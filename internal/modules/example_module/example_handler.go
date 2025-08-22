package example_module

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ExampleHandler handles HTTP requests for examples
type ExampleHandler struct {
	service ExampleService
}

// NewExampleHandler creates a new instance of ExampleHandler
func NewExampleHandler(service ExampleService) *ExampleHandler {
	return &ExampleHandler{service: service}
}

// CreateExample handles POST /examples - Auto-reload TEST! 🚀
func (h *ExampleHandler) CreateExample(c *fiber.Ctx) error {
	var req CreateExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.service.CreateExample(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetExample handles GET /examples/:id
func (h *ExampleHandler) GetExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	response, err := h.service.GetExample(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

// UpdateExample handles PUT /examples/:id
func (h *ExampleHandler) UpdateExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	var req UpdateExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.service.UpdateExample(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

// DeleteExample handles DELETE /examples/:id
func (h *ExampleHandler) DeleteExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	if err := h.service.DeleteExample(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListExamples handles GET /examples
func (h *ExampleHandler) ListExamples(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	responses, err := h.service.ListExamples(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":   responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

// RegisterRoutes registers all routes for the example module
func (h *ExampleHandler) RegisterRoutes(router fiber.Router) {
	examples := router.Group("/examples")
	examples.Post("", h.CreateExample)
	examples.Get("", h.ListExamples)
	examples.Get("/:id", h.GetExample)
	examples.Put("/:id", h.UpdateExample)
	examples.Delete("/:id", h.DeleteExample)
}
