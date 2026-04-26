package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// getEnv возвращает значение переменной окружения без пробелов по краям.
// Если значение пустое, возвращается defaultValue.
func getEnv(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt парсит целое значение из переменной окружения.
// При ошибке парсинга возвращается defaultValue.
func getEnvAsInt(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// getEnvAsBool парсит bool-значение из переменной окружения.
// При ошибке парсинга возвращается defaultValue.
func getEnvAsBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// getEnvAsDuration парсит длительность в формате Go (например, "15s")
// или как целое число секунд.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(value)
	if err == nil {
		return parsed
	}

	asInt, intErr := strconv.Atoi(value)
	if intErr == nil {
		return time.Duration(asInt) * time.Second
	}

	return defaultValue
}

// getEnvAsList парсит список значений, разделенных запятыми,
// и возвращает непустые элементы без пробелов по краям.
func getEnvAsList(key string, defaultValue []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return cloneStringSlice(defaultValue)
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		item := strings.TrimSpace(p)
		if item != "" {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return cloneStringSlice(defaultValue)
	}
	return out
}

// cloneStringSlice создает копию слайса, чтобы избежать побочных эффектов
// при изменении исходного значения.
func cloneStringSlice(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	return out
}
