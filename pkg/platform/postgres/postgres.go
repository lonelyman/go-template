package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-template/pkg/config"
)

// NewConnection คือ Public Function "เดียว" ที่เราจะเปิดให้ข้างนอกเรียกใช้
// มันรับพิมพ์เขียว (PostgresConfig) เข้ามา และคืนค่าเป็น DB Connection ที่พร้อมใช้งาน
func NewConnection(cfg config.PostgresConfig) (*gorm.DB, error) {
	// 1. สร้าง Connection String จาก Config ที่ได้รับมา
	dsn := cfg.BuildDSN() // เรียกใช้ Helper ที่เราสร้างไว้ใน pkg/config

	// 2. ตั้งค่า GORM Logger (สามารถปรับเปลี่ยนได้ตาม Environment)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log ทุก Query
	}

	// 3. เปิดการเชื่อมต่อ
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// 4. ⭐️ ตั้งค่า Connection Pool (นำไอเดียดีๆ ของน้องมาไว้ที่นี่) ⭐️
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("✅ Successfully connected to PostgreSQL: %s/%s", cfg.Host, cfg.DBName)

	return db, nil
}
