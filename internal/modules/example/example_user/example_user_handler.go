package example_user

import (
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"go-template/pkg/response"
	"go-template/pkg/validator"
	"time"

	govalidator "github.com/go-playground/validator/v10" // ⭐️ 1. ตั้งชื่อเล่นให้ไลบรารีเป็น "govalidator"
	"github.com/gofiber/fiber/v3"
)

// ====================================================================================
// DTOs (Data Transfer Objects)
// ====================================================================================

type CreateRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" vmsg:"required:กรุณาระบุรหัสผ่าน,min:รหัสผ่านต้องมีความยาวอย่างน้อย 8 ตัวอักษร"`
}

type GetUserByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"`
}

type ListUsersQuery struct {
	Limit  *int    `query:"limit"`
	Page   *int    `query:"page"`
	Offset *int    `query:"offset"`
	Cursor *string `query:"cursor"`
	Sort   *string `query:"sort" validate:"omitempty,sort_format"`
	Search *string `query:"search"`
}

type Response struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserByIDParams คือ DTO สำหรับ Path Parameter
type UpdateUserByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"`
}

// UpdateUserRequest คือ DTO สำหรับ JSON Body ที่จะอัปเดต
// ✨ ใช้ Pointer (*string) เพื่อให้เรารู้ว่า Client ส่ง field ไหนมาบ้าง ✨
type UpdateUserRequest struct {
	Name   *string `json:"name" validate:"omitempty,min=2"`
	Email  *string `json:"email" validate:"omitempty,email"`
	Status *string `json:"status" validate:"omitempty,oneof=active inactive" vmsg:"oneof:สถานะต้องเป็น active หรือ inactive"`
}

type DeleteUserByIDParams struct {
	ID uint `uri:"id" validate:"required,gte=1"`
}

// ChangePasswordRequest คือ DTO สำหรับเปลี่ยนรหัสผ่านโดยเฉพาะ
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}

// LoginRequest คือ DTO สำหรับรับ JSON body ตอนล็อกอิน
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ====================================================================================
// Handler
// ====================================================================================

// handler คือ struct ที่ทำงานจริง
type handler struct {
	service         Service
	log             logger.Logger
	bangkokLocation *time.Location
	validator       *govalidator.Validate // ⭐️ 2. ใช้ชื่อเล่นใหม่ในการอ้างอิง Type
}

// NewExampleUserHandler คือโรงงานสร้าง Handler
func NewExampleUserHandler(service Service, log logger.Logger, bangkokLocation *time.Location, validator *govalidator.Validate) *handler { // ⭐️ 3. ใช้ชื่อเล่นใหม่ในการอ้างอิง Type
	return &handler{
		service:         service,
		log:             log,
		bangkokLocation: bangkokLocation,
		validator:       validator,
	}
}

