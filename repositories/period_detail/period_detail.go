package period_detail

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/constant"
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

func (repo PeriodDetailRepository) getInstanceDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	if !ok {
		return repo.db
	}
	return tx
}

func (repo PeriodDetailRepository) Create(ctx context.Context, periodDetail *period_detail.PeriodDetail) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(periodDetail).Error
}

func (repo PeriodDetailRepository) GetByID(ctx context.Context, id uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("id = ?", id).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (repo PeriodDetailRepository) GetByPeriodAndUser(ctx context.Context, periodID, userID uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("periods_id = ? AND user_id = ?", periodID, userID).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (repo PeriodDetailRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period_detail.PeriodDetail{}).Where("id = ?", id).Updates(updates).Error
}

func (repo PeriodDetailRepository) Delete(ctx context.Context, id uint) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Delete(&period_detail.PeriodDetail{}, id).Error
}

func (repo PeriodDetailRepository) GetUsersByBatch(ctx context.Context, lastID uint, limit int) ([]uint, error) {
	var userIDs []uint
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Table("users").Select("id").Where("role = ?", "employee")

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	err := query.Order("id ASC").Limit(limit).Pluck("id", &userIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by batch: %w", err)
	}

	return userIDs, nil
}

func (repo PeriodDetailRepository) CreateBatch(ctx context.Context, periodDetails []period_detail.PeriodDetail) error {
	if len(periodDetails) == 0 {
		return nil
	}

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(&periodDetails).Error
}
