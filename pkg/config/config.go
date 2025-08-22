package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	External ExternalConfig `mapstructure:"external"`
	App      AppConfig      `mapstructure:"app"`
	Docker   DockerConfig   `mapstructure:"docker"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type DockerConfig struct {
	Env     string `mapstructure:"env"`
	DevPort string `mapstructure:"dev_port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	TimeZone string `mapstructure:"timezone"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type ExternalConfig struct {
	DHL DHLConfig `mapstructure:"dhl"`
}

type DHLConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

// LoadConfig loads configuration from multiple sources
func LoadConfig() (*Config, error) {
	// Enable environment variables first
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Map environment variables to config structure
	bindEnvVars()

	// Try to read .env file as environment variables (not as config file)
	if err := viper.ReadInConfig(); err != nil {
		// .env file is optional, continue without it
		fmt.Printf("Info: No config file found, using environment variables only\n")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

// bindEnvVars maps environment variables to config fields
func bindEnvVars() {
	// Server
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.mode", "FIBER_MODE")

	// App
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")

	// Docker
	viper.BindEnv("docker.env", "DOCKER_ENV")
	viper.BindEnv("docker.dev_port", "DEV_PORT")

	// Database (Primary - backward compatible)
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	viper.BindEnv("database.timezone", "DB_TIMEZONE")

	// Database - Primary (alternative path)
	viper.BindEnv("database.primary.host", "DB_HOST")
	viper.BindEnv("database.primary.port", "DB_PORT")
	viper.BindEnv("database.primary.user", "DB_USER")
	viper.BindEnv("database.primary.password", "DB_PASSWORD")
	viper.BindEnv("database.primary.name", "DB_NAME")
	viper.BindEnv("database.primary.ssl_mode", "DB_SSL_MODE")
	viper.BindEnv("database.primary.timezone", "DB_TIMEZONE")

	// Database - Analytics (separate server) - lowercase keys for viper
	viper.BindEnv("database.analytics.host", "ANALYTICS_DB_HOST")
	viper.BindEnv("database.analytics.port", "ANALYTICS_DB_PORT")
	viper.BindEnv("database.analytics.user", "ANALYTICS_DB_USER")
	viper.BindEnv("database.analytics.password", "ANALYTICS_DB_PASSWORD")
	viper.BindEnv("database.analytics.name", "ANALYTICS_DB_NAME")
	viper.BindEnv("database.analytics.ssl_mode", "ANALYTICS_DB_SSL_MODE")
	viper.BindEnv("database.analytics.timezone", "ANALYTICS_DB_TIMEZONE")

	// Database - ANALYTICS (uppercase for function compatibility)
	viper.BindEnv("ANALYTICS.host", "ANALYTICS_DB_HOST")
	viper.BindEnv("ANALYTICS.port", "ANALYTICS_DB_PORT")
	viper.BindEnv("ANALYTICS.user", "ANALYTICS_DB_USER")
	viper.BindEnv("ANALYTICS.password", "ANALYTICS_DB_PASSWORD")
	viper.BindEnv("ANALYTICS.dbname", "ANALYTICS_DB_NAME")
	viper.BindEnv("ANALYTICS.sslmode", "ANALYTICS_DB_SSL_MODE")
	viper.BindEnv("ANALYTICS.timezone", "ANALYTICS_DB_TIMEZONE")

	// Database - Logs (separate server)
	viper.BindEnv("database.logs.host", "LOGS_DB_HOST")
	viper.BindEnv("database.logs.port", "LOGS_DB_PORT")
	viper.BindEnv("database.logs.user", "LOGS_DB_USER")
	viper.BindEnv("database.logs.password", "LOGS_DB_PASSWORD")
	viper.BindEnv("database.logs.name", "LOGS_DB_NAME")
	viper.BindEnv("database.logs.ssl_mode", "LOGS_DB_SSL_MODE")
	viper.BindEnv("database.logs.timezone", "LOGS_DB_TIMEZONE")

	// Database - LOGS (uppercase for function compatibility)
	viper.BindEnv("LOGS.host", "LOGS_DB_HOST")
	viper.BindEnv("LOGS.port", "LOGS_DB_PORT")
	viper.BindEnv("LOGS.user", "LOGS_DB_USER")
	viper.BindEnv("LOGS.password", "LOGS_DB_PASSWORD")
	viper.BindEnv("LOGS.dbname", "LOGS_DB_NAME")
	viper.BindEnv("LOGS.sslmode", "LOGS_DB_SSL_MODE")
	viper.BindEnv("LOGS.timezone", "LOGS_DB_TIMEZONE")

	// Database - Reports (separate server)
	viper.BindEnv("database.reports.host", "REPORTS_DB_HOST")
	viper.BindEnv("database.reports.port", "REPORTS_DB_PORT")
	viper.BindEnv("database.reports.user", "REPORTS_DB_USER")
	viper.BindEnv("database.reports.password", "REPORTS_DB_PASSWORD")
	viper.BindEnv("database.reports.name", "REPORTS_DB_NAME")
	viper.BindEnv("database.reports.ssl_mode", "REPORTS_DB_SSL_MODE")
	viper.BindEnv("database.reports.timezone", "REPORTS_DB_TIMEZONE")

	// Database - REPORTS (uppercase for function compatibility)
	viper.BindEnv("REPORTS.host", "REPORTS_DB_HOST")
	viper.BindEnv("REPORTS.port", "REPORTS_DB_PORT")
	viper.BindEnv("REPORTS.user", "REPORTS_DB_USER")
	viper.BindEnv("REPORTS.password", "REPORTS_DB_PASSWORD")
	viper.BindEnv("REPORTS.dbname", "REPORTS_DB_NAME")
	viper.BindEnv("REPORTS.sslmode", "REPORTS_DB_SSL_MODE")
	viper.BindEnv("REPORTS.timezone", "REPORTS_DB_TIMEZONE")

	// Redis (default)
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.max_retries", "REDIS_MAX_RETRIES")
	viper.BindEnv("redis.pool_size", "REDIS_POOL_SIZE")
	viper.BindEnv("redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")

	// Redis - REDIS (uppercase for function compatibility)
	viper.BindEnv("REDIS.host", "REDIS_HOST")
	viper.BindEnv("REDIS.port", "REDIS_PORT")
	viper.BindEnv("REDIS.password", "REDIS_PASSWORD")
	viper.BindEnv("REDIS.db", "REDIS_DB")
	viper.BindEnv("REDIS.max_retries", "REDIS_MAX_RETRIES")
	viper.BindEnv("REDIS.pool_size", "REDIS_POOL_SIZE")
	viper.BindEnv("REDIS.min_idle_conns", "REDIS_MIN_IDLE_CONNS")

	// Redis - Cache (separate instance)
	viper.BindEnv("redis.cache.host", "CACHE_REDIS_HOST")
	viper.BindEnv("redis.cache.port", "CACHE_REDIS_PORT")
	viper.BindEnv("redis.cache.password", "CACHE_REDIS_PASSWORD")
	viper.BindEnv("redis.cache.db", "CACHE_REDIS_DB")
	viper.BindEnv("redis.cache.max_retries", "CACHE_REDIS_MAX_RETRIES")
	viper.BindEnv("redis.cache.pool_size", "CACHE_REDIS_POOL_SIZE")
	viper.BindEnv("redis.cache.min_idle_conns", "CACHE_REDIS_MIN_IDLE_CONNS")

	// Redis - CACHE (uppercase for function compatibility)
	viper.BindEnv("CACHE.host", "CACHE_REDIS_HOST")
	viper.BindEnv("CACHE.port", "CACHE_REDIS_PORT")
	viper.BindEnv("CACHE.password", "CACHE_REDIS_PASSWORD")
	viper.BindEnv("CACHE.db", "CACHE_REDIS_DB")
	viper.BindEnv("CACHE.max_retries", "CACHE_REDIS_MAX_RETRIES")
	viper.BindEnv("CACHE.pool_size", "CACHE_REDIS_POOL_SIZE")
	viper.BindEnv("CACHE.min_idle_conns", "CACHE_REDIS_MIN_IDLE_CONNS")

	// Auth
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")

	// External APIs
	viper.BindEnv("external.dhl.base_url", "DHL_API_URL")
	viper.BindEnv("external.dhl.api_key", "DHL_API_KEY")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("app.name", "Go Template")
	viper.SetDefault("app.version", "v1.0.0")
	viper.SetDefault("docker.env", "false")
	viper.SetDefault("docker.dev_port", "8080")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.timezone", "Asia/Bangkok")
	viper.SetDefault("database.primary.ssl_mode", "disable")
	viper.SetDefault("database.primary.timezone", "Asia/Bangkok")
}
