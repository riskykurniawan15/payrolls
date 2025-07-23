package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
)

type (
	JWTConfig struct {
		SecretKey string
	}

	JWTClaims struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.RegisteredClaims
	}

	UserContext struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
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
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.SecretKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid or expired token",
				}))
			}

			// Extract claims
			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, entities.ResponseFormater(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token claims",
				}))
			}

			// Create user context struct
			userCtx := UserContext{
				UserID:   claims.UserID,
				Username: claims.Username,
				Role:     claims.Role,
			}

			// Set user context to context
			c.Set("user", userCtx)

			return next(c)
		}
	}
}

// Helper function to get user context from context
func GetUserContext(c echo.Context) UserContext {
	userCtx, _ := c.Get("user").(UserContext)
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
