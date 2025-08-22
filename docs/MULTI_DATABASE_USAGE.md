# 🗄️ Multi-Database PostgreSQL Usage Guide

_Updated: August 2025 - พร้อม Viper Configuration, Repository Pattern และ Docker Support_

## 🚀 Overview

ระบบรองรับการเชื่อมต่อ PostgreSQL หลายฐานข้อมูลพร้อมกับ:

-  ✅ **Viper Configuration** - อ่านค่าจาก environment variables
-  ✅ **GORM Integration** - ORM สำหรับ Go
-  ✅ **Repository Pattern** - ใช้ Null Object Pattern สำหรับ optional services
-  ✅ **Docker Support** - ทำงานใน Docker containers
-  ✅ **Connection Pooling** - จัดการ connection อย่างมีประสิทธิภาพ
-  ✅ **Safe Migration** - AutoMigrate เฉพาะเมื่อต่อ database ได้จริง

## 📝 Available Database Connections:

### 1. **Primary Database** (Main Application)

```go
// Default connection (backward compatible)
primaryDB, err := platform.InitPostgres()
if err != nil {
    log.Fatal("Failed to connect to primary database:", err)
}

// Auto-migrate primary tables
primaryDB.AutoMigrate(&User{})

// ใช้งาน
var users []User
primaryDB.Find(&users)
```

### 2. **Analytics Database** (Repository Pattern)

```go
// ใช้ Repository Pattern - ปลอดภัยและ clean
var analyticsRepo IAnalyticsRepository
if analyticsDB, err := platform.InitPostgresWithName("ANALYTICS"); err != nil {
    log.Printf("⚠️ Analytics database not available: %v", err)
    analyticsRepo = NewNullAnalyticsRepository() // Null Object Pattern
} else {
    log.Println("✅ Analytics database connected")
    analyticsDB.AutoMigrate(&AnalyticsEvent{})
    analyticsRepo = NewPostgresAnalyticsRepository(analyticsDB)
}

// ใช้งาน - ไม่ต้องเช็ค nil
analyticsRepo.TrackEvent(userID, "user_login", map[string]interface{}{
    "timestamp": time.Now(),
    "ip":       "192.168.1.1",
})
```

### 3. **Logs Database** (Repository Pattern)

```go
// ใช้ Repository Pattern - ปลอดภัยและ clean
var logRepo ILogRepository
if logsDB, err := platform.InitPostgresWithName("LOGS"); err != nil {
    log.Printf("⚠️ Logs database not available: %v", err)
    logRepo = NewNullLogRepository() // Null Object Pattern
} else {
    log.Println("✅ Logs database connected")
    logsDB.AutoMigrate(&ApplicationLog{})
    logRepo = NewPostgresLogRepository(logsDB)
}

// ใช้งาน - ไม่ต้องเช็ค nil
logRepo.LogError("Something went wrong", map[string]interface{}{
    "user_id": 123,
    "service": "user-service",
})
```

### 4. **Reports Database** (Repository Pattern)

```go
// ใช้ Repository Pattern - ปลอดภัยและ clean
var reportsRepo IReportsRepository
if reportsDB, err := platform.InitPostgresWithName("REPORTS"); err != nil {
    log.Printf("⚠️ Reports database not available: %v", err)
    reportsRepo = NewNullReportsRepository() // Null Object Pattern
} else {
    log.Println("✅ Reports database connected")
    reportsRepo = NewPostgresReportsRepository(reportsDB)
}

// ใช้งาน - ไม่ต้องเช็ค nil
monthlyStats, err := reportsRepo.GetMonthlyStats(time.Now().AddDate(-1, 0, 0))
if err != nil {
    log.Printf("Failed to get monthly stats: %v", err)
}
```

## 🎯 Repository Interfaces

