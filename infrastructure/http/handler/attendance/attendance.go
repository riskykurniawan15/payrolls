package attendance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/attendance"
	attendanceServices "github.com/riskykurniawan15/payrolls/services/attendance"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IAttendanceHandler interface {
		GetAttendances(ctx echo.Context) error
		GetAttendanceByID(ctx echo.Context) error
		CheckIn(ctx echo.Context) error
		CheckOut(ctx echo.Context) error
		CheckOutByID(ctx echo.Context) error
	}

	AttendanceHandler struct {
		logger             logger.Logger
		attendanceServices attendanceServices.IAttendanceService
	}
)

func NewAttendanceHandlers(logger logger.Logger, attendanceServices attendanceServices.IAttendanceService) IAttendanceHandler {
	return &AttendanceHandler{
		logger:             logger,
		attendanceServices: attendanceServices,
	}
}

func (handler AttendanceHandler) GetAttendances(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	// Get pagination parameters
	pageStr := ctx.QueryParam("page")
	limitStr := ctx.QueryParam("limit")
	startDateStr := ctx.QueryParam("start_date")
	endDateStr := ctx.QueryParam("end_date")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse date parameters
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set end date to end of day
			endOfDay := parsed.Add(24*time.Hour - time.Second)
			endDate = &endOfDay
		}
	}

	handler.logger.InfoT("incoming get attendances request", requestID, map[string]interface{}{
		"user_id":     userID,
		"page":        page,
		"limit":       limit,
		"start_date":  startDateStr,
		"end_date":    endDateStr,
		"has_filters": startDate != nil || endDate != nil,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.attendanceServices.GetAttendances(serviceCtx, userID, page, limit, startDate, endDate)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("attendances retrieved successfully", requestID, map[string]interface{}{
		"user_id":     userID,
		"total":       response.Pagination.Total,
		"page":        response.Pagination.Page,
		"limit":       response.Pagination.Limit,
		"total_pages": response.Pagination.TotalPages,
		"count":       len(response.Data),
		"start_date":  startDateStr,
		"end_date":    endDateStr,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response.Data,
		"meta": response.Pagination,
	}))
}

func (handler AttendanceHandler) GetAttendanceByID(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	// Get attendance ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		handler.logger.WarningT("invalid attendance ID parameter", requestID, map[string]interface{}{
			"user_id": userID,
			"id_str":  idStr,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid attendance ID",
		}))
	}

	handler.logger.InfoT("incoming get attendance by ID request", requestID, map[string]interface{}{
		"user_id":       userID,
		"attendance_id": id,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.attendanceServices.GetAttendanceByID(serviceCtx, uint(id), userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id":       userID,
			"attendance_id": id,
			"error":         err.Error(),
		})
		return ctx.JSON(http.StatusNotFound, entities.ResponseFormater(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("attendance retrieved successfully", requestID, map[string]interface{}{
		"user_id":        userID,
		"attendance_id":  id,
		"check_in_date":  response.CheckInDate,
		"check_out_date": response.CheckOutDate,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler AttendanceHandler) CheckIn(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	var req attendance.AttendanceRequest

	if err := ctx.Bind(&req); err != nil {
		handler.logger.ErrorT("failed to bind request body", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming check-in request", requestID, map[string]interface{}{
		"user_id":  userID,
		"has_date": req.Date != nil,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.attendanceServices.CheckIn(serviceCtx, userID, req)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("check-in successful", requestID, map[string]interface{}{
		"user_id":       userID,
		"attendance_id": response.ID,
		"check_in_date": response.CheckInDate,
	})

	return ctx.JSON(http.StatusCreated, entities.ResponseFormater(http.StatusCreated, map[string]interface{}{
		"data": response,
	}))
}

func (handler AttendanceHandler) CheckOut(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	var req attendance.AttendanceRequest

	if err := ctx.Bind(&req); err != nil {
		handler.logger.ErrorT("failed to bind request body", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming check-out request", requestID, map[string]interface{}{
		"user_id":  userID,
		"has_date": req.Date != nil,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.attendanceServices.CheckOut(serviceCtx, userID, req)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("check-out successful", requestID, map[string]interface{}{
		"user_id":        userID,
		"attendance_id":  response.ID,
		"check_in_date":  response.CheckInDate,
		"check_out_date": response.CheckOutDate,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler AttendanceHandler) CheckOutByID(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	// Get attendance ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		handler.logger.WarningT("invalid attendance ID parameter", requestID, map[string]interface{}{
			"user_id": userID,
			"id_str":  idStr,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid attendance ID",
		}))
	}

	var req attendance.AttendanceRequest

	if err := ctx.Bind(&req); err != nil {
		handler.logger.ErrorT("failed to bind request body", requestID, map[string]interface{}{
			"user_id":       userID,
			"attendance_id": id,
			"error":         err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		}))
	}

	handler.logger.InfoT("incoming check-out by ID request", requestID, map[string]interface{}{
		"user_id":       userID,
		"attendance_id": id,
		"has_date":      req.Date != nil,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.attendanceServices.CheckOutByID(serviceCtx, uint(id), userID, req)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id":       userID,
			"attendance_id": id,
			"error":         err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("check-out by ID successful", requestID, map[string]interface{}{
		"user_id":        userID,
		"attendance_id":  response.ID,
		"check_in_date":  response.CheckInDate,
		"check_out_date": response.CheckOutDate,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
