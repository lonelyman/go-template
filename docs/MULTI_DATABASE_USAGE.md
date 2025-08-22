# 🗄️ Multi-Database PostgreSQL Usage Guide

_Updated: August 2025 - Simple & Practical Approach_

## 🚀 Overview

ระบบรองรับการเชื่อมต่อ PostgreSQL หลายฐานข้อมูลแบบง่ายๆ:

-  ✅ **Primary Database** - ฐานข้อมูลหลัก (Required)
-  ✅ **Optional Databases** - ฐานข้อมูลเสริม (Analytics, Logs, Reports)
-  ✅ **Graceful Degradation** - แอปทำงานได้แม้ฐานข้อมูลเสริมไม่มี
-  ✅ **Environment Variables** - ตั้งค่าผ่าน .env หรือ environment
-  ✅ **Docker Support** - ใช้งานใน Docker ได้

## 📝 Basic Usage:

### 1. **Primary Database** (Required)

```go
// เชื่อมต่อฐานข้อมูลหลัก - ต้องมี
primaryDB, err := postgres.NewConnection(cfg.Database)
if err != nil {
    log.Fatalf("Failed to connect to primary database: %v", err)
}
log.Println("✅ Primary database connected")

// ใช้งาน
var users []User
primaryDB.Find(&users)
```

### 2. **Optional Databases** (Simple Pattern)

```go
// Analytics Database (Optional)
analyticsDB, err := postgres.InitPostgresWithName("ANALYTICS")
if err != nil {
    log.Printf("⚠️ Analytics database not available: %v", err)
    analyticsDB = nil // ใช้ nil
} else {
    log.Println("✅ Analytics database connected")
}

// Logs Database (Optional)
logsDB, err := postgres.InitPostgresWithName("LOGS")
if err != nil {
    log.Printf("⚠️ Logs database not available: %v", err)
    logsDB = nil // ใช้ nil
} else {
    log.Println("✅ Logs database connected")
}

// ใช้งาน - เช็ค nil ก่อน
if analyticsDB != nil {
    analyticsDB.Create(&AnalyticsEvent{UserID: 123, Event: "login"})
}

if logsDB != nil {
    logsDB.Create(&ApplicationLog{Level: "INFO", Message: "User logged in"})
}
```

## 📊 Simple Data Structures

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
    Timestamp time.Time
    Details   map[string]interface{} `gorm:"type:jsonb"`
}
```

## 🏗️ Complete Example (main.go)

```go
package main

import (
    "log"
    "go-template/pkg/config"
    postgres "go-template/pkg/platform/postgres"
)

func main() {
    // 1. Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load configuration: %v", err)
    }

    // 2. Connect to primary database (required)
    primaryDB, err := postgres.NewConnection(cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to primary database: %v", err)
    }
    log.Println("✅ Primary database connected")

    // 3. Connect to optional databases
    analyticsDB, err := postgres.InitPostgresWithName("ANALYTICS")
    if err != nil {
        log.Printf("⚠️ Analytics database not available: %v", err)
        analyticsDB = nil
    } else {
        log.Println("✅ Analytics database connected")
    }

    logsDB, err := postgres.InitPostgresWithName("LOGS")
    if err != nil {
        log.Printf("⚠️ Logs database not available: %v", err)
        logsDB = nil
    } else {
        log.Println("✅ Logs database connected")
    }

    // 4. Use databases
    // Primary database - always available
    var users []User
    primaryDB.Find(&users)

    // Optional databases - check nil first
    if analyticsDB != nil {
        event := AnalyticsEvent{
            UserID: 123,
            Event:  "user_login",
        }
        analyticsDB.Create(&event)
    }

    if logsDB != nil {
        log := ApplicationLog{
            Level:   "INFO",
            Message: "Application started",
        }
        logsDB.Create(&log)
    }
}
```

## 🌍 Environment Variables (Simple)

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

### Optional (Additional Databases):

```env
# Analytics Database (Optional)
ANALYTICS_DB_HOST=analytics-server.com
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_USER=analytics_user
ANALYTICS_DB_PASSWORD=analytics_secret
ANALYTICS_DB_DBNAME=analytics_db
ANALYTICS_DB_SSLMODE=require
ANALYTICS_DB_TIMEZONE=Asia/Bangkok

# Logs Database (Optional)
LOGS_DB_HOST=logs-server.com
LOGS_DB_PORT=5432
LOGS_DB_USER=logs_user
LOGS_DB_PASSWORD=logs_secret
LOGS_DB_DBNAME=application_logs
LOGS_DB_SSLMODE=require
LOGS_DB_TIMEZONE=Asia/Bangkok
```

## 💡 Best Practices (Simple)

### ✅ DO - Simple Nil Check

```go
// ✅ Simple and clear
analyticsDB, err := postgres.InitPostgresWithName("ANALYTICS")
if err != nil {
    log.Printf("⚠️ Analytics database not available: %v", err)
    analyticsDB = nil
}

// Use with nil check
if analyticsDB != nil {
    analyticsDB.Create(&event)
}
```

### ❌ DON'T - Over-complicated Pattern

```go
// ❌ Too complex for simple usage
var analyticsRepo IAnalyticsRepository
if analyticsDB, err := postgres.InitPostgresWithName("ANALYTICS"); err != nil {
    analyticsRepo = NewNullAnalyticsRepository() // Too much abstraction
} else {
    analyticsRepo = NewPostgresAnalyticsRepository(analyticsDB)
}
```

## 📊 Use Cases

### 🎯 **Primary Database**: Core application data

-  Users, products, orders
-  **Always required**
-  App crashes if not available

### 📈 **Analytics Database**: Metrics and tracking

-  User behavior tracking
-  Performance metrics
-  **Optional** - app continues without it

### 📝 **Logs Database**: Application logging

-  Error logs, access logs
-  Debugging information
-  **Optional** - app continues without it

## 🛠️ Troubleshooting

### Connection Issues:

```bash
# Check if PostgreSQL is running
lsof -i :7430

# Test connection manually
psql -h localhost -p 7430 -U postgres -d go_template
```

### Docker Development:

```bash
# Run with Docker
make docker-dev

# Stop
make docker-dev-stop
```

---

_Last Updated: August 2025 - Simple & Practical_  
_For support, check the project's README.md or create an issue._
