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

	sort := "id:asc"
	if query.Sort != nil {
		sort = *query.Sort
	}

	if query.Cursor != nil {
		limit := 10
		if query.Limit != nil {
			limit = *query.Limit
		}

		userDomains, nextCursor, hasMore, serviceErr := h.service.ListUsersByCursor(*query.Cursor, limit, sort)
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

		userDomains, totalCount, serviceErr := h.service.ListUsersByPage(limit, offset, sort)
		if serviceErr != nil {
			return response.Error(c, serviceErr.(*custom_errors.AppError))
		}

		responsePayloads := h.toResponseList(userDomains)
		pagination := response.NewPagePagination(totalCount, limit, offset)
		return response.Success(c, fiber.StatusOK, "Users retrieved successfully", responsePayloads, pagination)
	}
}

// RegisterRoutes ลงทะเบียน routes ทั้งหมดของโมดูลนี้
func (h *handler) RegisterRoutes(router fiber.Router) {
	userRouter := router.Group("/users")
	userRouter.Post("", h.CreateUser)
	userRouter.Get("", h.ListUsers)
	userRouter.Get("/:id", h.GetUserByID)
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
