//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
	"gorm.io/gorm"

	healthServices "github.com/riskykurniawan15/payrolls/services/health"

	healthHandlers "github.com/riskykurniawan15/payrolls/infrastructure/http/handler/health"
)

type Dependencies struct {
	HealthHandlers healthHandlers.IHealthHandler
}

func InitializeHandler(db *gorm.DB) *Dependencies {
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
)

var ServicesSet = wire.NewSet(
	healthServices.NewHealthService,
)

var HandlerSet = wire.NewSet(
	healthHandlers.NewHealthHandlers,
)
