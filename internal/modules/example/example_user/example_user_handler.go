package example_user

import (
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"go-template/pkg/response"
	"go-template/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

// ====================================================================================
// DTOs (Data Transfer Objects)
// ====================================================================================

type CreateRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" vmsg:"required:กรุณาระบุรหัสผ่าน, min:กรุณาระบุรหัสผ่านที่มีความยาวอย่างน้อย 8 ตัวอักษร"`
}

type Response struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
	Role   string `json:"role"`
}

type GetUserByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"`
}

// ====================================================================================
// Handler
// ====================================================================================

// handler คือ struct ที่ทำงานจริง
type handler struct {
	service Service // Handler จะรู้จักแค่ "สัญญา" ของ Service
	log     logger.Logger
}

// NewExampleUserHandler คือโรงงานสร้าง Handler
func NewExampleUserHandler(service Service, log logger.Logger) *handler {
	return &handler{service: service, log: log}
}

// --- Handler Methods ---

func (h *handler) CreateUser(c fiber.Ctx) error {
	// 1. Bind
	req := new(CreateRequest)
	if err := c.Bind().Body(req); err != nil {
		appErr := custom_errors.NewWithDetails(fiber.StatusBadRequest, custom_errors.ErrInvalidFormat, "Request body is not valid JSON", err.Error())
		return response.Error(c, appErr)
	}
	// 2. ⭐️ เรียกใช้ Validator ฉบับใหม่ของเรา! ⭐️
	if validationResult := validator.Validate(req); !validationResult.IsValid {
		// เราสามารถส่ง Error Details ที่สวยงามกลับไปได้เลย
		appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}
	// 3. Translate DTO -> Domain
	domainData := &Domain{
		Name:  req.Name,
		Email: req.Email,
	}
	// 4. Call Service
	createdUserDomain, serviceErr := h.service.CreateUser(domainData, req.Password)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}
	// 5. Translate Domain -> Response DTO & Respond
	responsePayload := toResponse(createdUserDomain)
	// ⭐️ เรียกใช้ Success แบบใหม่! ⭐️
	return response.Success(c, fiber.StatusCreated, "User created successfully", responsePayload, nil)
}

// ใน example_user_handler.go

// GetUserByID handles GET /users/:id
func (h *handler) GetUserByID(c fiber.Ctx) error {
	// 1. Bind & Validate Path Parameter
	params := new(GetUserByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(params); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 2. Call Service
	userDomain, serviceErr := h.service.GetUserByID(params.ID)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// 3. Translate Domain -> Response DTO & Respond
	responsePayload := toResponse(userDomain)
	return response.Success(c, fiber.StatusOK, "User retrieved successfully", responsePayload, nil)
}

// RegisterRoutes ลงทะเบียน routes ทั้งหมดของโมดูลนี้
func (h *handler) RegisterRoutes(router fiber.Router) {
	userRouter := router.Group("/users")
	userRouter.Post("", h.CreateUser)
	userRouter.Get("/:id", h.GetUserByID)
	// ... (GET, PUT, DELETE routes) ...
}

// --- Private Helper ---
func toResponse(d *Domain) *Response {
	return &Response{
		ID:     d.ID,
		Name:   d.Name,
		Email:  d.Email,
		Status: d.Status,
		Role:   d.Role,
	}
}
