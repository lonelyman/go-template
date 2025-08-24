// pkg/auth/auth.go
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain password and returns its bcrypt hash.
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost คือความซับซ้อนในการ Hash (เลขยิ่งเยอะ ยิ่งช้า ยิ่งปลอดภัย)
	// DefaultCost (10) คือค่าที่เหมาะสมและปลอดภัยในปัจจุบัน
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain password.
// It returns nil on success, or an error if they don't match.
func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
