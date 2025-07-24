package data_tipes

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date",
			json:    `"2024-01-15"`,
			want:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
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
			json:    `"2024/01/15"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid date",
			json:    `"2024-13-45"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "with time",
			json:    `"2024-01-15 10:30:00"`,
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cd CustomDate
			err := json.Unmarshal([]byte(tt.json), &cd)

			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !cd.Time.Equal(tt.want) && !(cd.Time.IsZero() && tt.want.IsZero()) {
					t.Errorf("CustomDate.UnmarshalJSON() = %v, want %v", cd.Time, tt.want)
				}
			}
		})
	}
}

func TestCustomDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		date    CustomDate
		want    string
		wantErr bool
	}{
		{
			name:    "valid date",
			date:    CustomDate{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local)},
			want:    `"2024-01-15"`,
			wantErr: false,
		},
		{
			name:    "zero time",
			date:    CustomDate{Time: time.Time{}},
			want:    `null`,
			wantErr: false,
		},
		{
			name:    "different date",
			date:    CustomDate{Time: time.Date(2023, 12, 31, 0, 0, 0, 0, time.Local)},
			want:    `"2023-12-31"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDate.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(got) != tt.want {
					t.Errorf("CustomDate.MarshalJSON() = %v, want %v", string(got), tt.want)
				}
			}
		})
	}
}

func TestCustomDate_Integration(t *testing.T) {
	// Test round-trip marshaling and unmarshaling
	originalDate := CustomDate{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.Local)}

	// Marshal
	jsonData, err := json.Marshal(originalDate)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaledDate CustomDate
	err = json.Unmarshal(jsonData, &unmarshaledDate)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if !originalDate.Time.Equal(unmarshaledDate.Time) {
		t.Errorf("Round-trip failed: original = %v, unmarshaled = %v", originalDate.Time, unmarshaledDate.Time)
	}
}

func TestCustomDate_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "leap year date",
			json:    `"2024-02-29"`,
			wantErr: false,
		},
		{
			name:    "non-leap year february 29",
			json:    `"2023-02-29"`,
			wantErr: true,
		},
		{
			name:    "single digit month and day",
			json:    `"2024-1-5"`,
			wantErr: true,
		},
		{
			name:    "extra spaces",
			json:    `" 2024-01-15 "`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cd CustomDate
			err := json.Unmarshal([]byte(tt.json), &cd)

			if (err != nil) != tt.wantErr {
				t.Errorf("CustomDate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
