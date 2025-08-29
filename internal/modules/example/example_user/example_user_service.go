package example_user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"go-template/pkg/auth"
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Service คือ "สัญญา" ที่ Handler จะเรียกใช้
// ✨ 1. แก้ไข "สัญญา" ให้รับ Domain object และ password ✨
type Service interface {
	CreateUser(userToCreate *Domain, plainPassword string) (*Domain, error)
	GetUserByID(id uint) (*Domain, error)
	ListUsersByPage(limit, offset int, sort, search string) ([]*Domain, int, error)
	ListUsersByCursor(cursor string, limit int, sort, search string) ([]*Domain, string, bool, error)
	UpdateUser(id uint, req *UpdateUserRequest) (*Domain, error)
	DeleteUser(id uint) error
	ChangePassword(id uint, oldPassword, newPassword string) error
	Login(email, password string) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string) (newAccessToken string, err error)
}

// service คือ struct ที่ทำงานจริง
type service struct {
	repo Repository
	auth auth.Service
	log  logger.Logger
}

type cursorInfo struct {
	LastID    uint
	LastValue string // สามารถเก็บได้ทั้ง name หรือ created_at
}

// NewExampleUserService คือโรงงานสร้าง Service
func NewExampleUserService(repo Repository, auth auth.Service, log logger.Logger) Service {
	return &service{repo: repo, auth: auth, log: log}
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
	hashedPassword, err := s.auth.HashPassword(plainPassword)
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
	s.log.Dump("Full user object after creation:", userToCreate)
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
func (s *service) ListUsersByPage(limit, offset int, sort, search string) ([]*Domain, int, error) {
	// 1. "แปลภาษาเข็มทิศ" และตรวจสอบความปลอดภัย
	sortField, sortDirection, err := s.parseSortString(sort)
	if err != nil {
		// ถ้า Client ส่ง sort field ที่ไม่ได้รับอนุญาตมา ให้คืนค่า error
		return nil, 0, custom_errors.ValidationError("Sort parameter ไม่ถูกต้อง", err.Error())
	}

	// 2. เรียกใช้ Repository เพื่อดึงข้อมูลและจำนวนทั้งหมด
	userDomains, totalCount, repoErr := s.repo.ListByPage(limit, offset, sortField, sortDirection, search)
	if repoErr != nil {
		return nil, 0, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้", repoErr.Error())
	}

	return userDomains, totalCount, nil
}

// ListUsersByCursor handles cursor-based pagination and sorting.
func (s *service) ListUsersByCursor(cursor string, limit int, sort, search string) (results []*Domain, nextCursor string, hasMore bool, err error) {
	s.log.Debug("Service: ListUsersByCursor started", "cursor", cursor, "limit", limit, "sort", sort, "search", search)

	sortField, sortDirection, parseErr := s.parseSortString(sort)
	if parseErr != nil {
		appErr := custom_errors.ValidationError("Sort parameter ไม่ถูกต้อง", parseErr.Error())
		return nil, "", false, appErr
	}

	cursorData, decodeErr := s.decodeCursor(cursor)
	if decodeErr != nil {
		appErr := custom_errors.ValidationError("Cursor ไม่ถูกต้อง", decodeErr.Error())
		return nil, "", false, appErr
	}

	fetchLimit := limit + 1
	userDomains, repoErr := s.repo.ListByCursor(cursorData, fetchLimit, sortField, sortDirection, search)
	if repoErr != nil {
		appErr := custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้", repoErr.Error())
		return nil, "", false, appErr
	}

	if len(userDomains) > limit {
		hasMore = true
		results = userDomains[:limit]
		lastItemInResults := results[len(results)-1]
		nextCursor = s.encodeCursor(sortField, lastItemInResults)
	} else {
		hasMore = false
		results = userDomains
		nextCursor = ""
	}

	return results, nextCursor, hasMore, nil
}

// UpdateUser handles the business logic for updating a user.
func (s *service) UpdateUser(id uint, req *UpdateUserRequest) (*Domain, error) {
	// 1. ดึงข้อมูลผู้ใช้คนปัจจุบันจาก Repo มาก่อน เพื่อให้แน่ใจว่ามีตัวตนอยู่จริง
	existingUser, err := s.repo.GetByID(id)
	if err != nil {
		// ถ้า repo คืนค่า gorm.ErrRecordNotFound มา เราจะแปลงเป็น NotFoundError ของเรา
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NotFoundError(fmt.Sprintf("ไม่พบผู้ใช้งาน ID: %d", id))
		}
		// ถ้าเป็น error อื่น
		return nil, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ใช้", err.Error())
	}

	// 2. ⭐️ "รวมร่าง" ข้อมูล: เอาข้อมูลใหม่จาก req (ถ้ามี) ไปทับข้อมูลเก่า ⭐️
	// เราจะเช็คว่า pointer ไม่ใช่ nil ก่อนที่จะอัปเดต
	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	// 2. ⭐️⭐️⭐️ เพิ่ม Logic การตรวจสอบข้อมูลซ้ำ! ⭐️⭐️⭐️
	// เช็คว่ามีการส่งอีเมลใหม่มาหรือไม่ และอีเมลนั้นไม่ตรงกับของเดิม
	if req.Email != nil && *req.Email != existingUser.Email {
		// ถ้ามีการเปลี่ยนอีเมล ให้ไปเช็คว่าอีเมลใหม่นี้มีคนอื่นใช้แล้วหรือยัง
		userWithNewEmail, err := s.repo.GetByEmail(*req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการตรวจสอบอีเมลซ้ำ", err.Error())
		}
		// ถ้าเจอ user คนอื่นที่ใช้อีเมลนี้อยู่แล้ว
		if userWithNewEmail != nil {
			return nil, custom_errors.AlreadyExistsError("อีเมลนี้ถูกใช้งานโดยผู้ใช้อื่นแล้ว", nil)
		}
		// ถ้าไม่เจอ ก็อัปเดตอีเมลได้
		existingUser.Email = *req.Email
	}
	if req.Status != nil {
		existingUser.Status = *req.Status
	}

	// 3. สั่งให้ Repo บันทึกข้อมูลที่อัปเดตแล้ว
	if err := s.repo.Update(existingUser); err != nil {
		return nil, custom_errors.SystemErrorWithDetails("ไม่สามารถอัปเดตข้อมูลผู้ใช้ได้", err.Error())
	}

	// 4. คืนค่า Domain object ฉบับล่าสุดกลับไปให้ Handler
	return existingUser, nil
}

