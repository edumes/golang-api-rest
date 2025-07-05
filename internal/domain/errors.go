package domain

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Err     error       `json:"-"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewAppErrorWithDetails(code int, message string, err error, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: details,
	}
}

func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func NewConflictError(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, err)
}

func NewValidationAppError(message string, err error, details interface{}) *AppError {
	return NewAppErrorWithDetails(http.StatusBadRequest, message, err, details)
}

func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %d errors", len(v))
}

func NewValidationError(field, tag, value, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Tag:     tag,
		Value:   value,
		Message: message,
	}
}
