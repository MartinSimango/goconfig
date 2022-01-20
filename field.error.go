package goconfig

import (
	"fmt"
	"strings"
)

type Error interface {
	error
}

type FieldError struct {
	Field        string
	ErrorMessage string
}

// Check we implement interface
var _ Error = &FieldError{}

func (cfe *FieldError) Error() string {
	return fmt.Sprintf("error found with field '%s': %s", cfe.Field, cfe.ErrorMessage)
}

//FieldErrorsToString coverts an array of field arrays into a string
func FieldErrorsToString(fieldErrors []FieldError) string {
	var combinedFieldErrors []string
	for _, fieldError := range fieldErrors {
		combinedFieldErrors = append(combinedFieldErrors, fmt.Sprintf("%s - '%s'",
			fieldError.Field, fieldError.ErrorMessage))
	}
	return strings.Join(combinedFieldErrors, "\n")
}

// ToFieldError converts err into a field error
func ToFieldError(field string, err error) *FieldError {
	if err == nil {
		return nil
	}
	return &FieldError{
		Field:        field,
		ErrorMessage: err.Error(),
	}
}

// AppendFieldError appends fieldError to the configFieldErrors are and returns the resulting slice.
func AppendFieldError(configFieldErrors []FieldError, fieldError *FieldError) []FieldError {
	if fieldError != nil {
		return append(configFieldErrors, *fieldError)
	}
	return configFieldErrors
}
