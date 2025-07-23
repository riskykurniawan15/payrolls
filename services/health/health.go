package health

import (
	"context"

	healthModels "github.com/riskykurniawan15/payrolls/models/health"
	healthRepositories "github.com/riskykurniawan15/payrolls/repositories/health"
)

type (
	IHealthServices interface {
		HealthMetric(ctx context.Context) *healthModels.HealthMetric
	}
	HealthServices struct {
		healthRepositories healthRepositories.IHealthRepositories
	}
)

func NewHealthService(
	healthRepositories healthRepositories.IHealthRepositories,
) IHealthServices {
	return &HealthServices{
		healthRepositories: healthRepositories,
	}
}

func (svc HealthServices) HealthMetric(ctx context.Context) *healthModels.HealthMetric {
	var metric healthModels.HealthMetric
	status := map[string]interface{}{
		"database": "connected",
	}

	databaseHealth, err := svc.healthRepositories.DatabaseHealth(ctx)
	if err != nil {
		status["database"] = "refused"
		metric.DB = err.Error()
	} else {
		metric.DB = databaseHealth
	}

	metric.Status = status

	return &metric
}
