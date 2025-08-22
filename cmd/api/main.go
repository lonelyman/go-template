package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go-template/internal/adapters/primary/http/middleware"
	exampleModule "go-template/internal/modules/example-module"
	"go-template/pkg/config"
	"go-template/pkg/platform"
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

	// 2. เชื่อมต่อ Database (แบบ DI)
	db, err := platform.InitPostgres() // ✨ ใช้วิธีที่มีจริง
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Successfully connected to PostgreSQL")

	// 3. ประกอบร่าง Modules (Dependency Injection)
	exampleRepo := exampleModule.NewExampleRepository(db)
	exampleService := exampleModule.NewExampleService(exampleRepo)
	exampleHandler := exampleModule.NewExampleHandler(exampleService)

	// 4. ตั้งค่า Web Server (Fiber)
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s %s", cfg.App.Name, cfg.App.Version),
	})

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// --- ลงทะเบียน Routes ของแต่ละ Module (แบบ Modular) --- ✨
	apiV1 := app.Group("/api/v1")

	// Register example module routes
	examples := apiV1.Group("/examples")
	examples.Post("", exampleHandler.CreateExample)
	examples.Get("/:id", exampleHandler.GetExample)
	examples.Put("/:id", exampleHandler.UpdateExample)
	examples.Delete("/:id", exampleHandler.DeleteExample)

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
