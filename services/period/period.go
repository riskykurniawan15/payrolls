package period

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/models/period"
	periodRepo "github.com/riskykurniawan15/payrolls/repositories/period"
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
		periodRepo periodRepo.IPeriodRepository
	}
)

func NewPeriodService(periodRepo periodRepo.IPeriodRepository) IPeriodService {
	return &PeriodService{
		periodRepo: periodRepo,
	}
}

func (s *PeriodService) Create(ctx context.Context, req period.CreatePeriodRequest, userID uint) (*period.PeriodResponse, error) {
	// Generate code if not provided
	code := req.Code
	if code == nil || *code == "" {
		generatedCode, err := s.periodRepo.GenerateUniqueCode(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique code: %w", err)
		}
		code = &generatedCode
	} else {
		// Check if provided code already exists
		exists, err := s.periodRepo.IsCodeExists(ctx, *code)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("code '%s' already exists", *code)
		}
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD: %w", err)
	}

	// Parse end date and set to end of day (23:59:59)
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD: %w", err)
	}
	// Set end date to 23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Validate end date is after start date
	if !endDate.After(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	// Check for date conflicts with existing periods
	hasConflict, err := s.periodRepo.CheckDateConflict(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check date conflicts: %w", err)
	}
	if hasConflict {
		// Get conflicting periods for detailed error message
		conflictingPeriods, err := s.periodRepo.GetConflictingPeriods(ctx, startDate, endDate)
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

	if err := s.periodRepo.Create(ctx, period); err != nil {
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	// Convert to response
	response := s.toResponse(*period)
	return &response, nil
}

func (s *PeriodService) GetByID(ctx context.Context, id uint) (*period.PeriodResponse, error) {
	p, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Check if period is deleted (status = 9)
	if p.Status == constant.StatusDeleted {
		return nil, fmt.Errorf("period not found: (deleted)")
	}

	response := s.toResponse(*p)
	return &response, nil
}

func (s *PeriodService) Update(ctx context.Context, id uint, req period.UpdatePeriodRequest, userID uint) (*period.PeriodResponse, error) {
	// Check if period exists
	existingPeriod, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Check if period is deleted (status = 9)
	if existingPeriod.Status == constant.StatusDeleted {
		return nil, fmt.Errorf("cannot update period with status 9 (deleted)")
	}

	// Check if period has generated payroll
	if existingPeriod.UserExecutablePayroll != nil || existingPeriod.PayrollDate != nil {
		return nil, fmt.Errorf("cannot update period that has generated payroll")
	}

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
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD: %w", err)
		}
		newStartDate = startDate
		hasDateChanges = true
	} else {
		newStartDate = existingPeriod.StartDate
	}

	// Parse end date if provided
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD: %w", err)
		}
		// Set end date to 23:59:59
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
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
	// Check if period exists
	existingPeriod, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("period not found: %w", err)
	}

	// Check if period is already deleted (status = 9)
	if existingPeriod.Status == constant.StatusDeleted {
		return fmt.Errorf("period is already deleted")
	}

	// Check if period has generated payroll
	if existingPeriod.UserExecutablePayroll != nil || existingPeriod.PayrollDate != nil {
		return fmt.Errorf("cannot update period that has generated payroll")
	}

	// Soft delete
	if err := s.periodRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete period: %w", err)
	}

	return nil
}

func (s *PeriodService) List(ctx context.Context, req period.ListPeriodsRequest) (*period.ListPeriodsResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	response, err := s.periodRepo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list periods: %w", err)
	}

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
