package env

import (
	"os"
	"strconv"
)

// GetString fetches a string from the environment with a default value if it doesn't exist
func GetString(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

// GetUint parses a uint from the environment with a default value if it doesn't exist
func GetUint(key string, def uint) uint {
	if val, err := strconv.ParseUint(os.Getenv(key), 10, 32); err == nil {
		return uint(val)
	}

	return def
}
