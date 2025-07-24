package data_tipes

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomDateTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid datetime",
			json:    `"2024-01-15 10:30:45"`,
			want:    time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "empty string",
			json:    `""`,
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "null value",
			json:    `null`,
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "invalid format",
			json:    `"2024/01/15 10:30:45"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid datetime",
			json:    `"2024-13-45 25:70:80"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "date only",
			json:    `"2024-01-15"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "with milliseconds",
			json:    `"2024-01-15 10:30:45.123"`,
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cdt CustomDateTime
			err := json.Unmarshal([]byte(tt.json), &cdt)

			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDateTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !cdt.Time.Equal(tt.want) && !(cdt.Time.IsZero() && tt.want.IsZero()) {
					t.Errorf("CustomDateTime.UnmarshalJSON() = %v, want %v", cdt.Time, tt.want)
				}
			}
		})
	}
}

func TestCustomDateTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		datetime CustomDateTime
		want     string
		wantErr  bool
	}{
		{
			name:     "valid datetime",
			datetime: CustomDateTime{Time: time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)},
			want:     `"2024-01-15 10:30:45"`,
			wantErr:  false,
		},
		{
			name:     "zero time",
			datetime: CustomDateTime{Time: time.Time{}},
			want:     `null`,
			wantErr:  false,
		},
		{
			name:     "different datetime",
			datetime: CustomDateTime{Time: time.Date(2023, 12, 31, 23, 59, 59, 0, time.Local)},
			want:     `"2023-12-31 23:59:59"`,
			wantErr:  false,
		},
		{
			name:     "midnight",
			datetime: CustomDateTime{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.Local)},
			want:     `"2024-06-15 00:00:00"`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.datetime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDateTime.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(got) != tt.want {
					t.Errorf("CustomDateTime.MarshalJSON() = %v, want %v", string(got), tt.want)
				}
			}
		})
	}
}

func TestCustomDateTime_Integration(t *testing.T) {
	// Test round-trip marshaling and unmarshaling
	originalDateTime := CustomDateTime{Time: time.Date(2024, 6, 15, 14, 30, 25, 0, time.Local)}

	// Marshal
	jsonData, err := json.Marshal(originalDateTime)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaledDateTime CustomDateTime
	err = json.Unmarshal(jsonData, &unmarshaledDateTime)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if !originalDateTime.Time.Equal(unmarshaledDateTime.Time) {
		t.Errorf("Round-trip failed: original = %v, unmarshaled = %v", originalDateTime.Time, unmarshaledDateTime.Time)
	}
}

func TestCustomDateTime_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "leap year datetime",
			json:    `"2024-02-29 12:00:00"`,
			wantErr: false,
		},
		{
			name:    "non-leap year february 29",
			json:    `"2023-02-29 12:00:00"`,
			wantErr: true,
		},
		{
			name:    "single digit components",
			json:    `"2024-1-5 2:3:4"`,
			wantErr: true,
		},
		{
			name:    "extra spaces",
			json:    `" 2024-01-15 10:30:45 "`,
			wantErr: true,
		},
		{
			name:    "24 hour format",
			json:    `"2024-01-15 24:00:00"`,
			wantErr: true,
		},
		{
			name:    "60 minutes",
			json:    `"2024-01-15 10:60:45"`,
			wantErr: true,
		},
		{
			name:    "60 seconds",
			json:    `"2024-01-15 10:30:60"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cdt CustomDateTime
			err := json.Unmarshal([]byte(tt.json), &cdt)

			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDateTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomDateTime_TimeZoneHandling(t *testing.T) {
	// Test that the datetime is parsed in local timezone
	jsonStr := `"2024-01-15 10:30:45"`

	var cdt CustomDateTime
	err := json.Unmarshal([]byte(jsonStr), &cdt)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check that the timezone is local
	_, offset := cdt.Time.Zone()
	_, localOffset := time.Now().Zone()

	if offset != localOffset {
		t.Errorf("Timezone mismatch: got offset %d, want %d", offset, localOffset)
	}
}
