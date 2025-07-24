package code_generator

import (
	"math/rand"
	"time"
)

const (
	// Characters for code generation (uppercase letters only)
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// GenerateRandomCode generates a random 5-character uppercase code
func GenerateRandomCode(codeLength int) string {
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
func GeneratePeriodCode(prefix string, codeLength int) string {
	now := time.Now()
	dateStr := now.Format("20060102") // YYYYMMDD format
	randomCode := GenerateRandomCode(codeLength)

	return prefix + "/" + dateStr + "/" + randomCode
}
