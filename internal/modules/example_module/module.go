package example_module

import (
"github.com/gofiber/fiber/v2"
"gorm.io/gorm"
)

// ExampleModule represents the complete example module
type ExampleModule struct {
	handler *ExampleHandler
}

// NewExampleModule creates a new instance of ExampleModule with all dependencies
func NewExampleModule(db *gorm.DB) *ExampleModule {
	repo := NewExampleRepository(db)
	service := NewExampleService(repo)
	handler := NewExampleHandler(service)

	return &ExampleModule{
		handler: handler,
	}
}

// RegisterRoutes registers all routes for the example module
func (m *ExampleModule) RegisterRoutes(router fiber.Router) {
	examples := router.Group("/examples")
	examples.Post("", m.handler.CreateExample)
	examples.Get("/:id", m.handler.GetExample)
	examples.Put("/:id", m.handler.UpdateExample)
	examples.Delete("/:id", m.handler.DeleteExample)
}
