# ความแตกต่างระหว่าง Viper กับ os.Getenv()

## แบบเดิม (os.Getenv + defaults):
```go
host := getEnv("DB_HOST", "localhost")        // ใช้ default "localhost"
port := getEnv("DB_PORT", "5432")             // ใช้ default "5432"
user := getEnv("DB_USER", "postgres")         // ใช้ default "postgres"
password := getEnv("DB_PASSWORD", "password") // ใช้ default "password"

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback // แฮดโค้ด default
}
```

## แบบใหม่ (Viper + validation):
```go
// ใน .env:
DB_HOST=localhost
DB_PORT=7430
DB_USER=postgres
DB_PASSWORD=your_real_password

// ใน config binding:
viper.BindEnv("database.primary.host", "DB_HOST")
viper.BindEnv("database.primary.port", "DB_PORT")

// ใน code:
host := viper.GetString("database.primary.host")
// ถ้าไม่มี DB_HOST จะได้ "" และ error ทันที (ไม่ใช้ default)
```

# ข้อดีของ Viper:

## 1. ไม่มี hardcoded defaults
```go
// เก่า: ใช้ default แอบๆ
host := getEnv("DB_HOST", "localhost") // ถ้าไม่มี DB_HOST จะใช้ "localhost"

// ใหม่: error ทันทีถ้าไม่มี config
host := viper.GetString("database.primary.host") // ถ้าไม่มี DB_HOST จะได้ ""
if host == "" {
    return fmt.Errorf("DB_HOST is required")
}
```

## 2. Multiple databases
```go
// เก่า: ต่อได้แค่ DB เดียว
db, err := connectToDatabase()

// ใหม่: ต่อได้หลาย DB แยกกัน
primaryDB, err := platform.InitPostgresWithName("primary")
analyticsDB, err := platform.InitPostgresWithName("analytics")
logsDB, err := platform.InitPostgresWithName("logs")
```

## 3. Type safety
```go
// เก่า: ทุกอย่างเป็น string
maxConns := getEnv("DB_MAX_CONNS", "25")
maxConnsInt, _ := strconv.Atoi(maxConns) // ต้องแปลงเอง

// ใหม่: type safety
maxConns := viper.GetInt("database.primary.max_conns") // ได้ int ทันที
timeout := viper.GetDuration("database.primary.timeout") // ได้ time.Duration
enabled := viper.GetBool("database.primary.enabled") // ได้ bool
```

## 4. Environment mapping
```go
// เก่า: ใช้ env var name ตรงๆ
host := os.Getenv("DB_HOST")
analyticsHost := os.Getenv("ANALYTICS_DB_HOST")

// ใหม่: map เป็น structure ที่เข้าใจง่าย
host := viper.GetString("database.primary.host")
analyticsHost := viper.GetString("database.analytics.host")
```

## 5. Validation
```go
// เก่า: ไม่มี validation
host := getEnv("DB_HOST", "localhost") // อาจได้ค่าผิด แต่ไม่รู้

// ใหม่: มี validation
if !viper.IsSet("database.primary.host") {
    return fmt.Errorf("database configuration not found - check DB_HOST env var")
}
```

# การใช้งานปัจจุบัน:

## ใน .env (เหมือนเดิม):
```env
# Primary Database
DB_HOST=localhost
DB_PORT=7430
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_template

# Analytics Database
ANALYTICS_DB_HOST=analytics-server.com
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_USER=analytics_user
ANALYTICS_DB_PASSWORD=analytics_password
ANALYTICS_DB_NAME=analytics_db
```

## ใน code (เปลี่ยน):
```go
// แบบเดิม (ยังใช้ได้):
db, err := platform.InitPostgres()

// แบบใหม่ (ใช้ viper ข้างใน):
primaryDB, err := platform.InitPostgres() // เหมือนเดิม แต่ใช้ viper
analyticsDB, err := platform.InitPostgresWithName("analytics") // ใหม่!
```

# สรุป:
- **.env ยังใช้ได้เหมือนเดิม**
- **Viper แทนที่ os.Getenv() ข้างใน**
- **ไม่มี default values แฮดโค้ด**
- **รองรับ multiple databases คนละ server**
- **Type safety และ validation**
- **Error messages ชัดเจน**
