package example_module

import (
	"errors"
	"fmt"
)

// ExampleService interface defines the business logic contract
type ExampleService interface {
	CreateExample(req CreateExampleRequest) (*ExampleResponse, error)
	GetExample(id uint) (*ExampleResponse, error)
	UpdateExample(id uint, req UpdateExampleRequest) (*ExampleResponse, error)
	DeleteExample(id uint) error
	ListExamples(limit, offset int) ([]ExampleResponse, error)
}

// exampleService implements ExampleService
type exampleService struct {
	repo ExampleRepository
}

// NewExampleService creates a new instance of ExampleService
func NewExampleService(repo ExampleRepository) ExampleService {
	return &exampleService{repo: repo}
}

// CreateExample creates a new example
func (s *exampleService) CreateExample(req CreateExampleRequest) (*ExampleResponse, error) {
	// Check if email already exists
	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Create new example
	example := &ExampleDomain{
		Name:   req.Name,
		Email:  req.Email,
		Status: "active",
	}

	if err := s.repo.Create(example); err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	response := example.ToResponse()
	return &response, nil
}

// GetExample retrieves an example by ID
func (s *exampleService) GetExample(id uint) (*ExampleResponse, error) {
	example, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := example.ToResponse()
	return &response, nil
}

// UpdateExample updates an existing example
func (s *exampleService) UpdateExample(id uint, req UpdateExampleRequest) (*ExampleResponse, error) {
	example, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		example.Name = *req.Name
	}
	if req.Email != nil {
		// Check if email already exists for different record
		existing, _ := s.repo.GetByEmail(*req.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email already exists")
		}
		example.Email = *req.Email
	}
	if req.Status != nil {
		example.Status = *req.Status
	}

	if err := s.repo.Update(example); err != nil {
		return nil, fmt.Errorf("failed to update example: %w", err)
	}

	response := example.ToResponse()
	return &response, nil
}

// DeleteExample deletes an example
func (s *exampleService) DeleteExample(id uint) error {
	// Check if example exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// ListExamples retrieves all examples with pagination
func (s *exampleService) ListExamples(limit, offset int) ([]ExampleResponse, error) {
	examples, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list examples: %w", err)
	}

	responses := make([]ExampleResponse, len(examples))
	for i, example := range examples {
		responses[i] = example.ToResponse()
	}

	return responses, nil
}
