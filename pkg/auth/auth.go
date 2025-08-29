// pkg/auth/auth.go
package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service คือ "บริษัทรักษาความปลอดภัย" ของเรา
type Service interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, plainPassword string) error
	// ✨ อัปเกรด: เปลี่ยน GenerateToken เป็น GenerateTokens
	GenerateTokens(userID uint, role string) (accessToken string, refreshToken string, err error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	// ✨ เพิ่ม: ความสามารถในการตรวจสอบ Refresh Token
	ValidateRefreshToken(tokenString string) (userID uint, err error)
}

// JWTClaims คือ "ข้อมูล" ที่เราจะฝังเข้าไปใน "Access Token"
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// authService คือ struct ที่ทำงานจริง
type authService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewAuthService คือ "โรงงาน" สำหรับก่อตั้งบริษัทรักษาความปลอดภัย
func NewAuthService(secretKey string) Service {
	return &authService{
		secretKey:     []byte(secretKey),
		accessExpiry:  30 * time.Minute,    // 👈 กำหนดอายุ Access Token ที่นี่ (30 นาที)
		refreshExpiry: 15 * 24 * time.Hour, // 👈 กำหนดอายุ Refresh Token ที่นี่ (15 วัน)
	}
}

// --- Implementation of Service interface ---

func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *authService) ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// ✨ อัปเกรด: GenerateTokens จะสร้างกุญแจ 2 ดอกพร้อมกัน
func (s *authService) GenerateTokens(userID uint, role string) (accessToken string, refreshToken string, err error) {
	// 1. สร้าง Access Token (บัตรผ่านรายวัน - อายุสั้น)
	accessClaims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secretKey)
	if err != nil {
		return "", "", err
	}

	// 2. สร้าง Refresh Token (บัตรสมาชิก - อายุยาว)
	refreshClaims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID), // ใช้ Subject เก็บ UserID
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken ตรวจสอบความถูกต้องของ "Access Token"
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

// ✨ เพิ่ม: ValidateRefreshToken ตรวจสอบความถูกต้องของ "Refresh Token"
func (s *authService) ValidateRefreshToken(tokenString string) (userID uint, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ดึง UserID ที่เก็บไว้ใน 'sub' (Subject) claim
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return 0, fmt.Errorf("invalid user id in token")
		}

		userID_uint, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user id format in token")
		}
		return uint(userID_uint), nil
	}

	return 0, fmt.Errorf("invalid refresh token")
}
