package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver สำหรับ PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Driver สำหรับอ่านจากไฟล์
	"github.com/joho/godotenv"

	"go-template/pkg/config" // Import config loader ของเรา
)

func main() {
	// 1. โหลด .env สำหรับ Local Development
	// ตอนรันใน Docker Compose, env var จะถูกฉีดเข้ามาโดยตรงอยู่แล้ว
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found, using OS environment variables")
	}

	// 2. รับคำสั่งจาก Command Line (Flags)
	var dbName, migrationPath, action string
	flag.StringVar(&dbName, "db", "", "Name of the database to migrate (e.g., primary, logs)")
	flag.StringVar(&migrationPath, "path", "", "Path to the migration files (e.g., db/migrations/primary)")
	flag.StringVar(&action, "action", "up", "Migration action: up or down")
	flag.Parse()

	// ตรวจสอบว่าผู้ใช้ใส่ Flag ที่จำเป็นมาครบหรือไม่
	if dbName == "" || migrationPath == "" {
		log.Fatalf("❌ Both --db and --path flags are required! Usage: go run ./cmd/migrate/main.go --db=<name> --path=<path>")
	}

	log.Printf("🚀 Starting migration for database: '%s' | Action: '%s'", dbName, action)

	// 3. โหลด Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// 4. เลือก Connection String (DSN) ที่ถูกต้องตาม Flag ที่ได้รับมา
	var dsn string
	switch dbName {
	case "primary":
		dsn = cfg.Postgres.Primary.BuildDSN()
	case "logs":
		dsn = cfg.Postgres.Logs.BuildDSN()
	// ในอนาคตถ้ามี DB อื่นๆ ก็มาเพิ่ม case ที่นี่
	default:
		log.Fatalf("❌ Unknown database name: '%s'. Must be one of 'primary', 'logs'", dbName)
	}

	log.Printf("📁 Using migration files from: '%s'", migrationPath)

	// 5. สร้าง instance ของ migrate
	// สังเกตว่า path จะต้องขึ้นต้นด้วย file://
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dsn)
	if err != nil {
		log.Fatalf("❌ Failed to create migrate instance: %v", err)
	}

	// 6. รัน Action ตามที่ได้รับมา
	var migrationErr error
	switch action {
	case "up":
		migrationErr = m.Up()
	case "down":
		// สั่งให้ย้อนกลับไป 1 step
		migrationErr = m.Steps(-1)
	default:
		log.Fatalf("❌ Unknown action: '%s'. Must be 'up' or 'down'", action)
	}

	// 7. จัดการผลลัพธ์
	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		log.Fatalf("❌ Failed to apply migrations: %v", migrationErr)
	}

	if errors.Is(migrationErr, migrate.ErrNoChange) {
		log.Println("✅ No new migrations to apply.")
	} else {
		log.Println("✅ Database migration completed successfully!")
	}
}
