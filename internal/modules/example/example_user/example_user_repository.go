package example_user

import (
	"fmt"
	"go-template/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// Repository คือ "สัญญา" ที่ Service จะเรียกใช้
type Repository interface {
	Create(d *Domain) error
	GetByEmail(email string) (*Domain, error)
	GetByID(id uint) (*Domain, error)
}

// Model คือ "ชุดเกราะ" สำหรับ GORM
type Model struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Email       string `gorm:"uniqueIndex:idx_unique_active_email;not null"`
	Password    string `gorm:"not null"`
	Status      string `gorm:"not null;default:active"`
	Role        string `gorm:"not null;default:user"`
	LastLoginAt *time.Time
}

func (Model) TableName() string {
	return "example_users"
}

// repository คือ struct ที่ทำงานจริง
type repository struct {
	db  *gorm.DB
	log logger.Logger
}

// NewExampleRepository คือโรงงานสร้าง Repository
func NewExampleRepository(db *gorm.DB, log logger.Logger) Repository {
	return &repository{db: db, log: log}
}

// --- Implementation ---

func (r *repository) Create(d *Domain) error {
	gormModel := toGORM(d)
	result := r.db.Create(gormModel)
	if result.Error != nil {
		fmt.Print("Failed to create user in database")
		return result.Error
	}
	*d = *gormModel.toDomain() // อัปเดตค่าที่ DB สร้างให้กลับไปที่ Domain object
	return nil
}

func (r *repository) GetByEmail(email string) (*Domain, error) {
	var gormModel Model
	result := r.db.Where("email = ?", email).First(&gormModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return gormModel.toDomain(), nil
}

func (r *repository) GetByID(id uint) (*Domain, error) {
	var gormModel Model
	result := r.db.First(&gormModel, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return gormModel.toDomain(), nil
}

// --- Translators ---

func toGORM(d *Domain) *Model {
	return &Model{
		Name:     d.Name,
		Email:    d.Email,
		Password: d.PasswordHash,
		Status:   d.Status,
		Role:     d.Role,
	}
}

func (m *Model) toDomain() *Domain {
	return &Domain{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.Password,
		Status:       m.Status,
		Role:         m.Role,
		LastLoginAt:  m.LastLoginAt,
	}
}
