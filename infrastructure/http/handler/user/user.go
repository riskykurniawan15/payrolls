package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/user"
	userServices "github.com/riskykurniawan15/payrolls/services/user"
)

type (
	IUserHandler interface {
		Login(ctx echo.Context) error
		Profile(ctx echo.Context) error
	}

	UserHandler struct {
		userServices userServices.IUserService
	}
)

func NewUserHandlers(userServices userServices.IUserService) IUserHandler {
	return &UserHandler{
		userServices: userServices,
	}
}

func (handler UserHandler) Login(ctx echo.Context) error {
	var req user.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	// Validate request using custom validator
	if err := ctx.Validate(&req); err != nil {
		if validationErrors, ok := err.(*middleware.ValidationErrors); ok {
			return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
				"error":             "Validation failed",
				"validation_errors": validationErrors.GetValidationErrors(),
			}))
		}
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Call service
	response, err := handler.userServices.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler UserHandler) Profile(ctx echo.Context) error {
	// Get user ID from middleware context
	userID := middleware.GetUserID(ctx)

	// Call service
	response, err := handler.userServices.Profile(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
