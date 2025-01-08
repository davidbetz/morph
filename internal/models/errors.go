package models

import "fmt"

type NotImplementedError struct {
	Message string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("NotImplementedError: %s", e.Message)
}

func NewNotImplementedError(message string) *NotImplementedError {
	return &NotImplementedError{
		Message: message,
	}
}
