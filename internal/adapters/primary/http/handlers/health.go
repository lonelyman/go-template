package handlers

import "github.com/gofiber/fiber/v2"

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "go-template-api",
		"version": "1.0",
	})
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", h.HealthCheck)
	app.Get("/ping", h.HealthCheck) // Alternative endpoint
}
