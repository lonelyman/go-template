// cmd/migrate/main.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	// "go-template/pkg/config" // 👈 คอมเมนต์ config import ออกไปก่อน

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	var dbName, migrationPath string
	flag.StringVar(&dbName, "db", "primary", "Name of the database to migrate (e.g., primary)")
	flag.StringVar(&migrationPath, "path", "db/migrations/primary", "Path to the migration files (e.g., db/migrations/primary)")
	flag.Parse()

	log.Printf("🚀 Starting migration for database: '%s'", dbName)

	// 🛑🛑🛑 ลองคอมเมนต์โค้ดส่วนอ่าน Config นี้ทั้งหมดออกไปก่อนชั่วคราว 🛑🛑🛑
	/*
	   cfg, err := config.LoadConfig()
	   if err != nil {
	      log.Fatalf("❌ Could not load config: %v", err)
	   }

	   var dsn string
	   switch dbName {
	   case "primary":
	      dsn = cfg.Postgres.Primary.BuildDSN()
	   default:
	      log.Fatalf("❌ Unknown database name: '%s'", dbName)
	   }
	*/

	// ⭐️⭐️⭐️ แล้วลอง Hardcode DSN ที่ถูกต้อง 100% ลงไปตรงๆ แบบนี้เลย! ⭐️⭐️⭐️
	// (พี่เอาข้อมูลจาก .env ที่น้องเคยส่งให้พี่มาสร้างให้)
	dsn := "postgres://root:12345678@localhost:7430/go_template?sslmode=disable"

	log.Println("--- [DEBUG] FORCING CONNECTION TO:", "localhost:7430", "---")

	log.Printf("📁 Using migration files from: '%s'", migrationPath)

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dsn)
	if err != nil {
		log.Fatalf("❌ Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("❌ Failed to apply migrations: %v", err)
	}

	log.Println("✅ Database migration completed successfully!")
}
