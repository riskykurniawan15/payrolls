package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/attendance"
	"github.com/riskykurniawan15/payrolls/models/overtime"
	attendanceRepo "github.com/riskykurniawan15/payrolls/repositories/attendance"
	overtimeRepo "github.com/riskykurniawan15/payrolls/repositories/overtime"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IOvertimeService interface {
		Create(ctx context.Context, req overtime.CreateOvertimeRequest, userID uint) (*overtime.OvertimeResponse, error)
		GetByID(ctx context.Context, id uint, userID uint) (*overtime.OvertimeResponse, error)
		Update(ctx context.Context, id uint, req overtime.UpdateOvertimeRequest, userID uint) (*overtime.OvertimeResponse, error)
		Delete(ctx context.Context, id uint, userID uint) error
		List(ctx context.Context, req overtime.ListOvertimesRequest, userID uint) (*overtime.ListOvertimesResponse, error)
	}

	OvertimeService struct {
		logger         logger.Logger
		overtimeRepo   overtimeRepo.IOvertimeRepository
		attendanceRepo attendanceRepo.IAttendanceRepository
	}
)

func NewOvertimeService(logger logger.Logger, overtimeRepo overtimeRepo.IOvertimeRepository, attendanceRepo attendanceRepo.IAttendanceRepository) IOvertimeService {
	return &OvertimeService{
		logger:         logger,
		overtimeRepo:   overtimeRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (s *OvertimeService) Create(ctx context.Context, req overtime.CreateOvertimeRequest, userID uint) (*overtime.OvertimeResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing create overtime request", requestID, map[string]interface{}{
		"user_id":          userID,
		"overtimes_date":   req.OvertimesDate,
		"total_hours_time": req.TotalHoursTime,
		"created_by":       userID,
	})

	// Parse overtime date
	if req.OvertimesDate == nil {
		return nil, fmt.Errorf("overtimes_date is required")
	}
	overtimeDate := req.OvertimesDate.Time

	// Check if date is in the future
	if overtimeDate.After(time.Now()) {
		s.logger.WarningT("overtime date is in the future", requestID, map[string]interface{}{
			"overtimes_date": overtimeDate.Format("2006-01-02"),
		})
		return nil, fmt.Errorf("overtime date cannot be in the future")
	}

	// Check if it's a weekday
	if !attendance.IsWeekday(overtimeDate) {
		s.logger.InfoT("overtime on weekend is allowed", requestID, map[string]interface{}{
			"overtimes_date": overtimeDate.Format("2006-01-02"),
			"weekday":        overtimeDate.Weekday().String(),
		})
	} else {
		// For weekdays, check if user has checked out attendance
		s.logger.InfoT("checking attendance checkout for weekday", requestID, map[string]interface{}{
			"user_id":        userID,
			"overtimes_date": overtimeDate.Format("2006-01-02"),
		})

		// Get attendance for the specific date
		attendanceRecord, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, overtimeDate)
		if err != nil {
			if err.Error() == "record not found" {
				s.logger.WarningT("no attendance record found for weekday", requestID, map[string]interface{}{
					"user_id":        userID,
					"overtimes_date": overtimeDate.Format("2006-01-02"),
				})
				return nil, fmt.Errorf("attendance record not found for %s. Please check in and check out first", overtimeDate.Format("2006-01-02"))
			}
			s.logger.ErrorT("failed to get attendance record", requestID, map[string]interface{}{
				"error":          err.Error(),
				"user_id":        userID,
				"overtimes_date": overtimeDate.Format("2006-01-02"),
			})
			return nil, fmt.Errorf("failed to check attendance: %w", err)
		}

		// Check if user has checked out
		if attendanceRecord.CheckOutDate == nil {
			s.logger.WarningT("user has not checked out for weekday", requestID, map[string]interface{}{
				"user_id":        userID,
				"overtimes_date": overtimeDate.Format("2006-01-02"),
				"check_in_date":  attendanceRecord.CheckInDate.Format("2006-01-02 15:04:05"),
			})
			return nil, fmt.Errorf("please check out from attendance first before creating overtime for %s", overtimeDate.Format("2006-01-02"))
		}

		s.logger.InfoT("attendance checkout verified", requestID, map[string]interface{}{
			"user_id":        userID,
			"overtimes_date": overtimeDate.Format("2006-01-02"),
			"check_in_date":  attendanceRecord.CheckInDate.Format("2006-01-02 15:04:05"),
			"check_out_date": attendanceRecord.CheckOutDate.Format("2006-01-02 15:04:05"),
		})
	}

	// Check total hours limit (max 3 hours per day)
	totalHours, err := s.overtimeRepo.GetTotalHoursByUserAndDate(ctx, userID, overtimeDate)
	if err != nil {
		s.logger.ErrorT("failed to get total hours", requestID, map[string]interface{}{
			"error":          err.Error(),
			"user_id":        userID,
			"overtimes_date": overtimeDate.Format("2006-01-02"),
		})
		return nil, fmt.Errorf("failed to check total hours: %w", err)
	}

	newTotalHours := totalHours + req.TotalHoursTime
	if newTotalHours > 3.0 {
		s.logger.WarningT("total hours exceed limit", requestID, map[string]interface{}{
			"user_id":         userID,
			"overtimes_date":  overtimeDate.Format("2006-01-02"),
			"current_total":   totalHours,
			"requested_hours": req.TotalHoursTime,
			"new_total":       newTotalHours,
			"max_allowed":     3.0,
		})
		return nil, fmt.Errorf("total overtime hours for %s would exceed 3 hours limit (current: %.2f, requested: %.2f, total: %.2f)",
			overtimeDate.Format("2006-01-02"), totalHours, req.TotalHoursTime, newTotalHours)
	}

	s.logger.InfoT("overtime validation passed", requestID, map[string]interface{}{
		"user_id":          userID,
		"overtimes_date":   overtimeDate.Format("2006-01-02"),
		"total_hours_time": req.TotalHoursTime,
		"current_total":    totalHours,
		"new_total":        newTotalHours,
	})

	// Create overtime
	overtimeRecord := &overtime.Overtime{
		UserID:         userID,
		OvertimesDate:  overtimeDate,
		TotalHoursTime: req.TotalHoursTime,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedBy:      nil,
		UpdatedAt:      nil,
	}

	s.logger.InfoT("creating overtime in database", requestID, map[string]interface{}{
		"user_id":          userID,
		"overtimes_date":   overtimeDate.Format("2006-01-02"),
		"total_hours_time": req.TotalHoursTime,
		"created_by":       userID,
	})

	if err := s.overtimeRepo.Create(ctx, overtimeRecord); err != nil {
		s.logger.ErrorT("failed to create overtime", requestID, map[string]interface{}{
			"error":            err.Error(),
			"user_id":          userID,
			"overtimes_date":   overtimeDate.Format("2006-01-02"),
			"total_hours_time": req.TotalHoursTime,
		})
		return nil, fmt.Errorf("failed to create overtime: %w", err)
	}

	s.logger.InfoT("overtime created successfully", requestID, map[string]interface{}{
		"overtime_id":      overtimeRecord.ID,
		"user_id":          userID,
		"overtimes_date":   overtimeDate.Format("2006-01-02"),
		"total_hours_time": req.TotalHoursTime,
	})

	// Convert to response
	response := s.toResponse(*overtimeRecord)
	return &response, nil
}

