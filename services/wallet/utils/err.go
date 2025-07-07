// Package errors contains the error handler controller
//
// Handles internal errors (service, external api calls, etc)
package utils

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

// Error defines an application (domain) error, encapsulating error-related information.
type Error struct {
	// Message contains a human-readable description of the error.
	Message string

	// Cause represents the underlying cause of the error, if any. It can be another error.
	Cause error

	// Type is an optional field used to categorize the error type within the application.
	Type *Type

	// Stacktrace contains the stack trace associated with the error, typically used for debugging purposes.
	// It can be a string or an array of strings representing the call stack.
	Stacktrace string // It can be an array of strings later if needed.
}

// Type is a struct that can be used to categorize error types in a Go application.
// It provides fields to help identify and classify different types of errors that
// may occur during the execution of the program.
type Type struct {
	// Category is a string that represents the general category or group of the error.
	// For example, it can be "Validation", "Database", "API", or any other relevant category.
	Category string

	// Subcategory is an string field that can provide further refinement of the error type.
	// It can be used to specify a more specific subcategory within the main category.
	// For example, within the "Database" category, you might have subcategories like "Connection" or "Query".
	Subcategory string

	// Description is a human-readable string that describes the error in more detail.
	// It can provide additional information about what caused the error or where it occurred.
	Description string
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// GetCause returns the underlying cause of the error
func (e *Error) GetCause() error {
	return e.Cause
}

// Decorate decorates an error with a message and returns the decorated error.
// If the provided 'err' is nil, it returns nil.
//
// Parameters:
//   - err: The original error that you want to decorate.
//   - message: A string message that you want to add as a prefix to the error.
//   - errorType: An optional error type (can be nil) used to categorize the error.
//
// Returns:
//   - A decorated error, which includes the provided 'message' as a prefix to the original error message,
//     the original error as the cause, and an optional error type.
//
// Example:
//
//	originalError := someFunctionThatMayReturnAnError()
//	decoratedError := Decorate(originalError, "Failed to perform operation", nil)
//	if decoratedError != nil {
//	    log.Printf("Decorated Error: %s", decoratedError.Message)
//	}
func Decorate(err error, message string, errorType *Type) *Error {
	if err == nil {
		return nil
	}

	// Combine the original error message with the provided error type description.
	errorType.Description = fmt.Sprintf("%s: %s", errorType.Description, err.Error())

	// Create a new error with the decorated message and the original error as the cause.
	return &Error{
		Message:    message,
		Cause:      err,
		Type:       errorType,
		Stacktrace: GetStackTrace(),
	}
}

// getStackTrace captures the current call stack trace and returns it as a string.
//
// This function uses Go's built-in 'debug' package to capture the current call stack trace.
// The stack trace is then converted into a human-readable string format and returned.
//
// The call stack trace includes information about the sequence of function calls that
// led to the current point in the code. It is typically used for debugging and error
// reporting to provide insights into how the program reached a particular state or error.
//
// Returns:
//   - A string representation of the current call stack trace.
//
// Example Usage:
//
//	stackTrace := GetStackTrace()
//	fmt.Println("Stack Trace:")
//	fmt.Println(stackTrace)
func GetStackTrace() string {
	// Capture stack trace
	stackTrace := debug.Stack()

	// Convert stack trace to string
	stackString := string(stackTrace)

	return stackString
}

// GetRootErrorString extracts the root error message from an error value.
//
// This function takes an error value 'err' and retrieves the root cause of the error message,
// which is typically the last part of a colon-separated error message. If the error message
// doesn't contain a colon, it returns the original error message as the root cause.
//
// Parameters:
//   - err: The error value from which to extract the root cause message.
//
// Returns:
//   - A string representing the root cause of the error.
//
// Example:
//
//	err := fmt.Errorf("Failed to process data: unexpected EOF")
//	rootErrorMessage := GetRootErrorString(err) // Returns "unexpected EOF"
//
//	err = fmt.Errorf("An error occurred")
//	rootErrorMessage := GetRootErrorString(err) // Returns "An error occurred"
func GetRootErrorString(err error) string {
	// Start with the original error message
	errMsg := err.Error()

	// Split the error message by ':' to find the root cause
	splittedString := strings.Split(errMsg, ":")

	// If there are multiple segments (e.g., "prefix: root cause"), return the last one
	if len(splittedString) > 1 {
		return strings.TrimSpace(splittedString[len(splittedString)-1])
	}

	// If there's only one segment, return the original error message
	return errMsg
}

// Is checks if the error is equivalent to the target error.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
