package router

import (
	"github.com/labstack/echo/v4"
	dep "github.com/riskykurniawan15/payrolls/infrastructure/http"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/utils/validator"
)

func Routers(dep *dep.Dependencies, jwtSecret string) *echo.Echo {
	engine := echo.New()

	// Add custom validator
	engine.Validator = validator.NewCustomValidator()

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

	// Period routes
	periods := engine.Group("/periods", middleware.JWTMiddleware(jwtConfig), middleware.AdminOnlyMiddleware())
	{
		periods.POST("", dep.PeriodHandlers.Create)
		periods.GET("", dep.PeriodHandlers.List)
		periods.GET("/:id", dep.PeriodHandlers.GetByID)
		periods.PUT("/:id", dep.PeriodHandlers.Update)
		periods.DELETE("/:id", dep.PeriodHandlers.Delete)

		// Period detail routes
		periods.POST("/:id/run-payroll", dep.PeriodDetailHandlers.RunPayroll)
	}

	// Attendance routes (for all authenticated users)
	attendances := engine.Group("/attendances", middleware.JWTMiddleware(jwtConfig), middleware.EmployeeOnlyMiddleware())
	{
		attendances.GET("", dep.AttendanceHandlers.GetAttendances)
		attendances.GET("/:id", dep.AttendanceHandlers.GetAttendanceByID)
		attendances.POST("/check-in", dep.AttendanceHandlers.CheckIn)
		attendances.POST("/check-out", dep.AttendanceHandlers.CheckOut)
		attendances.POST("/check-out/:id", dep.AttendanceHandlers.CheckOutByID)
	}

	// Overtime routes (for all authenticated users)
	overtimes := engine.Group("/overtimes", middleware.JWTMiddleware(jwtConfig), middleware.EmployeeOnlyMiddleware())
	{
		overtimes.POST("", dep.OvertimeHandlers.Create)
		overtimes.GET("", dep.OvertimeHandlers.List)
		overtimes.GET("/:id", dep.OvertimeHandlers.GetByID)
		overtimes.PUT("/:id", dep.OvertimeHandlers.Update)
		overtimes.DELETE("/:id", dep.OvertimeHandlers.Delete)
	}

	// Reimbursement routes (for all authenticated users)
	reimbursements := engine.Group("/reimbursements", middleware.JWTMiddleware(jwtConfig), middleware.EmployeeOnlyMiddleware())
	{
		reimbursements.POST("", dep.ReimbursementHandlers.Create)
		reimbursements.GET("", dep.ReimbursementHandlers.List)
		reimbursements.GET("/:id", dep.ReimbursementHandlers.GetByID)
		reimbursements.PUT("/:id", dep.ReimbursementHandlers.Update)
		reimbursements.DELETE("/:id", dep.ReimbursementHandlers.Delete)
	}

	return engine
}
