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
	"go-template/pkg/auth"
	"go-template/pkg/config"
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"go-template/pkg/platform/postgres"
	"go-template/pkg/response"
	"go-template/pkg/validator"
)

// "ช่องรับ" ข้อมูล Build
var (
	AppVersion string
	BuildTime  string
	CommitHash string
)

// "นาฬิกาเทียบเวลา" Global
var (
	bangkokLocation *time.Location
)

func main() {
	// --- ตั้งค่า Time Zone ---
	var err error
	bangkokLocation, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("❌ Failed to load Bangkok time zone: %v", err)
	}

	// --- โหลด .env ---
	if os.Getenv("DOCKER_ENV") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	// --- โหลด Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// --- สร้าง Dependencies ส่วนกลาง ---
	var appLogger logger.Logger
	if cfg.Server.Mode == "development" {
		appLogger = logger.NewPrettyLogger()
	} else {
		appLogger = logger.NewSlogLogger()
	}
	appLogger.Info("Logger initialized", "mode", cfg.Server.Mode)

	appValidator := validator.New()
	authService := auth.NewAuthService(cfg.Auth.JWTSecret)

	authMiddleware := middleware.AuthMiddleware(authService)

	// --- เชื่อมต่อ Platforms ---
	primaryDB, err := postgres.NewConnection(cfg.Postgres.Primary, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to primary database", err)
		os.Exit(1)
	}
	var logsDB *gorm.DB
	if cfg.Postgres.Logs.Host != "" {
		logsDB, err = postgres.NewConnection(cfg.Postgres.Logs, appLogger)
		if err != nil {
			appLogger.Warn("Logs database configured but unavailable", "error", err)
			logsDB = nil
		}
	}

	// --- ประกอบร่าง Modules (Dependency Injection) ---
	_ = logsDB

	healthHandler := handlers.NewHealthHandler(primaryDB)

	// ⭐️⭐️⭐️ แก้ไขตรงนี้ให้ถูกต้อง! ⭐️⭐️⭐️
	exampleUserRepo := example_user.NewExampleUserRepository(primaryDB, appLogger)
	exampleUserService := example_user.NewExampleUserService(exampleUserRepo, authService, appLogger)
	exampleUserHandler := example_user.NewExampleUserHandler(exampleUserService, appLogger, bangkokLocation, appValidator)

	// --- ตั้งค่า Web Server (Fiber) ---
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

	// --- ติดตั้ง Middlewares & Routes ---
	app.Use(middleware.Logger(appLogger))
	app.Use(middleware.CORS())

	healthHandler.RegisterRoutes(app)

	apiV1 := app.Group("/api/v1")
	example := apiV1.Group("/example")
	exampleUserHandler.RegisterRoutes(example, authMiddleware)

	// --- เริ่มและปิดการทำงานของ Server ---
	go func() {
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

	if err := app.Shutdown(); err != nil {
		appLogger.Error("Server shutdown failed", err)
		os.Exit(1)
	}

	appLogger.Info("Server gracefully stopped")
}
