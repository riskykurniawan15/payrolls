package overtime

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/overtime"
	overtimeServices "github.com/riskykurniawan15/payrolls/services/overtime"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"github.com/riskykurniawan15/payrolls/utils/validator"
)

type (
	IOvertimeHandler interface {
		Create(ctx echo.Context) error
		GetByID(ctx echo.Context) error
		Update(ctx echo.Context) error
		Delete(ctx echo.Context) error
		List(ctx echo.Context) error
	}

	OvertimeHandler struct {
		logger           logger.Logger
		overtimeServices overtimeServices.IOvertimeService
	}
)

func NewOvertimeHandlers(logger logger.Logger, overtimeServices overtimeServices.IOvertimeService) IOvertimeHandler {
	return &OvertimeHandler{
		logger:           logger,
		overtimeServices: overtimeServices,
	}
}

func (handler OvertimeHandler) Create(ctx echo.Context) error {
	var req overtime.CreateOvertimeRequest
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
		"user_id":          req.UserID,
		"overtimes_date":   req.OvertimesDate,
		"total_hours_time": req.TotalHoursTime,
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
	handler.logger.InfoT("calling overtime service", requestID, map[string]interface{}{
		"user_id":          userID,
		"request_user_id":  req.UserID,
		"overtimes_date":   req.OvertimesDate,
		"total_hours_time": req.TotalHoursTime,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.overtimeServices.Create(serviceCtx, req, userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("overtime created successfully", requestID, map[string]interface{}{
		"overtime_id":      response.ID,
		"user_id":          response.UserID,
		"overtimes_date":   response.OvertimesDate.Format("2006-01-02"),
		"total_hours_time": response.TotalHoursTime,
	})

	return ctx.JSON(http.StatusCreated, entities.ResponseFormater(http.StatusCreated, map[string]interface{}{
		"data": response,
	}))
}

func (handler OvertimeHandler) GetByID(ctx echo.Context) error {
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
	response, err := handler.overtimeServices.GetByID(serviceCtx, uint(id))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler OvertimeHandler) Update(ctx echo.Context) error {
	// Get ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	var req overtime.UpdateOvertimeRequest
	requestID := middleware.GetRequestID(ctx)

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"id":               id,
		"overtimes_date":   req.OvertimesDate,
		"total_hours_time": req.TotalHoursTime,
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
	response, err := handler.overtimeServices.Update(serviceCtx, uint(id), req, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler OvertimeHandler) Delete(ctx echo.Context) error {
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
	err = handler.overtimeServices.Delete(serviceCtx, uint(id))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler OvertimeHandler) List(ctx echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	userIDStr := ctx.QueryParam("user_id")
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")
	sortBy := ctx.QueryParam("sort_by")
	sortDesc := ctx.QueryParam("sort_desc") == "true"

	requestID := middleware.GetRequestID(ctx)
	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"user_id":    userIDStr,
		"start_date": startDate,
		"end_date":   endDate,
		"sort_by":    sortBy,
		"sort_desc":  sortDesc,
	})

	// Build request
	req := overtime.ListOvertimesRequest{
		Page:      page,
		Limit:     limit,
		StartDate: &startDate,
		EndDate:   &endDate,
		SortBy:    sortBy,
		SortDesc:  sortDesc,
	}

	// Parse user_id if provided
	if userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			req.UserID = &userIDUint
		}
	}

	// Handle empty date strings
	if startDate == "" {
		req.StartDate = nil
	}
	if endDate == "" {
		req.EndDate = nil
	}

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.overtimeServices.List(serviceCtx, req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response.Data,
		"meta": response.Pagination,
	}))
}
