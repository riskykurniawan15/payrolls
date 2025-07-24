package data_tipes

import (
	"fmt"
	"strings"
	"time"
)

// CustomDateTime represents a custom time type that handles YYYY-MM-DD HH:MM:SS format
type CustomDateTime struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for CustomDateTime
func (ct *CustomDateTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string
	str := strings.Trim(string(data), `"`)

	// If empty string, return nil
	if str == "" || str == "null" {
		return nil
	}

	// Parse using the specific format: YYYY-MM-DD HH:MM:SS
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local); err == nil {
		ct.Time = t
		return nil
	}

	return fmt.Errorf("unable to parse date: %s. Expected format: YYYY-MM-DD HH:MM:SS (e.g., 2024-01-01 08:00:00)", str)
}

// MarshalJSON implements custom JSON marshaling for CustomDateTime
func (ct CustomDateTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, ct.Time.Format("2006-01-02 15:04:05"))), nil
}
