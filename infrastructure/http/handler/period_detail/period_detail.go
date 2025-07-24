package period_detail

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
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
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		}))
	}

	// Get user ID from middleware
	userID := middleware.GetUserID(ctx)

	// Call service
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), middleware.GetRequestID(ctx))
	response, err := handler.periodDetailServices.RunPayroll(serviceCtx, uint(periodID), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}
