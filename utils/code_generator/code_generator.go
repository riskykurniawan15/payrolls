package code_generator

import (
	"math/rand"
	"time"
)

const (
	// Characters for code generation (uppercase letters only)
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Code length
	codeLength = 5
)

// GenerateRandomCode generates a random 5-character uppercase code
func GenerateRandomCode() string {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Generate random code
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

// GeneratePeriodCode generates a period code with format PRD/YYYYMMDD/{random_code}
func GeneratePeriodCode() string {
	now := time.Now()
	dateStr := now.Format("20060102") // YYYYMMDD format
	randomCode := GenerateRandomCode()

	return "PRD/" + dateStr + "/" + randomCode
}

// IsCodeExists checks if a code already exists (case-insensitive)
// This function should be implemented in the repository layer
// For now, we'll just return the function signature
func IsCodeExists(code string) bool {
	// This will be implemented in the repository
	// Check if code exists where status != 9 (case-insensitive)
	return false
}
