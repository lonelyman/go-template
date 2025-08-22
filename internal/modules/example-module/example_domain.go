package example_module

import (
	"time"
)

// ExampleDomain represents the core business entity
type ExampleDomain struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Status    string    `json:"status" gorm:"default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ExampleDomain) TableName() string {
	return "examples"
}

// CreateExampleRequest represents the request to create an example
type CreateExampleRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateExampleRequest represents the request to update an example
type UpdateExampleRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Status *string `json:"status,omitempty"`
}

// ExampleResponse represents the response format
type ExampleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts domain to response
func (e *ExampleDomain) ToResponse() ExampleResponse {
	return ExampleResponse{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
