package validator

import (
	"regexp"
	"strings"
)

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
	// At least 8 characters, one uppercase, one lowercase, one digit
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

// IsValidName validates name format
func IsValidName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return false
	}

	// Only letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s\-']+$`)
	return nameRegex.MatchString(name)
}

// IsValidPhoneNumber validates phone number format
func IsValidPhoneNumber(phone string) bool {
	// Remove all non-digit characters
	phoneDigits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Check if it's between 10-15 digits
	return len(phoneDigits) >= 10 && len(phoneDigits) <= 15
}

// SanitizeString removes harmful characters from string
func SanitizeString(input string) string {
	// Remove control characters and trim whitespace
	cleaned := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}
