// pkg/validator/validator.go
package validator

import (
	"fmt"
	"reflect"
	"regexp" // 1. Import regexp เข้ามาเพื่อสร้างกฎของเราเอง
	"strings"

	"github.com/go-playground/validator/v10"
)

// ====================================================================================
// Structs for Validation Result
// ====================================================================================

// ValidationResult holds the complete validation result.
type ValidationResult struct {
	IsValid bool                    `json:"is_valid"`
	Errors  []ValidationErrorDetail `json:"errors,omitempty"`
}

// ValidationErrorDetail represents a single validation error.
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

// ====================================================================================
// Validator Factory & Custom Rules
// ====================================================================================

// New creates and configures a new validator instance.
// เราเปลี่ยนจาก var global มาเป็นฟังก์ชัน New() เพื่อให้เรา "ลงทะเบียน" กฎใหม่ๆ ได้
func New() *validator.Validate {
	v := validator.New()

	// ⭐️ ลงทะเบียน "กฎ" ใหม่ที่เราสร้างขึ้นเองที่นี่! ⭐️
	v.RegisterValidation("sort_format", validateSortFormat)

	return v
}

// validateSortFormat คือฟังก์ชันที่ทำงานเบื้องหลังของกฎ "sort_format"
// มันจะตรวจสอบว่า string อยู่ในรูปแบบ 'field:direction' หรือไม่
func validateSortFormat(fl validator.FieldLevel) bool {
	sortPattern := `^[a-zA-Z_]+:(asc|desc)$`
	matched, _ := regexp.MatchString(sortPattern, fl.Field().String())
	return matched
}

// ====================================================================================
// Main Validation Function
// ====================================================================================

// Validate validates a struct and returns our custom, detailed result.
// สังเกตว่าตอนนี้มันรับ validator instance (v) เข้ามาด้วย
func Validate(v *validator.Validate, s interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationErrorDetail{},
	}

	err := v.Struct(s)
	if err != nil {
		result.IsValid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			result.Errors = translateValidationErrors(s, validationErrors)
		}
		return result
	}

	return result
}

// ====================================================================================
// Error Translation & Parsing Helpers
// ====================================================================================

// translateValidationErrors converts validator.ValidationErrors to our custom format.
func translateValidationErrors(s interface{}, validationErrors validator.ValidationErrors) []ValidationErrorDetail {
	var details []ValidationErrorDetail
	structType := reflect.TypeOf(s)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	for _, ve := range validationErrors {
		var fieldName string
		var customMessage string

		if field, ok := structType.FieldByName(ve.Field()); ok {
			// Get field name from json tag
			jsonTag := field.Tag.Get("json")
			fieldName = strings.Split(jsonTag, ",")[0]
			if fieldName == "" {
				fieldName = ve.Field()
			}

			// Parse the vmsg tag to find a specific message for the failed rule
			vmsgTag := field.Tag.Get("vmsg")
			messageMap := parseCommaSeparatedVmsg(vmsgTag)
			customMessage = messageMap[ve.Tag()]

		} else {
			fieldName = ve.Field()
		}

		detail := ValidationErrorDetail{
			Field:   fieldName,
			Message: ternary(customMessage != "", customMessage, generateDefaultErrorMessage(ve.Tag(), ve.Param(), ve)),
			Value:   fmt.Sprintf("%v", ve.Value()),
		}
		details = append(details, detail)
	}
	return details
}

// parseCommaSeparatedVmsg is our smart parser for the vmsg tag.
func parseCommaSeparatedVmsg(tag string) map[string]string {
	// ... (โค้ดส่วนนี้เหมือนเดิมเป๊ะๆ) ...
	messageMap := make(map[string]string)
	if tag == "" {
		return messageMap
	}
	var parts []string
	var current strings.Builder
	escaped := false
	for _, char := range tag {
		if char == '\\' && !escaped {
			escaped = true
			continue
		}
		if char == ',' && !escaped {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(char)
		}
		escaped = false
	}
	parts = append(parts, current.String())
	for _, rule := range parts {
		keyValue := strings.SplitN(rule, ":", 2)
		if len(keyValue) == 2 {
			ruleName := strings.TrimSpace(keyValue[0])
			message := strings.ReplaceAll(strings.TrimSpace(keyValue[1]), `\,`, `,`)
			messageMap[ruleName] = message
		}
	}
	return messageMap
}

// generateDefaultErrorMessage creates a user-friendly error message if no custom message is provided.
func generateDefaultErrorMessage(tag, param string, originalError error) string {
	switch tag {
	// ... (case อื่นๆ เหมือนเดิม) ...
	case "email":
		return "ต้องเป็นรูปแบบอีเมลที่ถูกต้อง"

	// ⭐️ เพิ่ม Default Message สำหรับกฎใหม่ของเรา ⭐️
	case "sort_format":
		return "รูปแบบการเรียงข้อมูลต้องเป็น 'field:direction' (เช่น id:asc)"

	// --- Fallback ---
	default:
		return originalError.Error()
	}
}

// ternary is a small helper for conditional expressions.
func ternary(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	}
	return ifFalse
}
