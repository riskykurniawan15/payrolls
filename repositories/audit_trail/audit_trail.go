package audit_trail

import (
	"context"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/models/audit_trail"
	"gorm.io/gorm"
)

type (
	IAuditTrailRepository interface {
		Create(ctx context.Context, auditTrail *audit_trail.AuditTrail) error
	}

	AuditTrailRepository struct {
		db *gorm.DB
	}
)

func NewAuditTrailRepository(db *gorm.DB) IAuditTrailRepository {
	return &AuditTrailRepository{
		db: db,
	}
}

func (repo AuditTrailRepository) getInstanceDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	if !ok {
		return repo.db
	}
	return tx
}

func (repo AuditTrailRepository) Create(ctx context.Context, auditTrail *audit_trail.AuditTrail) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(auditTrail).Error
}
