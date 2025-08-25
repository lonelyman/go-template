package example_user

import (
	"errors"
	"fmt"
	"go-template/pkg/auth"
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

// Service คือ "สัญญา" ที่ Handler จะเรียกใช้
// ✨ 1. แก้ไข "สัญญา" ให้รับ Domain object และ password ✨
type Service interface {
	CreateUser(userToCreate *Domain, plainPassword string) (*Domain, error)
	GetUserByID(id uint) (*Domain, error)
	ListUsersByPage(limit, offset int, sort string) ([]*Domain, int, error)
	ListUsersByCursor(cursor string, limit int, sort string) ([]*Domain, string, bool, error)
}

// service คือ struct ที่ทำงานจริง
type service struct {
	repo      Repository
	jwtSecret string
	log       logger.Logger
}

// NewExampleUserService คือโรงงานสร้าง Service
func NewExampleUserService(repo Repository, jwtSecret string, log logger.Logger) Service {
	return &service{repo: repo, jwtSecret: jwtSecret, log: log}
}

// --- Implementation ---

// ✨ 2. แก้ไข "เมธอด" ให้รับ Domain object และ password ✨
func (s *service) CreateUser(userToCreate *Domain, plainPassword string) (*Domain, error) {
	// 1. ตรวจสอบ Logic ว่า email ซ้ำหรือไม่
	// (ใช้ Email จาก Domain object ที่รับเข้ามา)
	existingUser, err := s.repo.GetByEmail(userToCreate.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, custom_errors.SystemErrorWithDetails("ไม่สามารถตรวจสอบอีเมลได้", err.Error())
	}
	if existingUser != nil {
		return nil, custom_errors.AlreadyExistsError("อีเมลนี้ถูกใช้งานแล้ว", nil)
	}

	// 2. Hash Password (ใช้ password ดิบๆ ที่รับเข้ามา)
	hashedPassword, err := auth.HashPassword(plainPassword)
	if err != nil {
		return nil, custom_errors.SystemErrorWithDetails("ไม่สามารถเข้ารหัสรหัสผ่านได้", err.Error())
	}

	// 3. เติมข้อมูลที่เหลือให้ Domain object ที่ได้รับมา
	userToCreate.PasswordHash = hashedPassword
	userToCreate.Status = "active" // กำหนดค่าเริ่มต้นทางธุรกิจ
	userToCreate.Role = "user"     // กำหนดค่าเริ่มต้นทางธุรกิจ

	// 4. เรียกใช้ Repo เพื่อบันทึกข้อมูล
	if err := s.repo.Create(userToCreate); err != nil {
		return nil, custom_errors.SystemErrorWithDetails("ไม่สามารถสร้างผู้ใช้งานได้", err.Error())
	}
	s.log.Dumpf(logger.LevelSuccess, "Full user object after creation:", userToCreate)
	// 5. คืนค่า Domain object ที่สมบูรณ์แล้ว (ตอนนี้มี ID, CreatedAt แล้ว) กลับไป
	return userToCreate, nil
}

func (s *service) GetUserByID(id uint) (*Domain, error) {
	// 1. สั่งงาน Repository ให้ไปหาข้อมูล
	userDomain, err := s.repo.GetByID(id)

	// 2. ⭐️ Service ทำหน้าที่ "ตีความ" Error! ⭐️
	if err != nil {
		// ถ้า Error ที่ได้คือ "หาไม่เจอ" (gorm.ErrRecordNotFound)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ให้แปลงเป็น Business Error ของเรา คือ NotFoundError
			return nil, custom_errors.NotFoundError("ไม่พบผู้ใช้งาน ID: " + fmt.Sprintf("%d", id))
		}

		// ถ้าเป็น Error อื่นๆ (เช่น DB down)
		// ให้แปลงเป็น System Error
		return nil, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ใช้", err.Error())
	}

	// 3. ถ้าไม่มี Error ก็ส่งข้อมูลกลับไปให้ Handler
	return userDomain, nil
}

// ListUsersByPage handles page-based pagination and sorting.
func (s *service) ListUsersByPage(limit, offset int, sort string) ([]*Domain, int, error) {
	// 1. "แปลภาษาเข็มทิศ" และตรวจสอบความปลอดภัย
	sortField, sortDirection, err := parseSortString(sort)
	if err != nil {
		// ถ้า Client ส่ง sort field ที่ไม่ได้รับอนุญาตมา ให้คืนค่า error
		return nil, 0, custom_errors.ValidationError("Sort parameter ไม่ถูกต้อง", err.Error())
	}

	// 2. เรียกใช้ Repository เพื่อดึงข้อมูลและจำนวนทั้งหมด
	userDomains, totalCount, repoErr := s.repo.ListByPage(limit, offset, sortField, sortDirection)
	if repoErr != nil {
		return nil, 0, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้", repoErr.Error())
	}

	return userDomains, totalCount, nil
}

// ListUsersByCursor handles cursor-based pagination and sorting.
func (s *service) ListUsersByCursor(cursor string, limit int, sort string) ([]*Domain, string, bool, error) {
	// ⭐️⭐️⭐️ หมายเหตุ: การ Implement Cursor-based Pagination จริงๆ นั้นซับซ้อนมาก
	// จะต้องมีการเข้ารหัส/ถอดรหัส cursor (เช่น base64 ของ ID หรือ Timestamp)
	// และ Logic การ query ใน Repository ก็จะซับซ้อนกว่านี้มาก
	//
	// โค้ดด้านล่างนี้เป็นแค่ "โครงสร้างตัวอย่าง" เพื่อให้เห็นภาพรวมเท่านั้นนะ!
	// ⭐️⭐️⭐️

	sortField, sortDirection, err := parseSortString(sort)
	if err != nil {
		return nil, "", false, custom_errors.ValidationError("Sort parameter ไม่ถูกต้อง", err.Error())
	}

	// (ในชีวิตจริง เราจะต้องถอดรหัส cursor ก่อน)
	// lastID, _ := decodeCursor(cursor)

	userDomains, repoErr := s.repo.ListByCursor(0, limit, sortField, sortDirection) // ส่ง lastID เข้าไป
	if repoErr != nil {
		return nil, "", false, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้", repoErr.Error())
	}

	// (ในชีวิตจริง เราจะต้องสร้าง nextCursor และเช็ค hasMore จากข้อมูลที่ได้)
	nextCursor := "next_cursor_placeholder"
	hasMore := len(userDomains) == limit

	return userDomains, nextCursor, hasMore, nil
}

// --- Private Helper ---

// parseSortString คือ "นักแปลภาษาเข็มทิศ"
// มันจะแกะ string "field:direction" ออกมา และตรวจสอบกับ "แผนที่" (whitelist)
func parseSortString(sort string) (field string, direction string, err error) {
	// ⭐️ แผนที่: กำหนดว่า Client สามารถเรียงข้อมูลจาก field ไหนได้บ้าง
	// เพื่อป้องกันการยิง query ที่ไม่เหมาะสม (เช่น เรียงจาก password_hash)
	allowedSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}

	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid sort format, must be 'field:direction'")
	}

	field = parts[0]
	direction = parts[1]

	if !allowedSortFields[field] {
		return "", "", errors.New("sorting by this field is not allowed: " + field)
	}

	if direction != "asc" && direction != "desc" {
		return "", "", errors.New("invalid sort direction, must be 'asc' or 'desc'")
	}

	return field, direction, nil
}
