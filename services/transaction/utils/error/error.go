package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	// TypeInternal indicates an unexpected server-side error.
	TypeInternal ErrorType = "internal_server_error"
	// TypeInvalidInput indicates bad request data from the client.
	TypeInvalidInput ErrorType = "invalid_input"
	// TypeNotFound indicates a resource was not found.
	TypeNotFound ErrorType = "not_found"
	// TypeUnauthorized indicates authentication failure.
	TypeUnauthorized ErrorType = "unauthorized"
	// TypeForbidden indicates authorization failure.
	TypeForbidden ErrorType = "forbidden"
	// TypeConflict indicates a resource conflict (e.g., duplicate entry).
	TypeConflict ErrorType = "conflict"
	// TypeExternal indicates an error from an external service.
	TypeExternal ErrorType = "external_service_error"
)

type CustomError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the wrapped error, allowing for error chain inspection.
func (e *CustomError) Unwrap() error {
	return e.Err
}

// New creates a new CustomError with the given type and message.
func New(errType ErrorType, message string) error {
	return &CustomError{Type: errType, Message: message}
}

// Wrap creates a new CustomError wrapping an existing error.
func Wrap(errType ErrorType, message string, err error) error {
	return &CustomError{Type: errType, Message: message, Err: err}
}

// GetType attempts to extract the CustomErrorType from an error chain.
func GetType(err error) ErrorType {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type
	}
	return TypeInternal
}

// GetHTTPStatu.s maps a CustomErrorType to an HTTP status code
func GetHTTPStatus(errType ErrorType) int {
	switch errType {
	case TypeInvalidInput:
		return http.StatusBadRequest
	case TypeNotFound:
		return http.StatusNotFound
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeForbidden:
		return http.StatusForbidden
	case TypeConflict:
		return http.StatusConflict
	case TypeExternal:
		return http.StatusBadGateway
	case TypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
