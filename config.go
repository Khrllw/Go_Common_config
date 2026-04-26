package config

import (
	"net"
	"net/url"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAppName    = "user-service"
	defaultAppVersion = "1.0.0"
	defaultAppEnv     = "dev"

	defaultHTTPPort              = "6001"
	defaultGinMode               = "release"
	defaultHTTPReadTimeout       = 10 * time.Second
	defaultHTTPReadHeaderTimeout = 10 * time.Second
	defaultHTTPWriteTimeout      = 15 * time.Second
	defaultHTTPIdleTimeout       = 60 * time.Second
	defaultHTTPShutdownTimeout   = 20 * time.Second

	defaultDBHost            = "localhost"
	defaultDBPort            = "5432"
	defaultDBUser            = "postgres"
	defaultDBPassword        = "postgres"
	defaultDBName            = "user_service"
	defaultDBSchema          = "users"
	defaultDBSSLMode         = "disable"
	defaultDBMaxOpenConns    = 20
	defaultDBMaxIdleConns    = 10
	defaultDBConnMaxLifetime = 30 * time.Minute

	defaultLogLevel      = "info"
	defaultSwaggerPath   = "/swagger"
	defaultSwaggerEnable = false
)

// Config — корневая структура конфигурации приложения.
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	Swagger  SwaggerConfig
}

// AppConfig содержит метаданные сервиса.
type AppConfig struct {
	Name    string
	Version string
	Env     string
}

// HTTPConfig содержит runtime-настройки HTTP-сервера.
type HTTPConfig struct {
	Port              string
	GinMode           string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	AllowedOrigins    []string
}

// DatabaseConfig содержит параметры подключения к БД и пула соединений.
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	Schema          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoggerConfig содержит настройки логирования.
type LoggerConfig struct {
	Level string
}

// SwaggerConfig содержит настройки Swagger UI.
type SwaggerConfig struct {
	Enabled bool
	Path    string
}

// DSN формирует строку подключения к PostgreSQL в URI-формате.
func (c DatabaseConfig) DSN() string {
	query := url.Values{}
	query.Set("sslmode", c.SSLMode)

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Name,
		RawQuery: query.Encode(),
	}).String()
}

// LoadConfig загружает конфигурацию из переменных окружения.
// Файл .env используется как опциональный источник значений для локальной разработки.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", defaultAppName),
			Version: getEnv("APP_VERSION", defaultAppVersion),
			Env:     getEnv("APP_ENV", defaultAppEnv),
		},
		HTTP: HTTPConfig{
			Port:              getEnv("HTTP_PORT", defaultHTTPPort),
			GinMode:           getEnv("GIN_MODE", defaultGinMode),
			ReadTimeout:       getEnvAsDuration("HTTP_READ_TIMEOUT", defaultHTTPReadTimeout),
			ReadHeaderTimeout: getEnvAsDuration("HTTP_READ_HEADER_TIMEOUT", defaultHTTPReadHeaderTimeout),
			WriteTimeout:      getEnvAsDuration("HTTP_WRITE_TIMEOUT", defaultHTTPWriteTimeout),
			IdleTimeout:       getEnvAsDuration("HTTP_IDLE_TIMEOUT", defaultHTTPIdleTimeout),
			ShutdownTimeout:   getEnvAsDuration("HTTP_SHUTDOWN_TIMEOUT", defaultHTTPShutdownTimeout),
			AllowedOrigins:    getEnvAsList("HTTP_ALLOWED_ORIGINS", []string{"*"}),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", defaultDBHost),
			Port:            getEnv("DB_PORT", defaultDBPort),
			User:            getEnv("DB_USER", defaultDBUser),
			Password:        getEnv("DB_PASSWORD", defaultDBPassword),
			Name:            getEnv("DB_NAME", defaultDBName),
			Schema:          getEnv("DB_SCHEMA", defaultDBSchema),
			SSLMode:         getEnv("DB_SSLMODE", defaultDBSSLMode),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", defaultDBMaxOpenConns),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", defaultDBMaxIdleConns),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", defaultDBConnMaxLifetime),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", defaultLogLevel),
		},
		Swagger: SwaggerConfig{
			Enabled: getEnvAsBool("SWAGGER_ENABLED", defaultSwaggerEnable),
			Path:    getEnv("SWAGGER_PATH", defaultSwaggerPath),
		},
	}

	cfg.normalize()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