```go
// Analytics Repository Interface
type IAnalyticsRepository interface {
    TrackEvent(userID uint, event string, details map[string]interface{}) error
    GetUserEvents(userID uint, limit int) ([]AnalyticsEvent, error)
    GetEventCount(event string) (int64, error)
}

// Log Repository Interface
type ILogRepository interface {
    LogError(message string, details map[string]interface{}) error
    LogInfo(message string, details map[string]interface{}) error
    GetRecentLogs(limit int) ([]ApplicationLog, error)
}

// Reports Repository Interface
type IReportsRepository interface {
    GetMonthlyStats(since time.Time) ([]MonthlyReport, error)
    GetUserReport(userID uint) (*UserReport, error)
    GetTopProducts(limit int) ([]ProductStats, error)
}
```

## 🏗️ Repository Implementations

### PostgreSQL Implementations

```go
// Analytics Repository - PostgreSQL
type PostgresAnalyticsRepository struct {
    db *gorm.DB
}

func NewPostgresAnalyticsRepository(db *gorm.DB) IAnalyticsRepository {
    return &PostgresAnalyticsRepository{db: db}
}

func (r *PostgresAnalyticsRepository) TrackEvent(userID uint, event string, details map[string]interface{}) error {
    analyticsEvent := AnalyticsEvent{
        UserID:    userID,
        Event:     event,
        Timestamp: time.Now(),
        Details:   details,
    }
    return r.db.Create(&analyticsEvent).Error
}

func (r *PostgresAnalyticsRepository) GetUserEvents(userID uint, limit int) ([]AnalyticsEvent, error) {
    var events []AnalyticsEvent
    err := r.db.Where("user_id = ?", userID).Order("timestamp DESC").Limit(limit).Find(&events).Error
    return events, err
}

func (r *PostgresAnalyticsRepository) GetEventCount(event string) (int64, error) {
    var count int64
    err := r.db.Model(&AnalyticsEvent{}).Where("event = ?", event).Count(&count).Error
    return count, err
}

// Log Repository - PostgreSQL
type PostgresLogRepository struct {
    db *gorm.DB
}

func NewPostgresLogRepository(db *gorm.DB) ILogRepository {
    return &PostgresLogRepository{db: db}
}

func (r *PostgresLogRepository) LogError(message string, details map[string]interface{}) error {
    log := ApplicationLog{
        Level:     "ERROR",
        Message:   message,
        Details:   details,
        Timestamp: time.Now(),
    }
    return r.db.Create(&log).Error
}

func (r *PostgresLogRepository) LogInfo(message string, details map[string]interface{}) error {
    log := ApplicationLog{
        Level:     "INFO",
        Message:   message,
        Details:   details,
        Timestamp: time.Now(),
    }
    return r.db.Create(&log).Error
}

func (r *PostgresLogRepository) GetRecentLogs(limit int) ([]ApplicationLog, error) {
    var logs []ApplicationLog
    err := r.db.Order("timestamp DESC").Limit(limit).Find(&logs).Error
    return logs, err
}

// Reports Repository - PostgreSQL
type PostgresReportsRepository struct {
    db *gorm.DB
}

func NewPostgresReportsRepository(db *gorm.DB) IReportsRepository {
    return &PostgresReportsRepository{db: db}
}

func (r *PostgresReportsRepository) GetMonthlyStats(since time.Time) ([]MonthlyReport, error) {
    var stats []MonthlyReport
    err := r.db.Raw(`
        SELECT
            DATE_TRUNC('month', created_at) as month,
            COUNT(*) as total_users,
            AVG(order_amount) as avg_order
        FROM users
        WHERE created_at >= ?
        GROUP BY month
        ORDER BY month DESC
    `, since).Scan(&stats).Error
    return stats, err
}

func (r *PostgresReportsRepository) GetUserReport(userID uint) (*UserReport, error) {
    var report UserReport
    err := r.db.Raw(`
        SELECT user_id, COUNT(*) as total_orders, SUM(amount) as total_spent
        FROM orders WHERE user_id = ?
        GROUP BY user_id
    `, userID).Scan(&report).Error
    return &report, err
}

func (r *PostgresReportsRepository) GetTopProducts(limit int) ([]ProductStats, error) {
    var products []ProductStats
    err := r.db.Raw(`
        SELECT product_id, COUNT(*) as order_count, SUM(quantity) as total_sold
        FROM order_items
        GROUP BY product_id
        ORDER BY order_count DESC
        LIMIT ?
    `, limit).Scan(&products).Error
    return products, err
}
```

