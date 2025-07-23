package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/models/overtime"
	"gorm.io/gorm"
)

type (
	IOvertimeRepository interface {
		Create(ctx context.Context, overtime *overtime.Overtime) error
		GetByID(ctx context.Context, id uint) (*overtime.Overtime, error)
		Update(ctx context.Context, id uint, updates map[string]interface{}) error
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, req overtime.ListOvertimesRequest, userID uint) (*overtime.ListOvertimesResponse, error)

		GetTotalHoursByUserAndDate(ctx context.Context, userID uint, date time.Time) (float64, error)
	}

	OvertimeRepository struct {
		db *gorm.DB
	}
)

func NewOvertimeRepository(db *gorm.DB) IOvertimeRepository {
	return &OvertimeRepository{db: db}
}

func (r *OvertimeRepository) Create(ctx context.Context, overtime *overtime.Overtime) error {
	return r.db.WithContext(ctx).Create(overtime).Error
}

func (r *OvertimeRepository) GetByID(ctx context.Context, id uint) (*overtime.Overtime, error) {
	var o overtime.Overtime
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OvertimeRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&overtime.Overtime{}).Where("id = ?", id).Updates(updates).Error
}

func (r *OvertimeRepository) Delete(ctx context.Context, id uint) error {
	// Soft delete by setting status to 9 (if we add status field later)
	// For now, just delete the record
	return r.db.WithContext(ctx).Delete(&overtime.Overtime{}, id).Error
}

func (r *OvertimeRepository) List(ctx context.Context, req overtime.ListOvertimesRequest, userID uint) (*overtime.ListOvertimesResponse, error) {
	var overtimes []overtime.Overtime
	var total int64

	// Build query
	query := r.db.WithContext(ctx).Model(&overtime.Overtime{}).Where("user_id = ?", userID)

	// Apply date range filter
	if req.StartDate != nil && req.EndDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			endDate, err := time.Parse("2006-01-02", *req.EndDate)
			if err == nil {
				// Set end date to end of day
				endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				query = query.Where("overtimes_date BETWEEN ? AND ?", startDate, endDate)
			}
		}
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
	if err := query.Find(&overtimes).Error; err != nil {
		return nil, err
	}

	// Convert to response
	var responses []overtime.OvertimeResponse
	for _, o := range overtimes {
		responses = append(responses, r.toResponse(o))
	}

	// Calculate pagination info
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &overtime.ListOvertimesResponse{
		Data: responses,
		Pagination: overtime.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *OvertimeRepository) GetTotalHoursByUserAndDate(ctx context.Context, userID uint, date time.Time) (float64, error) {
	var totalHours float64
	// Get total hours for the specific date (from start of day to end of day)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.WithContext(ctx).
		Model(&overtime.Overtime{}).
		Where("user_id = ? AND overtimes_date BETWEEN ? AND ?", userID, startOfDay, endOfDay).
		Select("COALESCE(SUM(total_hours_time), 0)").
		Scan(&totalHours).Error

	return totalHours, err
}

// Helper function to convert Overtime to OvertimeResponse
func (r *OvertimeRepository) toResponse(o overtime.Overtime) overtime.OvertimeResponse {
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
