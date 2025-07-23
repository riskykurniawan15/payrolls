package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/user"
	userServices "github.com/riskykurniawan15/payrolls/services/user"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"github.com/riskykurniawan15/payrolls/utils/validator"
)

type (
	IUserHandler interface {
		Login(ctx echo.Context) error
		Profile(ctx echo.Context) error
	}

	UserHandler struct {
		logger       logger.Logger
		userServices userServices.IUserService
	}
)

func NewUserHandlers(logger logger.Logger, userServices userServices.IUserService) IUserHandler {
	return &UserHandler{
		logger:       logger,
		userServices: userServices,
	}
}

func (handler UserHandler) Login(ctx echo.Context) error {
	var req user.LoginRequest
	requestID := middleware.GetRequestID(ctx)

	if err := ctx.Bind(&req); err != nil {
		handler.logger.ErrorT("failed to bind request body", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"username": req.Username,
	})

	// Validate request using custom validator
	if err := ctx.Validate(&req); err != nil {
		if validationErrors, ok := err.(*validator.ValidationErrors); ok {
			handler.logger.WarningT("validation failed", requestID, map[string]interface{}{
				"validation_errors": validationErrors.GetValidationErrors(),
			})
			return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
				"error":             "Validation failed",
				"validation_errors": validationErrors.GetValidationErrors(),
			}))
		}
		handler.logger.ErrorT("validation error", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.userServices.Login(serviceCtx, req)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("login successful", requestID, map[string]interface{}{
		"username": req.Username,
		"user_id":  response.User.ID,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler UserHandler) Profile(ctx echo.Context) error {
	// Get user ID from middleware context
	userID := middleware.GetUserID(ctx)
	requestID := middleware.GetRequestID(ctx)

	handler.logger.InfoT("incoming profile request", requestID, map[string]interface{}{
		"user_id": userID,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.userServices.Profile(serviceCtx, userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("profile retrieved successfully", requestID, map[string]interface{}{
		"user_id":  userID,
		"username": response.Username,
		"role":     response.Role,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
