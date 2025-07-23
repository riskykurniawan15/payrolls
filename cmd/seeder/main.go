package main

import (
	"log"

	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/database/seeders"
	"github.com/riskykurniawan15/payrolls/driver"
)

func main() {
	// Load configuration
	cfg := config.Configuration()

	// Connect to database using existing driver
	db := driver.ConnectDB(cfg.PostgressDB)

	// Initialize seeder
	seeder := seeders.NewMainSeeder(db)

	// Run seeder
	if err := seeder.Run(); err != nil {
		log.Fatalf("Seeder error: %v", err)
	}

	log.Println("Seeder completed successfully!")
}