// RegisterRoutes ลงทะเบียน routes ทั้งหมดของโมดูลนี้
func (h *handler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {

	userRouter := router.Group("/users")

	// --- 🚪 Public Routes (ประตูสาธารณะ - ไม่ต้องใช้ Token) ---
	userRouter.Post("", h.CreateUser)
	userRouter.Post("/login", h.Login)
	userRouter.Get("", h.ListUsers) // เราอาจจะอยากให้ List เป็น public

	// --- 🔐 Protected Routes (โซนปลอดภัย - ต้องใช้ Token) ---
	// 1. สร้าง Group ใหม่ขึ้นมาจาก userRouter
	// 2. "ส่งยาม" (authMiddleware) เข้าไปเฝ้าที่ทางเข้าของ Group นี้
	protected := userRouter.Group("", authMiddleware)

	// 3. Route ทั้งหมดที่ถูกสร้างจาก `protected` group นี้
	// จะถูกป้องกันโดย "ยาม" ของเราโดยอัตโนมัติ!
	protected.Get("/:id", h.GetUserByID)
	protected.Put("/:id", h.UpdateUserByID)
	protected.Delete("/:id", h.DeleteUserByID)
	protected.Put("/:id/password", h.ChangePassword)
}

// --- Handler Methods ---

func (h *handler) CreateUser(c fiber.Ctx) error {
	req := new(CreateRequest)
	if err := c.Bind().Body(req); err != nil {
		appErr := custom_errors.InvalidFormatError("Request body is not valid JSON", err.Error())
		return response.Error(c, appErr)
	}

	if validationResult := validator.Validate(h.validator, req); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	domainData := &Domain{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUserDomain, serviceErr := h.service.CreateUser(domainData, req.Password)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	responsePayload := h.toResponse(createdUserDomain)
	return response.Success(c, fiber.StatusCreated, "User created successfully", responsePayload, nil)
}

func (h *handler) GetUserByID(c fiber.Ctx) error {
	params := new(GetUserByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}

	if validationResult := validator.Validate(h.validator, params); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	userDomain, serviceErr := h.service.GetUserByID(params.ID)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	responsePayload := h.toResponse(userDomain)
	return response.Success(c, fiber.StatusOK, "User retrieved successfully", responsePayload, nil)
}

func (h *handler) ListUsers(c fiber.Ctx) error {
	query := new(ListUsersQuery)
	if err := c.Bind().Query(query); err != nil {
		appErr := custom_errors.InvalidFormatError("Query parameter ไม่ถูกต้อง", err.Error())
		return response.Error(c, appErr)
	}

	if validationResult := validator.Validate(h.validator, query); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("Query parameter ไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	var search string
	if query.Search != nil {
		search = *query.Search
	}

	sort := "id:asc"
	if query.Sort != nil {
		sort = *query.Sort
	}

	if query.Cursor != nil {
		limit := 10
		if query.Limit != nil {
			limit = *query.Limit
		}

		userDomains, nextCursor, hasMore, serviceErr := h.service.ListUsersByCursor(*query.Cursor, limit, sort, search)
		if serviceErr != nil {
			return response.Error(c, serviceErr.(*custom_errors.AppError))
		}

		responsePayloads := h.toResponseList(userDomains)
		pagination := response.NewCursorPagination(nextCursor, hasMore)
		return response.Success(c, fiber.StatusOK, "Users retrieved successfully", responsePayloads, pagination)

	} else {
		limit := 10
		if query.Limit != nil {
			limit = *query.Limit
		}
		offset := 0
		if query.Offset != nil {
			offset = *query.Offset
		} else if query.Page != nil && *query.Page > 0 {
			offset = (*query.Page - 1) * limit
		}

		userDomains, totalCount, serviceErr := h.service.ListUsersByPage(limit, offset, sort, search)
		if serviceErr != nil {
			return response.Error(c, serviceErr.(*custom_errors.AppError))
		}

		responsePayloads := h.toResponseList(userDomains)
		pagination := response.NewPagePagination(totalCount, limit, offset)
		return response.Success(c, fiber.StatusOK, "Users retrieved successfully", responsePayloads, pagination)
	}
}

// UpdateUserByID handles PUT /users/:id
func (h *handler) UpdateUserByID(c fiber.Ctx) error {
	// 1. Bind & Validate Path Parameter
	params := new(UpdateUserByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, params); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 2. Bind & Validate Request Body
	req := new(UpdateUserRequest)
	if err := c.Bind().Body(req); err != nil {
		appErr := custom_errors.InvalidFormatError("Request body is not valid JSON", err.Error())
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, req); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 3. Call Service (ส่ง ID และ DTO ที่มี Pointer เข้าไป)
	updatedUserDomain, serviceErr := h.service.UpdateUser(params.ID, req)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// 4. Translate & Respond
	responsePayload := h.toResponse(updatedUserDomain)
	return response.Success(c, fiber.StatusOK, "User updated successfully", responsePayload, nil)
}

// DeleteUserByID handles DELETE /users/:id
func (h *handler) DeleteUserByID(c fiber.Ctx) error {
	// 1. Bind & Validate Path Parameter
	params := new(DeleteUserByIDParams)
	if err := c.Bind().URI(params); err != nil {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, params); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 2. Call Service
	if serviceErr := h.service.DeleteUser(params.ID); serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// 3. Respond with 204 No Content
	// นี่คือ Best Practice สำหรับ DELETE ที่สำเร็จ!
	return response.NoContent(c)
}

// ChangePassword handles PUT /users/:id/password
func (h *handler) ChangePassword(c fiber.Ctx) error {
	// 1. Bind & Validate Path Parameter (ID)
	params := new(UpdateUserByIDParams) // ใช้ DTO เดิมได้เลย
	if err := c.Bind().URI(params); err != nil {
		// สร้าง Error ที่เป็นมาตรฐานของเรา
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", fiber.Map{"id": "must be a positive integer"})
		// แล้วส่งกลับไปให้ "ผู้ช่วย" จัดการ
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, params); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ID ที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 2. Bind & Validate Request Body
	req := new(ChangePasswordRequest)
	if err := c.Bind().Body(req); err != nil {
		// สร้าง Error ที่เป็นมาตรฐานของเรา
		appErr := custom_errors.InvalidFormatError("Request body is not valid JSON", err.Error())
		// แล้วส่งกลับไปให้ "ผู้ช่วย" จัดการ
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, req); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ข้อมูลรหัสผ่านไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 3. Call Service
	if serviceErr := h.service.ChangePassword(params.ID, req.OldPassword, req.NewPassword); serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// 4. Respond with a simple success message
	return response.Message(c, fiber.StatusOK, "Password updated successfully")
}

// Login handles POST /users/login
func (h *handler) Login(c fiber.Ctx) error {
	// 1. Bind & Validate Request Body
	req := new(LoginRequest)
	if err := c.Bind().Body(req); err != nil {
		appErr := custom_errors.InvalidFormatError("Request body is not valid JSON", err.Error())
		return response.Error(c, appErr)
	}
	if validationResult := validator.Validate(h.validator, req); !validationResult.IsValid {
		appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
		return response.Error(c, appErr)
	}

	// 2. Call Service เพื่อขอ "กุญแจดิจิทัล" (JWT)
	// Service จะคืนค่าเป็น string (token) กลับมา
	token, serviceErr := h.service.Login(req.Email, req.Password)
	if serviceErr != nil {
		return response.Error(c, serviceErr.(*custom_errors.AppError))
	}

	// 3. ส่ง "กุญแจ" กลับไปให้ Client
	return response.Success(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token}, nil)
}

// --- Private Helpers ---
func (h *handler) toResponse(d *Domain) *Response {
	return &Response{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Status:    d.Status,
		Role:      d.Role,
		CreatedAt: d.CreatedAt.In(h.bangkokLocation),
		UpdatedAt: d.UpdatedAt.In(h.bangkokLocation),
	}
}

func (h *handler) toResponseList(domains []*Domain) []*Response {
	responses := make([]*Response, 0, len(domains))
	for _, d := range domains {
		responses = append(responses, h.toResponse(d))
	}
	return responses
}
