package example_user

/*
import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExampleRepository is a mock implementation of ExampleRepository
type MockExampleRepository struct {
	mock.Mock
}

func (m *MockExampleRepository) Create(example *ExampleDomain) error {
	args := m.Called(example)
	return args.Error(0)
}

func (m *MockExampleRepository) GetByID(id uint) (*ExampleDomain, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExampleDomain), args.Error(1)
}

func (m *MockExampleRepository) GetByEmail(email string) (*ExampleDomain, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExampleDomain), args.Error(1)
}

func (m *MockExampleRepository) Update(example *ExampleDomain) error {
	args := m.Called(example)
	return args.Error(0)
}

func (m *MockExampleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExampleRepository) GetAll(limit, offset int) ([]ExampleDomain, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]ExampleDomain), args.Error(1)
}

func TestExampleService_CreateExample(t *testing.T) {
	// Arrange
	mockRepo := new(MockExampleRepository)
	service := NewExampleService(mockRepo)

	req := CreateExampleRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Mock repository calls
	mockRepo.On("GetByEmail", req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*example_module.ExampleDomain")).Return(nil)

	// Act
	result, err := service.CreateExample(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, "active", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestExampleService_CreateExample_EmailExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockExampleRepository)
	service := NewExampleService(mockRepo)

	req := CreateExampleRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	existingExample := &ExampleDomain{
		ID:    1,
		Email: req.Email,
	}

	// Mock repository calls
	mockRepo.On("GetByEmail", req.Email).Return(existingExample, nil)

	// Act
	result, err := service.CreateExample(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email already exists")

	mockRepo.AssertExpectations(t)
}

func TestExampleService_GetExample(t *testing.T) {
	// Arrange
	mockRepo := new(MockExampleRepository)
	service := NewExampleService(mockRepo)

	exampleID := uint(1)
	expectedExample := &ExampleDomain{
		ID:     exampleID,
		Name:   "Test User",
		Email:  "test@example.com",
		Status: "active",
	}

	// Mock repository calls
	mockRepo.On("GetByID", exampleID).Return(expectedExample, nil)

	// Act
	result, err := service.GetExample(exampleID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.ID, result.ID)
	assert.Equal(t, expectedExample.Name, result.Name)
	assert.Equal(t, expectedExample.Email, result.Email)

	mockRepo.AssertExpectations(t)
}
*/
