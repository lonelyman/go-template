คู่มือการใช้งาน pkg/logger (สุดยอดนักข่าว)
แพ็กเกจ logger คือ "สำนักข่าว" กลางของแอปพลิชันเรา ถูกออกแบบมาให้มีความยืดหยุ่น, ใช้งานง่าย, และทรงพลัง โดยสามารถ "เปลี่ยนโฉม" การแสดงผลได้ตามสภาพแวดล้อม (Environment)

ปรัชญาหลัก:

ตอนพัฒนา (Development): แสดงผล Log ที่สวยงาม, มีสีสัน, และอ่านง่ายสำหรับมนุษย์ (Pretty Logger)

ตอนขึ้นโปรดักชัน (Production): แสดงผล Log ที่มีโครงสร้างชัดเจนเป็น JSON (Structured Logger หรือ slog) เพื่อให้เครื่องจักรสามารถนำไปวิเคราะห์ต่อได้

1. การตั้งค่า (Configuration)
"หัวหน้ากองบรรณาธิการ" (main.go) จะเป็นคนตัดสินใจเลือก "นักข่าว" ที่จะใช้งานโดยอัตโนมัติ โดยดูจากค่า Config server.mode ซึ่งถูกควบคุมโดย Environment Variable SERVER_MODE

ถ้า SERVER_MODE = "development" -> จะใช้ Pretty Logger

ถ้า SERVER_MODE = "production" (หรือค่าอื่นๆ) -> จะใช้ Slog Logger

2. การใช้งาน (Usage)
เราจะ ไม่ สร้าง Logger ขึ้นมาเองในแต่ละส่วนของโปรแกรม แต่เราจะใช้หลักการ Dependency Injection (DI) โดยให้ main.go เป็นคน "แจกจ่าย" logger.Logger instance ไปให้ทุกส่วนที่ต้องการใช้งาน (เช่น Repository, Service, Handler)

ตัวอย่างการติดตั้งใน Service:

// internal/modules/example_user/example_user_service.go

type service struct {
	repo Repository
	log  logger.Logger // 1. เพิ่มช่องสำหรับเก็บ "นักข่าว"
}

// 2. รับ "นักข่าว" เข้ามาตอนสร้าง
func NewExampleUserService(repo Repository, log logger.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

// 3. ใช้งาน "นักข่าว"
func (s *service) SomeMethod() {
    s.log.Info("Starting some method", "userID", 123)
    // ...
}

3. ความสามารถของนักข่าว (Logging Methods)
"นักข่าว" ของเรามีความสามารถ 3 รูปแบบหลักๆ

3.1 Structured Logging (สำหรับ Production และ Log ที่ต้องการรายละเอียด)
นี่คือวิธีหลักในการบันทึก Log เราจะให้ "พาดหัวข่าว" (msg) และตามด้วย "รายละเอียดข่าว" ที่เป็นคู่ key-value เสมอ

log.Debug(msg string, args ...any): สำหรับข้อมูลดีบักที่ละเอียดมากๆ

log.Info(msg string, args ...any): สำหรับเหตุการณ์ทั่วไป

log.Success(msg string, args ...any): สำหรับเหตุการณ์ที่บ่งบอกถึงความสำเร็จ

log.Warn(msg string, args ...any): สำหรับเหตุการณ์ที่น่าสงสัย แต่ยังไม่ถึงกับเป็น Error

log.Error(msg string, err error, args ...any): สำหรับข้อผิดพลาด (จะแนบ err เข้าไปใน Log โดยอัตโนมัติ)

ตัวอย่างการใช้งาน:

s.log.Info("User logged in successfully", "userID", 123, "ip_address", "127.0.0.1")
s.log.Error("Failed to connect to database", err, "retry_count", 3)

3.2 Simple Dumping (สำหรับ Debug ง่ายๆ ตอนพัฒนา)
นี่คือ "เครื่องมือพิเศษ" สำหรับตอนที่เราแค่อยากจะดูค่าในตัวแปรเร็วๆ โดยไม่ต้องมานั่งพิมพ์ key-value

log.Dump(a ...any): จะทำการ Dump ตัวแปรทั้งหมดที่ส่งเข้าไป ออกมาเป็น JSON สวยๆ

ตัวอย่างการใช้งาน:

userDomain, _ := s.repo.GetByID(1)
s.log.Dump(userDomain)

3.3 ✨ Highlighting (ปากกาไฮไลท์สำหรับ Debug)
นี่คือฟังก์ชันพิเศษสำหรับตอนที่เราต้องการ Debug Flow ที่ซับซ้อน และอยากจะ "ไฮไลท์" Log ของเราด้วยสีที่แตกต่างกัน

log.Highlight(color string, msg string, data ...any): จะทำการ Dump ข้อมูลเหมือน Dump แต่จะใช้สีที่เรากำหนด

ตัวอย่างการใช้งาน:

// เราสามารถเรียกใช้ Color Constants ที่อยู่ใน package logger ได้เลย
s.log.Highlight(logger.ColorPurple, "User data before update:", userBefore)
s.log.Highlight(logger.ColorCyan, "Update request payload:", updateRequest)

ผลลัพธ์ (Development - Pretty Logger):

🎨 HIGHLIGHT   service.go:80: User data before update:
{
  "ID": 1,
  "Name": "Nipon K.",
  ...
}
🎨 HIGHLIGHT   service.go:81: Update request payload:
{
  "Name": "Nipon Kaew",
  ...
}

ใน Production, Highlight() จะถูกแปลงเป็น Debug Level เพื่อความปลอดภัย