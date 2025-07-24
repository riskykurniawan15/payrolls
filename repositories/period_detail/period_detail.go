package period_detail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/models/payslip"
	"github.com/riskykurniawan15/payrolls/models/period_detail"
	periodRepo "github.com/riskykurniawan15/payrolls/repositories/period"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"gorm.io/gorm"
)

type (
	IPeriodDetailRepository interface {
		Create(ctx context.Context, periodDetail *period_detail.PeriodDetail) error
		GetByID(ctx context.Context, id uint) (*period_detail.PeriodDetail, error)
		GetByPeriodAndUser(ctx context.Context, periodID, userID uint) (*period_detail.PeriodDetail, error)
		Update(ctx context.Context, id uint, updates map[string]interface{}) error
		Delete(ctx context.Context, id uint) error
		DeleteByPeriodID(ctx context.Context, periodID uint) error
		GetUsersByBatch(ctx context.Context, lastID uint, limit int) ([]uint, error)
		CreateBatch(ctx context.Context, periodDetails []period_detail.PeriodDetail) error
		ListPayslip(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error)
		GetPayslipData(ctx context.Context, periodDetailID, userID uint) (*payslip.PayslipData, error)
		GetPayslipSummaryData(ctx context.Context, periodID uint) (*payslip.PayslipSummaryData, error)
	}

	PeriodDetailRepository struct {
		db         *gorm.DB
		periodRepo periodRepo.IPeriodRepository
		userRepo   userRepo.IUserRepository
	}
)

func NewPeriodDetailRepository(db *gorm.DB, periodRepo periodRepo.IPeriodRepository, userRepo userRepo.IUserRepository) IPeriodDetailRepository {
	return &PeriodDetailRepository{
		db:         db,
		periodRepo: periodRepo,
		userRepo:   userRepo,
	}
}

func (repo PeriodDetailRepository) getInstanceDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constant.TransactionKey).(*gorm.DB)
	if !ok {
		return repo.db
	}
	return tx
}

func (repo PeriodDetailRepository) Create(ctx context.Context, periodDetail *period_detail.PeriodDetail) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(periodDetail).Error
}

func (repo PeriodDetailRepository) GetByID(ctx context.Context, id uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("id = ?", id).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (repo PeriodDetailRepository) GetByPeriodAndUser(ctx context.Context, periodID, userID uint) (*period_detail.PeriodDetail, error) {
	var periodDetail period_detail.PeriodDetail
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	err := repo.getInstanceDB(ctx).WithContext(ctxWT).Where("periods_id = ? AND user_id = ?", periodID, userID).First(&periodDetail).Error
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (repo PeriodDetailRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Model(&period_detail.PeriodDetail{}).Where("id = ?", id).Updates(updates).Error
}

func (repo PeriodDetailRepository) Delete(ctx context.Context, id uint) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Delete(&period_detail.PeriodDetail{}, id).Error
}

func (repo PeriodDetailRepository) DeleteByPeriodID(ctx context.Context, periodID uint) error {
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	return repo.getInstanceDB(ctx).WithContext(ctxWT).Delete(&period_detail.PeriodDetail{}, "periods_id = ?", periodID).Error
}

func (repo PeriodDetailRepository) GetUsersByBatch(ctx context.Context, lastID uint, limit int) ([]uint, error) {
	var userIDs []uint
	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()
	query := repo.getInstanceDB(ctx).WithContext(ctxWT).Table("users").Select("id").Where("roles = ?", constant.EmployeeRole)

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	err := query.Order("id ASC").Limit(limit).Pluck("id", &userIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by batch: %w", err)
	}

	return userIDs, nil
}

func (repo PeriodDetailRepository) CreateBatch(ctx context.Context, periodDetails []period_detail.PeriodDetail) error {
	if len(periodDetails) == 0 {
		return nil
	}

	ctxWT, cancel := context.WithTimeout(ctx, constant.DBTimeout)
	defer cancel()

	return repo.getInstanceDB(ctx).WithContext(ctxWT).Create(&periodDetails).Error
}

func (repo PeriodDetailRepository) ListPayslip(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error) {
	// Build query to get period details with period information
	query := repo.getInstanceDB(ctx).WithContext(ctx).
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

func (repo PeriodDetailRepository) GetPayslipData(ctx context.Context, periodDetailID, userID uint) (*payslip.PayslipData, error) {
	// Get period detail
	periodDetail, err := repo.GetByID(ctx, periodDetailID)
	if err != nil {
		return nil, fmt.Errorf("period detail not found: %w", err)
	}

	// Check if period detail belongs to the user
	if periodDetail.UserID != userID {
		return nil, fmt.Errorf("period detail not found")
	}

	// Get period information
	period, err := repo.periodRepo.GetByID(ctx, periodDetail.PeriodsID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Get user information
	user, err := repo.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse overtime data
	var overtimeDetails []payslip.OvertimeData
	if periodDetail.Overtime != nil {
		json.Unmarshal(*periodDetail.Overtime, &overtimeDetails)
	}

	// Parse reimbursement data
	var reimbursements []payslip.ReimbursementData
	if periodDetail.Reimbursement != nil {
		json.Unmarshal(*periodDetail.Reimbursement, &reimbursements)
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

func (repo PeriodDetailRepository) GetPayslipSummaryData(ctx context.Context, periodID uint) (*payslip.PayslipSummaryData, error) {
	// Get period information
	period, err := repo.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Get all period details for this period
	query := repo.getInstanceDB(ctx).WithContext(ctx).
		Table("period_details").
		Select(`
			period_details.id,
			period_details.user_id,
			period_details.total_working,
			period_details.take_home_pay,
			users.username as employee_name
		`).
		Joins("JOIN users ON period_details.user_id = users.id").
		Where("period_details.periods_id = ?", periodID).
		Order("users.username ASC")

	var results []struct {
		ID           uint    `json:"id"`
		UserID       uint    `json:"user_id"`
		TotalWorking int     `json:"total_working"`
		TakeHomePay  float64 `json:"take_home_pay"`
		EmployeeName string  `json:"employee_name"`
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get period details: %w", err)
	}

	// Calculate summary data
	totalEmployees := len(results)
	totalTakeHomePay := float64(0)
	var employeeList []payslip.PayslipSummaryEmployee

	for i, result := range results {
		totalTakeHomePay += result.TakeHomePay
		employeeList = append(employeeList, payslip.PayslipSummaryEmployee{
			No:           i + 1,
			EmployeeName: result.EmployeeName,
			TakeHomePay:  result.TakeHomePay,
		})
	}

	// Calculate average working days (assuming all employees have same working days)
	totalWorkingDays := 0
	if totalEmployees > 0 {
		totalWorkingDays = results[0].TotalWorking
	}

	return &payslip.PayslipSummaryData{
		CompanyName:      "Company Name",
		PeriodName:       period.Name,
		TotalEmployees:   totalEmployees,
		TotalWorkingDays: totalWorkingDays,
		TotalTakeHomePay: totalTakeHomePay,
		EmployeeList:     employeeList,
		GeneratedAt:      time.Now(),
	}, nil
}
