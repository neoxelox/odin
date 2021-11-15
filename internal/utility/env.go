package utility

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvAsString(key string, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return def
}

func GetEnvAsInt(key string, def int) int {
	valueStr := GetEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return def
}

func GetEnvAsBool(key string, def bool) bool {
	valueStr := GetEnvAsString(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return def
}

func GetEnvAsSlice(key string, def []string) []string {
	valueStr := GetEnvAsString(key, "")
	if value := strings.Split(valueStr, ","); len(value) >= 1 {
		return value
	}

	return def
}
