package audit_trail

import (
	"time"
)

type AuditTrail struct {
	ID             uint      `json:"id" db:"id"`
	IP             string    `json:"ip" db:"ip"`
	Method         string    `json:"method" db:"method"`
	Path           string    `json:"path" db:"path"`
	UserID         *uint     `json:"user_id" db:"user_id"`
	Payload        *string   `json:"payload" db:"payload"`
	ResponseCode   int       `json:"response_code" db:"response_code"`
	ErrorResponse  *string   `json:"error_response" db:"error_response"`
	ResponseTimeMs *int      `json:"response_time_ms" db:"response_time_ms"`
	UserAgent      *string   `json:"user_agent" db:"user_agent"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type CreateAuditTrailRequest struct {
	IP             string  `json:"ip" validate:"required"`
	Method         string  `json:"method" validate:"required"`
	Path           string  `json:"path" validate:"required"`
	UserID         *uint   `json:"user_id"`
	Payload        *string `json:"payload"`
	ResponseCode   int     `json:"response_code" validate:"required"`
	ErrorResponse  *string `json:"error_response"`
	ResponseTimeMs *int    `json:"response_time_ms"`
	UserAgent      *string `json:"user_agent"`
}
