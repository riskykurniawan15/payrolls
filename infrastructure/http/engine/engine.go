package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riskykurniawan15/payrolls/config"
	dep "github.com/riskykurniawan15/payrolls/infrastructure/http"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/router"
	"github.com/riskykurniawan15/payrolls/utils/logger"
	"gorm.io/gorm"
)

type App struct {
	Config      config.Config
	Logger      logger.Logger
	PostgressDB *gorm.DB
}

func Start(app App) {
	// Initialize HTTP server
	e := router.Routers(dep.InitializeHandler(app.PostgressDB))

	// Start HTTP server in background
	go func() {
		e.HideBanner = true
		if err := e.Start(fmt.Sprintf("%s:%d", app.Config.Http.Server, app.Config.Http.Port)); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutdown signal received")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := e.Shutdown(ctx); err != nil {
		app.Logger.Error(fmt.Sprintf("Failed to shutdown HTTP server: %v", err))
	} else {
		app.Logger.Info("HTTP server shutdown successfully")
	}

	app.Logger.Info("Application shutdown completed")
}
