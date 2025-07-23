package user

import (
	"context"
	"errors"

	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/user"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"github.com/riskykurniawan15/payrolls/utils/bcrypt"
	"github.com/riskykurniawan15/payrolls/utils/jwt"
	"github.com/riskykurniawan15/payrolls/utils/logger"
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
		logger   logger.Logger
	}
)

func NewUserService(config config.Config, userRepo userRepo.IUserRepository, logger logger.Logger) IUserService {
	return &UserService{
		config:   config,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (service *UserService) Login(ctx context.Context, req user.LoginRequest) (response user.LoginResponse, err error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting login process", requestID, map[string]interface{}{
		"username": req.Username,
	})

	// Get user by username
	userData, err := service.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		service.logger.WarningT("user not found", requestID, map[string]interface{}{
			"username": req.Username,
			"error":    err.Error(),
		})
		return response, errors.New("invalid username or password")
	}

	service.logger.InfoT("user found, verifying password", requestID, map[string]interface{}{
		"username": req.Username,
		"user_id":  userData.ID,
		"role":     userData.Role,
	})

	// Verify password
	err = bcrypt.VerifyPassword(userData.Password, req.Password)
	if err != nil {
		service.logger.WarningT("password verification failed", requestID, map[string]interface{}{
			"username": req.Username,
			"user_id":  userData.ID,
			"error":    err.Error(),
		})
		return response, errors.New("invalid username or password")
	}

	service.logger.InfoT("password verified, generating JWT token", requestID, map[string]interface{}{
		"username": req.Username,
		"user_id":  userData.ID,
		"role":     userData.Role,
	})

	// Generate JWT token
	jwtConfig := jwt.JWTConfig{
		SecretKey: service.config.JWT.SecretKey,
		Expired:   service.config.JWT.Expired,
	}

	token, expiresAt, err := jwt.GenerateToken(jwtConfig, userData.ID, userData.Username, userData.Role)
	if err != nil {
		service.logger.ErrorT("failed to generate JWT token", requestID, map[string]interface{}{
			"username": req.Username,
			"user_id":  userData.ID,
			"error":    err.Error(),
		})
		return response, err
	}

	service.logger.InfoT("JWT token generated successfully", requestID, map[string]interface{}{
		"username":   req.Username,
		"user_id":    userData.ID,
		"role":       userData.Role,
		"expires_at": expiresAt,
	})

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
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting profile retrieval", requestID, map[string]interface{}{
		"user_id": userID,
	})

	// Get user data from database
	userData, err := service.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		service.logger.WarningT("user not found in database", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return response, errors.New("user not found")
	}

	service.logger.InfoT("user profile retrieved successfully", requestID, map[string]interface{}{
		"user_id":    userID,
		"username":   userData.Username,
		"role":       userData.Role,
		"created_at": userData.CreatedAt,
		"updated_at": userData.UpdatedAt,
	})

	return user.ProfileResponse{
		ID:        userData.ID,
		Username:  userData.Username,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}, nil
}

func (service *UserService) CreateUser(ctx context.Context, req user.CreateUserRequest) (response user.CreateUserResponse, err error) {
	requestID := middleware.GetRequestIDFromContext(ctx)

	service.logger.InfoT("starting user creation process", requestID, map[string]interface{}{
		"username": req.Username,
		"role":     req.Role,
	})

	// Check if username already exists
	existingUser, _ := service.userRepo.GetUserByUsername(ctx, req.Username)
	if existingUser.ID != 0 {
		service.logger.WarningT("username already exists", requestID, map[string]interface{}{
			"username":      req.Username,
			"existing_id":   existingUser.ID,
			"existing_role": existingUser.Role,
		})
		return response, errors.New("username already exists")
	}

	service.logger.InfoT("username is available, validating password", requestID, map[string]interface{}{
		"username": req.Username,
	})

	// Validate password strength
	if !bcrypt.IsValidPassword(req.Password, 6) {
		service.logger.WarningT("password validation failed", requestID, map[string]interface{}{
			"username": req.Username,
			"error":    "password must be at least 6 characters long",
		})
		return response, errors.New("password must be at least 6 characters long")
	}

	service.logger.InfoT("password validated, hashing password", requestID, map[string]interface{}{
		"username": req.Username,
	})

	// Hash password using environment cost
	hashedPassword, err := bcrypt.HashPasswordWithEnvCost(req.Password)
	if err != nil {
		service.logger.ErrorT("failed to hash password", requestID, map[string]interface{}{
			"username": req.Username,
			"error":    err.Error(),
		})
		return response, errors.New("failed to hash password")
	}

	service.logger.InfoT("password hashed successfully, creating user in database", requestID, map[string]interface{}{
		"username": req.Username,
		"role":     req.Role,
	})

	// Create user data
	userData := user.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Save user to database
	createdUser, err := service.userRepo.CreateUser(ctx, userData)
	if err != nil {
		service.logger.ErrorT("failed to create user in database", requestID, map[string]interface{}{
			"username": req.Username,
			"role":     req.Role,
			"error":    err.Error(),
		})
		return response, errors.New("failed to create user")
	}

	service.logger.InfoT("user created successfully", requestID, map[string]interface{}{
		"user_id":    createdUser.ID,
		"username":   createdUser.Username,
		"role":       createdUser.Role,
		"created_at": createdUser.CreatedAt,
	})

	return user.CreateUserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}