### Null Object Implementations (Safe Fallback)

```go
// Analytics Repository - Null Object
type NullAnalyticsRepository struct{}

func NewNullAnalyticsRepository() IAnalyticsRepository {
    return &NullAnalyticsRepository{}
}

func (r *NullAnalyticsRepository) TrackEvent(userID uint, event string, details map[string]interface{}) error {
    // ไม่ทำอะไรเลย - แต่ไม่ error
    return nil
}

func (r *NullAnalyticsRepository) GetUserEvents(userID uint, limit int) ([]AnalyticsEvent, error) {
    return []AnalyticsEvent{}, nil
}

func (r *NullAnalyticsRepository) GetEventCount(event string) (int64, error) {
    return 0, nil
}

// Log Repository - Null Object
type NullLogRepository struct{}

func NewNullLogRepository() ILogRepository {
    return &NullLogRepository{}
}

func (r *NullLogRepository) LogError(message string, details map[string]interface{}) error {
    // ไม่ทำอะไรเลย - แต่ไม่ error
    return nil
}

func (r *NullLogRepository) LogInfo(message string, details map[string]interface{}) error {
    return nil
}

func (r *NullLogRepository) GetRecentLogs(limit int) ([]ApplicationLog, error) {
    return []ApplicationLog{}, nil
}

// Reports Repository - Null Object
type NullReportsRepository struct{}

func NewNullReportsRepository() IReportsRepository {
    return &NullReportsRepository{}
}

func (r *NullReportsRepository) GetMonthlyStats(since time.Time) ([]MonthlyReport, error) {
    return []MonthlyReport{}, nil
}

func (r *NullReportsRepository) GetUserReport(userID uint) (*UserReport, error) {
    return &UserReport{}, nil
}

func (r *NullReportsRepository) GetTopProducts(limit int) ([]ProductStats, error) {
    return []ProductStats{}, nil
}
```

## 📊 Data Structures

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"unique"`
    CreatedAt time.Time
}

type AnalyticsEvent struct {
    ID        uint                   `gorm:"primaryKey"`
    UserID    uint
    Event     string
    Timestamp time.Time
    Details   map[string]interface{} `gorm:"type:jsonb"`
}

type ApplicationLog struct {
    ID        uint                   `gorm:"primaryKey"`
    Level     string
    Message   string
    Service   string
    Timestamp time.Time
    Details   map[string]interface{} `gorm:"type:jsonb"`
}

type MonthlyReport struct {
    Month      time.Time `json:"month"`
    TotalUsers int64     `json:"total_users"`
    AvgOrder   float64   `json:"avg_order"`
}

type UserReport struct {
    UserID      uint    `json:"user_id"`
    TotalOrders int64   `json:"total_orders"`
    TotalSpent  float64 `json:"total_spent"`
}

type ProductStats struct {
    ProductID  uint  `json:"product_id"`
    OrderCount int64 `json:"order_count"`
    TotalSold  int64 `json:"total_sold"`
}
```

## 🏗️ Complete Example

