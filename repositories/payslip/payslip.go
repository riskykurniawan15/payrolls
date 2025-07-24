package payslip

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/models/payslip"
	periodRepo "github.com/riskykurniawan15/payrolls/repositories/period"
	periodDetailRepo "github.com/riskykurniawan15/payrolls/repositories/period_detail"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"gorm.io/gorm"
)

type (
	IPayslipRepository interface {
		List(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error)
		GetPayslipData(ctx context.Context, periodDetailID, userID uint) (*payslip.PayslipData, error)
	}

	PayslipRepository struct {
		db               *gorm.DB
		periodDetailRepo periodDetailRepo.IPeriodDetailRepository
		periodRepo       periodRepo.IPeriodRepository
		userRepo         userRepo.IUserRepository
	}

	// OvertimeData for JSON unmarshaling
	OvertimeData struct {
		ID     uint    `json:"id"`
		Date   string  `json:"date"`
		Hours  float64 `json:"hours"`
		Amount float64 `json:"amount"`
	}

	// ReimbursementData for JSON unmarshaling
	ReimbursementData struct {
		ID     uint    `json:"id"`
		Title  string  `json:"title"`
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
	}
)

func NewPayslipRepository(
	db *gorm.DB,
	periodDetailRepo periodDetailRepo.IPeriodDetailRepository,
	periodRepo periodRepo.IPeriodRepository,
	userRepo userRepo.IUserRepository,
) IPayslipRepository {
	return &PayslipRepository{
		db:               db,
		periodDetailRepo: periodDetailRepo,
		periodRepo:       periodRepo,
		userRepo:         userRepo,
	}
}

func (r *PayslipRepository) List(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error) {
	// Build query to get period details with period information
	query := r.db.WithContext(ctx).
		Table("period_details").
		Select(`
			period_details.id,
			period_details.periods_id,
			periods.name as period_name,
			periods.start_date,
			periods.end_date,
			period_details.take_home_pay,
			period_details.created_at
		`).
		Joins("JOIN periods ON period_details.periods_id = periods.id").
		Where("period_details.user_id = ?", userID)

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count payslips: %w", err)
	}

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		// Default sort by start_date descending
		query = query.Order("periods.start_date DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	var summaries []payslip.PayslipSummary
	if err := query.Find(&summaries).Error; err != nil {
		return nil, fmt.Errorf("failed to get payslips: %w", err)
	}

	// Calculate pagination
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &payslip.PayslipListResponse{
		Data: summaries,
		Pagination: payslip.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *PayslipRepository) GetPayslipData(ctx context.Context, periodDetailID, userID uint) (*payslip.PayslipData, error) {
	// Get period detail
	periodDetail, err := r.periodDetailRepo.GetByID(ctx, periodDetailID)
	if err != nil {
		return nil, fmt.Errorf("period detail not found: %w", err)
	}

	// Check if period detail belongs to the user
	if periodDetail.UserID != userID {
		return nil, fmt.Errorf("period detail not found")
	}

	// Get period information
	period, err := r.periodRepo.GetByID(ctx, periodDetail.PeriodsID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Get user information
	user, err := r.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse overtime data
	var overtimeDetails []payslip.OvertimeData
	if periodDetail.Overtime != nil {
		var overtimeData []OvertimeData
		if err := json.Unmarshal(*periodDetail.Overtime, &overtimeData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal overtime data: %w", err)
		}

		// Convert to payslip format
		for _, ot := range overtimeData {
			overtimeDetails = append(overtimeDetails, payslip.OvertimeData{
				ID:     ot.ID,
				Date:   ot.Date,
				Hours:  ot.Hours,
				Amount: ot.Amount,
			})
		}
	}

	// Parse reimbursement data
	var reimbursements []payslip.ReimbursementData
	if periodDetail.Reimbursement != nil {
		var reimbursementData []ReimbursementData
		if err := json.Unmarshal(*periodDetail.Reimbursement, &reimbursementData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal reimbursement data: %w", err)
		}

		// Convert to payslip format
		for _, reimb := range reimbursementData {
			reimbursements = append(reimbursements, payslip.ReimbursementData{
				ID:     reimb.ID,
				Title:  reimb.Title,
				Date:   reimb.Date,
				Amount: reimb.Amount,
			})
		}
	}

	return &payslip.PayslipData{
		EmployeeName:       user.Username,
		PeriodName:         period.Name,
		StartDate:          period.StartDate,
		EndDate:            period.EndDate,
		TotalWorking:       periodDetail.TotalWorking,
		DailyRate:          periodDetail.DailyRate,
		BaseSalary:         periodDetail.AmountSalary,
		OvertimeDetails:    overtimeDetails,
		TotalOvertime:      periodDetail.AmountOvertime,
		Reimbursements:     reimbursements,
		TotalReimbursement: periodDetail.AmountReimbursement,
		TakeHomePay:        periodDetail.TakeHomePay,
		GeneratedAt:        time.Now(),
	}, nil
}
