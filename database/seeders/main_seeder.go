package seeders

import (
	"flag"
	"log"
	"os"

	"gorm.io/gorm"
)

type MainSeeder struct {
	db *gorm.DB
}

func NewMainSeeder(db *gorm.DB) *MainSeeder {
	return &MainSeeder{db: db}
}

func (s *MainSeeder) Run() error {
	// Parse command line flags
	seedType := flag.String("type", "all", "Type of seeder to run (all, users, etc.)")
	flag.Parse()

	log.Printf("Starting seeder with type: %s", *seedType)

	switch *seedType {
	case "all":
		return s.seedAll()
	case "users":
		return s.seedUsers()
	default:
		log.Printf("Unknown seeder type: %s", *seedType)
		flag.Usage()
		os.Exit(1)
	}

	return nil
}

func (s *MainSeeder) seedAll() error {
	log.Println("Running all seeders...")

	// Run seeders in order
	seeders := []func() error{
		s.seedUsers,
		// Add more seeders here
		// s.seedDepartments,
		// s.seedPositions,
	}

	for _, seeder := range seeders {
		if err := seeder(); err != nil {
			return err
		}
	}

	log.Println("All seeders completed successfully!")
	return nil
}

func (s *MainSeeder) seedUsers() error {
	userSeeder := NewUserSeeder(s.db)
	return userSeeder.Seed()
}

// Command line usage
func Usage() {
	log.Println("Database Seeder Usage:")
	log.Println("  go run database/seeders/main_seeder.go -type=all     # Run all seeders")
	log.Println("  go run database/seeders/main_seeder.go -type=users   # Run only user seeder")
	log.Println("")
	log.Println("Available seeder types:")
	log.Println("  - all: Run all seeders")
	log.Println("  - users: Run user seeder only")
}
