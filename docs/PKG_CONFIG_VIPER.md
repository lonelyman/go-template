คู่มือการใช้งาน pkg/config
แพ็กเกจ config คือ "ศูนย์บัญชาการ" ของแอปพลิเคชันเรา มีหน้าที่รับผิดชอบในการโหลด, จัดการ, และเข้าถึงค่าตั้งค่า (Configuration) ทั้งหมดของโปรเจกต์ โดยใช้ไลบรารี Viper ที่ทรงพลังเป็นแกนหลัก

ปรัชญาหลัก: ระบบ Config ของเราทำงานแบบ "3 ชั้น" โดยมีลำดับความสำคัญดังนี้:

Environment Variables (จาก .env หรือระบบ) > configs/config.yml > ค่า Default ในโค้ด

ซึ่งหมายความว่าค่าที่ตั้งใน Environment Variable จะ "ชนะ" และ "เขียนทับ" ค่าที่อยู่ในไฟล์ config.yml เสมอ

1. โครงสร้างของ Configuration
"พิมพ์เขียว" ของค่า Config ทั้งหมดถูกกำหนดไว้ใน struct ภายในไฟล์ config.go เพื่อให้เราสามารถเข้าถึงค่าต่างๆ ได้อย่างปลอดภัย (Type-Safe)

พิมพ์เขียวหลัก (Config struct)
// pkg/config/config.go

type Config struct {
   App      AppConfig    `mapstructure:"app"`
   Server   ServerConfig `mapstructure:"server"`
   Postgres PostgresDbs  `mapstructure:"postgres"`
   Auth     AuthConfig   `mapstructure:"auth"`
}

พิมพ์เขียวย่อย
type AppConfig struct {
   Name    string `mapstructure:"name"`
   Version string `mapstructure:"version"`
}

type ServerConfig struct {
   AppPort  string `mapstructure:"appport"`
   HostPort string `mapstructure:"hostport"`
}

type PostgresDbs struct {
   Primary PostgresConfig `mapstructure:"primary"`
   Logs    PostgresConfig `mapstructure:"logs"`
}

type AuthConfig struct {
   JWTSecret string `mapstructure:"jwtSecret"`
}

type PostgresConfig struct {
   Host     string `mapstructure:"host"`
   Port     string `mapstructure:"port"`
   User     string `mapstructure:"user"`
   Password string `mapstructure:"password"`
   DBName   string `mapstructure:"name"`
   SSLMode  string `mapstructure:"ssl_mode"`
}

2. แหล่งที่มาของ Configuration
2.1 configs/config.yml (เมนูมาตรฐาน)
นี่คือไฟล์ที่เก็บ ค่าเริ่มต้น (Defaults) ที่ไม่เป็นความลับ และควรจะถูก Commit เข้าไปใน Git Repository เพื่อให้ทุกคนในทีมมีค่าตั้งต้นเดียวกัน

ตัวอย่าง:

# configs/config.yml
app:
  name: "Go Template API"
  version: "v1.0.0"

server:
  appport: "9998"
  hostport: "9999"

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

2.2 .env (โพยลับ)
ไฟล์นี้คือที่สำหรับเก็บ ข้อมูลที่เป็นความลับ (Secrets) และใช้สำหรับ "เขียนทับ" (Override) ค่าเริ่มต้นสำหรับสภาพแวดล้อมนั้นๆ (เช่น Local Development)

สำคัญ: ไฟล์ .env จะต้องถูกใส่ไว้ใน .gitignore เสมอ!

Viper จะทำการ map Environment Variable ไปยัง struct ของเราโดยอัตโนมัติ โดยใช้กฎ UPPERCASE_WITH_UNDERSCORES

ตัวอย่าง:

POSTGRES_PRIMARY_HOST -> จะถูก map ไปที่ cfg.Postgres.Primary.Host

AUTH_JWTSECRET -> จะถูก map ไปที่ cfg.Auth.JWTSecret

# .env
POSTGRES_PRIMARY_HOST=localhost
POSTGRES_PRIMARY_PORT=7430
POSTGRES_PRIMARY_USER=root
POSTGRES_PRIMARY_PASSWORD=supersecretpassword
POSTGRES_PRIMARY_NAME=go_template

AUTH_JWTSECRET=another-super-secret-key

3. การใช้งาน
เราจะเรียกใช้ฟังก์ชัน config.LoadConfig() แค่ครั้งเดียวเท่านั้น ที่จุดเริ่มต้นของโปรแกรมใน cmd/api/main.go

// cmd/api/main.go

func main() {
    // โหลด .env ก่อน (ถ้ามี)
    if os.Getenv("DOCKER_ENV") != "true" {
        godotenv.Load()
    }

    // โหลด Config ทั้งหมด
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("❌ Failed to load configuration: %v", err)
    }

    // ตอนนี้เราสามารถเข้าถึงค่า Config ทั้งหมดผ่าน `cfg` object ได้แล้ว
    // เช่น cfg.Server.AppPort, cfg.Postgres.Primary, cfg.Auth.JWTSecret
    
    // ส่งต่อไปให้ส่วนอื่นๆ ของโปรแกรมผ่าน Dependency Injection
    db, err := postgres.NewConnection(cfg.Postgres.Primary)
    // ...
    service := example_user.NewService(repo, cfg.Auth.JWTSecret)
    // ...
}

Helper Method: BuildDSN()
PostgresConfig struct มี Helper Method ติดตัวมาด้วยเพื่อความสะดวกในการสร้าง Connection String

// ใช้งานใน cmd/migrate/main.go
dsn := cfg.Postgres.Primary.BuildDSN()
