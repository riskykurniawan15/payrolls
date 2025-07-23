package overtime

import (
	"time"
)

type (
	// Overtime model
	Overtime struct {
		ID             uint       `json:"id" gorm:"primaryKey"`
		UserID         uint       `json:"user_id" gorm:"not null"`
		OvertimesDate  time.Time  `json:"overtimes_date" gorm:"not null"`
		TotalHoursTime float64    `json:"total_hours_time" gorm:"type:decimal(5,2);not null;default:0.00"`
		CreatedBy      uint       `json:"created_by" gorm:"not null"`
		CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
		UpdatedBy      *uint      `json:"updated_by" gorm:"default:null"`
		UpdatedAt      *time.Time `json:"updated_at" gorm:"autoUpdateTime:false"`
	}

	// CreateOvertimeRequest for creating new overtime
	CreateOvertimeRequest struct {
		OvertimesDate  string  `json:"overtimes_date" validate:"required,datetime=2006-01-02"`
		TotalHoursTime float64 `json:"total_hours_time" validate:"required,min=0.01,max=3.00"`
	}

	// UpdateOvertimeRequest for updating overtime
	UpdateOvertimeRequest struct {
		OvertimesDate  *string  `json:"overtimes_date" validate:"omitempty,datetime=2006-01-02"`
		TotalHoursTime *float64 `json:"total_hours_time" validate:"omitempty,min=0.01,max=3.00"`
	}

	// OvertimeResponse for API responses
	OvertimeResponse struct {
		ID             uint       `json:"id"`
		UserID         uint       `json:"user_id"`
		OvertimesDate  time.Time  `json:"overtimes_date"`
		TotalHoursTime float64    `json:"total_hours_time"`
		CreatedBy      uint       `json:"created_by"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedBy      *uint      `json:"updated_by"`
		UpdatedAt      *time.Time `json:"updated_at"`
	}

	// ListOvertimesRequest for listing overtimes with filters
	ListOvertimesRequest struct {
		Page      int     `json:"page" validate:"min=1"`
		Limit     int     `json:"limit" validate:"min=1,max=100"`
		StartDate *string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
		EndDate   *string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
		SortBy    string  `json:"sort_by" validate:"omitempty,oneof=id user_id overtimes_date total_hours_time created_at"`
		SortDesc  bool    `json:"sort_desc"`
	}

	// ListOvertimesResponse for paginated response
	ListOvertimesResponse struct {
		Data       []OvertimeResponse `json:"data"`
		Pagination Pagination         `json:"pagination"`
	}

	// Pagination info
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
)

func (Overtime) TableName() string {
	return "overtimes"
}
