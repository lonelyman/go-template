คู่มือการใช้งาน PostgreSQL หลาย Database
"พิมพ์เขียว" ของเราถูกออกแบบมาให้รองรับการเชื่อมต่อกับ PostgreSQL Database หลายตัวได้อย่างยืดหยุ่นและทนทานต่อความผิดพลาด (Resilient) เอกสารนี้จะอธิบายหลักการและขั้นตอนการใช้งานทั้งหมด

เป้าหมายหลัก: เพื่อแยกประเภทของข้อมูลออกจากกัน ซึ่งช่วยเพิ่มประสิทธิภาพและความปลอดภัย เช่น:

Primary DB: เก็บข้อมูลหลักของธุรกิจ (Users, Orders) - ต้องพร้อมใช้งานเสมอ

Logs DB: เก็บข้อมูล Log การใช้งาน (Activity Logs) - ถ้าล่มไปชั่วคราว แอปหลักต้องไม่พังตาม

Analytics DB: เก็บข้อมูลสำหรับการวิเคราะห์ - เช่นเดียวกับ Logs DB

1. การตั้งค่า (Configuration)
   ระบบ Config ของเราใช้หลักการ 3 ชั้น (Defaults -> File -> Environment) เพื่อความยืดหยุ่นสูงสุด

1.1 configs/config.yml (เมนูมาตรฐาน)
นี่คือไฟล์ที่กำหนด "พิมพ์เขียว" และค่าเริ่มต้นสำหรับตอนพัฒนาที่เครื่องเรา

# configs/config.yml

postgres:
primary:
host: "localhost"
port: "7430"
user: "root"
password: "" # ไม่เก็บความลับที่นี่
name: "go_template"
ssl_mode: "disable"
logs:
host: "" # เว้นว่างไว้สำหรับ optional db
port: ""
user: ""
password: ""
name: ""
ssl_mode: "disable"

1.2 .env (โพยลับ)
ไฟล์นี้จะใช้ "ทับ" ค่าจาก config.yml และเก็บข้อมูลที่เป็นความลับ เราจะใช้ชื่อตัวแปรที่สอดคล้องกับโครงสร้างใน struct ของเรา

# .env (สำหรับ Local Development)

# --- Primary Database ---

POSTGRES_PRIMARY_HOST=localhost
POSTGRES_PRIMARY_PORT=7430
POSTGRES_PRIMARY_USER=root
POSTGRES_PRIMARY_PASSWORD=12345678
POSTGRES_PRIMARY_NAME=go_template

# --- Logs Database ---

POSTGRES_LOGS_HOST=localhost
POSTGRES_LOGS_PORT=7431 # (อาจจะเป็นคนละพอร์ตหรือคนละ instance)
POSTGRES_LOGS_USER=root
POSTGRES_LOGS_PASSWORD=12345678
POSTGRES_LOGS_NAME=go_template_logs

1.3 pkg/config/config.go (พิมพ์เขียวในโค้ด)
struct ใน Go ของเราจะต้องมีโครงสร้างที่ตรงกับไฟล์ .yml เพื่อให้ Viper สามารถ map ค่าได้อย่างถูกต้อง

// pkg/config/config.go

type Config struct {
// ...
Postgres PostgresDbs `mapstructure:"postgres"`
}

type PostgresDbs struct {
Primary PostgresConfig `mapstructure:"primary"`
Logs PostgresConfig `mapstructure:"logs"`
Analytics PostgresConfig `mapstructure:"analytics"`
}

type PostgresConfig struct {
Host string `mapstructure:"host"`
Port string `mapstructure:"port"`
User string `mapstructure:"user"`
Password string `mapstructure:"password"`
DBName string `mapstructure:"name"`
SSLMode string `mapstructure:"ssl_mode"`
}

2. การเชื่อมต่อในแอปพลิเคชัน (main.go)
   ใน main.go เราจะจัดการกับการเชื่อมต่อ DB แต่ละประเภทแตกต่างกันไปตามความสำคัญ

// cmd/api/main.go

// Primary Database (จำเป็นต้องมี)
// ใช้หลักการ "Fail Fast": ถ้าต่อไม่ได้ ให้โปรแกรมพังไปเลย
primaryDB, err := postgres.NewConnection(cfg.Postgres.Primary)
if err != nil {
log.Fatalf("❌ Failed to connect to primary database: %v", err)
}

// Logs Database (ไม่มีก็ได้)
// ใช้หลักการ "Fallback": ถ้าต่อไม่ได้ ก็แค่ Log เตือน แล้วทำงานต่อไป
var logsDB \*gorm.DB
if cfg.Postgres.Logs.Host != "" { // เช็คว่ามีการตั้งค่าหรือไม่
logsDB, err = postgres.NewConnection(cfg.Postgres.Logs)
if err != nil {
log.Printf("⚠️ Logs database configured but unavailable: %v", err)
logsDB = nil // ตั้งค่าเป็น nil แล้วทำงานต่อ
}
}

3. การจัดการโครงสร้าง (Database Migrations)
   เราจะแยก "พิมพ์เขียว" (.sql ไฟล์) และ "คำสั่ง" (Makefile) สำหรับแต่ละ Database ออกจากกันอย่างชัดเจน

3.1 โครงสร้างโฟลเดอร์
db/
└── migrations/
├── primary/
│ └── 000001_create_example_users_table.up.sql
└── logs/
└── 000001_create_activity_logs_table.up.sql

3.2 Makefile
"แผงควบคุม" ของเราจะมีคำสั่งแยกสำหรับแต่ละ Database ทำให้เราเลือกทำงานได้อย่างเจาะจง

# Makefile

# --- Shortcuts for convenience ---

db-migrate-primary:
@make db-migrate db=primary

db-migrate-logs:
@make db-migrate db=logs

# --- คำสั่งหลัก (ยืดหยุ่น) ---

db-migrate:
ifndef db
$(error db is not set. Usage: make db-migrate db=<primary|logs>)
endif
	@docker compose -f docker-compose.dev.yml run --rm migrate \
	  go run ./cmd/migrate/main.go --db=$(db) --path=db/migrations/$(db)

3.3 cmd/migrate/main.go
โปรแกรม migrate ของเราฉลาดพอที่จะรับ Flag --db และ --path เพื่อไปดึง Connection String ที่ถูกต้องจาก Config และทำงานกับไฟล์ Migration ที่ถูกต้อง

// cmd/migrate/main.go

// ...
var dsn string
switch dbName {
case "primary":
dsn = cfg.Postgres.Primary.BuildDSN()
case "logs":
dsn = cfg.Postgres.Logs.BuildDSN()
default:
log.Fatalf("❌ Unknown database name: '%s'", dbName)
}

m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dsn)
// ...
