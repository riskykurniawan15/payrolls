package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/attendance"
	attendanceRepo "github.com/riskykurniawan15/payrolls/repositories/attendance"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IAttendanceService interface {
		GetAttendances(ctx context.Context, userID uint, page, limit int, startDate, endDate *time.Time) (attendance.AttendanceListResponse, error)
		GetAttendanceByID(ctx context.Context, id, userID uint) (attendance.AttendanceResponse, error)
		CheckIn(ctx context.Context, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error)
		CheckOut(ctx context.Context, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error)
		CheckOutByID(ctx context.Context, id, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error)
	}

	AttendanceService struct {
		attendanceRepo attendanceRepo.IAttendanceRepository
		logger         logger.Logger
	}
)

func NewAttendanceService(logger logger.Logger, attendanceRepo attendanceRepo.IAttendanceRepository) IAttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		logger:         logger,
	}
}

func (service *AttendanceService) GetAttendances(ctx context.Context, userID uint, page, limit int, startDate, endDate *time.Time) (attendance.AttendanceListResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	// Set default date range to current month if not provided
	// Only set default if both startDate and endDate are nil
	if startDate == nil && endDate == nil {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second) // Last second of current month

		startDate = &startOfMonth
		endDate = &endOfMonth
	}

	service.logger.InfoT("starting get attendances", requestID, map[string]interface{}{
		"user_id":    userID,
		"page":       page,
		"limit":      limit,
		"start_date": startDate,
		"end_date":   endDate,
	})

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get attendances from repository
	attendances, total, err := service.attendanceRepo.GetAttendances(ctx, userID, page, limit, startDate, endDate)
	if err != nil {
		service.logger.ErrorT("failed to get attendances from repository", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return attendance.AttendanceListResponse{}, err
	}

	// Convert to response format
	var responses []attendance.AttendanceResponse
	for _, att := range attendances {
		response := attendance.AttendanceResponse{
			ID:           att.ID,
			UserID:       att.UserID,
			CheckInDate:  att.CheckInDate,
			CheckOutDate: att.CheckOutDate,
			CreatedAt:    att.CreatedAt,
		}
		if att.UpdatedAt != nil {
			response.UpdatedAt = att.UpdatedAt
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	service.logger.InfoT("attendances retrieved successfully", requestID, map[string]interface{}{
		"user_id":     userID,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"count":       len(responses),
		"start_date":  startDate,
		"end_date":    endDate,
	})

	return attendance.AttendanceListResponse{
		Data: responses,
		Pagination: attendance.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (service *AttendanceService) GetAttendanceByID(ctx context.Context, id, userID uint) (attendance.AttendanceResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting get attendance by ID", requestID, map[string]interface{}{
		"attendance_id": id,
		"user_id":       userID,
	})

	// Get attendance from repository
	att, err := service.attendanceRepo.GetAttendanceByID(ctx, id, userID)
	if err != nil {
		service.logger.ErrorT("failed to get attendance by ID", requestID, map[string]interface{}{
			"attendance_id": id,
			"user_id":       userID,
			"error":         err.Error(),
		})
		return attendance.AttendanceResponse{}, err
	}

	service.logger.InfoT("attendance retrieved successfully", requestID, map[string]interface{}{
		"attendance_id":  id,
		"user_id":        userID,
		"check_in_date":  att.CheckInDate,
		"check_out_date": att.CheckOutDate,
	})

	response := attendance.AttendanceResponse{
		ID:           att.ID,
		UserID:       att.UserID,
		CheckInDate:  att.CheckInDate,
		CheckOutDate: att.CheckOutDate,
		CreatedAt:    att.CreatedAt,
	}
	if att.UpdatedAt != nil {
		response.UpdatedAt = att.UpdatedAt
	}
	return response, nil
}

func (service *AttendanceService) CheckIn(ctx context.Context, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting check-in process", requestID, map[string]interface{}{
		"user_id":  userID,
		"has_date": req.Date != nil,
	})

	// Determine check-in date
	checkInDate := time.Now()
	if req.Date != nil {
		checkInDate = req.Date.Time
	}

	service.logger.InfoT("check-in date determined", requestID, map[string]interface{}{
		"user_id":       userID,
		"check_in_date": checkInDate,
		"weekday":       checkInDate.Weekday(),
	})

	// Validate weekday (Monday to Friday)
	if !attendance.IsWeekday(checkInDate) {
		service.logger.WarningT("check-in attempted on weekend", requestID, map[string]interface{}{
			"user_id":       userID,
			"check_in_date": checkInDate,
			"weekday":       checkInDate.Weekday(),
		})
		return attendance.AttendanceResponse{}, errors.New("check-in is only allowed on weekdays (Monday to Friday)")
	}

	// Check if user already has an active check-in (no check-out)
	existingAttendance, err := service.attendanceRepo.GetLatestCheckInByUserID(ctx, userID, &checkInDate)
	if err == nil && existingAttendance.CheckOutDate == nil {
		service.logger.WarningT("user already has active check-in", requestID, map[string]interface{}{
			"user_id":           userID,
			"existing_check_in": existingAttendance.CheckInDate,
			"existing_id":       existingAttendance.ID,
		})
		return attendance.AttendanceResponse{}, errors.New("you already have an active check-in. Please check-out first")
	}

	// Create new attendance record
	attendanceData := attendance.Attendance{
		UserID:      userID,
		CheckInDate: checkInDate,
		CreatedBy:   userID,
	}

	createdAttendance, err := service.attendanceRepo.CreateAttendance(ctx, attendanceData)
	if err != nil {
		service.logger.ErrorT("failed to create attendance record", requestID, map[string]interface{}{
			"user_id":       userID,
			"check_in_date": checkInDate,
			"error":         err.Error(),
		})
		return attendance.AttendanceResponse{}, errors.New("failed to create attendance record")
	}

	service.logger.InfoT("check-in successful", requestID, map[string]interface{}{
		"attendance_id": createdAttendance.ID,
		"user_id":       userID,
		"check_in_date": checkInDate,
	})

	response := attendance.AttendanceResponse{
		ID:           createdAttendance.ID,
		UserID:       createdAttendance.UserID,
		CheckInDate:  createdAttendance.CheckInDate,
		CheckOutDate: createdAttendance.CheckOutDate,
		CreatedAt:    createdAttendance.CreatedAt,
	}
	if createdAttendance.UpdatedAt != nil {
		response.UpdatedAt = createdAttendance.UpdatedAt
	}
	return response, nil
}