func (s *OvertimeService) GetByID(ctx context.Context, id uint, userID uint) (*overtime.OvertimeResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing get overtime by ID request", requestID, map[string]interface{}{
		"overtime_id": id,
	})

	o, err := s.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get overtime by ID", requestID, map[string]interface{}{
			"error":       err.Error(),
			"overtime_id": id,
		})
		return nil, fmt.Errorf("overtime not found: %w", err)
	}

	// Check if overtime belongs to the user
	if o.UserID != userID {
		s.logger.WarningT("user trying to access overtime of another user", requestID, map[string]interface{}{
			"overtime_id":      id,
			"overtime_user_id": o.UserID,
			"request_user_id":  userID,
		})
		return nil, fmt.Errorf("overtime not found")
	}

	s.logger.InfoT("overtime retrieved successfully", requestID, map[string]interface{}{
		"overtime_id":      id,
		"user_id":          o.UserID,
		"overtimes_date":   o.OvertimesDate.Format("2006-01-02"),
		"total_hours_time": o.TotalHoursTime,
	})

	response := s.toResponse(*o)
	return &response, nil
}

func (s *OvertimeService) Update(ctx context.Context, id uint, req overtime.UpdateOvertimeRequest, userID uint) (*overtime.OvertimeResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing update overtime request", requestID, map[string]interface{}{
		"overtime_id": id,
		"user_id":     userID,
	})

	// Check if overtime exists
	existingOvertime, err := s.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get overtime for update", requestID, map[string]interface{}{
			"error":       err.Error(),
			"overtime_id": id,
		})
		return nil, fmt.Errorf("overtime not found: %w", err)
	}

	// Check if overtime belongs to the user
	if existingOvertime.UserID != userID {
		s.logger.WarningT("user trying to update overtime of another user", requestID, map[string]interface{}{
			"overtime_id":      id,
			"overtime_user_id": existingOvertime.UserID,
			"request_user_id":  userID,
		})
		return nil, fmt.Errorf("overtime not found")
	}

	s.logger.InfoT("overtime validation passed", requestID, map[string]interface{}{
		"overtime_id":      id,
		"user_id":          existingOvertime.UserID,
		"overtimes_date":   existingOvertime.OvertimesDate.Format("2006-01-02"),
		"total_hours_time": existingOvertime.TotalHoursTime,
	})

	// Prepare updates
	updates := make(map[string]interface{})
	updates["updated_by"] = userID
	updates["updated_at"] = time.Now()

	// Update overtimes_date if provided
	if req.OvertimesDate != nil {
		overtimeDate := req.OvertimesDate.Time

		// Check if date is in the future
		if overtimeDate.After(time.Now()) {
			return nil, fmt.Errorf("overtime date cannot be in the future")
		}

		// Check if it's a weekday and validate attendance
		if attendance.IsWeekday(overtimeDate) {
			// Check if user has checked out attendance for the new date
			attendanceRecord, err := s.attendanceRepo.GetByUserAndDate(ctx, existingOvertime.UserID, overtimeDate)
			if err != nil {
				if err.Error() == "record not found" {
					return nil, fmt.Errorf("attendance record not found for %s. Please check in and check out first", overtimeDate.Format("2006-01-02"))
				}
				return nil, fmt.Errorf("failed to check attendance: %w", err)
			}

			if attendanceRecord.CheckOutDate == nil {
				return nil, fmt.Errorf("please check out from attendance first before updating overtime for %s", overtimeDate.Format("2006-01-02"))
			}
		}

		// Allow multiple overtime records per day, no duplicate check needed

		updates["overtimes_date"] = overtimeDate
	}

	// Update total_hours_time if provided
	if req.TotalHoursTime != nil {
		// Check total hours limit (max 3 hours per day)
		overtimeDate := existingOvertime.OvertimesDate
		if req.OvertimesDate != nil {
			overtimeDate = req.OvertimesDate.Time
		}

		totalHours, err := s.overtimeRepo.GetTotalHoursByUserAndDate(ctx, existingOvertime.UserID, overtimeDate)
		if err != nil {
			return nil, fmt.Errorf("failed to check total hours: %w", err)
		}

		// Subtract current overtime hours and add new hours
		newTotalHours := totalHours - existingOvertime.TotalHoursTime + *req.TotalHoursTime
		if newTotalHours > 3.0 {
			return nil, fmt.Errorf("total overtime hours for %s would exceed 3 hours limit (current: %.2f, requested: %.2f, total: %.2f)",
				overtimeDate.Format("2006-01-02"), totalHours-existingOvertime.TotalHoursTime, *req.TotalHoursTime, newTotalHours)
		}

		updates["total_hours_time"] = *req.TotalHoursTime
	}

	// Apply updates
	if err := s.overtimeRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update overtime: %w", err)
	}

	// Get updated overtime
	updatedOvertime, err := s.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated overtime: %w", err)
	}

	response := s.toResponse(*updatedOvertime)
	return &response, nil
}

