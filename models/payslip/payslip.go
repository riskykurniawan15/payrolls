package payslip

import (
	"time"
)

type (
	// PayslipListRequest for listing payslips with filters
	PayslipListRequest struct {
		Page     int    `json:"page" validate:"min=1"`
		Limit    int    `json:"limit" validate:"min=1,max=100"`
		SortBy   string `json:"sort_by" validate:"omitempty,oneof=id periods_id start_date end_date take_home_pay created_at"`
		SortDesc bool   `json:"sort_desc"`
	}

	// PayslipListResponse for paginated response
	PayslipListResponse struct {
		Data       []PayslipSummary `json:"data"`
		Pagination Pagination       `json:"pagination"`
	}

	// PayslipSummary for list view
	PayslipSummary struct {
		ID          uint      `json:"id"`
		PeriodsID   uint      `json:"periods_id"`
		PeriodName  string    `json:"period_name"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
		TakeHomePay float64   `json:"take_home_pay"`
		CreatedAt   time.Time `json:"created_at"`
	}

	// GeneratePayslipRequest for generating payslip
	GeneratePayslipRequest struct {
		PeriodDetailID uint `json:"period_detail_id" validate:"required"`
	}

	// GeneratePayslipResponse for API response
	GeneratePayslipResponse struct {
		PrintURL  string    `json:"print_url"`
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	// PayslipData for HTML template
	PayslipData struct {
		CompanyName        string              `json:"company_name"`
		EmployeeName       string              `json:"employee_name"`
		PeriodName         string              `json:"period_name"`
		StartDate          time.Time           `json:"start_date"`
		EndDate            time.Time           `json:"end_date"`
		TotalWorking       int                 `json:"total_working"`
		DailyRate          float64             `json:"daily_rate"`
		BaseSalary         float64             `json:"base_salary"`
		OvertimeDetails    []OvertimeData      `json:"overtime_details"`
		TotalOvertime      float64             `json:"total_overtime"`
		Reimbursements     []ReimbursementData `json:"reimbursements"`
		TotalReimbursement float64             `json:"total_reimbursement"`
		TakeHomePay        float64             `json:"take_home_pay"`
		GeneratedAt        time.Time           `json:"generated_at"`
	}

	// OvertimeDetail for payslip
	OvertimeData struct {
		ID     uint    `json:"id"`
		Date   string  `json:"date"`
		Hours  float64 `json:"hours"`
		Amount float64 `json:"amount"`
	}

	// ReimbursementDetail for payslip
	ReimbursementData struct {
		ID     uint    `json:"id"`
		Title  string  `json:"title"`
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
	}

	// Pagination info
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
)
