package example_module

import (
	"errors"

	"gorm.io/gorm"
)

// ExampleRepository interface defines the contract for data access
type ExampleRepository interface {
	Create(example *ExampleDomain) error
	GetByID(id uint) (*ExampleDomain, error)
	GetByEmail(email string) (*ExampleDomain, error)
	Update(example *ExampleDomain) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]ExampleDomain, error)
}

// exampleRepository implements ExampleRepository
type exampleRepository struct {
	db *gorm.DB
}

// NewExampleRepository creates a new instance of ExampleRepository
func NewExampleRepository(db *gorm.DB) ExampleRepository {
	return &exampleRepository{db: db}
}

// Create creates a new example
func (r *exampleRepository) Create(example *ExampleDomain) error {
	return r.db.Create(example).Error
}

// GetByID retrieves an example by ID
func (r *exampleRepository) GetByID(id uint) (*ExampleDomain, error) {
	var example ExampleDomain
	err := r.db.First(&example, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("example not found")
	}
	return &example, err
}

// GetByEmail retrieves an example by email
func (r *exampleRepository) GetByEmail(email string) (*ExampleDomain, error) {
	var example ExampleDomain
	err := r.db.Where("email = ?", email).First(&example).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("example not found")
	}
	return &example, err
}

// Update updates an existing example
func (r *exampleRepository) Update(example *ExampleDomain) error {
	return r.db.Save(example).Error
}

// Delete deletes an example by ID
func (r *exampleRepository) Delete(id uint) error {
	return r.db.Delete(&ExampleDomain{}, id).Error
}

// GetAll retrieves all examples with pagination
func (r *exampleRepository) GetAll(limit, offset int) ([]ExampleDomain, error) {
	var examples []ExampleDomain
	err := r.db.Limit(limit).Offset(offset).Find(&examples).Error
	return examples, err
}
