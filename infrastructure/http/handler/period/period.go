package period

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/period"
	periodServices "github.com/riskykurniawan15/payrolls/services/period"
	"github.com/riskykurniawan15/payrolls/utils/validator"
)

type (
	IPeriodHandler interface {
		Create(ctx echo.Context) error
		GetByID(ctx echo.Context) error
		Update(ctx echo.Context) error
		Delete(ctx echo.Context) error
		List(ctx echo.Context) error
	}

	PeriodHandler struct {
		periodServices periodServices.IPeriodService
	}
)

func NewPeriodHandlers(periodServices periodServices.IPeriodService) IPeriodHandler {
	return &PeriodHandler{
		periodServices: periodServices,
	}
}

func (handler PeriodHandler) Create(ctx echo.Context) error {
	var req period.CreatePeriodRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		if validationErrors, ok := err.(*validator.ValidationErrors); ok {
			return ctx.JSON(http.StatusBadRequest, entities.Response{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Error:   "Validation failed",
				Meta: map[string]interface{}{
					"validation_errors": validationErrors.GetValidationErrors(),
				},
			})
		}
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Call service
	response, err := handler.periodServices.Create(ctx.Request().Context(), req, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusCreated, entities.ResponseFormater(http.StatusCreated, map[string]interface{}{
		"data": response,
	}))
}

func (handler PeriodHandler) GetByID(ctx echo.Context) error {
	// Get ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	// Call service
	response, err := handler.periodServices.GetByID(ctx.Request().Context(), uint(id))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler PeriodHandler) Update(ctx echo.Context) error {
	// Get ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	var req period.UpdatePeriodRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		if validationErrors, ok := err.(*validator.ValidationErrors); ok {
			return ctx.JSON(http.StatusBadRequest, entities.Response{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Error:   "Validation failed",
				Meta: map[string]interface{}{
					"validation_errors": validationErrors.GetValidationErrors(),
				},
			})
		}
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Call service
	response, err := handler.periodServices.Update(ctx.Request().Context(), uint(id), req, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler PeriodHandler) Delete(ctx echo.Context) error {
	// Get ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	// Call service
	err = handler.periodServices.Delete(ctx.Request().Context(), uint(id))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler PeriodHandler) List(ctx echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	search := ctx.QueryParam("search")
	statusStr := ctx.QueryParam("status")
	sortBy := ctx.QueryParam("sort_by")
	sortDesc := ctx.QueryParam("sort_desc") == "true"

	// Build request
	req := period.ListPeriodsRequest{
		Page:     page,
		Limit:    limit,
		Search:   search,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}

	// Parse status if provided
	if statusStr != "" {
		if status, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusInt8 := int8(status)
			req.Status = &statusInt8
		}
	}

	// Call service
	response, err := handler.periodServices.List(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
