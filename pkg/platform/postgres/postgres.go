package platform

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-template/pkg/config"
)

// ConnectionConfig holds database connection parameters
type ConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// setDefaults applies default values to connection config
func (c *ConnectionConfig) setDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.TimeZone == "" {
		c.TimeZone = "UTC"
	}
	if c.User == "" {
		c.User = "postgres"
	}
	if c.Password == "" {
		c.Password = "postgres"
	}
	if c.DBName == "" {
		c.DBName = "postgres"
	}
}

// buildDSN creates PostgreSQL connection string
func (c *ConnectionConfig) buildDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.TimeZone)
}

// createConnection establishes database connection with shared logic
func createConnection(connConfig ConnectionConfig, connectionName string) (*gorm.DB, error) {
	connConfig.setDefaults()
	dsn := connConfig.buildDSN()

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", connectionName, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB (%s): %w", connectionName, err)
	}

	// Set connection pool defaults
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("✅ Successfully connected to PostgreSQL database (%s) at %s:%d/%s",
		connectionName, connConfig.Host, connConfig.Port, connConfig.DBName)

	return db, nil
}

// NewConnection creates a new postgres connection using config struct
func NewConnection(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	// Convert port string to int
	port := 5432
	if dbConfig.Port != "" {
		if p, err := strconv.Atoi(dbConfig.Port); err == nil {
			port = p
		}
	}

	connConfig := ConnectionConfig{
		Host:     dbConfig.Host,
		Port:     port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		DBName:   dbConfig.Name,
		SSLMode:  dbConfig.SSLMode,
		TimeZone: dbConfig.TimeZone,
	}

	return createConnection(connConfig, "config-based")
}

// InitPostgres initializes default postgres connection
func InitPostgres() (*gorm.DB, error) {
	return InitPostgresWithName("database")
}

// InitPostgresWithName initializes postgres connection with custom configuration name
func InitPostgresWithName(connectionName string) (*gorm.DB, error) {
	envPrefix := strings.ToUpper(connectionName)

	// Get configuration values based on connection name
	var connConfig ConnectionConfig

	if connectionName == "database" {
		// Default database connection
		connConfig = ConnectionConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.name"),
			SSLMode:  viper.GetString("database.ssl_mode"),
			TimeZone: viper.GetString("database.timezone"),
		}
	} else {
		// Named connections use uppercase keys
		connConfig = ConnectionConfig{
			Host:     viper.GetString(fmt.Sprintf("%s.host", envPrefix)),
			Port:     viper.GetInt(fmt.Sprintf("%s.port", envPrefix)),
			User:     viper.GetString(fmt.Sprintf("%s.user", envPrefix)),
			Password: viper.GetString(fmt.Sprintf("%s.password", envPrefix)),
			DBName:   viper.GetString(fmt.Sprintf("%s.dbname", envPrefix)),
			SSLMode:  viper.GetString(fmt.Sprintf("%s.sslmode", envPrefix)),
			TimeZone: viper.GetString(fmt.Sprintf("%s.timezone", envPrefix)),
		}
	}

	return createConnection(connConfig, envPrefix)
}
