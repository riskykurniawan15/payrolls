package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	return godotenv.Load(file)
}

func GetEnv[T any](key string, defaultVal T) T {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	switch any(defaultVal).(type) {
	case int:
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return any(parsed).(T)
	case uint64:
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return defaultVal
		}
		return any(parsed).(T)
	case int64:
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return defaultVal
		}
		return any(parsed).(T)
	case bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return defaultVal
		}
		return any(parsed).(T)
	case float64:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return defaultVal
		}
		return any(parsed).(T)
	case string:
		return any(val).(T)
	default:
		return defaultVal
	}
}
