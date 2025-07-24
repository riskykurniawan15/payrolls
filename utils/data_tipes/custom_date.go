package data_tipes

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime represents a custom time type that handles YYYY-MM-DD format
type CustomDate struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for CustomDate
func (ct *CustomDate) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string
	str := strings.Trim(string(data), `"`)

	// If empty string, return nil
	if str == "" || str == "null" {
		return nil
	}

	// Parse using the specific format: YYYY-MM-DD
	if t, err := time.ParseInLocation("2006-01-02", str, time.Local); err == nil {
		ct.Time = t
		return nil
	}

	return fmt.Errorf("unable to parse date: %s. Expected format: YYYY-MM-DD (e.g., 2024-01-01)", str)
}

// MarshalJSON implements custom JSON marshaling for CustomDate
func (ct CustomDate) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, ct.Time.Format("2006-01-02"))), nil
}
