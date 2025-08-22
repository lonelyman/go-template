package response

import (
	"go-template/pkg/custom_errors"

	"github.com/gofiber/fiber/v3"
)

// Success - สำหรับส่ง Response เมื่อทำงานสำเร็จ
func Success(c fiber.Ctx, httpStatus int, data interface{}) error {
	return c.Status(httpStatus).JSON(fiber.Map{
		"data": data,
	})
}

// Error - สำหรับส่ง Error Response ที่เป็นมาตรฐานของเรา
func Error(c fiber.Ctx, err *custom_errors.AppError) error {
	return c.Status(err.HTTPStatus).JSON(fiber.Map{
		"error": err,
	})
}

// NoContent - สำหรับส่ง Response เมื่อไม่มีข้อมูล
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