```go
package main

import (
    "log"
    "time"
    "go-template/pkg/platform"
    "go-template/pkg/config"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Connect to primary database (required)
    primaryDB, err := platform.InitPostgres()
    if err != nil {
        log.Fatal("Failed to connect to primary database:", err)
    }
    log.Println("✅ Primary database connected")
    primaryDB.AutoMigrate(&User{})

    // Connect to analytics database (optional) - Repository Pattern
    var analyticsRepo IAnalyticsRepository
    if analyticsDB, err := platform.InitPostgresWithName("ANALYTICS"); err != nil {
        log.Printf("⚠️ Analytics database not available: %v", err)
        analyticsRepo = NewNullAnalyticsRepository()
    } else {
        log.Println("✅ Analytics database connected")
        analyticsDB.AutoMigrate(&AnalyticsEvent{})
        analyticsRepo = NewPostgresAnalyticsRepository(analyticsDB)
    }

    // Connect to logs database (optional) - Repository Pattern
    var logRepo ILogRepository
    if logsDB, err := platform.InitPostgresWithName("LOGS"); err != nil {
        log.Printf("⚠️ Logs database not available: %v", err)
        logRepo = NewNullLogRepository()
    } else {
        log.Println("✅ Logs database connected")
        logsDB.AutoMigrate(&ApplicationLog{})
        logRepo = NewPostgresLogRepository(logsDB)
    }

    // Example usage
    user := User{
        Email:     "test@example.com",
        CreatedAt: time.Now(),
    }

    // Save to primary database
    if err := primaryDB.Create(&user).Error; err != nil {
        log.Printf("Error creating user: %v", err)

        // Log error using Repository Pattern
        logRepo.LogError("Failed to create user", map[string]interface{}{
            "email": user.Email,
            "error": err.Error(),
        })
        return
    }

    log.Printf("✅ User created: %d", user.ID)

    // Track analytics event using Repository Pattern
    analyticsRepo.TrackEvent(user.ID, "user_registered", map[string]interface{}{
        "timestamp": time.Now(),
        "source":    "api",
    })

    // Graceful shutdown
    defer func() {
        log.Println("🔄 Closing database connections...")
        platform.CloseAllConnections()
        log.Println("✅ All connections closed")
    }()
}
```

## 🌍 Environment Variables Configuration

### Required (Primary Database):

```env
DB_HOST=localhost
DB_PORT=7430
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_template
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Bangkok
```

### Optional (Analytics Database):

```env
ANALYTICS_DB_HOST=analytics-server.company.com
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_USER=analytics_user
ANALYTICS_DB_PASSWORD=analytics_secret
ANALYTICS_DB_NAME=analytics_db
ANALYTICS_DB_SSL_MODE=require
ANALYTICS_DB_TIMEZONE=Asia/Bangkok
```

### Optional (Logs Database):

```env
LOGS_DB_HOST=logs-server.company.com
LOGS_DB_PORT=5432
LOGS_DB_USER=logs_user
LOGS_DB_PASSWORD=logs_secret
LOGS_DB_NAME=application_logs
LOGS_DB_SSL_MODE=require
LOGS_DB_TIMEZONE=Asia/Bangkok
```

### Optional (Reports Database):

```env
REPORTS_DB_HOST=replica-server.company.com
REPORTS_DB_PORT=5432
REPORTS_DB_USER=readonly_user
REPORTS_DB_PASSWORD=readonly_secret
REPORTS_DB_NAME=go_template
REPORTS_DB_SSL_MODE=require
REPORTS_DB_TIMEZONE=Asia/Bangkok
```

## 🚨 Best Practices

### ✅ DO - Use Repository Pattern

```go
// ✅ Safe: Repository Pattern with Null Object
var logRepo ILogRepository
if logsDB, err := platform.InitPostgresWithName("LOGS"); err != nil {
    logRepo = NewNullLogRepository() // Safe fallback
} else {
    logsDB.AutoMigrate(&ApplicationLog{}) // Only migrate when connected
    logRepo = NewPostgresLogRepository(logsDB)
}

// ใช้งาน - ไม่ต้องเช็ค nil
logRepo.LogError("Something went wrong", details)
```

### ❌ DON'T - Direct Database Fallback

