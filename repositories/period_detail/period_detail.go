package period_detail

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/models/period_detail"
	"gorm.io/gorm"
)

type (
	IPeriodDetailRepository interface {
		Create(ctx context.Context, periodDetail *period_detail.PeriodDetail) error
		GetByID(ctx context.Context, id uint) (*period_detail.PeriodDetail, error)
		GetByPeriodAndUser(ctx context.Context, periodID, userID uint) (*period_detail.PeriodDetail, error)
		Update(ctx context.Context, id uint, updates map[string]interface{}) error
		Delete(ctx context.Context, id uint) error
		GetUsersByBatch(ctx context.Context, lastID uint, limit int) ([]uint, error)
		CreateBatch(ctx context.Context, periodDetails []period_detail.PeriodDetail) error
	}

	PeriodDetailRepository struct {
		db *gorm.DB
	}
)

func NewPeriodDetailRepository(db *gorm.DB) IPeriodDetailRepository {
	return &PeriodDetailRepository{
		db: db,
	}
}

func (r *PeriodDetailRepository) Create(ctx context.Context, periodDetail *period_detail.PeriodDetail) error {
	return r.db.WithContext(ctx).Create(periodDetail).Error
}

func (r *PeriodDetailRepository) GetByID(ctx context.Context, id uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (r *PeriodDetailRepository) GetByPeriodAndUser(ctx context.Context, periodID, userID uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	err := r.db.WithContext(ctx).Where("periods_id = ? AND user_id = ?", periodID, userID).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (r *PeriodDetailRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&period_detail.PeriodDetail{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PeriodDetailRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&period_detail.PeriodDetail{}, id).Error
}

func (r *PeriodDetailRepository) GetUsersByBatch(ctx context.Context, lastID uint, limit int) ([]uint, error) {
	var userIDs []uint
	query := r.db.WithContext(ctx).Table("users").Select("id").Where("role = ?", "employee")

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	err := query.Order("id ASC").Limit(limit).Pluck("id", &userIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by batch: %w", err)
	}

	return userIDs, nil
}

func (r *PeriodDetailRepository) CreateBatch(ctx context.Context, periodDetails []period_detail.PeriodDetail) error {
	if len(periodDetails) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&periodDetails).Error
}
