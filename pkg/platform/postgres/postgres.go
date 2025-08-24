package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"go-template/pkg/config"
	"go-template/pkg/logger"
)

// gormLoggerAdapter คือ "หัวแปลงปลั๊ก" ที่ทำให้ Logger ของเราคุยกับ GORM ได้
type gormLoggerAdapter struct {
	appLogger logger.Logger
}

// Implement gormlogger.Interface
func (l *gormLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l // เราจะจัดการ level เองข้างล่าง
}
func (l *gormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	l.appLogger.Info(fmt.Sprintf(msg, data...))
}
func (l *gormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.appLogger.Warn(fmt.Sprintf(msg, data...))
}
func (l *gormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	l.appLogger.Error(fmt.Sprintf(msg, data...), nil)
}
func (l *gormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// เราสามารถเพิ่ม Logic การ log SQL query ที่นี่ได้ถ้าต้องการ
}

// NewConnection คือ Public Function ของเรา
// ✨ 2. แก้ไขให้รับ appLogger เข้ามาด้วย ✨
func NewConnection(cfg config.PostgresConfig, appLogger logger.Logger) (*gorm.DB, error) {
	dsn := cfg.BuildDSN()

	// ✨ 3. สร้าง GORM Logger ที่ใช้ "หัวแปลงปลั๊ก" ของเรา ✨
	newLogger := &gormLoggerAdapter{appLogger: appLogger}

	gormConfig := &gorm.Config{
		Logger: newLogger.LogMode(gormlogger.Info), // ตั้งค่าให้ GORM ใช้ Logger ใหม่ของเรา
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	appLogger.Info("Successfully connected to PostgreSQL", "host", cfg.Host, "dbName", cfg.DBName)

	return db, nil
}
