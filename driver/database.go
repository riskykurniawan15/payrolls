package driver

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/riskykurniawan15/payrolls/config"
)

func ConnectDB(cfg config.PostgressDB) *gorm.DB {
	defer func() {
		if r := recover(); r != nil {
			log.Panic(fmt.Sprint(r))
		}
	}()

	log.Println("Connection to database")

	// Determine SSL mode string
	sslMode := "disable"
	if cfg.SSLMode {
		sslMode = "require"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.DBServer,
			cfg.DBUser,
			cfg.DBPass,
			cfg.DBName,
			cfg.DBPort,
			sslMode,
			cfg.DBTimeZone,
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic("Failed to Connect Postgress")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get underlying sql.DB")
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleCon)                                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenCon)                                  // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifeTime) * time.Minute) // Maximum lifetime of connections

	log.Printf("Database connection pool configured: MaxIdle=%d, MaxOpen=%d, MaxLifetime=%dmin",
		cfg.DBMaxIdleCon, cfg.DBMaxOpenCon, cfg.DBMaxLifeTime)

	return db
}
