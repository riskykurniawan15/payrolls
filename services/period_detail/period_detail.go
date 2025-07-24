package period_detail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/constant"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/period"
	"github.com/riskykurniawan15/payrolls/models/period_detail"
	attendanceRepo "github.com/riskykurniawan15/payrolls/repositories/attendance"
	instanceRepo "github.com/riskykurniawan15/payrolls/repositories/instance"
	overtimeRepo "github.com/riskykurniawan15/payrolls/repositories/overtime"
	periodRepo "github.com/riskykurniawan15/payrolls/repositories/period"
	periodDetailRepo "github.com/riskykurniawan15/payrolls/repositories/period_detail"
	reimbursementRepo "github.com/riskykurniawan15/payrolls/repositories/reimbursement"
	userRepo "github.com/riskykurniawan15/payrolls/repositories/user"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IPeriodDetailService interface {
		RunPayroll(ctx context.Context, periodID uint, userID uint) (*period_detail.RunPayrollResponse, error)
	}

	PeriodDetailService struct {
		logger            logger.Logger
		periodDetailRepo  periodDetailRepo.IPeriodDetailRepository
		periodRepo        periodRepo.IPeriodRepository
		userRepo          userRepo.IUserRepository
		attendanceRepo    attendanceRepo.IAttendanceRepository
		overtimeRepo      overtimeRepo.IOvertimeRepository
		reimbursementRepo reimbursementRepo.IReimbursementRepository
		instanceRepo      instanceRepo.IInstanceRepository
	}

	// PayrollData for storing calculation results
	PayrollData struct {
		UserID              uint                `json:"user_id"`
		DailyRate           float64             `json:"daily_rate"`
		TotalWorking        int                 `json:"total_working"`
		AmountSalary        float64             `json:"amount_salary"`
		Overtime            []OvertimeData      `json:"overtime"`
		AmountOvertime      float64             `json:"amount_overtime"`
		Reimbursement       []ReimbursementData `json:"reimbursement"`
		AmountReimbursement float64             `json:"amount_reimbursement"`
		TakeHomePay         float64             `json:"take_home_pay"`
	}

	OvertimeData struct {
		ID     uint    `json:"id"`
		Date   string  `json:"date"`
		Hours  float64 `json:"hours"`
		Amount float64 `json:"amount"`
	}

	ReimbursementData struct {
		ID     uint    `json:"id"`
		Title  string  `json:"title"`
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
	}
)

func NewPeriodDetailService(
	logger logger.Logger,
	periodDetailRepo periodDetailRepo.IPeriodDetailRepository,
	periodRepo periodRepo.IPeriodRepository,
	userRepo userRepo.IUserRepository,
	attendanceRepo attendanceRepo.IAttendanceRepository,
	overtimeRepo overtimeRepo.IOvertimeRepository,
	reimbursementRepo reimbursementRepo.IReimbursementRepository,
	instanceRepo instanceRepo.IInstanceRepository,
) IPeriodDetailService {
	return &PeriodDetailService{
		logger:            logger,
		periodDetailRepo:  periodDetailRepo,
		periodRepo:        periodRepo,
		userRepo:          userRepo,
		attendanceRepo:    attendanceRepo,
		overtimeRepo:      overtimeRepo,
		reimbursementRepo: reimbursementRepo,
		instanceRepo:      instanceRepo,
	}
}

func (s *PeriodDetailService) RunPayroll(ctx context.Context, periodID, userID uint) (*period_detail.RunPayrollResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing run payroll request", requestID, map[string]interface{}{
		"period_id": periodID,
		"user_id":   userID,
	})

	// Get period data
	periodData, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		s.logger.ErrorT("failed to get period data", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": periodID,
		})
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Check if period is active
	if periodData.Status != constant.StatusActive {
		s.logger.ErrorT("period is not active", requestID, map[string]interface{}{
			"period_id": periodID,
			"status":    periodData.Status,
		})
		return nil, fmt.Errorf("period is not active")
	}

	// Update period status to processing
	err = s.periodRepo.Update(ctx, periodID, map[string]interface{}{
		"status":                  constant.StatusProcessing,
		"user_executable_payroll": userID,
		"payroll_date":            time.Now(),
		"updated_by":              userID,
	})
	if err != nil {
		s.logger.ErrorT("failed to update period status", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": periodID,
		})
		return nil, fmt.Errorf("failed to update period status: %w", err)
	}

	// Generate job ID
	jobID := fmt.Sprintf("payroll_%d_%s", periodID, requestID)

	// Start background processing
	go s.processPayrollBackground(context.Background(), periodID, periodData, userID, jobID)

	s.logger.InfoT("payroll job started", requestID, map[string]interface{}{
		"period_id": periodID,
		"job_id":    jobID,
	})

	return &period_detail.RunPayrollResponse{
		Status: "Payroll processing started",
		JobID:  jobID,
	}, nil
}

