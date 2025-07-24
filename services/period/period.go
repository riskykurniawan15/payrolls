package period

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/period"
	periodRepo "github.com/riskykurniawan15/payrolls/repositories/period"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IPeriodService interface {
		Create(ctx context.Context, req period.CreatePeriodRequest, userID uint) (*period.PeriodResponse, error)
		GetByID(ctx context.Context, id uint) (*period.PeriodResponse, error)
		Update(ctx context.Context, id uint, req period.UpdatePeriodRequest, userID uint) (*period.PeriodResponse, error)
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, req period.ListPeriodsRequest) (*period.ListPeriodsResponse, error)
	}

	PeriodService struct {
		logger     logger.Logger
		periodRepo periodRepo.IPeriodRepository
	}
)

func NewPeriodService(logger logger.Logger, periodRepo periodRepo.IPeriodRepository) IPeriodService {
	return &PeriodService{
		logger:     logger,
		periodRepo: periodRepo,
	}
}

func (s *PeriodService) Create(ctx context.Context, req period.CreatePeriodRequest, userID uint) (*period.PeriodResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing create period request", requestID, map[string]interface{}{
		"user_id": userID,
		"name":    req.Name,
		"code":    req.Code,
	})

	// Generate code if not provided
	code := req.Code
	if code == nil || *code == "" {
		s.logger.InfoT("generating unique code", requestID, map[string]interface{}{
			"name": req.Name,
		})
		generatedCode, err := s.periodRepo.GenerateUniqueCode(ctx)
		if err != nil {
			s.logger.ErrorT("failed to generate unique code", requestID, map[string]interface{}{
				"error": err.Error(),
				"name":  req.Name,
			})
			return nil, fmt.Errorf("failed to generate unique code: %w", err)
		}
		code = &generatedCode
		s.logger.InfoT("code generated successfully", requestID, map[string]interface{}{
			"generated_code": *code,
		})
	} else {
		// Check if provided code already exists
		s.logger.InfoT("checking code existence", requestID, map[string]interface{}{
			"code": *code,
		})
		exists, err := s.periodRepo.IsCodeExists(ctx, *code)
		if err != nil {
			s.logger.ErrorT("failed to check code existence", requestID, map[string]interface{}{
				"error": err.Error(),
				"code":  *code,
			})
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			s.logger.WarningT("code already exists", requestID, map[string]interface{}{
				"code": *code,
			})
			return nil, fmt.Errorf("code '%s' already exists", *code)
		}
		s.logger.InfoT("code is unique", requestID, map[string]interface{}{
			"code": *code,
		})
	}

	// Parse start and end date
	if req.StartDate == nil {
		return nil, fmt.Errorf("start_date is required")
	}
	if req.EndDate == nil {
		return nil, fmt.Errorf("end_date is required")
	}
	startDate := req.StartDate.Time
	endDate := req.EndDate.Time

	// Set end date to 23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	s.logger.InfoT("parsed dates", requestID, map[string]interface{}{
		"start_date": startDate.Format("2006-01-02 15:04:05"),
		"end_date":   endDate.Format("2006-01-02 15:04:05"),
	})

	// Validate end date is after start date
	if !endDate.After(startDate) {
		s.logger.WarningT("end_date is not after start_date", requestID, map[string]interface{}{
			"start_date": startDate.Format("2006-01-02 15:04:05"),
			"end_date":   endDate.Format("2006-01-02 15:04:05"),
		})
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	// Check for date conflicts with existing periods
	s.logger.InfoT("checking date conflicts", requestID, map[string]interface{}{
		"start_date": startDate.Format("2006-01-02 15:04:05"),
		"end_date":   endDate.Format("2006-01-02 15:04:05"),
	})
	hasConflict, err := s.periodRepo.CheckDateConflict(ctx, startDate, endDate)
	if err != nil {
		s.logger.ErrorT("failed to check date conflicts", requestID, map[string]interface{}{
			"error":      err.Error(),
			"start_date": startDate.Format("2006-01-02 15:04:05"),
			"end_date":   endDate.Format("2006-01-02 15:04:05"),
		})
		return nil, fmt.Errorf("failed to check date conflicts: %w", err)
	}
	if hasConflict {
		// Get conflicting periods for detailed error message
		conflictingPeriods, err := s.periodRepo.GetConflictingPeriods(ctx, startDate, endDate)
		if err != nil {
			s.logger.ErrorT("failed to get conflicting periods", requestID, map[string]interface{}{
				"error":      err.Error(),
				"start_date": startDate.Format("2006-01-02 15:04:05"),
				"end_date":   endDate.Format("2006-01-02 15:04:05"),
			})
			return nil, fmt.Errorf("date range conflicts with existing periods")
		}

		// Build detailed error message
		var conflictDetails []string
		for _, p := range conflictingPeriods {
			conflictDetails = append(conflictDetails, fmt.Sprintf("Period '%s' (%s to %s)", p.Name, p.StartDate.Format("2006-01-02"), p.EndDate.Format("2006-01-02")))
		}

		s.logger.WarningT("date conflict detected", requestID, map[string]interface{}{
			"start_date":        startDate.Format("2006-01-02 15:04:05"),
			"end_date":          endDate.Format("2006-01-02 15:04:05"),
			"conflicting_count": len(conflictingPeriods),
			"conflicts":         conflictDetails,
		})
		return nil, fmt.Errorf("date range conflicts with existing periods: %s", strings.Join(conflictDetails, ", "))
	}

	s.logger.InfoT("no date conflicts found", requestID, map[string]interface{}{
		"start_date": startDate.Format("2006-01-02 15:04:05"),
		"end_date":   endDate.Format("2006-01-02 15:04:05"),
	})

	// Create period with status always 1
	period := &period.Period{
		Code:      *code,
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    constant.StatusActive, // Always active for new periods
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedBy: nil, // Explicitly set to nil
		UpdatedAt: nil, // Explicitly set to nil
	}

	s.logger.InfoT("creating period in database", requestID, map[string]interface{}{
		"code":       *code,
		"name":       req.Name,
		"start_date": startDate.Format("2006-01-02 15:04:05"),
		"end_date":   endDate.Format("2006-01-02 15:04:05"),
		"created_by": userID,
	})

	if err := s.periodRepo.Create(ctx, period); err != nil {
		s.logger.ErrorT("failed to create period", requestID, map[string]interface{}{
			"error": err.Error(),
			"code":  *code,
			"name":  req.Name,
		})
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	s.logger.InfoT("period created successfully", requestID, map[string]interface{}{
		"period_id": period.ID,
		"code":      *code,
		"name":      req.Name,
	})

	// Convert to response
	response := s.toResponse(*period)
	return &response, nil
}

func (s *PeriodService) GetByID(ctx context.Context, id uint) (*period.PeriodResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing get period by ID request", requestID, map[string]interface{}{
		"period_id": id,
	})

	p, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get period by ID", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": id,
		})
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Check if period is deleted (status = 9)
	if p.Status == constant.StatusDeleted {
		s.logger.WarningT("period is deleted", requestID, map[string]interface{}{
			"period_id": id,
			"status":    p.Status,
		})
		return nil, fmt.Errorf("period not found: (deleted)")
	}

	s.logger.InfoT("period retrieved successfully", requestID, map[string]interface{}{
		"period_id": id,
		"code":      p.Code,
		"name":      p.Name,
		"status":    p.Status,
	})

	response := s.toResponse(*p)
	return &response, nil
}

