package period

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/period"
	periodServices "github.com/riskykurniawan15/payrolls/services/period"
	"github.com/riskykurniawan15/payrolls/utils/logger"
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
		logger         logger.Logger
		periodServices periodServices.IPeriodService
	}
)

func NewPeriodHandlers(logger logger.Logger, periodServices periodServices.IPeriodService) IPeriodHandler {
	return &PeriodHandler{
		logger:         logger,
		periodServices: periodServices,
	}
}

func (handler PeriodHandler) Create(ctx echo.Context) error {
	var req period.CreatePeriodRequest
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
		"name":       req.Name,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
	})

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		if validationErrors, ok := err.(*validator.ValidationErrors); ok {
			handler.logger.WarningT("validation failed", requestID, map[string]interface{}{
				"validation_errors": validationErrors.GetValidationErrors(),
			})
			return ctx.JSON(http.StatusBadRequest, entities.Response{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Error:   "Validation failed",
				Meta: map[string]interface{}{
					"validation_errors": validationErrors.GetValidationErrors(),
				},
			})
		}
		handler.logger.ErrorT("validation error", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)
	handler.logger.InfoT("calling period service", requestID, map[string]interface{}{
		"user_id": userID,
		"name":    req.Name,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.periodServices.Create(serviceCtx, req, userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("period created successfully", requestID, map[string]interface{}{
		"period_id": response.ID,
		"code":      response.Code,
		"name":      response.Name,
	})

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

	requestID := middleware.GetRequestID(ctx)
	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"id": id,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.periodServices.GetByID(serviceCtx, uint(id))
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
	requestID := middleware.GetRequestID(ctx)

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"id":         id,
		"name":       req.Name,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
	})

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

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.periodServices.Update(serviceCtx, uint(id), req, userID)
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

	requestID := middleware.GetRequestID(ctx)
	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"id": id,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	err = handler.periodServices.Delete(serviceCtx, uint(id))
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

	requestID := middleware.GetRequestID(ctx)
	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"page":      page,
		"limit":     limit,
		"search":    search,
		"status":    statusStr,
		"sort_by":   sortBy,
		"sort_desc": sortDesc,
	})

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

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.periodServices.List(serviceCtx, req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
