//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/riskykurniawan15/payrolls/config"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
	periodRepositories "github.com/riskykurniawan15/payrolls/repositories/period"
	periodDetailRepositories "github.com/riskykurniawan15/payrolls/repositories/period_detail"
	userRepositories "github.com/riskykurniawan15/payrolls/repositories/user"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"gorm.io/gorm"

	attendanceRepositories "github.com/riskykurniawan15/payrolls/repositories/attendance"
	auditTrailRepositories "github.com/riskykurniawan15/payrolls/repositories/audit_trail"
	overtimeRepositories "github.com/riskykurniawan15/payrolls/repositories/overtime"
	reimbursementRepositories "github.com/riskykurniawan15/payrolls/repositories/reimbursement"
	attendanceServices "github.com/riskykurniawan15/payrolls/services/attendance"
	auditTrailServices "github.com/riskykurniawan15/payrolls/services/audit_trail"
	healthServices "github.com/riskykurniawan15/payrolls/services/health"
	overtimeServices "github.com/riskykurniawan15/payrolls/services/overtime"
	payslipServices "github.com/riskykurniawan15/payrolls/services/payslip"
	periodServices "github.com/riskykurniawan15/payrolls/services/period"
	periodDetailServices "github.com/riskykurniawan15/payrolls/services/period_detail"
	reimbursementServices "github.com/riskykurniawan15/payrolls/services/reimbursement"
	userServices "github.com/riskykurniawan15/payrolls/services/user"

	attendanceHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/attendance"
	healthHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/health"
	overtimeHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/overtime"
	payslipHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/payslip"
	periodHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/period"
	periodDetailHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/period_detail"
	reimbursementHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/reimbursement"
	userHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/user"
)

type Dependencies struct {
	HealthHandlers        healthHandlers.IHealthHandler
	UserHandlers          userHandlers.IUserHandler
	PeriodHandlers        periodHandlers.IPeriodHandler
	PeriodDetailHandlers  periodDetailHandlers.IPeriodDetailHandler
	AttendanceHandlers    attendanceHandlers.IAttendanceHandler
	OvertimeHandlers      overtimeHandlers.IOvertimeHandler
	ReimbursementHandlers reimbursementHandlers.IReimbursementHandler
	PayslipHandlers       payslipHandlers.IPayslipHandler
	AuditTrailService     auditTrailServices.IAuditTrailService
}

func InitializeHandler(db *gorm.DB, cfg config.Config, logger logger.Logger) *Dependencies {
	wire.Build(
		RepositorySet,
		ServicesSet,
		HandlerSet,
		wire.Struct(new(Dependencies), "*"),
	)
	return nil
}

var RepositorySet = wire.NewSet(
	healthRepositories.NewHealthRepositories,
	userRepositories.NewUserRepository,
	periodRepositories.NewPeriodRepository,
	periodDetailRepositories.NewPeriodDetailRepository,
	attendanceRepositories.NewAttendanceRepository,
	auditTrailRepositories.NewAuditTrailRepository,
	overtimeRepositories.NewOvertimeRepository,
	reimbursementRepositories.NewReimbursementRepository,
)

var ServicesSet = wire.NewSet(
	healthServices.NewHealthService,
	userServices.NewUserService,
	periodServices.NewPeriodService,
	periodDetailServices.NewPeriodDetailService,
	attendanceServices.NewAttendanceService,
	auditTrailServices.NewAuditTrailService,
	overtimeServices.NewOvertimeService,
	reimbursementServices.NewReimbursementService,
	payslipServices.NewPayslipService,
)

var HandlerSet = wire.NewSet(
	healthHandlers.NewHealthHandlers,
	userHandlers.NewUserHandlers,
	periodHandlers.NewPeriodHandlers,
	periodDetailHandlers.NewPeriodDetailHandlers,
	attendanceHandlers.NewAttendanceHandlers,
	overtimeHandlers.NewOvertimeHandlers,
	reimbursementHandlers.NewReimbursementHandlers,
	payslipHandlers.NewPayslipHandlers,
)
