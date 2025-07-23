package router

import (
	"github.com/labstack/echo/v4"
	dep "github.com/riskykurniawan15/payrolls/infrastructure/http"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
)

func Routers(dep *dep.Dependencies, jwtSecret string) *echo.Echo {
	engine := echo.New()

	// Add custom validator
	engine.Validator = middleware.NewCustomValidator()

	// Add request ID middleware globally
	engine.Use(middleware.RequestIDMiddleware())

	// Public routes
	engine.GET("/health", dep.HealthHandlers.Metric)
	engine.POST("/auth/login", dep.UserHandlers.Login)

	// Protected routes with JWT middleware
	jwtConfig := middleware.JWTConfig{SecretKey: jwtSecret}
	protected := engine.Group("/user", middleware.JWTMiddleware(jwtConfig))
	{
		protected.GET("/profile", dep.UserHandlers.Profile)
	}

	return engine
}
