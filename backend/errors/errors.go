package errors

import "fmt"

// APIError es la estructura estándar para errores API
type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Details string `json:"details,omitempty"`
}

// Error implementa la interfaz error
func (e *APIError) Error() string {
	return e.Message
}

// Errores predefinidos comunes
var (
	// HTTP 4xx Errors
	ErrBadRequest = &APIError{
		Code:    400,
		Message: "Bad request",
	}

	ErrUnauthorized = &APIError{
		Code:    401,
		Message: "Unauthorized - Please login",
	}

	ErrForbidden = &APIError{
		Code:    403,
		Message: "Forbidden - Access denied",
	}

	ErrNotFound = &APIError{
		Code:    404,
		Message: "Resource not found",
	}

	ErrConflict = &APIError{
		Code:    409,
		Message: "Resource already exists",
	}

	ErrUnprocessable = &APIError{
		Code:    422,
		Message: "Validation failed",
	}

	// HTTP 5xx Errors
	ErrInternal = &APIError{
		Code:    500,
		Message: "Internal server error",
	}

	ErrUnavailable = &APIError{
		Code:    503,
		Message: "Service unavailable",
	}
)

// Constructor functions with details

// NewValidationError crea un error de validación con detalles
func NewValidationError(field, message string) *APIError {
	return &APIError{
		Code:    422,
		Message: "Validation failed",
		Details: fmt.Sprintf("%s: %s", field, message),
	}
}

// NewBadRequest crea un error de bad request con detalles
func NewBadRequest(details string) *APIError {
	return &APIError{
		Code:    400,
		Message: "Bad request",
		Details: details,
	}
}

// NewNotFound crea un error de no encontrado con detalles
func NewNotFound(resourceType string, id interface{}) *APIError {
	return &APIError{
		Code:    404,
		Message: "Resource not found",
		Details: fmt.Sprintf("%s with ID %v not found", resourceType, id),
	}
}

// NewConflict crea un error de conflicto con detalles
func NewConflict(details string) *APIError {
	return &APIError{
		Code:    409,
		Message: "Resource already exists",
		Details: details,
	}
}

// NewInternal crea un error interno con detalles
func NewInternal(details string) *APIError {
	return &APIError{
		Code:    500,
		Message: "Internal server error",
		Details: details,
	}
}

// NewDatabaseError crea un error de base de datos
func NewDatabaseError(operation string, err error) *APIError {
	return &APIError{
		Code:    500,
		Message: "Database error",
		Details: fmt.Sprintf("%s failed: %v", operation, err),
	}
}

// NewValidationErrors agrupa múltiples errores de validación
type ValidationErrors map[string]string

// Error implementa la interfaz error para ValidationErrors
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	details := ""
	for field, msg := range ve {
		details += fmt.Sprintf("%s: %s; ", field, msg)
	}
	return fmt.Sprintf("validation failed: %s", details[:len(details)-2])
}

// ToAPIError convierte ValidationErrors a APIError
func (ve ValidationErrors) ToAPIError() *APIError {
	details := ""
	for field, msg := range ve {
		details += fmt.Sprintf("%s: %s; ", field, msg)
	}
	if details != "" {
		details = details[:len(details)-2]
	}

	return &APIError{
		Code:    422,
		Message: "Validation failed",
		Details: details,
	}
}
