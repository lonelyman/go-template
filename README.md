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

โปรเจกต์ Go ที่ใช้ **Hexagonal Architecture** พร้อมด้วย modules, adapters และระบบ testing ที่สมบูรณ์

## 🏗️ โครงสร้างโปรเจกต์

```
.
├── cmd/                         # 🚀 Entry Points - จุดเริ่มต้นของแอปพลิเคชัน
│   └── api/
│       └── main.go              # ❤️ Main application - จุดเริ่มต้นหลัก
│
├── internal/                    # 🏢 Core Business Logic - หัวใจของระบบ
│   ├── modules/                 # 📦 Business Modules - โมดูลธุรกิจแต่ละฟีเจอร์
│   │   └── example_module/      # 📋 Example Module - ตัวอย่างโมดูลสมบูรณ์
│   │       ├── example_domain.go     # 🏛️ Domain Objects - โครงสร้างข้อมูลหลัก
│   │       ├── example_handler.go    # 🚪 HTTP Handlers - จัดการ HTTP requests
│   │       ├── example_repository.go # 🗃️ Data Repository - จัดการข้อมูล
│   │       ├── example_service.go    # ⚙️ Business Logic - ตรรกะทางธุรกิจ
│   │       └── example_test.go       # 🧪 Unit Tests - การทดสอบหน่วย
│   │
│   └── adapters/                # 🔌 External Adapters - ตัวเชื่อมต่อภายนอก
│       ├── primary/             # 📥 Inbound Adapters - รับข้อมูลเข้า
│       │   └── http/            # 🌐 HTTP Layer
│       │       ├── handlers/    # 🎯 HTTP Handlers
│       │       │   └── health.go     # ❤️ Health Check Endpoint
│       │       └── middleware/  # 🛡️ HTTP Middleware
│       │           └── middleware.go # 🔒 CORS, Auth, Logger
│       │
│       └── secondary/           # 📤 Outbound Adapters - ส่งข้อมูลออก
│           └── dhl/             # 📦 DHL Integration
│               └── dhl_adapter.go    # 🚚 DHL API Adapter
│
├── pkg/                         # 🧰 Shared Packages - แพ็คเกจใช้ร่วม
│   ├── auth/                    # 🔐 Authentication
│   │   └── jwt.go               # 🎫 JWT Token Management
│   ├── config/                  # ⚙️ Configuration
│   │   └── config.go            # 📋 App Configuration
│   ├── platform/                # 🏗️ Platform Integrations
│   │   ├── postgres/            # 🐘 PostgreSQL
│   │   │   └── postgres.go      # 💾 Database Connection
│   │   └── redis/               # 🔴 Redis
│   │       └── redis.go         # ⚡ Cache Connection
│   ├── utils/                   # 🛠️ Utilities
│   │   └── string.go            # 📝 String Helpers
│   └── validator/               # ✅ Validation
│       └── validator.go         # 🔍 Input Validation
│
├── tests/                       # 🧪 Integration Tests - การทดสอบรวม
│   └── main_test.go             # 🎯 Test Setup & Integration Tests
│
├── build/                       # 🏭 Build & Deployment
│   ├── dev/                     # 🔧 Development
│   │   └── Dockerfile           # 🐳 Dev Docker Image (with hot reload)
│   └── prod/                    # 🚀 Production
│       └── Dockerfile           # 🐳 Prod Docker Image (optimized)
│
├── assets/                      # 📁 Static Assets - ไฟล์คงที่
│   ├── fonts/                   # 🔤 Font Files
│   ├── images/                  # 🖼️ Image Files
│   └── templates/               # 📄 Templates
│       └── email.html           # 📧 Email Template
│
├── configs/                     # 📋 Configuration Files
│   ├── config.yml               # ⚙️ App Configuration
│   └── database-example.env     # 🗃️ Database Config Example
│
├── docs/                        # 📚 Documentation
│   ├── DOCKER.md                # 🐳 Docker Guide
│   ├── MULTI_DATABASE_USAGE.md  # 🗃️ Database Guide
│   ├── VIPER_EXPLANATION.md     # 📋 Config Management
│   └── VIPER_QUICK_GUIDE.md     # ⚡ Quick Config Guide
│
├── scripts/                     # 📜 Utility Scripts
│   └── clean-cache.sh           # 🧹 Cache Cleanup
│
├── docker-compose.yml           # 🐳 Production Docker Compose
├── docker-compose.dev.yml       # 🔧 Development Docker Compose
├── Makefile                     # 🎯 Build Commands
├── go.mod                       # 📦 Go Dependencies
├── go.sum                       # 🔒 Dependency Checksums
└── README.md                    # 📖 Project Documentation
```

