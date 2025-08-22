# 🔍 Viper Configuration Mapping คืออะไร?

## 📋 Environment Variables → Viper Keys

| .env File | Viper Key | อธิบาย |
|-----------|-----------|---------|
| `DB_HOST=localhost` | `database.primary.host` | Primary Database hostname |
| `DB_PORT=7430` | `database.primary.port` | Primary Database port |
| `ANALYTICS_DB_HOST=server.com` | `database.analytics.host` | Analytics Database hostname |
| `LOGS_DB_HOST=logs.com` | `database.logs.host` | Logs Database hostname |
| `PORT=9998` | `server.port` | Server port |
| `JWT_SECRET=abc123` | `auth.jwt_secret` | JWT secret key |

## 🔧 การทำงานของ Viper

### 1. Binding ใน `pkg/config/config.go`:
```go
// Primary Database
viper.BindEnv("database.primary.host", "DB_HOST")
viper.BindEnv("database.primary.port", "DB_PORT")

// Analytics Database  
viper.BindEnv("database.analytics.host", "ANALYTICS_DB_HOST")
viper.BindEnv("database.analytics.port", "ANALYTICS_DB_PORT")
```

### 2. ใน `.env` file:
```env
# Primary Database
DB_HOST=localhost
DB_PORT=7430

# Analytics Database
ANALYTICS_DB_HOST=analytics-server.com
ANALYTICS_DB_PORT=5432
```

### 3. ใน code:
```go
// ✅ ถูกต้อง - ใช้ viper key ที่มี binding
host := viper.GetString("database.primary.host")        // ได้ "localhost"
analyticsHost := viper.GetString("database.analytics.host") // ได้ "analytics-server.com"

// ❌ ผิด - ใช้ key ที่ไม่มี binding
host := viper.GetString("database.host")                // ได้ ""
```

## 🎯 Multi-Database Example:

### Environment Variables:
```env
# Primary Database (Local)
DB_HOST=localhost
DB_PORT=7430
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_template

# Analytics Database (Remote)
ANALYTICS_DB_HOST=analytics.company.com
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_USER=analytics_user
ANALYTICS_DB_PASSWORD=analytics_secret
ANALYTICS_DB_NAME=analytics_db

# Logs Database (Remote)
LOGS_DB_HOST=logs.company.com
LOGS_DB_PORT=5432
LOGS_DB_USER=logs_user
LOGS_DB_PASSWORD=logs_secret
LOGS_DB_NAME=application_logs
```

### Code Usage:
```go
// Primary database (backward compatible)
primaryDB, err := platform.InitPostgres()
// หรือ
primaryDB, err := platform.InitPostgresWithName("primary")

// Analytics database (คนละ server)
analyticsDB, err := platform.InitPostgresWithName("analytics")

// Logs database (คนละ server)
logsDB, err := platform.InitPostgresWithName("logs")
```

## 🚀 ข้อดี:
1. **ไม่มี hardcoded defaults** - รู้ทันทีถ้าลืมตั้งค่า
2. **Multiple databases** - ต่อ DB หลายตัวได้ คนละ server
3. **Type safety** - `viper.GetInt()`, `viper.GetBool()`
4. **Clear error messages** - บอกชัดเจนว่าขาด env var ไหน
5. **Structured config** - จัดกลุ่มได้ เช่น `database.primary.*`, `database.analytics.*`

## 📝 สรุป:
- **Viper ≠ แทนที่ .env**
- **Viper = เครื่องมืออ่าน .env แบบ smart**
- **แต่ละ ENV VAR ต้อง bind ก่อน**
- **ใช้ viper key แทน env var name**
- **รองรับหลาย databases แยกกันได้**
