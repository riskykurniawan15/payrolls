package user

import (
	"context"

	"github.com/riskykurniawan15/payrolls/models/user"
	"gorm.io/gorm"
)

type (
	IUserRepository interface {
		GetUserByUsername(ctx context.Context, username string) (user.User, error)
		GetUserByID(ctx context.Context, id uint) (user.User, error)
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

func (repo UserRepository) GetUserByUsername(ctx context.Context, username string) (user user.User, err error) {
	if err = repo.db.WithContext(ctx).Where("LOWER(username) = LOWER(?)", username).First(&user).Error; err != nil {
		return
	}

	return user, nil
}

func (repo UserRepository) GetUserByID(ctx context.Context, id uint) (user user.User, err error) {
	if err = repo.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return
	}

	return user, nil
}
