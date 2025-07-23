package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/utils/jwt"
)

type (
	JWTConfig struct {
		SecretKey string
	}
)

func JWTMiddleware(config JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "Authorization header is required",
				}))
			}

			// Extract token from "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid authorization header format. Use 'Bearer <token>'",
				}))
			}

			tokenString := tokenParts[1]

			// Parse and validate JWT token
			claims, err := jwt.ParseToken(tokenString, config.SecretKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid or expired token",
				}))
			}

			// Extract user context
			userCtx := jwt.ExtractUserContext(claims)

			// Set user context to context
			c.Set("user", userCtx)

			return next(c)
		}
	}
}

// Helper function to get user context from context
func GetUserContext(c echo.Context) jwt.UserContext {
	userCtx, _ := c.Get("user").(jwt.UserContext)
	return userCtx
}

// Helper function to get user ID from context
func GetUserID(c echo.Context) uint {
	userCtx := GetUserContext(c)
	return userCtx.UserID
}

// Helper function to get username from context
func GetUsername(c echo.Context) string {
	userCtx := GetUserContext(c)
	return userCtx.Username
}

// Helper function to get role from context
func GetRole(c echo.Context) string {
	userCtx := GetUserContext(c)
	return userCtx.Role
}
