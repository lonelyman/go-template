package logger

// Logger คือ "สัญญาใจ" หรือ Interface ที่นักข่าวทุกคนต้องทำตาม
// ไม่ว่าจะเป็นนักข่าวสายสวยงาม หรือสายโปรดักชัน ก็ต้องมี 4 ความสามารถนี้
type Logger interface {
	// Debug ใช้สำหรับ Log ที่มีรายละเอียดเยอะๆ สำหรับตอนพัฒนาเท่านั้น
	Debug(msg string, args ...any)

	// Info ใช้สำหรับ Log เหตุการณ์ทั่วไปที่เกิดขึ้น
	Info(msg string, args ...any)

	// Warn ใช้สำหรับ Log เหตุการณ์ที่น่าสงสัย แต่ยังไม่ถึงกับเป็น Error
	Warn(msg string, args ...any)

	// Error ใช้สำหรับ Log ข้อผิดพลาดที่เกิดขึ้นในระบบ
	Error(msg string, err error, args ...any)
	// Success ใช้สำหรับ Log ข้อความที่บ่งบอกว่าการทำงานสำเร็จ
	Success(msg string, args ...any)

	Print(msg string)
	
	Dump(data interface{})

	Dumpf(level string, msg string, data interface{})
}

const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarn    = "WARN"
	LevelError   = "ERROR"
	LevelSuccess = "SUCCESS"
)
