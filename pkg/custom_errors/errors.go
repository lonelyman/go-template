package custom_errors

// ⭐️ 1. เอา Error Code ที่เป็นมาตรฐานของน้องมาไว้ที่นี่ ⭐️
const (
	// Authentication
	ErrUnauthorized     = "UNAUTHORIZED"
	ErrInvalidToken     = "INVALID_TOKEN"
	ErrTokenExpired     = "TOKEN_EXPIRED"
	ErrPermissionDenied = "PERMISSION_DENIED"

	// Validation
	ErrValidation    = "VALIDATION_ERROR"
	ErrMissingParam  = "MISSING_PARAMETER"
	ErrInvalidFormat = "INVALID_FORMAT"

	// Resource
	ErrNotFound      = "NOT_FOUND"
	ErrAlreadyExists = "ALREADY_EXISTS"

	// System
	ErrSystem      = "SYSTEM_ERROR"
	ErrExternalAPI = "EXTERNAL_API_ERROR"
)

// ⭐️ 2. เอา SimpleError struct ของน้องมาใช้เป็น AppError ของเรา ⭐️
type AppError struct {
	HTTPStatus int         `json:"-"` // เราจะเพิ่ม HTTPStatus เข้าไปเพื่อให้จัดการง่ายขึ้น และซ่อนจาก JSON
	Message    string      `json:"message"`
	Code       string      `json:"code"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// ⭐️ 3. ย้าย Helper Functions ทั้งหมดของน้องมาไว้ที่นี่! ⭐️
// ทำให้การสร้าง Error สวยงามและเป็นมาตรฐานเดียวกันทั้งโปรเจกต์

// --- Helper Functions ---

func New(status int, code, message string) *AppError {
	// (อาจจะเพิ่ม checker เหมือนของน้องก็ได้)
	return &AppError{HTTPStatus: status, Code: code, Message: message}
}

func NewWithDetails(status int, code, message string, details interface{}) *AppError {
	return &AppError{HTTPStatus: status, Code: code, Message: message, Details: details}
}

// Helper ที่ใช้งานบ่อยๆ
func NotFoundError(message string) *AppError {
	return New(404, ErrNotFound, message)
}

func ValidationError(message string, details interface{}) *AppError {
	return NewWithDetails(400, ErrValidation, message, details)
}

func UnauthorizedError(message string) *AppError {
	return New(401, ErrUnauthorized, message)
}

func SystemError(message string) *AppError {
	return New(500, ErrSystem, message)
}