## 🚀 เทคโนโลยีที่ใช้

### 🌟 Core Technologies

-  **Go 1.24** - Programming Language
-  **Fiber v3** - High-performance HTTP Framework
-  **PostgreSQL** - Primary Database
-  **Redis** - Caching & Session Storage

### 🏗️ Architecture & Patterns

-  **Hexagonal Architecture** - Clean Architecture Pattern
-  **Domain-Driven Design (DDD)** - Business Logic Organization
-  **Dependency Injection** - Loose Coupling via Interfaces
-  **Repository Pattern** - Data Access Abstraction

### 🛠️ Development Tools

-  **Docker** - Containerization
-  **Docker Compose** - Multi-container Development
-  **CompileDaemon** - Hot Reload for Development
-  **Viper** - Configuration Management
-  **JWT** - Authentication & Authorization

### 📊 Database & Caching

-  **GORM** - ORM for Database Operations
-  **PostgreSQL Driver** - Database Connectivity
-  **Redis Client** - Cache Operations

### 🧪 Testing & Quality

-  **Go Testing** - Built-in Testing Framework
-  **Integration Tests** - End-to-end Testing
-  **Unit Tests** - Component Testing

## ⚡ Quick Start

### 1. 📥 Clone & Setup

```bash
# Clone repository
git clone <repository-url>
cd go-template

# Install dependencies
go mod tidy
```

### 2. 🗃️ Setup Database

**Option A: Using Docker (Recommended)**

```bash
# Start PostgreSQL with Docker Compose
make docker-dev-up
```

**Option B: Local PostgreSQL**

```bash
# Install PostgreSQL and create database
createdb go_template

# Or using psql
psql -U postgres -c "CREATE DATABASE go_template;"
```

### 3. ⚙️ Configure Environment

```bash
# Copy environment template
cp configs/database-example.env .env

# Edit .env file with your settings
```

### 4. 🏃‍♂️ Run Application

```bash
# Development mode (hot reload)
make run

# Or run directly
go run cmd/api/main.go

# Using Docker Compose
make docker-dev-up
```

### 5. 🧪 Test Application

```bash
# Health check
curl http://localhost:9998/health

# Test example API
curl -X POST http://localhost:9998/api/v1/examples \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

## 🛠️ Development Commands

```bash
# Development
make run                 # Run application
make docker-dev-up       # Start development environment
make docker-dev-down     # Stop development environment

# Testing
make test                # Run all tests
make test-coverage       # Run tests with coverage
make test-integration    # Run integration tests

# Building
make build               # Build binary
make docker-build        # Build Docker image

# Utilities
make clean               # Clean build artifacts
make lint                # Run linter
```

## 🏗️ Project Architecture

### 🔄 Request Flow

```
HTTP Request → Middleware → Handler → Service → Repository → Database
                    ↓
HTTP Response ← Handler ← Service ← Repository ← Database
```

### 📦 Module Structure

แต่ละ module ใน `internal/modules/` จะมีโครงสร้าง:

```go
// domain.go - ข้อมูลหลักและ business rules
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// repository.go - interface และ implementation สำหรับข้อมูล
type UserRepository interface {
    Create(user *User) error
    GetByID(id int) (*User, error)
}

// service.go - ตรรกะทางธุรกิจ
type UserService interface {
    CreateUser(name, email string) (*User, error)
    GetUser(id int) (*User, error)
}

// handler.go - HTTP handlers
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // Handle HTTP request/response
}
```

## 🔧 Environment Variables

| Variable      | Description       | Default       | Required |
| ------------- | ----------------- | ------------- | -------- |
| `DB_HOST`     | PostgreSQL Host   | `localhost`   | ✅       |
| `DB_PORT`     | PostgreSQL Port   | `5432`        | ✅       |
| `DB_USER`     | Database User     | `postgres`    | ✅       |
| `DB_PASSWORD` | Database Password | `password`    | ✅       |
| `DB_NAME`     | Database Name     | `go_template` | ✅       |
| `PORT`        | Server Port       | `9998`        | ✅       |
| `JWT_SECRET`  | JWT Secret Key    | -             | ✅       |
| `REDIS_HOST`  | Redis Host        | `localhost`   | ❌       |
| `REDIS_PORT`  | Redis Port        | `6379`        | ❌       |

## 📚 API Documentation

### 🏥 Health Check

```bash
GET /health
```

### 📋 Example Module APIs

```bash
# Create Example
POST /api/v1/examples
Content-Type: application/json
{
  "name": "John Doe",
  "email": "john@example.com"
}

