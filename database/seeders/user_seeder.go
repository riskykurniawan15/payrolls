package seeders

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct {
	db *gorm.DB
}

const defaultCreatedBy = 1
const defaultCreatedEmployee = 100

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

func (s *UserSeeder) Seed() error {
	log.Println("Seeding users...")

	// Create admin user first
	if err := s.createAdminUser(); err != nil {
		return err
	}

	// Create employee users
	if err := s.createEmployeeUsers(); err != nil {
		return err
	}

	log.Println("Users seeding completed!")
	return nil
}

func (s *UserSeeder) createAdminUser() error {
	// Hash password for admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := User{
		ID:        defaultCreatedBy,
		Username:  "admin",
		Password:  string(hashedPassword),
		Roles:     "admin",
		Salary:    0,
		CreatedBy: defaultCreatedBy, // Self-reference for admin
		CreatedAt: time.Now(),
	}

	// Check if admin already exists
	var count int64
	s.db.Model(&User{}).Where("username = ?", adminUser.Username).Count(&count)

	if count == 0 {
		if err := s.db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Printf("Created admin user: %s", adminUser.Username)
	} else {
		log.Printf("Admin user already exists: %s", adminUser.Username)
	}

	return nil
}

func (s *UserSeeder) createEmployeeUsers() error {
	log.Println("Creating", defaultCreatedEmployee, "employee users...")

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Salary ranges for employees
	salaryRanges := []struct {
		min, max float64
	}{
		{3000000, 5000000},   // Junior
		{5000000, 8000000},   // Mid-level
		{8000000, 12000000},  // Senior
		{12000000, 15000000}, // Lead/Manager
	}

	// Generate employees
	for i := 1; i <= defaultCreatedEmployee; i++ {
		// Generate fake name
		firstName := faker.FirstName()
		lastName := faker.LastName()

		// Create username from name (lowercase, no spaces)
		username := strings.ToLower(fmt.Sprintf("%s_%s", firstName, lastName))
		username = strings.ReplaceAll(username, " ", "")
		username = strings.ReplaceAll(username, "-", "")
		username = strings.ReplaceAll(username, "'", "")

		// Hash password (same as username)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Random salary from ranges
		salaryRange := salaryRanges[i%len(salaryRanges)]
		salary := salaryRange.min + rand.Float64()*(salaryRange.max-salaryRange.min)
		salary = float64(int(salary*100)) / 100 // Round to 2 decimal places

		employeeUser := User{
			Username:  username,
			Password:  string(hashedPassword),
			Roles:     "employee",
			Salary:    salary,
			CreatedBy: defaultCreatedBy,
			CreatedAt: time.Now(),
		}

		// Check if user already exists
		var count int64
		s.db.Model(&User{}).Where("username = ?", employeeUser.Username).Count(&count)

		if count == 0 {
			if err := s.db.Create(&employeeUser).Error; err != nil {
				return err
			}
			if i%10 == 0 { // Log every 10th user
				log.Printf("Created employee %d: %s (salary: %.2f)", i, username, salary)
			}
		} else {
			if i%10 == 0 { // Log every 10th user
				log.Printf("Employee already exists: %s", username)
			}
		}
	}

	log.Printf("Created %d employee users", defaultCreatedEmployee)
	return nil
}

// User model for seeder
type User struct {
	ID        uint    `gorm:"primaryKey"`
	Username  string  `gorm:"unique;not null"`
	Password  string  `gorm:"not null"`
	Roles     string  `gorm:"not null"`
	Salary    float64 `gorm:"not null"`
	CreatedBy uint
	CreatedAt time.Time
	UpdatedBy *uint
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false"`
}

func (User) TableName() string {
	return "users"
}
