package middleware

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/models/audit_trail"
	auditTrailService "github.com/riskykurniawan15/payrolls/services/audit_trail"
	"github.com/riskykurniawan15/payrolls/utils/jwt"
)

type AuditTrailMiddleware struct {
	auditTrailService auditTrailService.IAuditTrailService
}

func NewAuditTrailMiddleware(auditTrailService auditTrailService.IAuditTrailService) *AuditTrailMiddleware {
	return &AuditTrailMiddleware{
		auditTrailService: auditTrailService,
	}
}

var sensitivePaths = []string{
	"/auth/login",
	"/health",
}

func (m *AuditTrailMiddleware) AuditTrail() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Capture request body
			var requestBody []byte
			if c.Request().Body != nil {
				requestBody, _ = io.ReadAll(c.Request().Body)
				// Restore body for handler
				c.Request().Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			// Create custom response writer to capture response
			responseWriter := &responseCapture{
				ResponseWriter: c.Response().Writer,
				body:           &bytes.Buffer{},
			}
			c.Response().Writer = responseWriter

			// Execute handler
			err := next(c)

			// Calculate response time
			responseTime := int(time.Since(start).Milliseconds())

			// Get user context from JWT if available
			var userID *uint
			if user := c.Get("user"); user != nil {
				if userContext, ok := user.(jwt.UserContext); ok {
					userID = &userContext.UserID
				}
			}

			// Prepare audit trail request
			auditReq := &audit_trail.CreateAuditTrailRequest{
				IP:             getClientIP(c),
				Method:         c.Request().Method,
				Path:           c.Request().URL.Path,
				UserID:         userID,
				ResponseCode:   c.Response().Status,
				ResponseTimeMs: &responseTime,
				UserAgent:      getStringPtr(c.Request().UserAgent()),
			}

			// Handle request payload (exclude sensitive data)
			if len(requestBody) > 0 && !isSensitiveEndpoint(c.Request().URL.Path) {
				payloadStr := string(requestBody)
				auditReq.Payload = &payloadStr
			}

			// Handle error response
			if err != nil {
				errorMsg := err.Error()
				auditReq.ErrorResponse = &errorMsg
			} else if c.Response().Status >= 400 {
				// Capture response body for error status codes
				responseBody := responseWriter.body.String()
				if responseBody != "" {
					auditReq.ErrorResponse = &responseBody
				}
			}

			// Log audit trail asynchronously to avoid blocking response
			go func() {
				ctx := context.Background()
				_ = m.auditTrailService.LogRequest(ctx, auditReq)
			}()

			return err
		}
	}
}

// responseCapture captures response body
type responseCapture struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rc *responseCapture) Write(b []byte) (int, error) {
	rc.body.Write(b)
	return rc.ResponseWriter.Write(b)
}

// Helper functions
func getClientIP(c echo.Context) string {
	// Check for forwarded headers first
	if ip := c.Request().Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := c.Request().Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return c.RealIP()
}

func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func isSensitiveEndpoint(path string) bool {
	for _, sensitivePath := range sensitivePaths {
		if strings.Contains(path, sensitivePath) {
			return true
		}
	}
	return false
}
