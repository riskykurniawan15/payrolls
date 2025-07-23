package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestIDMiddleware adds a request ID to the context
// If X-Request-ID header exists, use it; otherwise generate a new UUID
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID exists in header
			requestID := c.Request().Header.Get(RequestIDHeader)

			// If no request ID in header, generate a new UUID
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Set request ID to context
			c.Set(RequestIDKey, requestID)

			// Add request ID to response header for client reference
			c.Response().Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c echo.Context) string {
	requestID, _ := c.Get(RequestIDKey).(string)
	return requestID
}

// GetRequestIDFromContext retrieves the request ID from context.Context
// This can be used in repository layer
func GetRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Try to get request ID from context
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// AddRequestIDToContext adds request ID to a context.Context
// This is a helper function to avoid direct access to RequestIDKey
func AddRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}
