package attendance

import (
	"fmt"
	"strings"
	"time"
)

type (
	// AttendanceRequest represents the request for check-in
	AttendanceRequest struct {
		Date *CustomTime `json:"date,omitempty"` // Optional, if not provided use current time
	}

	// CustomTime represents a custom time type that handles YYYY-MM-DD HH:MM:SS format
	CustomTime struct {
		time.Time
	}

	// AttendanceResponse represents the attendance response
	AttendanceResponse struct {
		ID           uint       `json:"id"`
		UserID       uint       `json:"user_id"`
		CheckInDate  time.Time  `json:"check_in_date"`
		CheckOutDate *time.Time `json:"check_out_date,omitempty"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	}

	// AttendanceListResponse represents the list of attendances
	AttendanceListResponse struct {
		Data       []AttendanceResponse `json:"data"`
		Pagination Pagination           `json:"pagination"`
	}

	// Pagination info
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}

	// Attendance represents the attendance model
	Attendance struct {
		ID           uint       `json:"id" gorm:"primaryKey;column:id"`
		UserID       uint       `json:"user_id" gorm:"column:user_id;not null"`
		CheckInDate  time.Time  `json:"check_in_date" gorm:"column:check_in_date;not null"`
		CheckOutDate *time.Time `json:"check_out_date,omitempty" gorm:"column:check_out_date"`
		CreatedBy    uint       `json:"created_by" gorm:"column:created_by;not null"`
		CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
		UpdatedBy    *uint      `json:"updated_by,omitempty" gorm:"column:updated_by"`
		UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime:false"`
	}
)

// TableName specifies the table name for Attendance model
func (Attendance) TableName() string {
	return "attendances"
}

// UnmarshalJSON implements custom JSON unmarshaling for CustomTime
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string
	str := strings.Trim(string(data), `"`)

	// If empty string, return nil
	if str == "" || str == "null" {
		return nil
	}

	// Parse using the specific format: YYYY-MM-DD HH:MM:SS
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local); err == nil {
		ct.Time = t
		return nil
	}

	return fmt.Errorf("unable to parse date: %s. Expected format: YYYY-MM-DD HH:MM:SS (e.g., 2024-01-01 08:00:00)", str)
}

// MarshalJSON implements custom JSON marshaling for CustomTime
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, ct.Time.Format("2006-01-02 15:04:05"))), nil
}

// IsWeekday checks if the given date is a weekday (Monday to Friday)
func IsWeekday(date time.Time) bool {
	weekday := date.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// IsValidCheckOutTime checks if the check-out time is valid (after check-in and before 11 PM)
func IsValidCheckOutTime(checkInDate, checkOutDate time.Time) bool {
	// Check-out must be after check-in

	return checkOutDate.After(checkInDate)
}