// DeleteUser handles the business logic for deleting a user.
func (s *service) DeleteUser(id uint) error {
	// 1. ตรวจสอบให้แน่ใจก่อนว่ามี User ID นี้อยู่จริงหรือไม่
	// เพื่อที่เราจะสามารถคืนค่า 404 Not Found ที่ถูกต้องกลับไปได้
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_errors.NotFoundError(fmt.Sprintf("ไม่พบผู้ใช้งาน ID: %d ที่ต้องการลบ", id))
		}
		return custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ใช้ก่อนลบ", err.Error())
	}

	// 2. ถ้ามีอยู่จริง ก็สั่งให้ Repo ทำการลบ
	if err := s.repo.Delete(id); err != nil {
		return custom_errors.SystemErrorWithDetails("ไม่สามารถลบข้อมูลผู้ใช้ได้", err.Error())
	}

	// 3. ถ้าสำเร็จ ก็ไม่ต้องคืนค่าอะไรกลับไป (return nil)
	return nil
}

// ChangePassword handles the business logic for changing a user's password.
func (s *service) ChangePassword(id uint, oldPassword, newPassword string) error {
	// 1. ดึงข้อมูลผู้ใช้คนปัจจุบันจาก Repo มาก่อน
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_errors.NotFoundError(fmt.Sprintf("ไม่พบผู้ใช้งาน ID: %d", id))
		}
		return custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ใช้", err.Error())
	}

	// 2. ⭐️ "ตรวจสอบกุญแจเก่า": เปรียบเทียบรหัสผ่านเก่าที่ส่งมากับ Hash ใน DB ⭐️
	// เราจะเรียกใช้ "บอดี้การ์ด" (pkg/auth) ของเรามาช่วย
	if err := s.auth.ComparePassword(user.PasswordHash, oldPassword); err != nil {
		// ถ้า err ไม่ใช่ nil แสดงว่ารหัสผ่านเก่าไม่ถูกต้อง!
		return custom_errors.UnauthorizedError("รหัสผ่านเดิมไม่ถูกต้อง")
	}

	// 3. "สร้างกุญแจใหม่": Hash รหัสผ่านใหม่
	newHashedPassword, err := s.auth.HashPassword(newPassword)
	if err != nil {
		return custom_errors.SystemErrorWithDetails("ไม่สามารถเข้ารหัสรหัสผ่านใหม่ได้", err.Error())
	}

	// 4. อัปเดต "กุญแจ" ใน object ของเรา
	user.PasswordHash = newHashedPassword

	// 5. สั่งให้ Repo บันทึกการเปลี่ยนแปลง
	if err := s.repo.Update(user); err != nil {
		return custom_errors.SystemErrorWithDetails("ไม่สามารถบันทึกรหัสผ่านใหม่ได้", err.Error())
	}

	// 6. ถ้าสำเร็จ ก็ไม่ต้องคืนค่าอะไรกลับไป (return nil)
	return nil
}

