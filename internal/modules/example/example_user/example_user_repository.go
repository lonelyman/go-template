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
	ListByPage(limit, offset int, sortField, sortDirection string) ([]*Domain, int, error)
	ListByCursor(lastID uint, limit int, sortField, sortDirection string) ([]*Domain, error)
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
	r.log.Dumpf(logger.LevelDebug, "result", gormModel)
	loc, _ := time.LoadLocation("Asia/Bangkok")
	fmt.Println("DB time (raw):", gormModel.CreatedAt)
	fmt.Println("DB time (Bangkok):", gormModel.CreatedAt.In(loc))
	//fmt.Println("------ggggg------", id, gormModel)
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
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ListByPage handles page-based pagination
func (r *repository) ListByPage(limit, offset int, sortField, sortDirection string) ([]*Domain, int, error) {
	var gormModels []Model
	var totalCount int64

	// 1. นับจำนวนทั้งหมดก่อน (สำหรับ Pagination)
	if err := r.db.Model(&Model{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 2. สร้างคำสั่ง Order By
	orderClause := fmt.Sprintf("%s %s", sortField, sortDirection)

	// 3. ดึงข้อมูลตามหน้า
	result := r.db.Order(orderClause).Limit(limit).Offset(offset).Find(&gormModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// 4. แปลง GORM Models กลับเป็น Domain Structs
	domains := make([]*Domain, 0, len(gormModels))
	for _, model := range gormModels {
		domains = append(domains, model.toDomain())
	}

	return domains, int(totalCount), nil
}

// ListByCursor handles cursor-based pagination
func (r *repository) ListByCursor(lastID uint, limit int, sortField, sortDirection string) ([]*Domain, error) {
	var gormModels []Model

	// สร้าง query เริ่มต้น
	query := r.db.Model(&Model{})

	// สร้างคำสั่ง Order By
	orderClause := fmt.Sprintf("%s %s", sortField, sortDirection)
	query = query.Order(orderClause)

	// ⭐️ Logic ของ Cursor: ดึงข้อมูลที่ "ถัดจาก" รายการสุดท้ายที่เห็น ⭐️
	if lastID > 0 {
		// ตัวอย่างง่ายๆ คือดึง ID ที่มากกว่า ID สุดท้าย (ถ้าเรียงแบบ asc)
		// ในชีวิตจริง Logic ตรงนี้จะซับซ้อนกว่านี้มาก
		if sortField == "id" && sortDirection == "asc" {
			query = query.Where("id > ?", lastID)
		}
	}

	// ดึงข้อมูล
	result := query.Limit(limit).Find(&gormModels)
	if result.Error != nil {
		return nil, result.Error
	}

	// แปลง GORM Models กลับเป็น Domain Structs
	domains := make([]*Domain, 0, len(gormModels))
	for _, model := range gormModels {
		domains = append(domains, model.toDomain())
	}

	return domains, nil
}
