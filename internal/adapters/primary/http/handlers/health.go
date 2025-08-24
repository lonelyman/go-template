package handlers

import (
	"context"
	"time"

	"go-template/pkg/response"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ⭐️ 1. สร้าง "พิมพ์เขียว" (Structs) สำหรับ Health Response โดยเฉพาะ ⭐️
type HealthResponse struct {
	Status       string           `json:"status"`
	Service      string           `json:"service"`
	Dependencies DependencyStatus `json:"dependencies"`
}

type DependencyStatus struct {
	Database string `json:"database"`
}

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	dbStatus := "ok"

	// พยายาม Ping DB
	sqlDB, err := h.db.DB()
	if err != nil {
		dbStatus = "error"
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			dbStatus = "error"
		}
	}

	// ⭐️ 2. สร้างข้อมูล Response โดยใช้ Struct ที่เราเพิ่งสร้าง ⭐️
	healthData := HealthResponse{
		Status:  "ok",
		Service: "go-template-api", // (อาจจะดึงมาจาก config ก็ได้นะ)
		Dependencies: DependencyStatus{
			Database: dbStatus,
		},
	}

	// ถ้า DB มีปัญหา ให้ตอบกลับด้วย Status 503
	if dbStatus == "error" {
		healthData.Status = "error"
		// ⭐️ 3. เรียกใช้ response.Success แม้กระทั่งตอน Error! ⭐️
		// เพราะเรายังอยากให้โครงสร้างเป็น {"data": ...} แต่บอกสถานะว่า error
		return response.Success(c, fiber.StatusServiceUnavailable, "Database connection error", healthData, nil)
	}

	// ⭐️ 4. เรียกใช้ "ผู้ช่วย" ของเราเพื่อส่ง Response ที่เป็นมาตรฐาน! ⭐️
	return response.Success(c, fiber.StatusOK, "Health check passed", healthData, nil)
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(app fiber.Router) {
	app.Get("/health", h.HealthCheck)
}
