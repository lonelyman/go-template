package tests

import (
	"log"
	"os"
	"testing"

	"go-template/pkg/platform"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "7430")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "12345678")
	os.Setenv("DB_NAME", "postgres")

	// Initialize test database
	db, err := platform.InitPostgres()
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Run migrations
	// db.AutoMigrate(&example_module.ExampleDomain{})

	// Run tests
	code := m.Run()

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Exit(code)
}
