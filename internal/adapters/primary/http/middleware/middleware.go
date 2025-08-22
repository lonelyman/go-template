package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Logger middleware for Fiber
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Continue processing request
		err := c.Next()

		// Log after processing
		fmt.Printf("%s - [%s] \"%s %s %s %d %s\"\n",
			c.IP(),
			start.Format(time.RFC3339),
			c.Method(),
			c.Path(),
			c.Protocol(),
			c.Response().StatusCode(),
			time.Since(start),
		)

		return err
	}
}

// CORS middleware
func CORS() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}

// Auth middleware (placeholder)
func Auth() fiber.Handler {
	return func(c fiber.Ctx) error {
		// TODO: Implement JWT authentication
		return c.Next()
	}
}
