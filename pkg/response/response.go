// pkg/response/response.go
package response

import (
	"math"

	"go-template/pkg/custom_errors"

	"github.com/gofiber/fiber/v3"
)

// ====================================================================================
// Response Structs (พิมพ์เขียวของ Response)
// ====================================================================================

// Pagination คือพิมพ์เขียวสำหรับข้อมูลการแบ่งหน้า
// รองรับทั้ง Page-based และ Cursor-based โดยใช้ Pointer และ omitempty
// เพื่อให้แสดงผลเฉพาะ field ที่จำเป็นในแต่ละสถานการณ์
type Pagination struct {
	// --- Page-based fields ---
	TotalRecords *int `json:"total_records,omitempty"`
	Limit        *int `json:"limit,omitempty"`
	Offset       *int `json:"offset,omitempty"`
	TotalPages   *int `json:"total_pages,omitempty"`
	CurrentPage  *int `json:"current_page,omitempty"`

	// --- Cursor-based fields ---
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    *bool   `json:"has_more,omitempty"`
}

// ====================================================================================
// Constructors & Helpers (โรงงานสร้าง Response)
// ====================================================================================

// NewPagePagination คือโรงงานสำหรับสร้าง Page-based Pagination object
// มันจะคำนวณ TotalPages และ CurrentPage ให้เราโดยอัตโนมัติ!
func NewPagePagination(totalRecords, limit, offset int) *Pagination {
	if limit <= 0 {
		limit = 1 // ป้องกันการหารด้วยศูนย์
	}
	if offset < 0 {
		offset = 0
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	currentPage := int(math.Floor(float64(offset)/float64(limit))) + 1

	return &Pagination{
		TotalRecords: &totalRecords,
		Limit:        &limit,
		Offset:       &offset,
		TotalPages:   &totalPages,
		CurrentPage:  &currentPage,
	}
}

// Success คือ "ผู้ช่วย" หลักสำหรับส่ง Response เมื่อทำงานสำเร็จ
// สามารถรับ pagination ที่เป็น optional ได้
func Success(c fiber.Ctx, httpStatus int, message string, data interface{}, pagination *Pagination) error {
	// สร้าง Body พื้นฐาน
	body := fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}

	// ถ้ามีข้อมูล pagination ส่งเข้ามาด้วย ก็เพิ่มเข้าไปใน Body
	if pagination != nil {
		body["pagination"] = pagination
	}

	return c.Status(httpStatus).JSON(body)
}

// Error คือ "ผู้ช่วย" หลักสำหรับส่ง Error Response ที่เป็นมาตรฐานของเรา
func Error(c fiber.Ctx, err *custom_errors.AppError) error {
	return c.Status(err.HTTPStatus).JSON(fiber.Map{
		"success": false,
		"message": err.Message,
		"error": fiber.Map{
			"code":    err.Code,
			"details": err.Details,
		},
	})
}

// Message คือ "ผู้ช่วย" สำหรับส่งแค่ข้อความกลับไป, ไม่มี data
func Message(c fiber.Ctx, httpStatus int, message string) error {
	return c.Status(httpStatus).JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

// NoContent คือ "ผู้ช่วย" สำหรับส่ง Response ที่ไม่มี Body กลับไป (HTTP 204)
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
