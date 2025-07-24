package reimbursement

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/reimbursement"
	reimbursementServices "github.com/riskykurniawan15/payrolls/services/reimbursement"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"github.com/riskykurniawan15/payrolls/utils/validator"
)

type (
	IReimbursementHandler interface {
		Create(ctx echo.Context) error
		GetByID(ctx echo.Context) error
		Update(ctx echo.Context) error
		Delete(ctx echo.Context) error
		List(ctx echo.Context) error
	}

	ReimbursementHandler struct {
		logger                logger.Logger
		reimbursementServices reimbursementServices.IReimbursementService
	}
)

func NewReimbursementHandlers(logger logger.Logger, reimbursementServices reimbursementServices.IReimbursementService) IReimbursementHandler {
	return &ReimbursementHandler{
		logger:                logger,
		reimbursementServices: reimbursementServices,
	}
}

func (handler ReimbursementHandler) Create(ctx echo.Context) error {
	var req reimbursement.CreateReimbursementRequest
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
		"title":       req.Title,
		"date":        req.Date,
		"amount":      req.Amount,
		"description": req.Description,
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
	handler.logger.InfoT("calling reimbursement service", requestID, map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"date":    req.Date,
		"amount":  req.Amount,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.reimbursementServices.Create(serviceCtx, req, userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("reimbursement created successfully", requestID, map[string]interface{}{
		"reimbursement_id": response.ID,
		"user_id":          response.UserID,
		"title":            response.Title,
		"date":             response.Date.Format("2006-01-02"),
		"amount":           response.Amount,
	})

	return ctx.JSON(http.StatusCreated, entities.ResponseFormater(http.StatusCreated, map[string]interface{}{
		"data": response,
	}))
}

func (handler ReimbursementHandler) GetByID(ctx echo.Context) error {
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

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.reimbursementServices.GetByID(serviceCtx, uint(id), userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler ReimbursementHandler) Update(ctx echo.Context) error {
	// Get ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	var req reimbursement.UpdateReimbursementRequest
	requestID := middleware.GetRequestID(ctx)

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"id":          id,
		"title":       req.Title,
		"date":        req.Date,
		"amount":      req.Amount,
		"description": req.Description,
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
	response, err := handler.reimbursementServices.Update(serviceCtx, uint(id), req, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler ReimbursementHandler) Delete(ctx echo.Context) error {
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

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	err = handler.reimbursementServices.Delete(serviceCtx, uint(id), userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler ReimbursementHandler) List(ctx echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")
	sortBy := ctx.QueryParam("sort_by")
	sortDesc := ctx.QueryParam("sort_desc") == "true"

	requestID := middleware.GetRequestID(ctx)
	handler.logger.InfoT("incoming request", requestID, map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"start_date": startDate,
		"end_date":   endDate,
		"sort_by":    sortBy,
		"sort_desc":  sortDesc,
	})

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Build request - user can only see their own reimbursements
	req := reimbursement.ListReimbursementsRequest{
		Page:      page,
		Limit:     limit,
		StartDate: &startDate,
		EndDate:   &endDate,
		SortBy:    sortBy,
		SortDesc:  sortDesc,
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
	response, err := handler.reimbursementServices.List(serviceCtx, req, userID)
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