# Get All Examples
GET /api/v1/examples

# Get Example by ID
GET /api/v1/examples/{id}

# Update Example
PUT /api/v1/examples/{id}
Content-Type: application/json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}

# Delete Example
DELETE /api/v1/examples/{id}
```

## 🎯 สร้าง Module ใหม่

### 1. สร้างโครงสร้างไฟล์

```bash
mkdir -p internal/modules/your_module
cd internal/modules/your_module

# สร้างไฟล์หลัก
touch your_domain.go
touch your_repository.go
touch your_service.go
touch your_handler.go
touch your_test.go
```

### 2. ตัวอย่างโครงสร้าง Module

```go
// your_domain.go
type YourEntity struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" validate:"required"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// your_repository.go
type YourRepository interface {
    Create(entity *YourEntity) error
    GetByID(id int) (*YourEntity, error)
    GetAll() ([]*YourEntity, error)
    Update(entity *YourEntity) error
    Delete(id int) error
}

// your_service.go
type YourService interface {
    CreateEntity(name string) (*YourEntity, error)
    GetEntity(id int) (*YourEntity, error)
    GetAllEntities() ([]*YourEntity, error)
    UpdateEntity(id int, name string) (*YourEntity, error)
    DeleteEntity(id int) error
}

// your_handler.go
type YourHandler struct {
    service YourService
}

func (h *YourHandler) SetupRoutes(app *fiber.App) {
    api := app.Group("/api/v1/your-entities")
    api.Post("/", h.Create)
    api.Get("/", h.GetAll)
    api.Get("/:id", h.GetByID)
    api.Put("/:id", h.Update)
    api.Delete("/:id", h.Delete)
}
```

### 3. Register ใน main.go

```go
// cmd/api/main.go
yourRepo := your_module.NewYourRepository(db)
yourService := your_module.NewYourService(yourRepo)
yourHandler := your_module.NewYourHandler(yourService)
yourHandler.SetupRoutes(app)
```

## 🐳 Docker Usage

### Development Environment

```bash
# Start full development stack
make docker-dev-up

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop development stack
make docker-dev-down
```

### Production Build

```bash
# Build production image
make docker-build

# Run production container
make docker-run
```

## 🧪 Testing Strategy

### Unit Tests

-  ทดสอบ Service layer logic
-  Mock dependencies (repository, external services)
-  ใช้ table-driven tests

```bash
# Run unit tests
go test ./internal/modules/...

# With coverage
go test -cover ./internal/modules/...
```

### Integration Tests

-  ทดสอบ API endpoints จริง
-  ใช้ test database
-  ทดสอบ database operations

```bash
# Run integration tests
go test ./tests/...
```

## 🚀 Deployment

### Build Binary

```bash
# Build for current platform
make build

# Build for Linux (production)
GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go
```

### Docker Production

```bash
# Build production image
docker build -f build/prod/Dockerfile -t go-template:latest .

# Run production container
docker run -p 9998:9998 --env-file .env go-template:latest
```

## 🛡️ Security Features

-  **JWT Authentication** - Token-based auth
-  **CORS Middleware** - Cross-origin protection
-  **Input Validation** - Request validation
-  **SQL Injection Protection** - GORM ORM protection
-  **Environment Variables** - Sensitive data protection

## 🔍 Monitoring & Logging

-  **Health Check Endpoint** - `/health`
-  **Structured Logging** - JSON format logs
-  **Request Logging** - HTTP request/response logs
-  **Error Handling** - Centralized error handling

## 🤝 Contributing

1. **Fork** the repository
2. **Create** feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** changes (`git commit -m 'Add amazing feature'`)
4. **Push** to branch (`git push origin feature/amazing-feature`)
5. **Create** Pull Request

### Code Style

-  Follow **Go conventions**
-  Use **gofmt** for formatting
-  Add **comments** for public functions
-  Write **tests** for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

If you have any questions or need help, please:

-  Open an **Issue** on GitHub
-  Check the **Documentation** in `/docs`
-  Review **Example Module** for reference

---

**Happy Coding! 🚀**
