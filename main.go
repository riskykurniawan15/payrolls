package main

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/driver"
	http "github.com/riskykurniawan15/payrolls/infrastructure/http/engine"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

func main() {
	// Load configuration
	cfg := config.Configuration()

	// Connect to database using existing driver
	db := driver.ConnectDB(cfg.PostgressDB)

	// Initialize logger with config
	log := logger.NewLoggerWithConfig(logger.LoggerConfig{
		OutputMode: cfg.Logger.OutputMode,
		LogLevel:   cfg.Logger.LogLevel,
		LogDir:     cfg.Logger.LogDir,
	})

	// Defer logger close to ensure all logs are written
	defer log.Close()

	// Defer database cleanup with timeout
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sqlDB, err := db.WithContext(ctx).DB()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get DB connection: %v", err))
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Error(fmt.Sprintf("Failed to close DB connection: %v", err))
			} else {
				log.Info("Database connection closed successfully")
			}
		}
	}()

	// Start HTTP server with graceful shutdown
	http.Start(http.App{
		Config:      cfg,
		Logger:      log,
		PostgressDB: db,
	})
}
