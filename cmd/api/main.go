package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"go-template/internal/adapters/primary/http/handlers"
	"go-template/internal/adapters/primary/http/middleware"
	"go-template/internal/modules/example/example_user"
	"go-template/pkg/config"
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"go-template/pkg/platform/postgres"
	"go-template/pkg/response"
	"go-template/pkg/validator"
)

// "ช่องรับ" ข้อมูล Build ที่จะถูกยิงเข้ามาโดย Linker Flags (-ldflags)
var (
	AppVersion string
	BuildTime  string
	CommitHash string
)

// ⭐️ 1. สร้าง "นาฬิกาเทียบเวลา" ของเราเป็นตัวแปร Global ⭐️
var (
	bangkokLocation *time.Location
)

func main() {

	// --- ⭐️ 2. ตั้งค่า "นาฬิกาเทียบเวลา" เป็นอย่างแรกสุด! ⭐️ ---
	var err error
	bangkokLocation, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("❌ Failed to load Bangkok time zone: %v", err)
	}

	// --- 0. โหลด Environment Variables (สำหรับ Local Dev) ---
	if os.Getenv("DOCKER_ENV") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	// --- 1. โหลด Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// --- 2. สร้าง Logger ---
	var appLogger logger.Logger
	if cfg.Server.Mode == "development" {
		appLogger = logger.NewPrettyLogger()
	} else {
		appLogger = logger.NewSlogLogger()
	}
	appLogger.Info("Logger initialized", "mode", cfg.Server.Mode)

	appValidator := validator.New()

	// --- 3. เชื่อมต่อ Platforms (Databases) ---
	primaryDB, err := postgres.NewConnection(cfg.Postgres.Primary, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to primary database", err)
		os.Exit(1)
	}
	// เชื่อมต่อ Database สำหรับเก็บ Logs (ถ้ามีการตั้งค่า)
	var logsDB *gorm.DB
	if cfg.Postgres.Logs.Host != "" {
		logsDB, err = postgres.NewConnection(cfg.Postgres.Logs, appLogger)
		if err != nil {
			appLogger.Warn("Logs database configured but unavailable", "error", err)
			logsDB = nil
		}
	}

	// --- 4. ประกอบร่าง Modules (Dependency Injection) ---
	_ = logsDB // ป้องกัน unused variable

	healthHandler := handlers.NewHealthHandler(primaryDB)

	exampleUserRepo := example_user.NewExampleRepository(primaryDB, appLogger)
	exampleUserService := example_user.NewExampleUserService(exampleUserRepo, cfg.Auth.JWTSecret, appLogger)
	exampleUserHandler := example_user.NewExampleUserHandler(exampleUserService, appLogger, bangkokLocation, appValidator)

	// --- 5. ตั้งค่า Web Server (Fiber) ---
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s %s", cfg.App.Name, AppVersion),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*custom_errors.AppError); ok {
				return response.Error(c, appErr)
			}
			systemErr := custom_errors.SystemErrorWithDetails("An unexpected error occurred", err.Error())
			appLogger.Error("Unhandled error has occurred", systemErr)
			return response.Error(c, systemErr)
		},
	})

	// --- 6. ติดตั้ง Middlewares & Routes ---
	app.Use(middleware.Logger(appLogger))
	app.Use(middleware.CORS())

	healthHandler.RegisterRoutes(app)

	apiV1 := app.Group("/api/v1")
	example := apiV1.Group("/example")
	exampleUserHandler.RegisterRoutes(example)

	// --- 7. เริ่มและปิดการทำงานของ Server ---
	go func() {
		// แอปของเราจะ Listen ที่ AppPort "ข้างใน" Container เสมอ
		listenAddr := fmt.Sprintf(":%s", cfg.Server.AppPort)

		appLogger.Info("Server starting...",
			"version", AppVersion,
			"buildTime", BuildTime,
			"commit", CommitHash,
		)
		appLogger.Info("Application running",
			"internalPort", cfg.Server.AppPort,
			"externalUrl", fmt.Sprintf("http://localhost:%s", cfg.Server.HostPort),
		)

		if err := app.Listen(listenAddr); err != nil {
			appLogger.Error("Server failed to start", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	appLogger.Info("Shutting down server...")

	// ใช้ app.Shutdown() แบบไม่มี context ตามเวอร์ชัน Fiber ที่เราใช้
	if err := app.Shutdown(); err != nil {
		appLogger.Error("Server shutdown failed", err)
		os.Exit(1)
	}

	appLogger.Info("Server gracefully stopped")
}
