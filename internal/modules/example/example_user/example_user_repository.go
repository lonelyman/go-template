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
	ListByPage(limit, offset int, sortField, sortDirection, search string) ([]*Domain, int, error)
	ListByCursor(cursorData *cursorInfo, limit int, sortField, sortDirection, search string) ([]*Domain, error)
	Update(d *Domain) error
	Delete(id uint) error
	HardDelete(id uint) error
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

// NewExampleUserRepository คือโรงงานสร้าง Repository
func NewExampleUserRepository(db *gorm.DB, log logger.Logger) Repository {
	return &repository{db: db, log: log}
}

// --- Implementation ---

func (r *repository) Create(d *Domain) error {
	gormModel := toGORMForCreate(d)
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
	r.log.Dump("result", gormModel)
	loc, _ := time.LoadLocation("Asia/Bangkok")
	fmt.Println("DB time (raw):", gormModel.CreatedAt)
	fmt.Println("DB time (Bangkok):", gormModel.CreatedAt.In(loc))
	//fmt.Println("------ggggg------", id, gormModel)
	return gormModel.toDomain(), nil
}

// ListByPage handles page-based pagination
func (r *repository) ListByPage(limit, offset int, sortField, sortDirection, search string) ([]*Domain, int, error) {
	var gormModels []Model
	var totalCount int64

	// 1. สร้าง query เริ่มต้น
	query := r.db.Model(&Model{})

	// 2. ⭐️ "ติดตั้งแว่นขยาย": เพิ่มเงื่อนไขการค้นหาถ้ามี ⭐️
	if search != "" {
		// ใช้ LIKE เพื่อค้นหาบางส่วนของคำใน field name หรือ email
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// 3. นับจำนวนทั้งหมด (หลังจากที่กรองด้วย search แล้ว)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 4. "ติดตั้งเข็มทิศ": สร้างคำสั่ง Order By
	orderClause := fmt.Sprintf("%s %s", sortField, sortDirection)

	// 5. ดึงข้อมูลตามหน้า
	result := query.Order(orderClause).Limit(limit).Offset(offset).Find(&gormModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// 6. แปล GORM Models กลับเป็น Domain Structs
	domains := make([]*Domain, 0, len(gormModels))
	for _, model := range gormModels {
		domains = append(domains, model.toDomain())
	}

	return domains, int(totalCount), nil
}

// ListByCursor handles cursor-based pagination, sorting, and searching.
func (r *repository) ListByCursor(cursorData *cursorInfo, limit int, sortField, sortDirection, search string) ([]*Domain, error) {
	var gormModels []Model

	query := r.db.Model(&Model{})

	// "ติดตั้งแว่นขยาย"
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// ⭐️⭐️⭐️ อัปเกรด Logic การสร้าง Query ให้ฉลาดขึ้น! ⭐️⭐️⭐️
	if sortField == "id" {
		// --- กรณีเรียงด้วย ID (ท่าธรรมดา) ---
		orderClause := fmt.Sprintf("id %s", sortDirection)
		query = query.Order(orderClause)

		if cursorData.LastID > 0 {
			if sortDirection == "asc" {
				query = query.Where("id > ?", cursorData.LastID)
			} else {
				query = query.Where("id < ?", cursorData.LastID)
			}
		}
	} else {
		// --- กรณีเรียงด้วย Field อื่น (ท่าไม้ตาย Keyset) ---
		orderClause := fmt.Sprintf("%s %s, id %s", sortField, sortDirection, sortDirection)
		query = query.Order(orderClause)

		if cursorData.LastID > 0 {
			var operator string
			if sortDirection == "asc" {
				operator = ">"
			} else {
				operator = "<"
			}
			query = query.Where(
				fmt.Sprintf("(%s %s ?) OR (%s = ? AND id %s ?)", sortField, operator, sortField, operator),
				cursorData.LastValue, cursorData.LastValue, cursorData.LastID,
			)
		}
	}

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

// Update saves the changes of an existing user to the database.
func (r *repository) Update(d *Domain) error {
	// 1. แปลง Domain object ที่อัปเดตแล้ว กลับไปเป็น GORM Model
	gormModel := toGORMForUpdate(d)

	// 2. ใช้ Save() เพื่อบันทึกการเปลี่ยนแปลง
	// GORM จะใช้ ID ที่มีอยู่ในการหาแถวที่จะ UPDATE
	// และจะอัปเดต "ทุก" field รวมถึง UpdatedAt ให้เป็นเวลาปัจจุบันโดยอัตโนมัติ
	result := r.db.Save(gormModel)
	if result.Error != nil {
		return result.Error
	}

	// (Optional) อัปเดตค่า UpdatedAt กลับไปที่ Domain object เพื่อความสมบูรณ์
	d.UpdatedAt = gormModel.UpdatedAt

	return nil
}

// Delete performs a soft delete on a user by their ID.
func (r *repository) Delete(id uint) error {
	// GORM's Delete function will automatically set the 'deleted_at' field
	// because our 'Model' struct embeds gorm.Model.
	result := r.db.Delete(&Model{}, id)
	if result.Error != nil {
		return result.Error
	}

	// GORM v2 returns an error if no record is found to delete.
	// We can check if any rows were actually affected.
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // คืนค่า error มาตรฐานเพื่อให้ Service ตีความ
	}

	return nil
}

// HardDelete performs a permanent delete from the database.
// 💥 WARNING: This action is irreversible! 💥
func (r *repository) HardDelete(id uint) error {
	// ⭐️ ใช้ .Unscoped() นำหน้า .Delete() ⭐️
	// เพื่อบอก GORM ให้ทำการ Hard Delete
	result := r.db.Unscoped().Delete(&Model{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// --- Translators ---

func toGORMForCreate(d *Domain) *Model {
	return &Model{
		Name:        d.Name,
		Email:       d.Email,
		Password:    d.PasswordHash,
		Status:      d.Status,
		Role:        d.Role,
		LastLoginAt: d.LastLoginAt,
	}
}

func toGORMForUpdate(d *Domain) *Model {
	return &Model{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:        d.Name,
		Email:       d.Email,
		Password:    d.PasswordHash,
		Status:      d.Status,
		Role:        d.Role,
		LastLoginAt: d.LastLoginAt,
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
