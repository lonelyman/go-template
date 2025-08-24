package middleware

import (
	"time"

	"go-template/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

// Logger is a middleware that logs HTTP requests.
// ✨ 2. แก้ไขให้รับ "นักข่าว" (Logger) เข้ามา ✨
func Logger(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// ไปทำงานใน Handler ต่อไปก่อน
		err := c.Next()

		stop := time.Now()
		latency := stop.Sub(start)

		// ✨ 3. ใช้ "นักข่าว" ของเราบันทึก Log! ✨
		log.Info("Request handled",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", latency.String(),
			"ip", c.IP(),
		)

		return err
	}
}

// CORS is a middleware for Cross-Origin Resource Sharing
func CORS() fiber.Handler {
	// ... (โค้ดส่วนนี้เหมือนเดิม) ...
	return func(c fiber.Ctx) error {
		// ...
		return c.Next()
	}
}
