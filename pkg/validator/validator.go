package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldValidationError represents a single validation error with enhanced details
type FieldValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorMessage string `json:"error_message"`
	ActualValue  string `json:"actual_value"`
	Constraint   string `json:"constraint,omitempty"`
}

// StructValidationResult holds the complete validation result
type StructValidationResult struct {
	IsValid bool                   `json:"is_valid"`
	Errors  []FieldValidationError `json:"errors,omitempty"`
}

// Enhanced validator instance
var structValidator = validator.New()

// ValidateStructWithDetails validates a struct and returns detailed error information
// with support for custom field names and error messages via struct tags
func ValidateStructWithDetails(structToValidate interface{}) *StructValidationResult {
	result := &StructValidationResult{
		IsValid: true,
		Errors:  []FieldValidationError{},
	}

	// Validate required pointer fields first (common with query parameter parsing)
	requiredPointerErrors := validateRequiredPointerFields(structToValidate)
	if len(requiredPointerErrors) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, requiredPointerErrors...)
		return result
	}

	// Perform standard struct validation
	err := structValidator.Struct(structToValidate)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := extractFieldValidationErrors(structToValidate, validationErrors)
			if len(fieldErrors) > 0 {
				result.IsValid = false
				result.Errors = fieldErrors
			}
		}
	}

	return result
}

// validateRequiredPointerFields checks for nil pointer fields that are marked as required
// This is particularly useful when parsing query parameters into pointer fields
func validateRequiredPointerFields(structToValidate interface{}) []FieldValidationError {
	var errors []FieldValidationError

	structValue := reflect.ValueOf(structToValidate)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	if structValue.Kind() != reflect.Struct {
		return errors
	}

	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Check if field is a pointer and has required validation
		if field.Kind() == reflect.Ptr && isFieldRequired(fieldType) && field.IsNil() {
			fieldName := extractFieldName(fieldType)
			errorMessage := extractRequiredErrorMessage(fieldType, fieldName)

			error := FieldValidationError{
				FieldName:    fieldName,
				ErrorMessage: errorMessage,
				ActualValue:  "<nil>",
				Constraint:   "required",
			}
			errors = append(errors, error)
		}
	}

	return errors
}

// extractFieldValidationErrors converts validator.ValidationErrors to our custom error format
func extractFieldValidationErrors(structToValidate interface{}, validationErrors validator.ValidationErrors) []FieldValidationError {
	var errors []FieldValidationError

	structType := reflect.TypeOf(structToValidate)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	for _, ve := range validationErrors {
		fieldType, found := structType.FieldByName(ve.Field())

		var fieldName, errorMessage string
		if found {
			fieldName = extractFieldName(fieldType)
			errorMessage = extractTagErrorMessage(fieldType, ve.Tag())
		} else {
			// Fallback if field is not found
			fieldName = ve.Field()
			errorMessage = generateDefaultErrorMessage(ve.Tag())
		}

		error := FieldValidationError{
			FieldName:    fieldName,
			ErrorMessage: errorMessage,
			ActualValue:  fmt.Sprintf("%v", ve.Value()),
			Constraint:   ve.Param(),
		}
		errors = append(errors, error)
	}

	return errors
}

// Helper functions for better code organization

// isFieldRequired checks if a struct field has required validation tag
func isFieldRequired(fieldType reflect.StructField) bool {
	validateTag := fieldType.Tag.Get("validate")
	return strings.Contains(validateTag, "required")
}

// extractFieldName gets the field name from struct tag or uses the struct field name as fallback
func extractFieldName(fieldType reflect.StructField) string {
	if fieldName := fieldType.Tag.Get("field_name"); fieldName != "" {
		return fieldName
	}
	return fieldType.Name
}

// extractRequiredErrorMessage gets the custom required error message or generates a default one
func extractRequiredErrorMessage(fieldType reflect.StructField, fieldName string) string {
	if errorMessage := fieldType.Tag.Get("error_required"); errorMessage != "" {
		return errorMessage
	}
	return fmt.Sprintf("%s is required", fieldName)
}

// extractTagErrorMessage gets the custom error message for a specific validation tag
func extractTagErrorMessage(fieldType reflect.StructField, tag string) string {
	errorTagKey := "error_" + tag
	if errorMessage := fieldType.Tag.Get(errorTagKey); errorMessage != "" {
		return errorMessage
	}
	return generateDefaultErrorMessage(tag)
}

// generateDefaultErrorMessage creates a user-friendly error message for validation tags
func generateDefaultErrorMessage(tag string) string {
	switch tag {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "value is too small"
	case "max":
		return "value is too large"
	case "len":
		return "invalid length"
	case "numeric":
		return "must be numeric"
	case "alpha":
		return "must contain only alphabetic characters"
	case "alphanum":
		return "must contain only alphanumeric characters"
	default:
		return fmt.Sprintf("validation failed for rule: %s", tag)
	}
}

// Legacy function for backward compatibility
// Deprecated: Use ValidateStructWithDetails instead
func ValidateStructDetails(s interface{}) ([]ValidationDetail, error) {
	result := ValidateStructWithDetails(s)

	// Convert to old format
	var details []ValidationDetail
	for _, err := range result.Errors {
		detail := ValidationDetail{
			Field: err.FieldName,
			Tag:   err.ErrorMessage,
			Value: err.ActualValue,
			Param: err.Constraint,
		}
		details = append(details, detail)
	}

	if !result.IsValid {
		return details, nil
	}

	return nil, nil
}

// ValidationDetail - legacy type for backward compatibility
// Deprecated: Use FieldValidationError instead
type ValidationDetail struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
	Param string `json:"param,omitempty"`
}