func (s *OvertimeService) Delete(ctx context.Context, id uint, userID uint) error {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing delete overtime request", requestID, map[string]interface{}{
		"overtime_id": id,
	})

	// Check if overtime exists
	existingOvertime, err := s.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get overtime for delete", requestID, map[string]interface{}{
			"error":       err.Error(),
			"overtime_id": id,
		})
		return fmt.Errorf("overtime not found: %w", err)
	}

	// Check if overtime belongs to the user
	if existingOvertime.UserID != userID {
		s.logger.WarningT("user trying to delete overtime of another user", requestID, map[string]interface{}{
			"overtime_id":      id,
			"overtime_user_id": existingOvertime.UserID,
			"request_user_id":  userID,
		})
		return fmt.Errorf("overtime not found")
	}

	s.logger.InfoT("overtime validation passed for delete", requestID, map[string]interface{}{
		"overtime_id":      id,
		"user_id":          existingOvertime.UserID,
		"overtimes_date":   existingOvertime.OvertimesDate.Format("2006-01-02"),
		"total_hours_time": existingOvertime.TotalHoursTime,
	})

	// Delete overtime
	if err := s.overtimeRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorT("failed to delete overtime", requestID, map[string]interface{}{
			"error":       err.Error(),
			"overtime_id": id,
		})
		return fmt.Errorf("failed to delete overtime: %w", err)
	}

	s.logger.InfoT("overtime deleted successfully", requestID, map[string]interface{}{
		"overtime_id":      id,
		"user_id":          existingOvertime.UserID,
		"overtimes_date":   existingOvertime.OvertimesDate.Format("2006-01-02"),
		"total_hours_time": existingOvertime.TotalHoursTime,
	})

	return nil
}

