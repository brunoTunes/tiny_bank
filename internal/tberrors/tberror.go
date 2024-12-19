package tberrors

import "fmt"

type ValidationError struct {
	message string
	field   string
}

func NewValidationError(message, field string) error {
	return ValidationError{
		message: message,
		field:   field,
	}
}

func (validationError ValidationError) Error() string {
	return fmt.Sprintf("%s\n", validationError.message)
}
