//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/riskykurniawan15/payrolls/config"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
	periodRepositories "github.com/riskykurniawan15/payrolls/repositories/period"
	userRepositories "github.com/riskykurniawan15/payrolls/repositories/user"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"gorm.io/gorm"

	attendanceRepositories "github.com/riskykurniawan15/payrolls/repositories/attendance"
	overtimeRepositories "github.com/riskykurniawan15/payrolls/repositories/overtime"
	attendanceServices "github.com/riskykurniawan15/payrolls/services/attendance"
	healthServices "github.com/riskykurniawan15/payrolls/services/health"
	overtimeServices "github.com/riskykurniawan15/payrolls/services/overtime"
	periodServices "github.com/riskykurniawan15/payrolls/services/period"
	userServices "github.com/riskykurniawan15/payrolls/services/user"

	attendanceHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/attendance"
	healthHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/health"
	overtimeHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/overtime"
	periodHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/period"
	userHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/user"
)

type Dependencies struct {
	HealthHandlers     healthHandlers.IHealthHandler
	UserHandlers       userHandlers.IUserHandler
	PeriodHandlers     periodHandlers.IPeriodHandler
	AttendanceHandlers attendanceHandlers.IAttendanceHandler
	OvertimeHandlers   overtimeHandlers.IOvertimeHandler
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
	attendanceRepositories.NewAttendanceRepository,
	overtimeRepositories.NewOvertimeRepository,
)

var ServicesSet = wire.NewSet(
	healthServices.NewHealthService,
	userServices.NewUserService,
	periodServices.NewPeriodService,
	attendanceServices.NewAttendanceService,
	overtimeServices.NewOvertimeService,
)

var HandlerSet = wire.NewSet(
	healthHandlers.NewHealthHandlers,
	userHandlers.NewUserHandlers,
	periodHandlers.NewPeriodHandlers,
	attendanceHandlers.NewAttendanceHandlers,
	overtimeHandlers.NewOvertimeHandlers,
)