func (s *PeriodService) Update(ctx context.Context, id uint, req period.UpdatePeriodRequest, userID uint) (*period.PeriodResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing update period request", requestID, map[string]interface{}{
		"period_id": id,
		"user_id":   userID,
		"name":      req.Name,
		"code":      req.Code,
	})

	// Check if period exists
	existingPeriod, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get period for update", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": id,
		})
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Check if period is deleted (status = 9)
	if existingPeriod.Status == constant.StatusDeleted {
		s.logger.WarningT("cannot update deleted period", requestID, map[string]interface{}{
			"period_id": id,
			"status":    existingPeriod.Status,
		})
		return nil, fmt.Errorf("cannot update period with status 9 (deleted)")
	}

	// Check if period has generated payroll
	if existingPeriod.UserExecutablePayroll != nil || existingPeriod.PayrollDate != nil {
		s.logger.WarningT("cannot update period with generated payroll", requestID, map[string]interface{}{
			"period_id":               id,
			"user_executable_payroll": existingPeriod.UserExecutablePayroll,
			"payroll_date":            existingPeriod.PayrollDate,
		})
		return nil, fmt.Errorf("cannot update period that has generated payroll")
	}

	s.logger.InfoT("period validation passed", requestID, map[string]interface{}{
		"period_id": id,
		"code":      existingPeriod.Code,
		"name":      existingPeriod.Name,
		"status":    existingPeriod.Status,
	})

	// Prepare updates
	updates := make(map[string]interface{})
	updates["updated_by"] = userID
	updates["updated_at"] = time.Now()

	// Update code if provided
	if req.Code != nil && *req.Code != "" {
		// Check if new code already exists (excluding current period)
		exists, err := s.periodRepo.IsCodeExists(ctx, *req.Code, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("code '%s' already exists", *req.Code)
		}
		updates["code"] = *req.Code
	}

	// Update name if provided
	if req.Name != nil {
		updates["name"] = *req.Name
	}

	// Handle date updates and conflict checking
	var newStartDate, newEndDate time.Time
	var hasDateChanges bool

	// Parse start date if provided
	if req.StartDate != nil {
		newStartDate = req.StartDate.Time
		hasDateChanges = true
	} else {
		newStartDate = existingPeriod.StartDate
	}

	// Parse end date if provided
	if req.EndDate != nil {
		// Set end date to 23:59:59
		endDate := req.EndDate.Time.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		newEndDate = endDate
		hasDateChanges = true
	} else {
		newEndDate = existingPeriod.EndDate
	}

	// Check for date conflicts if dates are being changed
	if hasDateChanges {
		// Validate end date is after start date
		if !newEndDate.After(newStartDate) {
			return nil, fmt.Errorf("end_date must be after start_date")
		}

		// Check for date conflicts with existing periods (excluding current period)
		hasConflict, err := s.periodRepo.CheckDateConflict(ctx, newStartDate, newEndDate, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check date conflicts: %w", err)
		}
		if hasConflict {
			// Get conflicting periods for detailed error message
			conflictingPeriods, err := s.periodRepo.GetConflictingPeriods(ctx, newStartDate, newEndDate, id)
			if err != nil {
				return nil, fmt.Errorf("date range conflicts with existing periods")
			}

			// Build detailed error message
			var conflictDetails []string
			for _, p := range conflictingPeriods {
				conflictDetails = append(conflictDetails, fmt.Sprintf("Period '%s' (%s to %s)", p.Name, p.StartDate.Format("2006-01-02"), p.EndDate.Format("2006-01-02")))
			}

			return nil, fmt.Errorf("date range conflicts with existing periods: %s", strings.Join(conflictDetails, ", "))
		}

		// Add date updates to the updates map
		if req.StartDate != nil {
			updates["start_date"] = newStartDate
		}
		if req.EndDate != nil {
			updates["end_date"] = newEndDate
		}
	}

	// Apply updates
	if err := s.periodRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update period: %w", err)
	}

	// Get updated period
	updatedPeriod, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated period: %w", err)
	}

	response := s.toResponse(*updatedPeriod)
	return &response, nil
}

