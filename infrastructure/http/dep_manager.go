//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/riskykurniawan15/payrolls/config"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
	userRepositories "github.com/riskykurniawan15/payrolls/repositories/user"
	"gorm.io/gorm"

	healthServices "github.com/riskykurniawan15/payrolls/services/health"
	userServices "github.com/riskykurniawan15/payrolls/services/user"

	healthHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/health"
	userHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/user"
)

type Dependencies struct {
	HealthHandlers healthHandlers.IHealthHandler
	UserHandlers   userHandlers.IUserHandler
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
)

var ServicesSet = wire.NewSet(
	healthServices.NewHealthService,
	userServices.NewUserService,
)

var HandlerSet = wire.NewSet(
	healthHandlers.NewHealthHandlers,
	userHandlers.NewUserHandlers,
)
