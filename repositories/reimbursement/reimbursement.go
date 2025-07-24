package reimbursement

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/models/reimbursement"
	"gorm.io/gorm"
)

type (
	IReimbursementRepository interface {
		Create(ctx context.Context, reimbursement *reimbursement.Reimbursement) error
		GetByID(ctx context.Context, id uint) (*reimbursement.Reimbursement, error)
		Update(ctx context.Context, id uint, updates map[string]interface{}) error
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, req reimbursement.ListReimbursementsRequest, userID uint) (*reimbursement.ListReimbursementsResponse, error)
		GetByUserAndDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]reimbursement.Reimbursement, error)
	}

	ReimbursementRepository struct {
		db *gorm.DB
	}
)

func NewReimbursementRepository(db *gorm.DB) IReimbursementRepository {
	return &ReimbursementRepository{db: db}
}

func (r *ReimbursementRepository) Create(ctx context.Context, reimbursement *reimbursement.Reimbursement) error {
	return r.db.WithContext(ctx).Create(reimbursement).Error
}

func (r *ReimbursementRepository) GetByID(ctx context.Context, id uint) (*reimbursement.Reimbursement, error) {
	var reimb reimbursement.Reimbursement
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&reimb).Error
	if err != nil {
		return nil, err
	}
	return &reimb, nil
}

func (r *ReimbursementRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&reimbursement.Reimbursement{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ReimbursementRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&reimbursement.Reimbursement{}, id).Error
}

func (r *ReimbursementRepository) List(ctx context.Context, req reimbursement.ListReimbursementsRequest, userID uint) (*reimbursement.ListReimbursementsResponse, error) {
	var reimbursements []reimbursement.Reimbursement
	var total int64

	// Build query - user can only see their own reimbursements
	query := r.db.WithContext(ctx).Model(&reimbursement.Reimbursement{}).Where("user_id = ?", userID)

	// Apply date range filter
	if req.StartDate != nil && req.EndDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			endDate, err := time.Parse("2006-01-02", *req.EndDate)
			if err == nil {
				// Set end date to end of day
				endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
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
	if err := query.Find(&reimbursements).Error; err != nil {
		return nil, err
	}

	// Convert to response
	var responses []reimbursement.ReimbursementResponse
	for _, reimb := range reimbursements {
		responses = append(responses, r.toResponse(reimb))
	}

	// Calculate pagination info
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &reimbursement.ListReimbursementsResponse{
		Data: responses,
		Pagination: reimbursement.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

// Helper function to convert Reimbursement to ReimbursementResponse
func (r *ReimbursementRepository) GetByUserAndDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]reimbursement.Reimbursement, error) {
	var reimbursements []reimbursement.Reimbursement
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&reimbursements).Error
	return reimbursements, err
}

func (r *ReimbursementRepository) toResponse(reimb reimbursement.Reimbursement) reimbursement.ReimbursementResponse {
	return reimbursement.ReimbursementResponse{
		ID:          reimb.ID,
		UserID:      reimb.UserID,
		Title:       reimb.Title,
		Date:        reimb.Date,
		Amount:      reimb.Amount,
		Description: reimb.Description,
		CreatedBy:   reimb.CreatedBy,
		CreatedAt:   reimb.CreatedAt,
		UpdatedBy:   reimb.UpdatedBy,
		UpdatedAt:   reimb.UpdatedAt,
	}
}
