package config

import (
	"os"
	"strconv"
)

// GetEnvAsString gets the environment variable defined by key as an string
func GetEnvAsString(key string, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

// GetEnvAsInt gets the environment variable defined by key as an integer
func GetEnvAsInt(key string, def int) int {
	valueStr := GetEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return def
}

// GetEnvAsBool gets the environment variable defined by key as a boolean
func GetEnvAsBool(key string, def bool) bool {
	valueStr := GetEnvAsString(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return def
}
