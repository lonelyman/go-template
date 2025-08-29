package middleware

import (
	"time"

	"go-template/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
// ⭐️⭐️⭐️ 2. อัปเกรด "ยามตรวจคนเข้าเมือง" ของเรา! ⭐️⭐️⭐️
// CORS configures Cross-Origin Resource Sharing.
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: []string{"https://gemini.google.com"},
		// AllowMethods: []string{
		// 	fiber.MethodGet,
		// 	fiber.MethodPost,
		// 	fiber.MethodHead,
		// 	fiber.MethodPut,
		// 	fiber.MethodDelete,
		// 	fiber.MethodPatch,
		// },
		// AllowHeaders:     []string{},
		AllowCredentials: false,

		// AllowOrigins: "http://localhost:3000", // 👈 ใน Production จริง เราจะระบุโดเมนของ Frontend แบบนี้

		// สำหรับ Development เราจะใช้ "*" เพื่ออนุญาต "ทุกเมือง" ไปก่อน
		// AllowOrigins: []string{"*"},

		// // อนุญาตให้ส่ง "บัตรประจำตัว" (Cookie) ข้ามเมืองได้
		// AllowCredentials: true,

		// // อนุญาตให้ใช้ Method เหล่านี้ได้
		// AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// // อนุญาตให้มี Header เหล่านี้ได้ (สำคัญมากสำหรับ Auth)
		// AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	})
}