```go
// ❌ Dangerous: Direct database fallback
logsDB, err := platform.InitPostgresWithName("LOGS")
if err != nil {
    logsDB = primaryDB // อันตราย!
}

// อาจไปสร้าง table ใน primary DB โดยไม่ตั้งใจ
logsDB.AutoMigrate(&ApplicationLog{}) // 💥 Dangerous!
```

## 🔗 Redis Integration

```go
// Default Redis connection
rdb, err := platform.InitRedis()
if err != nil {
    log.Printf("Redis not available: %v", err)
}

// Named Redis connections
cacheRedis, err := platform.InitRedisWithName("CACHE")
sessionRedis, err := platform.InitRedisWithName("SESSION")
```

### Redis Environment Variables:

```env
# Default Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Cache Redis
CACHE_REDIS_HOST=cache-server.company.com
CACHE_REDIS_PORT=6379
CACHE_REDIS_PASSWORD=cache_secret
CACHE_REDIS_DB=1

# Session Redis
SESSION_REDIS_HOST=session-server.company.com
SESSION_REDIS_PORT=6379
SESSION_REDIS_PASSWORD=session_secret
SESSION_REDIS_DB=2
```

## 🐳 Docker Support

### Development Mode:

```bash
# รัน PostgreSQL + Redis + Application
make docker-dev

# ดู logs
make docker-dev-logs

# หยุด
make docker-dev-stop
```

### Docker Environment Variables:

```env
# ใน docker-compose.dev.yml
DB_HOST=host.docker.internal  # เชื่อมต่อ PostgreSQL บน host machine
REDIS_HOST=host.docker.internal  # เชื่อมต่อ Redis บน host machine
```

## ⚙️ Viper Configuration Details

### Configuration Keys:

```go
// Primary Database (lowercase)
viper.GetString("database.host")
viper.GetString("database.port")
viper.GetString("database.user")

// Named Databases (uppercase)
viper.GetString("ANALYTICS.host")
viper.GetString("LOGS.port")
viper.GetString("REPORTS.user")
```

### Auto-binding Environment Variables:

```go
// Primary database
viper.BindEnv("database.host", "DB_HOST")
viper.BindEnv("database.port", "DB_PORT")

// Analytics database
viper.BindEnv("ANALYTICS.host", "ANALYTICS_DB_HOST")
viper.BindEnv("ANALYTICS.port", "ANALYTICS_DB_PORT")
```

## 📊 Use Cases

### 🎯 **Primary Database**: Core application data

-  Users, products, orders
-  Critical business logic
-  Transactional consistency required

### 📈 **Analytics Database**: Metrics and tracking

-  User behavior tracking
-  Performance metrics
-  Business intelligence data
-  Can be on separate high-performance server

### 📝 **Logs Database**: Application logging

-  Error logs, access logs
-  Audit trails
-  Debugging information
-  High write volume, separate from main DB

### 📊 **Reports Database**: Read-only queries

-  Complex reporting queries
-  Data warehousing
-  Read replica to reduce load on primary
-  Historical data analysis

## 🛠️ Troubleshooting

### Common Issues:

1. **Connection Refused**:

   ```bash
   # Check if PostgreSQL is running
   lsof -i :7430

   # Check Docker containers
   docker ps
   ```

2. **Vendoring Issues**:

   ```bash
   # Clean and rebuild
   rm -rf vendor/
   go clean -modcache
   go mod download
   ```

3. **VS Code Go Extension Issues**:
   -  Restart VS Code completely
   -  Command Palette: "Go: Restart Language Server"

### Docker Development:

```bash
# Rebuild containers
make docker-dev-stop
docker compose -f docker-compose.dev.yml build --no-cache
make docker-dev
```

## 📚 Further Reading

-  [Viper Configuration Library](https://github.com/spf13/viper)
-  [GORM PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
-  [Redis Go Client](https://github.com/redis/go-redis)
-  [Docker Compose for Development](https://docs.docker.com/compose/)

---

_Last Updated: August 2025_
_For support, check the project's README.md or create an issue._
