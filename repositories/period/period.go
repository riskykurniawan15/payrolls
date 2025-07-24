package period

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riskykurniawan15/payrolls/utils/code_generator"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/models/period"
	"gorm.io/gorm"
)

type (
	IPeriodRepository interface {
		Create(ctx context.Context, period *period.Period) error
		GetByID(ctx context.Context, id uint) (*period.Period, error)
		GetByCode(ctx context.Context, code string) (*period.Period, error)
		Update(ctx context.Context, id uint, updates map[string]interface{}) error
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, req period.ListPeriodsRequest) (*period.ListPeriodsResponse, error)
		IsCodeExists(ctx context.Context, code string, excludeID ...uint) (bool, error)
		CheckDateConflict(ctx context.Context, startDate, endDate time.Time, excludeID ...uint) (bool, error)
		GetConflictingPeriods(ctx context.Context, startDate, endDate time.Time, excludeID ...uint) ([]period.Period, error)
		GenerateUniqueCode(ctx context.Context) (string, error)
	}

	PeriodRepository struct {
		db *gorm.DB
	}
)

func NewPeriodRepository(db *gorm.DB) IPeriodRepository {
	return &PeriodRepository{db: db}
}

func (repo PeriodRepository) getInstanceDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	if !ok {
		return repo.db
	}
	return tx
}

func (repo PeriodRepository) Create(ctx context.Context, period *period.Period) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(period).Error
}

func (repo PeriodRepository) GetByID(ctx context.Context, id uint) (*period.Period, error) {
	var p period.Period

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo PeriodRepository) GetByCode(ctx context.Context, code string) (*period.Period, error) {
	var p period.Period

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("LOWER(code) = LOWER(?)", code).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo PeriodRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).Where("id = ?", id).Updates(updates).Error
}

func (repo PeriodRepository) Delete(ctx context.Context, id uint) error {
	// Soft delete by setting status to 9
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).Where("id = ?", id).Update("status", 9).Error
}

func (repo PeriodRepository) List(ctx context.Context, req period.ListPeriodsRequest) (*period.ListPeriodsResponse, error) {
	var periods []period.Period
	var total int64

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	// Build query
	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).Where("status != ?", constant.StatusDeleted)

	// Apply search filter
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", searchTerm, searchTerm)
	}

	// Apply status filter
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if req.SortBy != "" {
		sortOrder := "ASC"
		if req.SortDesc {
			sortOrder = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, sortOrder))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	if err := query.Find(&periods).Error; err != nil {
		return nil, err
	}

	// Convert to response
	var responses []period.PeriodResponse
	for _, p := range periods {
		responses = append(responses, repo.toResponse(p))
	}

	// Calculate pagination info
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &period.ListPeriodsResponse{
		Data: responses,
		Pagination: period.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (repo PeriodRepository) IsCodeExists(ctx context.Context, code string, excludeID ...uint) (bool, error) {
	var count int64
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).
		Where("LOWER(code) = LOWER(?) AND status != ?", code, constant.StatusDeleted)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CheckDateConflict checks if the given date range conflicts with existing periods
func (repo PeriodRepository) CheckDateConflict(ctx context.Context, startDate, endDate time.Time, excludeID ...uint) (bool, error) {
	var count int64

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).
		Where("status != ?", constant.StatusDeleted).
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
			startDate, startDate,
			endDate, endDate,
			startDate, endDate,
		)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetConflictingPeriods returns periods that conflict with the given date range
func (repo PeriodRepository) GetConflictingPeriods(ctx context.Context, startDate, endDate time.Time, excludeID ...uint) ([]period.Period, error) {
	var periods []period.Period

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period.Period{}).
		Where("status != ?", constant.StatusDeleted).
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
			startDate, startDate,
			endDate, endDate,
			startDate, endDate,
		)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Find(&periods).Error
	return periods, err
}

func (repo PeriodRepository) GenerateUniqueCode(ctx context.Context) (string, error) {
	// Try up to 10 times to generate a unique code
	for i := 0; i < 10; i++ {
		randomCode := code_generator.GeneratePeriodCode("PRD", 3)

		exists, err := repo.IsCodeExists(ctx, randomCode)
		if err != nil {
			return "", err
		}

		if !exists {
			return randomCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}

// Helper function to convert Period to PeriodResponse
func (repo PeriodRepository) toResponse(p period.Period) period.PeriodResponse {
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
