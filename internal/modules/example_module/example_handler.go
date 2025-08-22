package example_module

import (
	"go-template/pkg/custom_errors"
	"go-template/pkg/response"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ====================================================================================
// API Data Transfer Objects (DTOs)
// - Structs ที่ใช้คุยกับโลกภายนอก (JSON, URI)
// - จะถูกประกาศและใช้งานเฉพาะในไฟล์ Handler นี้เท่านั้น
// ====================================================================================

type CreateExampleRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type GetExampleByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"` // gte=1 คือต้องมากกว่าหรือเท่ากับ 1
}

type UpdateExampleByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"` // gte=1 คือต้องมากกว่าหรือเท่ากับ 1
}

type UpdateExampleRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Status *string `json:"status,omitempty"`
}

type DeleteExampleByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"` // gte=1 คือต้องมากกว่าหรือเท่ากับ 1
}

type ExampleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExampleHandler handles HTTP requests for examples
type ExampleHandler struct {
	service ExampleService
}

// NewExampleHandler creates a new instance of ExampleHandler
func NewExampleHandler(service ExampleService) *ExampleHandler {
	return &ExampleHandler{service: service}
}

// CreateExample handles POST /examples - Auto-reload TEST! 🚀
func (h *ExampleHandler) CreateExample(c fiber.Ctx) error {
	var req CreateExampleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// error_message,err := validator.ValidateStructDetails(req)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }
	exampleDomain := ExampleDomain{
		Name:  req.Name,
		Email: req.Email,
	}
	resp, err := h.service.CreateExample(exampleDomain)
	if err != nil {
		return response.Error(c, err.(*custom_errors.AppError))
	}
	return response.Success(c, fiber.StatusCreated, resp)
}

/*
// GetExample handles GET /examples/:id
func (h *ExampleHandler) GetExampleByID(c fiber.Ctx) error {
	// 1. สร้าง struct มารอรับค่า
	params := new(GetExampleByIDParams)

	// 2. ⭐️ ใช้ c.Bind().URI() เพื่อดึงและแปลงค่า! ⭐️
	if err := c.Bind().URI(params); err != nil {
		// Error นี้จะเกิดขึ้นถ้า user ส่ง string ที่ไม่ใช่ตัวเลขเข้ามา
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{
			"id": "must be a positive integer",
		})
		return response.Error(c, appErr)
	}

	domainObject, err := h.service.GetExample(params.ID)
	if err != nil {
		return response.Error(c, err.(*custom_errors.AppError))
	}
	responsePayload := toExampleResponse(domainObject)
	return response.Success(c, fiber.StatusOK, responsePayload)
}

// UpdateExample handles PUT /examples/:id
func (h *ExampleHandler) UpdateExampleByID(c fiber.Ctx) error {
	// === 1. Bind & Validate Path Parameters ===
	params := new(UpdateExampleByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}
	// if validationDetails := validator.ValidateStruct(params); validationDetails != nil {
	// 	appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationDetails)
	// 	return response.Error(c, appErr)
	// }

	// === 2. Bind & Validate Request Body ===
	req := new(UpdateExampleRequest)
	if err := c.Bind().Body(&req); err != nil {
		// ✨ ดันเข้ากรอบ #1: ใช้มาตรฐาน Error ของเรา ✨
		appErr := custom_errors.New(fiber.StatusBadRequest, custom_errors.ErrInvalidFormat, "Request body is not valid JSON", err.Error())
		return response.Error(c, appErr)
	}
	// if validationDetails := validator.ValidateStruct(req); validationDetails != nil {
	// 	appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationDetails)
	// 	return response.Error(c, appErr)
	// }

	// === 3. Call Service (ยังไม่ต้อง Translate ก็ได้สำหรับ Update) ===
	// ✨ ดันเข้ากรอบ #2: Service ควรจะคืนค่า AppError ที่มีโครงสร้างชัดเจน ✨
	// เราส่ง ID และ Request DTO เข้าไปให้ Service จัดการต่อ
	updatedExample, serviceErr := h.service.UpdateExample(params.ID, req)
	if serviceErr != nil {
		// ถ้า serviceErr เป็น *custom_errors.AppError อยู่แล้ว ก็จะทำงานได้เลย
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// === 4. Send Success Response ===
	// ✨ ดันเข้ากรอบ #3: ใช้มาตรฐาน Success Response ของเรา ✨
	return response.Success(c, fiber.StatusOK, updatedExample)
}

// DeleteExample handles DELETE /examples/:id
func (h *ExampleHandler) DeleteExampleByID(c fiber.Ctx) error {
	// ⭐️ ใช้ Bind() กับ DTO ที่สร้างขึ้นมาใหม่
	params := new(DeleteExampleByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}
	// เรียก Service
	if serviceErr := h.service.DeleteExample(params.ID); serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}
	// ⭐️ ใช้ response helper ตอบกลับ
	return response.NoContent(c)
}

// ListExamples handles GET /examples
func (h *ExampleHandler) ListExamples(c fiber.Ctx) error {
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
*/

// RegisterRoutes registers all routes for the example module
func (h *ExampleHandler) RegisterRoutes(router fiber.Router) {
	examples := router.Group("/examples")
	examples.Post("", h.CreateExample)
	// examples.Get("", h.ListExamples)
	// examples.Get("/:id", h.GetExampleByID)
	// examples.Put("/:id", h.UpdateExampleByID)
	// examples.Delete("/:id", h.DeleteExampleByID)
}

// ====================================================================================
// Private Helper Functions
// ====================================================================================

// toExampleResponse คือ "ผู้ช่วย" ที่ทำหน้าที่แปลง Domain -> Response DTO
// จะถูกเรียกใช้ได้แค่ภายในไฟล์ example_handler.go นี้เท่านั้น
func toExampleResponse(domain *ExampleDomain) *ExampleResponse {
	return &ExampleResponse{
		ID:        domain.ID,
		Name:      domain.Name,
		Email:     domain.Email,
		Status:    domain.Status,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func toExampleResponseList(domains []ExampleDomain) []ExampleResponse {
	responses := make([]ExampleResponse, 0, len(domains))
	for _, domain := range domains {
		responses = append(responses, *toExampleResponse(&domain))
	}
	return responses
}
