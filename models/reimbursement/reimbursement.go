package reimbursement

import (
	"time"
)

type (
	// Reimbursement model
	Reimbursement struct {
		ID          uint       `json:"id" gorm:"primaryKey"`
		UserID      uint       `json:"user_id" gorm:"not null"`
		Title       string     `json:"title" gorm:"not null"`
		Date        time.Time  `json:"date" gorm:"not null"`
		Amount      float64    `json:"amount" gorm:"type:decimal(10,2);not null"`
		Description *string    `json:"description" gorm:"default:null"`
		CreatedBy   uint       `json:"created_by" gorm:"not null"`
		CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
		UpdatedBy   *uint      `json:"updated_by" gorm:"default:null"`
		UpdatedAt   *time.Time `json:"updated_at" gorm:"autoUpdateTime:false"`
	}

	// CreateReimbursementRequest for creating new reimbursement
	CreateReimbursementRequest struct {
		Title       string  `json:"title" validate:"required,min=3,max=100"`
		Date        *string `json:"date" validate:"omitempty,datetime=2006-01-02"`
		Amount      float64 `json:"amount" validate:"required,min=0.01"`
		Description *string `json:"description" validate:"omitempty,max=500"`
	}

	// UpdateReimbursementRequest for updating reimbursement
	UpdateReimbursementRequest struct {
		Title       *string  `json:"title" validate:"omitempty,min=3,max=100"`
		Date        *string  `json:"date" validate:"omitempty,datetime=2006-01-02"`
		Amount      *float64 `json:"amount" validate:"omitempty,min=0.01"`
		Description *string  `json:"description" validate:"omitempty,max=500"`
	}

	// ReimbursementResponse for API responses
	ReimbursementResponse struct {
		ID          uint       `json:"id"`
		UserID      uint       `json:"user_id"`
		Title       string     `json:"title"`
		Date        time.Time  `json:"date"`
		Amount      float64    `json:"amount"`
		Description *string    `json:"description"`
		CreatedBy   uint       `json:"created_by"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedBy   *uint      `json:"updated_by"`
		UpdatedAt   *time.Time `json:"updated_at"`
	}

	// ListReimbursementsRequest for listing reimbursements with filters
	ListReimbursementsRequest struct {
		Page      int     `json:"page" validate:"min=1"`
		Limit     int     `json:"limit" validate:"min=1,max=100"`
		StartDate *string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
		EndDate   *string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
		SortBy    string  `json:"sort_by" validate:"omitempty,oneof=id user_id title date amount created_at"`
		SortDesc  bool    `json:"sort_desc"`
	}

	// ListReimbursementsResponse for paginated response
	ListReimbursementsResponse struct {
		Data       []ReimbursementResponse `json:"data"`
		Pagination Pagination              `json:"pagination"`
	}

	// Pagination info
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
)

func (Reimbursement) TableName() string {
	return "reimbursements"
}
