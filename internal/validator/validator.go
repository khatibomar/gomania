package validator

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error for a specific field.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of ValidationError.
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors.
func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range ve {
		sb.WriteString(err.Error())
		if i < len(ve)-1 {
			sb.WriteString("; ")
		}
	}
	return sb.String()
}

// IsValidationError checks if the error is of type ValidationErrors.
func IsValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}
