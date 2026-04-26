package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var schemaNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
var allowedGinModes = map[string]struct{}{
	"debug":   {},
	"release": {},
	"test":    {},
}
var allowedDBSSLModes = map[string]struct{}{
	"disable":     {},
	"allow":       {},
	"prefer":      {},
	"require":     {},
	"verify-ca":   {},
	"verify-full": {},
}

// normalize приводит значения конфигурации к каноничному виду перед валидацией.
func (c *Config) normalize() {
	c.App.Env = strings.ToLower(strings.TrimSpace(c.App.Env))
	c.HTTP.GinMode = strings.ToLower(strings.TrimSpace(c.HTTP.GinMode))
	c.Logger.Level = strings.ToLower(strings.TrimSpace(c.Logger.Level))
	c.Database.Schema = strings.TrimSpace(c.Database.Schema)
	c.Database.SSLMode = strings.ToLower(strings.TrimSpace(c.Database.SSLMode))
	c.Swagger.Path = normalizeSwaggerPath(c.Swagger.Path)
	c.HTTP.AllowedOrigins = normalizeOrigins(c.HTTP.AllowedOrigins)
}

// Validate выполняет раннюю валидацию обязательных и ограниченных полей.
func (c *Config) Validate() error {
	if c.Database.Schema == "" {
		return errors.New("DB_SCHEMA is required")
	}
	if !schemaNameRegex.MatchString(c.Database.Schema) {
		return fmt.Errorf("DB_SCHEMA contains invalid characters: %s", c.Database.Schema)
	}
	if c.Database.Name == "" {
		return errors.New("DB_NAME is required")
	}
	if c.HTTP.Port == "" {
		return errors.New("HTTP_PORT is required")
	}
	port, err := strconv.Atoi(c.HTTP.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("HTTP_PORT must be in range 1..65535: %s", c.HTTP.Port)
	}

	if _, ok := allowedGinModes[c.HTTP.GinMode]; !ok {
		return fmt.Errorf("GIN_MODE must be one of: debug, release, test: %s", c.HTTP.GinMode)
	}

	if c.HTTP.ReadTimeout <= 0 || c.HTTP.ReadHeaderTimeout <= 0 || c.HTTP.WriteTimeout <= 0 || c.HTTP.IdleTimeout <= 0 || c.HTTP.ShutdownTimeout <= 0 {
		return errors.New("HTTP timeouts must be positive")
	}

	if c.Database.Host == "" {
		return errors.New("DB_HOST is required")
	}
	if c.Database.User == "" {
		return errors.New("DB_USER is required")
	}
	if c.Database.Port == "" {
		return errors.New("DB_PORT is required")
	}
	dbPort, dbPortErr := strconv.Atoi(c.Database.Port)
	if dbPortErr != nil || dbPort < 1 || dbPort > 65535 {
		return fmt.Errorf("DB_PORT must be in range 1..65535: %s", c.Database.Port)
	}
	if _, ok := allowedDBSSLModes[c.Database.SSLMode]; !ok {
		return fmt.Errorf("DB_SSLMODE has unsupported value: %s", c.Database.SSLMode)
	}
	if c.Database.MaxOpenConns < 0 {
		return errors.New("DB_MAX_OPEN_CONNS must be >= 0")
	}
	if c.Database.MaxIdleConns < 0 {
		return errors.New("DB_MAX_IDLE_CONNS must be >= 0")
	}
	if c.Database.MaxOpenConns > 0 && c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return errors.New("DB_MAX_IDLE_CONNS must be <= DB_MAX_OPEN_CONNS")
	}
	if c.Database.ConnMaxLifetime < 0 {
		return errors.New("DB_CONN_MAX_LIFETIME must be >= 0")
	}
	if len(c.HTTP.AllowedOrigins) == 0 {
		return errors.New("HTTP_ALLOWED_ORIGINS must contain at least one origin")
	}
	if c.Swagger.Path == "" || !strings.HasPrefix(c.Swagger.Path, "/") {
		return errors.New("SWAGGER_PATH must start with '/'")
	}

	return nil
}

// normalizeSwaggerPath нормализует путь Swagger:
// добавляет ведущий "/", убирает хвостовые "/" и подставляет дефолт при пустом значении.
func normalizeSwaggerPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return defaultSwaggerPath
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	trimmed = strings.TrimRight(trimmed, "/")
	if trimmed == "" {
		return defaultSwaggerPath
	}
	return trimmed
}

// normalizeOrigins очищает список origins:
// удаляет пустые значения и дубликаты, при необходимости подставляет "*".
func normalizeOrigins(origins []string) []string {
	if len(origins) == 0 {
		return []string{"*"}
	}

	out := make([]string, 0, len(origins))
	seen := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
