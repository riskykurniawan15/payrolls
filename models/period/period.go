package period

import (
	"time"
)

type (
	// Period model
	Period struct {
		ID                    uint       `json:"id" gorm:"primaryKey"`
		Code                  string     `json:"code" gorm:"uniqueIndex;not null"`
		Name                  string     `json:"name" gorm:"not null"`
		StartDate             time.Time  `json:"start_date" gorm:"not null"`
		EndDate               time.Time  `json:"end_date" gorm:"not null"`
		Status                int8       `json:"status" gorm:"default:1"`
		UserExecutablePayroll *uint      `json:"user_executable_payroll" gorm:"column:user_executable_payroll"`
		PayrollDate           *time.Time `json:"payroll_date" gorm:"column:payroll_date"`
		CreatedBy             uint       `json:"created_by" gorm:"not null"`
		CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`
		UpdatedBy             *uint      `json:"updated_by" gorm:"default:null"`
		UpdatedAt             *time.Time `json:"updated_at" gorm:"autoUpdateTime:false"`
	}

	// CreatePeriodRequest for creating new period
	CreatePeriodRequest struct {
		Code      *string `json:"code" validate:"omitempty,min=3,max=50"`
		Name      string  `json:"name" validate:"required,min=3,max=100"`
		StartDate string  `json:"start_date" validate:"required,datetime=2006-01-02"`
		EndDate   string  `json:"end_date" validate:"required,datetime=2006-01-02"`
	}

	// UpdatePeriodRequest for updating period
	UpdatePeriodRequest struct {
		Code      *string `json:"code" validate:"omitempty,min=3,max=50"`
		Name      *string `json:"name" validate:"omitempty,min=3,max=100"`
		StartDate *string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
		EndDate   *string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	}

	// PeriodResponse for API responses
	PeriodResponse struct {
		ID                    uint       `json:"id"`
		Code                  string     `json:"code"`
		Name                  string     `json:"name"`
		StartDate             time.Time  `json:"start_date"`
		EndDate               time.Time  `json:"end_date"`
		Status                int8       `json:"status"`
		UserExecutablePayroll *uint      `json:"user_executable_payroll"`
		PayrollDate           *time.Time `json:"payroll_date"`
		CreatedBy             uint       `json:"created_by"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedBy             *uint      `json:"updated_by"`
		UpdatedAt             *time.Time `json:"updated_at"`
	}

	// ListPeriodsRequest for listing periods with filters
	ListPeriodsRequest struct {
		Page     int    `json:"page" validate:"min=1"`
		Limit    int    `json:"limit" validate:"min=1,max=100"`
		Search   string `json:"search"`
		Status   *int8  `json:"status"`
		SortBy   string `json:"sort_by" validate:"omitempty,oneof=id code name start_date end_date status created_at"`
		SortDesc bool   `json:"sort_desc"`
	}

	// ListPeriodsResponse for paginated response
	ListPeriodsResponse struct {
		Data       []PeriodResponse `json:"data"`
		Pagination Pagination       `json:"pagination"`
	}

	// Pagination info
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
)

func (Period) TableName() string {
	return "periods"
}