func (s *PeriodDetailService) processPayrollBackground(c context.Context, periodID uint, periodData *period.Period, userExecutablePayroll uint, jobID string) {
	requestID := fmt.Sprintf("bg_%s", jobID)
	status := constant.StatusFailed
	s.logger.InfoT("starting background payroll processing", requestID, map[string]interface{}{
		"period_id": periodID,
		"job_id":    jobID,
	})

	ctx, tx, err := s.instanceRepo.BeginTransactionWithContext(c)
	if err != nil {
		s.logger.ErrorT("failed to begin transaction", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer tx.Rollback()

	defer func() {
		// Update period status to completed
		err := s.periodRepo.Update(ctx, periodID, map[string]interface{}{
			"status":     status,
			"updated_by": userExecutablePayroll,
		})
		if err != nil {
			s.logger.ErrorT("failed to update period status to completed", requestID, map[string]interface{}{
				"error":     err.Error(),
				"period_id": periodID,
			})
		} else {
			s.logger.InfoT("payroll processing completed", requestID, map[string]interface{}{
				"period_id": periodID,
				"job_id":    jobID,
			})
		}

		if status == constant.StatusCompleted {
			tx.Commit()
		}
	}()

	// Parse period dates
	startDate := periodData.StartDate
	endDate := periodData.EndDate

	// Process users in batches
	lastID := uint(0)
	batchSize := 50

	s.periodDetailRepo.DeleteByPeriodID(ctx, periodID)

	for {
		userIDs, err := s.periodDetailRepo.GetUsersByBatch(ctx, lastID, batchSize)
		if err != nil {
			s.logger.ErrorT("failed to get users batch", requestID, map[string]interface{}{
				"error":      err.Error(),
				"last_id":    lastID,
				"batch_size": batchSize,
			})
			return
		}

		if len(userIDs) == 0 {
			break // No more users to process
		}

		// Process batch
		err = s.processUserBatch(ctx, periodID, userIDs, startDate, endDate, userExecutablePayroll, requestID)
		if err != nil {
			s.logger.ErrorT("failed to process user batch", requestID, map[string]interface{}{
				"error":    err.Error(),
				"user_ids": userIDs,
			})
			continue
		}

		lastID = userIDs[len(userIDs)-1]
		s.logger.InfoT("processed user batch", requestID, map[string]interface{}{
			"batch_size": len(userIDs),
			"last_id":    lastID,
		})
	}

	status = constant.StatusCompleted
}

func (s *PeriodDetailService) processUserBatch(ctx context.Context, periodID uint, userIDs []uint, startDate, endDate time.Time, userExecutablePayroll uint, requestID string) error {
	// Use database transaction
	var periodDetails []period_detail.PeriodDetail

	for _, userID := range userIDs {
		payrollData, err := s.calculatePayroll(ctx, userID, startDate, endDate, requestID)
		if err != nil {
			s.logger.ErrorT("failed to calculate payroll for user", requestID, map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			continue
		}

		// Convert overtime data to JSON
		overtimeJSON, err := json.Marshal(payrollData.Overtime)
		if err != nil {
			s.logger.ErrorT("failed to marshal overtime data", requestID, map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			continue
		}

		// Convert reimbursement data to JSON
		reimbursementJSON, err := json.Marshal(payrollData.Reimbursement)
		if err != nil {
			s.logger.ErrorT("failed to marshal reimbursement data", requestID, map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			continue
		}

		periodDetail := period_detail.PeriodDetail{
			PeriodsID:           periodID,
			UserID:              userID,
			DailyRate:           payrollData.DailyRate,
			TotalWorking:        payrollData.TotalWorking,
			AmountSalary:        payrollData.AmountSalary,
			Overtime:            (*period_detail.JSON)(&overtimeJSON),
			AmountOvertime:      payrollData.AmountOvertime,
			Reimbursement:       (*period_detail.JSON)(&reimbursementJSON),
			AmountReimbursement: payrollData.AmountReimbursement,
			TakeHomePay:         payrollData.TakeHomePay,
			CreatedBy:           userExecutablePayroll,
			CreatedAt:           time.Now(),
		}

		periodDetails = append(periodDetails, periodDetail)
	}

	// Create batch records
	if len(periodDetails) > 0 {
		err := s.periodDetailRepo.CreateBatch(ctx, periodDetails)
		if err != nil {
			return fmt.Errorf("failed to create period details batch: %w", err)
		}
	}

	return nil
}

func (s *PeriodDetailService) calculatePayroll(ctx context.Context, userID uint, startDate, endDate time.Time, requestID string) (*PayrollData, error) {
	// Get user data
	userData, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	// Calculate working days from attendance data (weekdays only)
	payDay, totalWorking := 0, 0
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			// Check if user has checked in attendance for this date
			hasCheckedIn, err := s.hasCheckedInAttendance(ctx, userID, currentDate)
			if err != nil {
				s.logger.ErrorT("failed to check attendance", requestID, map[string]interface{}{
					"error":   err.Error(),
					"user_id": userID,
					"date":    currentDate.Format("2006-01-02"),
				})
			} else if hasCheckedIn {
				totalWorking++
			}
			payDay++
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Calculate daily rate (salary / working days)
	dailyRate := float64(0)
	dailyRate = userData.Salary / float64(payDay)

	// Calculate salary amount
	amountSalary := dailyRate * float64(totalWorking)

	// Get overtime data for the period
	overtimeData, amountOvertime, err := s.calculateOvertime(ctx, userID, startDate, endDate, dailyRate, requestID)
	if err != nil {
		s.logger.ErrorT("failed to calculate overtime", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		overtimeData = []OvertimeData{}
		amountOvertime = 0
	}

	// Get reimbursement data for the period
	reimbursementData, amountReimbursement, err := s.calculateReimbursement(ctx, userID, startDate, endDate, requestID)
	if err != nil {
		s.logger.ErrorT("failed to calculate reimbursement", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		reimbursementData = []ReimbursementData{}
		amountReimbursement = 0
	}

	// Calculate take home pay
	takeHomePay := amountSalary + amountOvertime + amountReimbursement

	return &PayrollData{
		UserID:              userID,
		DailyRate:           dailyRate,
		TotalWorking:        totalWorking,
		AmountSalary:        amountSalary,
		Overtime:            overtimeData,
		AmountOvertime:      amountOvertime,
		Reimbursement:       reimbursementData,
		AmountReimbursement: amountReimbursement,
		TakeHomePay:         takeHomePay,
	}, nil
}

func (s *PeriodDetailService) calculateOvertime(ctx context.Context, userID uint, startDate, endDate time.Time, dailyRate float64, requestID string) ([]OvertimeData, float64, error) {
	// Get overtime records for the period
	overtimes, err := s.overtimeRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get overtime data: %w", err)
	}

	var overtimeData []OvertimeData
	totalAmount := float64(0)

	// Calculate hourly rate (daily rate / 8 hours)
	hourlyRate := dailyRate / 8.0

	for _, ot := range overtimes {
		// Overtime rate is 2x hourly rate
		overtimeAmount := ot.TotalHoursTime * hourlyRate * 2.0
		totalAmount += overtimeAmount

		overtimeData = append(overtimeData, OvertimeData{
			ID:     ot.ID,
			Date:   ot.OvertimesDate.Format("2006-01-02"),
			Hours:  ot.TotalHoursTime,
			Amount: overtimeAmount,
		})
	}

	return overtimeData, totalAmount, nil
}

func (s *PeriodDetailService) calculateReimbursement(ctx context.Context, userID uint, startDate, endDate time.Time, requestID string) ([]ReimbursementData, float64, error) {
	// Get reimbursement records for the period
	reimbursements, err := s.reimbursementRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reimbursement data: %w", err)
	}

	var reimbursementData []ReimbursementData
	totalAmount := float64(0)

	for _, reimb := range reimbursements {
		totalAmount += reimb.Amount

		reimbursementData = append(reimbursementData, ReimbursementData{
			ID:     reimb.ID,
			Title:  reimb.Title,
			Date:   reimb.Date.Format("2006-01-02"),
			Amount: reimb.Amount,
		})
	}

	return reimbursementData, totalAmount, nil
}

func (s *PeriodDetailService) hasCheckedInAttendance(ctx context.Context, userID uint, date time.Time) (bool, error) {
	// Get attendance record for the specific date
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, startOfDay)
	if err != nil {
		// If no attendance found, return false (not an error)
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, fmt.Errorf("failed to get attendance data: %w", err)
	}

	// Check if attendance record exists (has check-in)
	return attendance != nil, nil
}
