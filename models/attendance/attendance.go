package attendance

import (
	"time"

	"github.com/riskykurniawan15/payrolls/utils/data_tipes"
)

type (
	// AttendanceRequest represents the request for check-in
	AttendanceRequest struct {
		Date *data_tipes.CustomDateTime `json:"date,omitempty"` // Optional, if not provided use current time
	}

	// AttendanceResponse represents the attendance response
	AttendanceResponse struct {
		ID           uint       `json:"id"`
		UserID       uint       `json:"user_id"`
		CheckInDate  time.Time  `json:"check_in_date"`
		CheckOutDate *time.Time `json:"check_out_date"`
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