func (s *OvertimeService) List(ctx context.Context, req overtime.ListOvertimesRequest, userID uint) (*overtime.ListOvertimesResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing list overtimes request", requestID, map[string]interface{}{
		"page":       req.Page,
		"limit":      req.Limit,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"sort_by":    req.SortBy,
		"sort_desc":  req.SortDesc,
	})

	// Set default values if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	s.logger.InfoT("list overtimes with filters", requestID, map[string]interface{}{
		"page":       req.Page,
		"limit":      req.Limit,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"sort_by":    req.SortBy,
		"sort_desc":  req.SortDesc,
	})

	response, err := s.overtimeRepo.List(ctx, req, userID)
	if err != nil {
		s.logger.ErrorT("failed to list overtimes", requestID, map[string]interface{}{
			"error": err.Error(),
			"page":  req.Page,
			"limit": req.Limit,
		})
		return nil, fmt.Errorf("failed to list overtimes: %w", err)
	}

	s.logger.InfoT("overtimes listed successfully", requestID, map[string]interface{}{
		"total_count":  response.Pagination.Total,
		"total_pages":  response.Pagination.TotalPages,
		"current_page": response.Pagination.Page,
		"limit":        response.Pagination.Limit,
		"data_count":   len(response.Data),
	})

	return response, nil
}

// Helper function to convert Overtime to OvertimeResponse
func (s *OvertimeService) toResponse(o overtime.Overtime) overtime.OvertimeResponse {
	return overtime.OvertimeResponse{
		ID:             o.ID,
		UserID:         o.UserID,
		OvertimesDate:  o.OvertimesDate,
		TotalHoursTime: o.TotalHoursTime,
		CreatedBy:      o.CreatedBy,
		CreatedAt:      o.CreatedAt,
		UpdatedBy:      o.UpdatedBy,
		UpdatedAt:      o.UpdatedAt,
	}
}
