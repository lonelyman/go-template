package platform

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitPostgres initializes default postgres connection
func InitPostgres() (*gorm.DB, error) {
	return InitPostgresWithName("database")
}

// InitPostgresWithName initializes postgres connection with custom configuration name
func InitPostgresWithName(connectionName string) (*gorm.DB, error) {
	envPrefix := strings.ToUpper(connectionName)

	// Get configuration values based on connection name
	var host string
	var port int
	var user, password, dbname, sslmode, timezone string

	if connectionName == "database" {
		// Default database connection
		host = viper.GetString("database.host")
		port = viper.GetInt("database.port")
		user = viper.GetString("database.user")
		password = viper.GetString("database.password")
		dbname = viper.GetString("database.name")
		sslmode = viper.GetString("database.ssl_mode")
		timezone = viper.GetString("database.timezone")
	} else {
		// Named connections use uppercase keys
		host = viper.GetString(fmt.Sprintf("%s.host", envPrefix))
		port = viper.GetInt(fmt.Sprintf("%s.port", envPrefix))
		user = viper.GetString(fmt.Sprintf("%s.user", envPrefix))
		password = viper.GetString(fmt.Sprintf("%s.password", envPrefix))
		dbname = viper.GetString(fmt.Sprintf("%s.dbname", envPrefix))
		sslmode = viper.GetString(fmt.Sprintf("%s.sslmode", envPrefix))
		timezone = viper.GetString(fmt.Sprintf("%s.timezone", envPrefix))
	}

	// Set default values if not provided
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 5432
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	if timezone == "" {
		timezone = "UTC"
	}

	// Build connection string for GORM
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open database connection with GORM
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection (%s): %w", envPrefix, err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB (%s): %w", envPrefix, err)
	}

	// Configure connection pool
	maxOpenConns := viper.GetInt(fmt.Sprintf("%s.max_open_conns", connectionName))
	maxIdleConns := viper.GetInt(fmt.Sprintf("%s.max_idle_conns", connectionName))
	maxLifetime := viper.GetDuration(fmt.Sprintf("%s.max_lifetime", connectionName))

	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	if maxIdleConns == 0 {
		maxIdleConns = 25
	}
	if maxLifetime == 0 {
		maxLifetime = 5 * time.Minute
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	log.Printf("Successfully connected to PostgreSQL database (%s) at %s:%d/%s", envPrefix, host, port, dbname)

	return db, nil
}
