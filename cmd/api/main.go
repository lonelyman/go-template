package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	"go-template/internal/adapters/primary/http/handlers"
	"go-template/internal/adapters/primary/http/middleware"
	"go-template/internal/modules/example_module"
	"go-template/pkg/config"
	postgres "go-template/pkg/platform/postgres"

	"gorm.io/gorm"
)

func main() {
	// 0. โหลด Environment Variables
	if os.Getenv("DOCKER_ENV") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	// 1. โหลด Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Primary Database (จำเป็นต้องมี)
	primaryDB, err := postgres.NewConnection(cfg.Postgres.Primary)
	if err != nil {
		log.Fatalf("Failed to connect to primary database: %v", err)
	}

	// Logs Database (ไม่มีก็ได้)
	var logsDB *gorm.DB               // ประกาศเป็น nil ไว้ก่อน
	if cfg.Postgres.Logs.Host != "" { // เช็คว่ามี config ของ logs db ไหม
		logsDB, err = postgres.NewConnection(cfg.Postgres.Logs)
		if err != nil {
			log.Printf("⚠️ Logs database configured but unavailable: %v", err)
			logsDB = nil // ถ้าต่อไม่ได้ก็ให้เป็น nil เหมือนเดิม
		}
	}

	// 5. ประกอบร่าง Modules (Dependency Injection)
	exampleRepo := example_module.NewExampleRepository(primaryDB)
	exampleService := example_module.NewExampleService(exampleRepo)
	exampleHandler := example_module.NewExampleHandler(exampleService)

	// ตัวอย่างการใช้งาน multiple databases
	_ = logsDB // ป้องกัน unused variable

	// Health handler
	healthHandler := handlers.NewHealthHandler()

	// 4. ตั้งค่า Web Server (Fiber)
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s %s", cfg.App.Name, cfg.App.Version),
	})

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Register health routes
	healthHandler.RegisterRoutes(app)

	// --- ลงทะเบียน Routes ของแต่ละ Module (แบบ Modular) --- ✨
	apiV1 := app.Group("/api/v1")

	// Register example module routes
	exampleHandler.RegisterRoutes(apiV1)

	// 5. เริ่มและปิดการทำงานของ Server (Start & Graceful Shutdown)
	go func() {
		listenAddr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		if err := app.Listen(listenAddr); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil { // ✨ ถูกต้องสำหรับ Fiber v2
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("✅ Server gracefully stopped")
}
