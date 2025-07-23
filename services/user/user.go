package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/models/user"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type (
	IUserService interface {
		Login(ctx context.Context, req user.LoginRequest) (user.LoginResponse, error)
		Profile(ctx context.Context, userID uint) (user.ProfileResponse, error)
	}

	UserService struct {
		config   config.Config
		userRepo userRepo.IUserRepository
	}
)

func NewUserService(config config.Config, userRepo userRepo.IUserRepository) IUserService {
	return &UserService{
		config:   config,
		userRepo: userRepo,
	}
}

func (service UserService) Login(ctx context.Context, req user.LoginRequest) (response user.LoginResponse, err error) {
	// Get user by username
	userData, err := service.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return response, errors.New("invalid username or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password))
	if err != nil {
		return response, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, expiresAt, err := service.generateJWT(userData)
	if err != nil {
		return response, err
	}

	return user.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: user.UserInfo{
			ID:       userData.ID,
			Username: userData.Username,
			Role:     userData.Role,
		},
	}, nil
}

func (service UserService) Profile(ctx context.Context, userID uint) (response user.ProfileResponse, err error) {
	// Get user data from database
	userData, err := service.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return response, errors.New("user not found")
	}

	return user.ProfileResponse{
		ID:        userData.ID,
		Username:  userData.Username,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}, nil
}

func (service UserService) generateJWT(user user.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(service.config.JWT.Expired) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(service.config.JWT.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
