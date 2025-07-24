package code_generator

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateRandomCode(t *testing.T) {
	tests := []struct {
		name       string
		codeLength int
	}{
		{
			name:       "5 character code",
			codeLength: 5,
		},
		{
			name:       "10 character code",
			codeLength: 10,
		},
		{
			name:       "1 character code",
			codeLength: 1,
		},
		{
			name:       "0 character code",
			codeLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GenerateRandomCode(tt.codeLength)

			// Check length
			if len(code) != tt.codeLength {
				t.Errorf("GenerateRandomCode() length = %v, want %v", len(code), tt.codeLength)
			}

			// Check if code contains only valid characters
			for _, char := range code {
				if !strings.Contains(chars, string(char)) {
					t.Errorf("GenerateRandomCode() contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestGenerateRandomCodeUniqueness(t *testing.T) {
	// Test that multiple calls generate different codes
	codes := make(map[string]bool)

	for i := 0; i < 100; i++ {
		code := GenerateRandomCode(5)
		if codes[code] {
			t.Errorf("GenerateRandomCode() generated duplicate code: %s", code)
		}
		codes[code] = true
	}
}

func TestGeneratePeriodCode(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		codeLength int
	}{
		{
			name:       "PRD prefix",
			prefix:     "PRD",
			codeLength: 5,
		},
		{
			name:       "ATT prefix",
			prefix:     "ATT",
			codeLength: 3,
		},
		{
			name:       "empty prefix",
			prefix:     "",
			codeLength: 5,
		},
		{
			name:       "long prefix",
			prefix:     "VERYLONGPREFIX",
			codeLength: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GeneratePeriodCode(tt.prefix, tt.codeLength)

			// Check if code starts with prefix
			if tt.prefix != "" && !strings.HasPrefix(code, tt.prefix) {
				t.Errorf("GeneratePeriodCode() does not start with prefix: %s", tt.prefix)
			}

			// Check if code contains date format (YYYYMMDD)
			parts := strings.Split(code, "/")
			if len(parts) < 3 {
				t.Errorf("GeneratePeriodCode() does not have expected format: %s", code)
			}

			// Check date part
			datePart := parts[1]
			if len(datePart) != 8 {
				t.Errorf("GeneratePeriodCode() date part length = %v, want 8", len(datePart))
			}

			// Try to parse date
			_, err := time.Parse("20060102", datePart)
			if err != nil {
				t.Errorf("GeneratePeriodCode() date part is not valid: %s", datePart)
			}

			// Check random code part
			randomPart := parts[2]
			if len(randomPart) != tt.codeLength {
				t.Errorf("GeneratePeriodCode() random part length = %v, want %v", len(randomPart), tt.codeLength)
			}

			// Check if random part contains only valid characters
			for _, char := range randomPart {
				if !strings.Contains(chars, string(char)) {
					t.Errorf("GeneratePeriodCode() random part contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestGeneratePeriodCodeFormat(t *testing.T) {
	prefix := "TEST"
	codeLength := 5
	code := GeneratePeriodCode(prefix, codeLength)

	// Check format: PREFIX/YYYYMMDD/RANDOMCODE
	parts := strings.Split(code, "/")
	if len(parts) != 3 {
		t.Errorf("GeneratePeriodCode() format is incorrect: %s", code)
	}

	// Check that parts[0] is the prefix
	if parts[0] != prefix {
		t.Errorf("GeneratePeriodCode() prefix = %s, want %s", parts[0], prefix)
	}

	// Check that parts[1] is a valid date (8 digits)
	if len(parts[1]) != 8 {
		t.Errorf("GeneratePeriodCode() date part length = %d, want 8", len(parts[1]))
	}

	// Check that parts[2] is the random code
	if len(parts[2]) != codeLength {
		t.Errorf("GeneratePeriodCode() random part length = %d, want %d", len(parts[2]), codeLength)
	}
}

func TestGeneratePeriodCodeUniqueness(t *testing.T) {
	// Test that multiple calls generate different codes
	codes := make(map[string]bool)
	prefix := "TEST"

	for i := 0; i < 50; i++ {
		code := GeneratePeriodCode(prefix, 5)
		if codes[code] {
			t.Errorf("GeneratePeriodCode() generated duplicate code: %s", code)
		}
		codes[code] = true
	}
}

func TestCharsConstant(t *testing.T) {
	// Test that chars constant contains expected characters
	expectedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	if chars != expectedChars {
		t.Errorf("chars constant = %s, want %s", chars, expectedChars)
	}

	// Test that chars contains only uppercase letters and numbers
	for _, char := range chars {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			t.Errorf("chars constant contains invalid character: %c", char)
		}
	}
}
