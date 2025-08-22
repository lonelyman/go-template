package utils

import (
	"fmt"
	"strings"
	"time"
)

// FormatTime formats time to Thailand timezone
func FormatTime(t time.Time) string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return t.In(loc).Format("2006-01-02 15:04:05")
}

// ParseThaiTime parses time string in Thailand timezone
func ParseThaiTime(timeStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
}

// TruncateString truncates string to specified length
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// SlugifyString converts string to URL-friendly slug
func SlugifyString(str string) string {
	// Convert to lowercase
	str = strings.ToLower(str)

	// Replace spaces and special characters with hyphens
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.ReplaceAll(str, "_", "-")

	// Remove multiple consecutive hyphens
	for strings.Contains(str, "--") {
		str = strings.ReplaceAll(str, "--", "-")
	}

	// Trim hyphens from start and end
	str = strings.Trim(str, "-")

	return str
}

// GenerateID generates a simple ID from timestamp
func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// IsEmpty checks if string is empty or only whitespace
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// Contains checks if slice contains item (generic function would be better in Go 1.18+)
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveString removes item from slice
func RemoveString(slice []string, item string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
