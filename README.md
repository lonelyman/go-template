# Go Template Project

โปรเจกต์ Go ที่ใช้ Hexagonal Architecture พร้อม### เตรียม Database (PostgreSQL)

**ตัวเลือก 1: ใช้ PostgreSQL ที่มีอยู่แล้วในเครื่อง**

```bash
# ตรวจสอบว่า PostgreSQL รันอยู่หรือไม่
pg_isready -h localhost -p 5432

# หากยังไม่รัน ให้เริ่ม PostgreSQL service
# สำหรับ macOS (Homebrew)
brew services start postgresql

# สำหรับ Linux (systemd)
sudo systemctl start postgresql

# สร้าง database
createdb go_template

# หรือใช้ psql
psql -h localhost -U postgres -c "CREATE DATABASE go_template;"
```

**ตัวเลือก 2: ใช้ Docker PostgreSQL**

```bash
# รัน PostgreSQL แยกต่างหาก
docker run --name postgres-go-template \
  -e POSTGRES_DB=go_template \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15-alpine

# หรือ uncomment postgres service ใน docker-compose.yml แล้วรัน
docker-compose up -d postgres
```

2. Copy environment variables:ters และระบบ testing ที่สมบูรณ์

# Go Template Project

โปรเจกต์ Go ที่ใช้ Hexagonal Architecture พร้อมด้วย modules, adapters และระบบ testing ที่สมบูรณ์

## โครงสร้างโปรเจกต์

```
.
├── tests/                     # (Lab ทดสอบ) บ้านของการทดสอบแบบ Integration Test
│   ├── example_module_api_test.go # - เทสต์ API ของ Example Module (ยิง HTTP จริง -> เช็ค DB จริง)
│   └── main_test.go           # - ไฟล์ตั้งค่ากลางสำหรับ Test (เช่น ต่อ DB Test)
│
├── assets/                      # (คลังทรัพย์สิน) บ้านของไฟล์ Static (Fonts, Images, Templates)
│   ├── fonts/
│   ├── images/
│   └── templates/
│
├── build/                       # (โรงงานประกอบร่าง) บ้านของ Dockerfile และไฟล์เกี่ยวกับการ Build
│   └── package/
│       └── Dockerfile           # - Dockerfile แบบ Multi-stage ที่ดีที่สุด
│
├── cmd/                         # (ปุ่มสตาร์ท) ที่อยู่ของโปรแกรมที่สั่งรันได้
│   └── api/
│       └── main.go              # ❤️ จุดเริ่มต้นเดียวของแอปเรา, ที่ประกอบร่างทุกอย่างเข้าด้วยกัน
│
├── configs/                     # (ห้องสำรอง) สำหรับไฟล์ config เพิ่มเติมในอนาคต (ปัจจุบันใช้ .env)
│
├── internal/                    # === ⭐️ หัวใจและสมองของโปรเจกต์ (โค้ดหลักทั้งหมด) ⭐️ ===
│   ├── modules/                 # (แผนกต่างๆ) บ้านของ Business Logic แต่ละฟีเจอร์
│   │   └── example-module/        #   - ตัวอย่าง 1 ฟีเจอร์ที่สมบูรณ์
│   │       ├── example_handler.go      #     - Handler: ประตูหน้าด่านของฟีเจอร์, คุยด้วยภาษา HTTP/JSON
│   │       ├── example_service.go      #     - Service: สมองของฟีเจอร์, ที่อยู่ของ Business Logic
│   │       ├── example_repository.go   #     - Repository: แขนขาของฟีเจอร์, คุยกับ Database
│   │       ├── example_domain.go       #     - Domain: พิมพ์เขียวข้อมูลของฟีเจอร์, บริสุทธิ์ที่สุด
│   │       ├── example_service_test.go #     - Unit Test: หน่วยตรวจสอบคุณภาพของ Service
│   │       └── module.go               #     - Module: ประกอบร่างทุกอย่างเข้าด้วยกัน
│   │
│   └── adapters/                # (ประตูเชื่อมต่อกลาง) ที่อยู่ของ Adapters ที่ "ใช้ร่วมกัน"
│       ├── primary/               #   - ประตูทางเข้า (สำหรับ Request ที่เข้ามา)
│       │   └── http/                #     - ประตูสำหรับภาษา HTTP
│       │       └── middleware/      #       - ยามเฝ้าประตู (Logger, CORS, Auth Middleware)
│       │
│       └── secondary/             #   - ประตูทางออก (สำหรับเรียกไปข้างนอก)
│           └── dhl/               #     - Adapter สำหรับคุยกับ DHL API
│
├── pkg/                         # === 🧰 กล่องเครื่องมือช่าง (Reusable Code) 🧰 ===
│   ├── platform/                # เครื่องมือเชื่อมต่อ Platform (Postgres, Redis)
│   ├── auth/                    # เครื่องมือจัดการ Auth (JWT, Password Hashing)
│   └── utils/                   # เครื่องมือจิปาถะ (String, Time)
│
├── api/                         # (ห้องสมุด) สำหรับไฟล์ Document/Spec ของ API เช่น OpenAPI/Swagger
├── .env.example                 # ตัวอย่างไฟล์ Environment Variables สำหรับนักพัฒนาคนอื่น
├── docker-compose.yml           # คู่มือรันโปรเจกต์และ Services อื่นๆ (เช่น DB) ที่เครื่องเรา
├── go.mod                       # รายการ Library ที่โปรเจกต์เราใช้
├── Makefile                     # รวมคำสั่งสั้นๆ ที่ใช้บ่อย (run, test, build)
└── README.md                    # ป้ายหน้าบ้าน, คำอธิบายโปรเจกต์
```

