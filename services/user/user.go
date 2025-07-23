package user

import (
	"context"
	"errors"

	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/models/user"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"github.com/riskykurniawan15/payrolls/utils/bcrypt"
	"github.com/riskykurniawan15/payrolls/utils/jwt"
)

type (
	IUserService interface {
		Login(ctx context.Context, req user.LoginRequest) (user.LoginResponse, error)
		Profile(ctx context.Context, userID uint) (user.ProfileResponse, error)
		CreateUser(ctx context.Context, req user.CreateUserRequest) (user.CreateUserResponse, error)
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

func (service *UserService) Login(ctx context.Context, req user.LoginRequest) (response user.LoginResponse, err error) {
	// Get user by username
	userData, err := service.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return response, errors.New("invalid username or password")
	}

	// Verify password
	err = bcrypt.VerifyPassword(userData.Password, req.Password)
	if err != nil {
		return response, errors.New("invalid username or password")
	}

	// Generate JWT token
	jwtConfig := jwt.JWTConfig{
		SecretKey: service.config.JWT.SecretKey,
		Expired:   service.config.JWT.Expired,
	}

	token, expiresAt, err := jwt.GenerateToken(jwtConfig, userData.ID, userData.Username, userData.Role)
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

func (service *UserService) Profile(ctx context.Context, userID uint) (response user.ProfileResponse, err error) {
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

func (service *UserService) CreateUser(ctx context.Context, req user.CreateUserRequest) (response user.CreateUserResponse, err error) {
	// Check if username already exists
	existingUser, _ := service.userRepo.GetUserByUsername(ctx, req.Username)
	if existingUser.ID != 0 {
		return response, errors.New("username already exists")
	}

	// Validate password strength
	if !bcrypt.IsValidPassword(req.Password, 6) {
		return response, errors.New("password must be at least 6 characters long")
	}

	// Hash password using environment cost
	hashedPassword, err := bcrypt.HashPasswordWithEnvCost(req.Password)
	if err != nil {
		return response, errors.New("failed to hash password")
	}

	// Create user data
	userData := user.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Save user to database
	createdUser, err := service.userRepo.CreateUser(ctx, userData)
	if err != nil {
		return response, errors.New("failed to create user")
	}

	return user.CreateUserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}
