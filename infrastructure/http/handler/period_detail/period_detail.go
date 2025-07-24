package period_detail

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	periodDetailModel "github.com/riskykurniawan15/payrolls/models/period_detail"
	periodDetailService "github.com/riskykurniawan15/payrolls/services/period_detail"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IPeriodDetailHandler interface {
		RunPayroll(ctx echo.Context) error
	}

	PeriodDetailHandler struct {
		periodDetailServices periodDetailService.IPeriodDetailService
		logger               logger.Logger
	}
)

func NewPeriodDetailHandlers(
	periodDetailServices periodDetailService.IPeriodDetailService,
	logger logger.Logger,
) IPeriodDetailHandler {
	return &PeriodDetailHandler{
		periodDetailServices: periodDetailServices,
		logger:               logger,
	}
}

func (handler PeriodDetailHandler) RunPayroll(ctx echo.Context) error {
	// Get period ID from URL parameter
	periodIDStr := ctx.Param("id")
	periodID, err := strconv.ParseUint(periodIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid period ID",
		})
	}

	// Bind request body
	var req periodDetailModel.RunPayrollRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"errors":  err,
		})
	}

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Call service
	serviceCtx := ctx.Request().Context()
	response, err := handler.periodDetailServices.RunPayroll(serviceCtx, uint(periodID), req, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
