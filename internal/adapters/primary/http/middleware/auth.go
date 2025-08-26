// internal/adapters/primary/http/middleware/auth.go
package middleware

import (
	"strings"

	"go-template/pkg/auth"
	"go-template/pkg/custom_errors"
	"go-template/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware คือ "โรงงาน" ที่สร้างยามหน้าด่านของเรา
// มันจะรับ "บริษัทรักษาความปลอดภัย" (auth.Service) เข้ามาเป็นเครื่องมือ
func AuthMiddleware(authService auth.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. ขอดู "บัตรผ่าน" (Authorization Header)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// ถ้าไม่มีบัตรมาเลย ก็ไม่ต้องให้เข้า
			appErr := custom_errors.UnauthorizedError("Authorization header is required")
			return response.Error(c, appErr)
		}

		// 2. ตรวจสอบรูปแบบของบัตร (ต้องเป็น "Bearer [token]")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := custom_errors.New(fiber.StatusUnauthorized, custom_errors.ErrInvalidToken, "Invalid token format")
			return response.Error(c, appErr)
		}
		tokenString := parts[1]

		// 3. ส่งบัตรไปให้ "บริษัท" ตรวจสอบ
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			// ถ้าบริษัทบอกว่าบัตรปลอมหรือหมดอายุ ก็ไม่ต้องให้เข้า
			appErr := custom_errors.New(fiber.StatusUnauthorized, custom_errors.ErrInvalidToken, "Invalid or expired token")
			return response.Error(c, appErr)
		}

		// 4. ถ้าบัตรถูกต้อง...
		// เราจะเก็บข้อมูลที่ถอดรหัสได้ (claims) ไว้ใน c.Locals()
		// เพื่อให้ Handler ที่อยู่ถัดไปสามารถหยิบไปใช้ได้ (เช่น รู้ว่าใครคือเจ้าของ Request นี้)
		c.Locals("user_id", claims.UserID)
		c.Locals("user_role", claims.Role)

		// 5. "เปิดประตู" ให้ Request วิ่งเข้าไปทำงานต่อ
		return c.Next()
	}
}