func (s *PeriodService) Delete(ctx context.Context, id uint) error {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing delete period request", requestID, map[string]interface{}{
		"period_id": id,
	})

	// Check if period exists
	existingPeriod, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get period for delete", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": id,
		})
		return fmt.Errorf("period not found: %w", err)
	}

	// Check if period is already deleted (status = 9)
	if existingPeriod.Status == constant.StatusDeleted {
		s.logger.WarningT("period is already deleted", requestID, map[string]interface{}{
			"period_id": id,
			"status":    existingPeriod.Status,
		})
		return fmt.Errorf("period is already deleted")
	}

	// Check if period has generated payroll
	if existingPeriod.UserExecutablePayroll != nil || existingPeriod.PayrollDate != nil {
		s.logger.WarningT("cannot delete period with generated payroll", requestID, map[string]interface{}{
			"period_id":               id,
			"user_executable_payroll": existingPeriod.UserExecutablePayroll,
			"payroll_date":            existingPeriod.PayrollDate,
		})
		return fmt.Errorf("cannot delete period that has generated payroll")
	}

	s.logger.InfoT("period validation passed for delete", requestID, map[string]interface{}{
		"period_id": id,
		"code":      existingPeriod.Code,
		"name":      existingPeriod.Name,
		"status":    existingPeriod.Status,
	})

	// Soft delete by setting status to 9
	if err := s.periodRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorT("failed to delete period", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": id,
		})
		return fmt.Errorf("failed to delete period: %w", err)
	}

	s.logger.InfoT("period deleted successfully", requestID, map[string]interface{}{
		"period_id": id,
		"code":      existingPeriod.Code,
		"name":      existingPeriod.Name,
	})

	return nil
}

func (s *PeriodService) List(ctx context.Context, req period.ListPeriodsRequest) (*period.ListPeriodsResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing list periods request", requestID, map[string]interface{}{
		"page":      req.Page,
		"limit":     req.Limit,
		"search":    req.Search,
		"status":    req.Status,
		"sort_by":   req.SortBy,
		"sort_desc": req.SortDesc,
	})

	// Set default values if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	s.logger.InfoT("list periods with filters", requestID, map[string]interface{}{
		"page":      req.Page,
		"limit":     req.Limit,
		"search":    req.Search,
		"status":    req.Status,
		"sort_by":   req.SortBy,
		"sort_desc": req.SortDesc,
	})

	response, err := s.periodRepo.List(ctx, req)
	if err != nil {
		s.logger.ErrorT("failed to list periods", requestID, map[string]interface{}{
			"error": err.Error(),
			"page":  req.Page,
			"limit": req.Limit,
		})
		return nil, fmt.Errorf("failed to list periods: %w", err)
	}

	s.logger.InfoT("periods listed successfully", requestID, map[string]interface{}{
		"total_count":  response.Pagination.Total,
		"total_pages":  response.Pagination.TotalPages,
		"current_page": response.Pagination.Page,
		"limit":        response.Pagination.Limit,
		"data_count":   len(response.Data),
	})

	return response, nil
}

// Helper function to convert Period to PeriodResponse
func (s *PeriodService) toResponse(p period.Period) period.PeriodResponse {
	return period.PeriodResponse{
		ID:                    p.ID,
		Code:                  p.Code,
		Name:                  p.Name,
		StartDate:             p.StartDate,
		EndDate:               p.EndDate,
		Status:                p.Status,
		UserExecutablePayroll: p.UserExecutablePayroll,
		PayrollDate:           p.PayrollDate,
		CreatedBy:             p.CreatedBy,
		CreatedAt:             p.CreatedAt,
		UpdatedBy:             p.UpdatedBy,
		UpdatedAt:             p.UpdatedAt,
	}
}