## เริ่มต้นใช้งาน

### ติดตั้ง Dependencies

```bash
go mod tidy
```

### เตรียม Database (PostgreSQL)

1. ติดตั้ง PostgreSQL หรือใช้ Docker:

```bash
# ใช้ Docker Compose (แนะนำ)
docker-compose up -d postgres

# หรือรัน PostgreSQL แยก
docker run --name postgres-go-template
  -e POSTGRES_DB=go_template
  -e POSTGRES_USER=postgres
  -e POSTGRES_PASSWORD=password
  -p 5432:5432
  -d postgres:15-alpine
```

2. Copy environment variables:

```bash
cp .env.example .env
```

### รันโปรเจกต์

```bash
# รันด้วย Makefile
make run

# หรือรันตรงๆ
go run cmd/api/main.go

# รันด้วย Docker Compose (ใช้ PostgreSQL ที่มีอยู่แล้ว - postgres service ถูกคอมเมนต์ไว้)
docker-compose up

# หรือรันเฉพาะ Redis และ App
docker-compose up redis app
```

### ทดสอบ API

```bash
# Health check
curl http://localhost:8080/health

# สร้าง Example
curl -X POST http://localhost:8080/api/v1/examples
  -H "Content-Type: application/json"
  -d '{"name": "John Doe", "email": "john@example.com"}'

# ดู Examples
curl http://localhost:8080/api/v1/examples
```

## การพัฒนา

### Architecture

-  **Hexagonal Architecture**: แยก Business Logic ออกจาก Infrastructure
-  **Domain-Driven Design**: แต่ละ module มี domain ของตัวเอง
-  **Dependency Injection**: ใช้ interface สำหรับ loose coupling

### Testing

```bash
# Unit Tests
make test

# Integration Tests
make test-integration

# Test Coverage
make test-coverage
```

### เพิ่ม Module ใหม่

1. สร้างโฟลเดอร์ใน `internal/modules/your-module/`
2. สร้างไฟล์: `domain.go`, `repository.go`, `service.go`, `handler.go`, `module.go`
3. Register routes ใน `cmd/api/main.go`

### Docker

```bash
# Build image
make docker-build

# Run with Docker
make docker-run
```

## Environment Variables

| Variable      | Description       | Default     |
| ------------- | ----------------- | ----------- |
| `DB_HOST`     | Database host     | localhost   |
| `DB_PORT`     | Database port     | 5432        |
| `DB_USER`     | Database user     | postgres    |
| `DB_PASSWORD` | Database password | password    |
| `DB_NAME`     | Database name     | go_template |
| `PORT`        | Server port       | 8080        |
| `JWT_SECRET`  | JWT secret key    | -           |

## API Documentation

ดู [API Documentation](./api/README.md) สำหรับรายละเอียดของ endpoints

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## วิธีการใช้งาน

1. Clone โปรเจกต์
2. ติดตั้ง dependencies: `go mod tidy`
3. รันโปรเจกต์: `make run`
4. ทดสอบ: `make test`

## การพัฒนา

-  ใช้ Hexagonal Architecture
-  แยก Business Logic ออกจาก Infrastructure
-  มี Integration Test และ Unit Test ครบถ้วน
