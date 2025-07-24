package user

import (
	"context"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/models/user"
	"gorm.io/gorm"
)

type (
	IUserRepository interface {
		GetUserByUsername(ctx context.Context, username string) (user.User, error)
		GetUserByID(ctx context.Context, id uint) (user.User, error)
		CreateUser(ctx context.Context, user user.User) (user.User, error)
	}

	UserRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo UserRepository) getInstanceDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	if !ok {
		return repo.getInstanceDB(ctx)
	}
	return tx
}

func (repo UserRepository) GetUserByUsername(ctx context.Context, username string) (user user.User, err error) {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	if err = repo.getInstanceDB(ctx).WithContext(ctxWT).Where("LOWER(username) = LOWER(?)", username).First(&user).Error; err != nil {
		return
	}

	return user, nil
}

func (repo UserRepository) GetUserByID(ctx context.Context, id uint) (user user.User, err error) {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	if err = repo.getInstanceDB(ctx).WithContext(ctxWT).Where("id = ?", id).First(&user).Error; err != nil {
		return
	}

	return user, nil
}

func (repo UserRepository) CreateUser(ctx context.Context, user user.User) (createdUser user.User, err error) {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	if err = repo.getInstanceDB(ctx).WithContext(ctxWT).Create(&user).Error; err != nil {
		return
	}

	return user, nil
}
