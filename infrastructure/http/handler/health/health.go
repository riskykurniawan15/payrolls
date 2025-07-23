package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	healthServices "github.com/riskykurniawan15/payrolls/services/health"
)

type (
	IHealthHandler interface {
		Metric(ctx echo.Context) error
	}
	HealthHandler struct {
		healthServices healthServices.IHealthServices
	}
)

func NewHealthHandlers(healthServices healthServices.IHealthServices) IHealthHandler {
	return &HealthHandler{
		healthServices: healthServices,
	}
}

func (handler HealthHandler) Metric(ctx echo.Context) error {
	metric := handler.healthServices.HealthMetric(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": metric,
	}))
}
