//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/riskykurniawan15/payrolls/config"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
	periodRepositories "github.com/riskykurniawan15/payrolls/repositories/period"
	userRepositories "github.com/riskykurniawan15/payrolls/repositories/user"
	"gorm.io/gorm"

	healthServices "github.com/riskykurniawan15/payrolls/services/health"
	periodServices "github.com/riskykurniawan15/payrolls/services/period"
	userServices "github.com/riskykurniawan15/payrolls/services/user"

	healthHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/health"
	periodHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/period"
	userHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/user"
)

type Dependencies struct {
	HealthHandlers healthHandlers.IHealthHandler
	UserHandlers   userHandlers.IUserHandler
	PeriodHandlers periodHandlers.IPeriodHandler
}

func InitializeHandler(db *gorm.DB, cfg config.Config) *Dependencies {
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
)

var ServicesSet = wire.NewSet(
	healthServices.NewHealthService,
	userServices.NewUserService,
	periodServices.NewPeriodService,
)

var HandlerSet = wire.NewSet(
	healthHandlers.NewHealthHandlers,
	userHandlers.NewUserHandlers,
	periodHandlers.NewPeriodHandlers,
)
