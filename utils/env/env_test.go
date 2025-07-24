package env

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "non-existent file",
			filename: "non_existent.env",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  false, // godotenv.Load("") doesn't return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadEnv(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEnv_String(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal string
		want       string
		clearEnv   bool
	}{
		{
			name:       "existing env var",
			key:        "TEST_STRING",
			value:      "test_value",
			defaultVal: "default_value",
			want:       "test_value",
		},
		{
			name:       "non-existing env var",
			key:        "NON_EXISTENT",
			value:      "",
			defaultVal: "default_value",
			want:       "default_value",
			clearEnv:   true,
		},
		{
			name:       "empty env var",
			key:        "EMPTY_STRING",
			value:      "",
			defaultVal: "default_value",
			want:       "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_Int(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal int
		want       int
		clearEnv   bool
	}{
		{
			name:       "valid integer",
			key:        "TEST_INT",
			value:      "42",
			defaultVal: 0,
			want:       42,
		},
		{
			name:       "invalid integer",
			key:        "INVALID_INT",
			value:      "not_a_number",
			defaultVal: 100,
			want:       100,
		},
		{
			name:       "non-existing env var",
			key:        "NON_EXISTENT_INT",
			value:      "",
			defaultVal: 50,
			want:       50,
			clearEnv:   true,
		},
		{
			name:       "negative integer",
			key:        "NEGATIVE_INT",
			value:      "-10",
			defaultVal: 0,
			want:       -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_Uint64(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal uint64
		want       uint64
		clearEnv   bool
	}{
		{
			name:       "valid uint64",
			key:        "TEST_UINT64",
			value:      "18446744073709551615",
			defaultVal: 0,
			want:       18446744073709551615,
		},
		{
			name:       "invalid uint64",
			key:        "INVALID_UINT64",
			value:      "not_a_number",
			defaultVal: 100,
			want:       100,
		},
		{
			name:       "negative number",
			key:        "NEGATIVE_UINT64",
			value:      "-10",
			defaultVal: 50,
			want:       50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_Int64(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal int64
		want       int64
		clearEnv   bool
	}{
		{
			name:       "valid int64",
			key:        "TEST_INT64",
			value:      "9223372036854775807",
			defaultVal: 0,
			want:       9223372036854775807,
		},
		{
			name:       "negative int64",
			key:        "NEGATIVE_INT64",
			value:      "-9223372036854775808",
			defaultVal: 0,
			want:       -9223372036854775808,
		},
		{
			name:       "invalid int64",
			key:        "INVALID_INT64",
			value:      "not_a_number",
			defaultVal: 100,
			want:       100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_Bool(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal bool
		want       bool
		clearEnv   bool
	}{
		{
			name:       "true value",
			key:        "TEST_BOOL_TRUE",
			value:      "true",
			defaultVal: false,
			want:       true,
		},
		{
			name:       "false value",
			key:        "TEST_BOOL_FALSE",
			value:      "false",
			defaultVal: true,
			want:       false,
		},
		{
			name:       "1 as true",
			key:        "TEST_BOOL_ONE",
			value:      "1",
			defaultVal: false,
			want:       true,
		},
		{
			name:       "0 as false",
			key:        "TEST_BOOL_ZERO",
			value:      "0",
			defaultVal: true,
			want:       false,
		},
		{
			name:       "invalid bool",
			key:        "INVALID_BOOL",
			value:      "maybe",
			defaultVal: true,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_Float64(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal float64
		want       float64
		clearEnv   bool
	}{
		{
			name:       "valid float",
			key:        "TEST_FLOAT",
			value:      "3.14159",
			defaultVal: 0.0,
			want:       3.14159,
		},
		{
			name:       "negative float",
			key:        "NEGATIVE_FLOAT",
			value:      "-2.718",
			defaultVal: 0.0,
			want:       -2.718,
		},
		{
			name:       "integer as float",
			key:        "INT_AS_FLOAT",
			value:      "42",
			defaultVal: 0.0,
			want:       42.0,
		},
		{
			name:       "invalid float",
			key:        "INVALID_FLOAT",
			value:      "not_a_number",
			defaultVal: 1.5,
			want:       1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.value)
			}

			got := GetEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv_UnknownType(t *testing.T) {
	// Test with an unknown type (should return default value)
	type CustomType struct {
		Value string
	}

	defaultVal := CustomType{Value: "default"}

	// Set an environment variable
	os.Setenv("TEST_CUSTOM", "some_value")

	got := GetEnv("TEST_CUSTOM", defaultVal)
	if got != defaultVal {
		t.Errorf("GetEnv() = %v, want %v", got, defaultVal)
	}
}

func TestGetEnv_Cleanup(t *testing.T) {
	// Clean up environment variables set during tests
	envVars := []string{
		"TEST_STRING", "EMPTY_STRING", "TEST_INT", "INVALID_INT", "NEGATIVE_INT",
		"TEST_UINT64", "INVALID_UINT64", "NEGATIVE_UINT64", "TEST_INT64", "NEGATIVE_INT64", "INVALID_INT64",
		"TEST_BOOL_TRUE", "TEST_BOOL_FALSE", "TEST_BOOL_ONE", "TEST_BOOL_ZERO", "INVALID_BOOL",
		"TEST_FLOAT", "NEGATIVE_FLOAT", "INT_AS_FLOAT", "INVALID_FLOAT", "TEST_CUSTOM",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
