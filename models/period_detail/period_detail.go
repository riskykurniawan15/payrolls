package period_detail

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type (
	// PeriodDetail model
	PeriodDetail struct {
		ID                  uint       `json:"id" gorm:"primaryKey"`
		PeriodsID           uint       `json:"periods_id" gorm:"not null"`
		UserID              uint       `json:"user_id" gorm:"not null"`
		DailyRate           float64    `json:"daily_rate" gorm:"type:decimal(15,2);not null;default:0.00"`
		TotalWorking        int        `json:"total_working" gorm:"not null;default:0"`
		AmountSalary        float64    `json:"amount_salary" gorm:"type:decimal(15,2);not null;default:0.00"`
		Overtime            *JSON      `json:"overtime" gorm:"type:jsonb"`
		AmountOvertime      float64    `json:"amount_overtime" gorm:"type:decimal(15,2);not null;default:0.00"`
		Reimbursement       *JSON      `json:"reimbursement" gorm:"type:jsonb"`
		AmountReimbursement float64    `json:"amount_reimbursement" gorm:"type:decimal(15,2);not null;default:0.00"`
		TakeHomePay         float64    `json:"take_home_pay" gorm:"type:decimal(15,2);not null;default:0.00"`
		CreatedBy           uint       `json:"created_by" gorm:"not null"`
		CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
		UpdatedBy           *uint      `json:"updated_by" gorm:"default:null"`
		UpdatedAt           *time.Time `json:"updated_at" gorm:"autoUpdateTime:false"`
	}

	// JSON type for handling JSONB fields
	JSON json.RawMessage

	// RunPayrollRequest for running payroll
	RunPayrollRequest struct {
		UserExecutablePayroll uint `json:"user_executable_payroll" validate:"required"`
	}

	// RunPayrollResponse for API response
	RunPayrollResponse struct {
		Message string `json:"message"`
		JobID   string `json:"job_id"`
	}
)

func (PeriodDetail) TableName() string {
	return "period_details"
}

// Value implements the driver.Valuer interface for JSON
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

// Scan implements the sql.Scanner interface for JSON
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
	case string:
		*j = JSON([]byte(v))
	default:
		return json.Unmarshal([]byte(v.(string)), j)
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = JSON(data)
	return nil
}
