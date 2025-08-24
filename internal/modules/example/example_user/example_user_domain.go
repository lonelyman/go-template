package example_user

import "time"

// Domain คือพิมพ์เขียวหลักของข้อมูล User ในระบบของเรา
// จะต้องบริสุทธิ์ ไม่มี gorm tags หรือ json tags
type Domain struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Status       string
	Role         string
	LastLoginAt  *time.Time
}