func (service *AttendanceService) CheckOut(ctx context.Context, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting check-out process", requestID, map[string]interface{}{
		"user_id":  userID,
		"has_date": req.Date != nil,
	})

	// Get the latest active check-in for the user
	attendanceData, err := service.attendanceRepo.GetLatestCheckInByUserID(ctx, userID, nil)
	if err != nil {
		service.logger.ErrorT("no active check-in found for check-out", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return attendance.AttendanceResponse{}, errors.New("no active check-in found. Please check-in first")
	}

	// Determine check-out date
	checkOutDate := time.Now()
	if req.Date != nil {
		checkOutDate = req.Date.Time
	}

	service.logger.InfoT("check-out date determined", requestID, map[string]interface{}{
		"user_id":        userID,
		"attendance_id":  attendanceData.ID,
		"check_in_date":  attendanceData.CheckInDate,
		"check_out_date": checkOutDate,
	})

	// Validate check-out time
	if !attendance.IsValidCheckOutTime(attendanceData.CheckInDate, checkOutDate) {
		service.logger.WarningT("invalid check-out time", requestID, map[string]interface{}{
			"user_id":        userID,
			"attendance_id":  attendanceData.ID,
			"check_in_date":  attendanceData.CheckInDate,
			"check_out_date": checkOutDate,
		})
		return attendance.AttendanceResponse{}, errors.New("invalid check-out time. Check-out must be after check-in")
	}

	// Update attendance with check-out date
	attendanceData.CheckOutDate = &checkOutDate
	attendanceData.UpdatedBy = &userID
	now := time.Now()
	attendanceData.UpdatedAt = &now

	updatedAttendance, err := service.attendanceRepo.UpdateAttendance(ctx, attendanceData)
	if err != nil {
		service.logger.ErrorT("failed to update attendance with check-out", requestID, map[string]interface{}{
			"user_id":        userID,
			"attendance_id":  attendanceData.ID,
			"check_out_date": checkOutDate,
			"error":          err.Error(),
		})
		return attendance.AttendanceResponse{}, errors.New("failed to update attendance record")
	}

	service.logger.InfoT("check-out successful", requestID, map[string]interface{}{
		"attendance_id":  updatedAttendance.ID,
		"user_id":        userID,
		"check_in_date":  updatedAttendance.CheckInDate,
		"check_out_date": updatedAttendance.CheckOutDate,
	})

	response := attendance.AttendanceResponse{
		ID:           updatedAttendance.ID,
		UserID:       updatedAttendance.UserID,
		CheckInDate:  updatedAttendance.CheckInDate,
		CheckOutDate: updatedAttendance.CheckOutDate,
		CreatedAt:    updatedAttendance.CreatedAt,
	}
	if updatedAttendance.UpdatedAt != nil {
		response.UpdatedAt = updatedAttendance.UpdatedAt
	}
	return response, nil
}

