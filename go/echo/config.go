package main

import (
	"os"
	"strconv"
)

func getEnvAsString(key string, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return def
}

func getEnvAsBool(key string, def bool) bool {
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return def
}
