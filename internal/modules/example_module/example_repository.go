package example_module

import (
	"errors"

	"gorm.io/gorm"
)

type ExampleModel struct {
	gorm.Model
	Name   string `gorm:"not null"`
	Email  string `gorm:"unique;not null"`
	Status string `gorm:"default:active"`
}

// TableName specifies the table name for GORM
func (ExampleModel) TableName() string {
	return "examples"
}

// ExampleRepository interface defines the contract for data access
type ExampleRepository interface {
	Create(example *ExampleDomain) error
	// GetByID(id uint) (*ExampleModel, error)
	GetByEmail(email string) (*ExampleDomain, error)
	// Update(example *ExampleModel) error
	// Delete(id uint) error
	// GetAll(limit, offset int) ([]ExampleModel, error)
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
func (r *exampleRepository) Create(domain *ExampleDomain) error {
	gormModel := toGORM(domain) // 👈 ใช้ toGORM ฉบับใหม่ที่สั้นลง

	result := r.db.Create(gormModel) // GORM ส่งข้อมูลไป DB
	if result.Error != nil {
		return result.Error
	}
	// ⭐️ พอถึงบรรทัดนี้ GORM จะเอา ID ที่ DB สร้างให้ มาใส่ใน gormModel.ID อัตโนมัติ

	// เราก็แค่เอา ID นั้นมาอัปเดตกลับไปที่ Domain ของเรา
	domain.ID = gormModel.ID
	domain.CreatedAt = gormModel.CreatedAt
	domain.UpdatedAt = gormModel.UpdatedAt

	return nil
}

// GetByEmail retrieves an example by email
func (r *exampleRepository) GetByEmail(email string) (*ExampleDomain, error) {
	var gormModel ExampleModel
	err := r.db.Where("email = ?", email).First(&gormModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	domain := gormModel.toDomain()
	return domain, err
}

/*
// GetByID retrieves an example by ID
func (r *exampleRepository) GetByID(id uint) (*ExampleDomain, error) {
	var example ExampleDomain
	err := r.db.First(&example, id).Error
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
*/

// ====================================================================================
// Private Helper Functions (นักแปลภาษา)
// ====================================================================================

// ⭐️ toGORM: คือ "นักแปล" ขาไป (Domain -> GORM) ⭐️
// รับ Domain object ที่บริสุทธิ์เข้ามา แล้วแปลงเป็น GORM model เพื่อเตรียมคุยกับ DB
func toGORM(domain *ExampleDomain) *ExampleModel {
	return &ExampleModel{
		Name:   domain.Name,
		Email:  domain.Email,
		Status: domain.Status,
	}
}

// ⭐️ toDomain: คือ "นักแปล" ขากลับ (GORM -> Domain) ⭐️
// เป็น method ของ GORM model รับหน้าที่แปลงตัวเองกลับไปเป็น Domain object ที่บริสุทธิ์
// เพื่อส่งกลับไปให้ Service ใช้งานต่อ
func (g *ExampleModel) toDomain() *ExampleDomain {
	return &ExampleDomain{
		ID:        g.ID,
		Name:      g.Name,
		Email:     g.Email,
		Status:    g.Status,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}
