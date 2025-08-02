package goapi

import (
	"fmt"
)

// APIError representa un error de la API
type APIError struct {
	Code    int
	Message string
	Details interface{}
}

// Error implementa la interfaz error
func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
	}
	return e.Message
}

// NewAPIError crea un nuevo error de API
func NewAPIError(code int, message string, details ...interface{}) *APIError {
	var detailsData interface{}
	if len(details) > 0 {
		detailsData = details[0]
	}

	return &APIError{
		Code:    code,
		Message: message,
		Details: detailsData,
	}
}

// NotFoundError crea un error de recurso no encontrado
func NotFoundError(resource string, id interface{}) *APIError {
	return NewAPIError(
		404,
		fmt.Sprintf("%s con id %v no encontrado", resource, id),
	)
}

// BadRequestError crea un error de solicitud incorrecta
func BadRequestError(message string) *APIError {
	return NewAPIError(400, message)
}

// ValidationError crea un error de validación
func ValidationError(message string, details interface{}) *APIError {
	return NewAPIError(422, message, details)
}

// InternalError crea un error interno del servidor
func InternalError(err error) *APIError {
	return NewAPIError(500, "Error interno del servidor: "+err.Error())
}