// Login handles the business logic for user authentication.
func (s *service) Login(email, password string) (accessToken string, refreshToken string, err error) {
	// 1. ค้นหาผู้ใช้ด้วยอีเมล
	user, repoErr := s.repo.GetByEmail(email)
	if repoErr != nil {
		if !errors.Is(repoErr, gorm.ErrRecordNotFound) {
			return "", "", custom_errors.SystemErrorWithDetails("เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ใช้", repoErr.Error())
		}
		return "", "", custom_errors.UnauthorizedError("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	// 2. "ตรวจสอบตัวตน"
	if err := s.auth.ComparePassword(user.PasswordHash, password); err != nil {
		return "", "", custom_errors.UnauthorizedError("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	// 3. ✨ "สร้างกุญแจสองชั้น"! ✨
	accessToken, refreshToken, tokenErr := s.auth.GenerateTokens(user.ID, user.Role)
	if tokenErr != nil {
		return "", "", custom_errors.SystemErrorWithDetails("ไม่สามารถสร้าง Token ได้", tokenErr.Error())
	}

	// 4. คืนค่า "กุญแจ" ทั้งสองดอกกลับไปให้ Handler
	return accessToken, refreshToken, nil
}

// RefreshToken handles the logic for refreshing an access token.
func (s *service) RefreshToken(refreshToken string) (newAccessToken string, err error) {
	// 1. ส่ง "บัตรสมาชิก" ไปให้ "บริษัทรักษาความปลอดภัย" ตรวจสอบ
	// ถ้าบัตรถูกต้อง จะได้ UserID กลับมา
	userID, err := s.auth.ValidateRefreshToken(refreshToken)
	if err != nil {

		return "", custom_errors.SystemErrorWithDetails("ไม่สามารถตรวจสอบ Refresh Token ได้", err.Error())
	}

	// 2. (Pro-Tip!) ตรวจสอบอีกชั้น: เช็คว่า User คนนี้ยังมีตัวตนอยู่ในระบบจริงๆ หรือไม่
	user, err := s.repo.GetByID(userID)
	if err != nil {
		// ถ้าหาไม่เจอ (อาจจะถูกลบไปแล้ว) ก็ไม่ควรออกบัตรใหม่ให้
		return "", custom_errors.UnauthorizedError("User not found")
	}

	// 3. ถ้าทุกอย่างถูกต้อง... ให้ออก "บัตรผ่านรายวัน" (Access Token) ใบใหม่!
	// สังเกตว่าเราจะใช้ GenerateToken ตัวเดิมที่สร้างแค่ Access Token ได้เลย
	_, newAccessToken, err = s.auth.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return "", custom_errors.SystemErrorWithDetails("ไม่สามารถสร้าง Access Token ใหม่ได้", err.Error())
	}

	// 4. คืนค่า "บัตรผ่านรายวัน" ใบใหม่กลับไป
	return newAccessToken, nil
}

// --- Private Helpers ---

func (s *service) parseSortString(sort string) (field string, direction string, err error) {
	allowedSortFields := map[string]bool{
		"id": true, "name": true, "email": true, "created_at": true, "updated_at": true,
	}
	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid sort format")
	}
	field, direction = parts[0], parts[1]
	if !allowedSortFields[field] {
		return "", "", errors.New("sorting by this field is not allowed: " + field)
	}
	if direction != "asc" && direction != "desc" {
		return "", "", errors.New("invalid sort direction")
	}
	return field, direction, nil
}

// ✨ อัปเกรด "เครื่องมือ" สร้างและอ่านที่คั่นหนังสือ ✨
func (s *service) encodeCursor(sortField string, lastItem *Domain) string {
	var valueToEncode string
	switch sortField {
	case "name":
		valueToEncode = fmt.Sprintf("%s,%d", lastItem.Name, lastItem.ID)
	case "created_at":
		valueToEncode = fmt.Sprintf("%s,%d", lastItem.CreatedAt.Format(time.RFC3339Nano), lastItem.ID)
	default: // id
		valueToEncode = fmt.Sprintf("%d", lastItem.ID)
	}
	return base64.StdEncoding.EncodeToString([]byte(valueToEncode))
}

func (s *service) decodeCursor(encodedCursor string) (*cursorInfo, error) {
	if encodedCursor == "" {
		return &cursorInfo{LastID: 0, LastValue: ""}, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, errors.New("invalid base64 format")
	}
	parts := strings.Split(string(decoded), ",")

	if len(parts) == 1 { // id only
		id, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid cursor id content")
		}
		return &cursorInfo{LastID: uint(id)}, nil
	}
	if len(parts) == 2 { // value + id
		id, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid cursor id content")
		}
		return &cursorInfo{LastValue: parts[0], LastID: uint(id)}, nil
	}

	return nil, errors.New("invalid cursor format")
}
