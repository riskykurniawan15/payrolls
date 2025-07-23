package router

import (
	"github.com/labstack/echo/v4"
	dep "github.com/riskykurniawan15/payrolls/infrastructure/http"
)

func Routers(dep *dep.Dependencies) *echo.Echo {
	engine := echo.New()

	engine.GET("/health", dep.HealthHandlers.Metric)

	return engine
}
