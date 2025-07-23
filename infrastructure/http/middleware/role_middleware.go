package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
)

// RoleMiddleware creates middleware to check if user has required role
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from JWT context
			userRole := GetRole(c)
			if userRole == "" {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "User role not found",
				}))
			}

			// Check if user has required role
			if userRole != requiredRole {
				return c.JSON(http.StatusForbidden, entities.ResponseFormater(http.StatusForbidden, map[string]interface{}{
					"error": "Access denied. Required role: " + requiredRole,
				}))
			}

			return next(c)
		}
	}
}

// AdminOnlyMiddleware creates middleware that only allows admin access
func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return RoleMiddleware(constant.AdminRole)
}

// EmployeeOnlyMiddleware creates middleware that only allows employee access
func EmployeeOnlyMiddleware() echo.MiddlewareFunc {
	return RoleMiddleware(constant.EmployeeRole)
}

// MultipleRolesMiddleware creates middleware that allows access for multiple roles
func MultipleRolesMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from JWT context
			userRole := GetRole(c)
			if userRole == "" {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "User role not found",
				}))
			}

			// Check if user has any of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return c.JSON(http.StatusForbidden, entities.ResponseFormater(http.StatusForbidden, map[string]interface{}{
					"error": "Access denied. Required roles: " + stringSliceToString(allowedRoles),
				}))
			}

			return next(c)
		}
	}
}

// Helper function to convert string slice to comma-separated string
func stringSliceToString(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}

	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += ", " + slice[i]
	}
	return result
}