func (service *AttendanceService) CheckOutByID(ctx context.Context, id, userID uint, req attendance.AttendanceRequest) (attendance.AttendanceResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting check-out by ID process", requestID, map[string]interface{}{
		"attendance_id": id,
		"user_id":       userID,
		"has_date":      req.Date != nil,
	})

	// Get the specific attendance record
	attendanceData, err := service.attendanceRepo.GetAttendanceByIDForUpdate(ctx, id, userID)
	if err != nil {
		service.logger.ErrorT("attendance not found for check-out", requestID, map[string]interface{}{
			"attendance_id": id,
			"user_id":       userID,
			"error":         err.Error(),
		})
		return attendance.AttendanceResponse{}, err
	}

	// Check if already checked out
	if attendanceData.CheckOutDate != nil {
		service.logger.WarningT("attendance already checked out", requestID, map[string]interface{}{
			"attendance_id":  id,
			"user_id":        userID,
			"check_out_date": attendanceData.CheckOutDate,
		})
		return attendance.AttendanceResponse{}, errors.New("attendance already checked out")
	}

	// Determine check-out date
	checkOutDate := time.Now()
	if req.Date != nil {
		checkOutDate = req.Date.Time
	}

	service.logger.InfoT("check-out date determined", requestID, map[string]interface{}{
		"attendance_id":  id,
		"user_id":        userID,
		"check_in_date":  attendanceData.CheckInDate,
		"check_out_date": checkOutDate,
	})

	// Validate check-out time
	if !attendance.IsValidCheckOutTime(attendanceData.CheckInDate, checkOutDate) {
		service.logger.WarningT("invalid check-out time", requestID, map[string]interface{}{
			"attendance_id":  id,
			"user_id":        userID,
			"check_in_date":  attendanceData.CheckInDate,
			"check_out_date": checkOutDate,
		})
		return attendance.AttendanceResponse{}, errors.New("invalid check-out time. Check-out must be after check-in")
	}

	// Update attendance with check-out date
	attendanceData.CheckOutDate = &checkOutDate
	attendanceData.UpdatedBy = &userID
	now := time.Now()
	attendanceData.UpdatedAt = &now

	updatedAttendance, err := service.attendanceRepo.UpdateAttendance(ctx, attendanceData)
	if err != nil {
		service.logger.ErrorT("failed to update attendance with check-out", requestID, map[string]interface{}{
			"attendance_id":  id,
			"user_id":        userID,
			"check_out_date": checkOutDate,
			"error":          err.Error(),
		})
		return attendance.AttendanceResponse{}, errors.New("failed to update attendance record")
	}

	service.logger.InfoT("check-out by ID successful", requestID, map[string]interface{}{
		"attendance_id":  updatedAttendance.ID,
		"user_id":        userID,
		"check_in_date":  updatedAttendance.CheckInDate,
		"check_out_date": updatedAttendance.CheckOutDate,
	})

	response := attendance.AttendanceResponse{
		ID:           updatedAttendance.ID,
		UserID:       updatedAttendance.UserID,
		CheckInDate:  updatedAttendance.CheckInDate,
		CheckOutDate: updatedAttendance.CheckOutDate,
		CreatedAt:    updatedAttendance.CreatedAt,
	}
	if updatedAttendance.UpdatedAt != nil {
		response.UpdatedAt = updatedAttendance.UpdatedAt
	}
	return response, nil
}
