// pkg/auth/auth.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service คือ "บริษัทรักษาความปลอดภัย" ของเรา
// เป็น Interface ที่กำหนดว่าบริษัทนี้ต้องทำอะไรได้บ้าง
type Service interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, plainPassword string) error
	GenerateToken(userID uint, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// JWTClaims คือ "ข้อมูล" ที่เราจะฝังเข้าไปใน "กุญแจดิจิทัล" (Token)
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// authService คือ struct ที่ทำงานจริงๆ และเก็บ "รหัสลับ" ไว้ข้างใน
type authService struct {
	secretKey []byte
}

// NewAuthService คือ "โรงงาน" สำหรับก่อตั้งบริษัทรักษาความปลอดภัย
// เราจะเรียกใช้ฟังก์ชันนี้แค่ครั้งเดียวใน main.go
func NewAuthService(secretKey string) Service {
	return &authService{
		secretKey: []byte(secretKey),
	}
}

// --- Implementation of Service interface ---

// HashPassword เข้ารหัสรหัสผ่านด้วย bcrypt
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword เปรียบเทียบรหัสผ่านกับ hash
func (s *authService) ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// GenerateToken สร้าง JWT Token ใหม่
func (s *authService) GenerateToken(userID uint, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Token มีอายุ 7 วัน
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken ตรวจสอบความถูกต้องของ JWT Token
func (s *authService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
