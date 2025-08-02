// Package validation provides validation functionality for GoAPI
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// ValidateStruct validates a struct using tags
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// FormatValidationErrors formats validator errors into a more readable format
func FormatValidationErrors(err error) ValidationErrors {
	var validationErrors ValidationErrors
	
	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErrors {
			validationError := ValidationError{
				Field: fieldError.Field(),
				Tag:   fieldError.Tag(),
				Value: fmt.Sprintf("%v", fieldError.Value()),
			}
			
			// Generate human-readable messages
			switch fieldError.Tag() {
			case "required":
				validationError.Message = fmt.Sprintf("El campo '%s' es requerido", fieldError.Field())
			case "min":
				validationError.Message = fmt.Sprintf("El campo '%s' debe tener un valor mínimo de %s", fieldError.Field(), fieldError.Param())
			case "max":
				validationError.Message = fmt.Sprintf("El campo '%s' debe tener un valor máximo de %s", fieldError.Field(), fieldError.Param())
			case "email":
				validationError.Message = fmt.Sprintf("El campo '%s' debe ser un email válido", fieldError.Field())
			case "url":
				validationError.Message = fmt.Sprintf("El campo '%s' debe ser una URL válida", fieldError.Field())
			case "len":
				validationError.Message = fmt.Sprintf("El campo '%s' debe tener exactamente %s caracteres", fieldError.Field(), fieldError.Param())
			case "gte":
				validationError.Message = fmt.Sprintf("El campo '%s' debe ser mayor o igual a %s", fieldError.Field(), fieldError.Param())
			case "lte":
				validationError.Message = fmt.Sprintf("El campo '%s' debe ser menor o igual a %s", fieldError.Field(), fieldError.Param())
			default:
				validationError.Message = fmt.Sprintf("El campo '%s' no cumple con la validación '%s'", fieldError.Field(), fieldError.Tag())
			}
			
			validationErrors = append(validationErrors, validationError)
		}
	}
	
	return validationErrors
}

// QueryParam represents a query parameter with validation
type QueryParam struct {
	Name         string
	Type         string
	Required     bool
	DefaultValue interface{}
	Description  string
	Example      interface{}
}

// PathParam represents a path parameter with validation
type PathParam struct {
	Name        string
	Type        string
	Description string
	Example     interface{}
}

// ParseQueryParams parses and validates query parameters from a request
func ParseQueryParams(queryValues map[string][]string, params []QueryParam) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var validationErrors ValidationErrors

	for _, param := range params {
		values, exists := queryValues[param.Name]
		
		// Check if required parameter is missing
		if param.Required && (!exists || len(values) == 0 || values[0] == "") {
			validationErrors = append(validationErrors, ValidationError{
				Field:   param.Name,
				Tag:     "required",
				Message: fmt.Sprintf("El parámetro de consulta '%s' es requerido", param.Name),
			})
			continue
		}
		
		// Use default value if parameter is not provided
		if !exists || len(values) == 0 || values[0] == "" {
			if param.DefaultValue != nil {
				result[param.Name] = param.DefaultValue
			}
			continue
		}
		
		// Parse the value based on type
		value := values[0]
		parsedValue, err := parseValue(value, param.Type)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field:   param.Name,
				Tag:     "type",
				Value:   value,
				Message: fmt.Sprintf("El parámetro '%s' debe ser de tipo %s", param.Name, param.Type),
			})
			continue
		}
		
		result[param.Name] = parsedValue
	}
	
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}
	
	return result, nil
}

// parseValue parses a string value to the specified type
func parseValue(value, valueType string) (interface{}, error) {
	switch valueType {
	case "string":
		return value, nil
	case "int":
		return strconv.Atoi(value)
	case "int64":
		return strconv.ParseInt(value, 10, 64)
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "bool":
		return strconv.ParseBool(value)
	default:
		return value, nil
	}
}

// BindAndValidate binds request data to a struct and validates it
func BindAndValidate(data interface{}, target interface{}) error {
	// First, try to bind the data
	if err := bindData(data, target); err != nil {
		return fmt.Errorf("error binding data: %w", err)
	}
	
	// Then validate the struct
	validator := NewValidator()
	if err := validator.ValidateStruct(target); err != nil {
		return FormatValidationErrors(err)
	}
	
	return nil
}

// bindData binds data to a target struct (simplified version)
func bindData(_, target interface{}) error {
	// This is a simplified implementation
	// In a real implementation, you would use reflection to properly bind data
	targetValue := reflect.ValueOf(target)
	
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}
	
	// For now, we assume data is already in the correct format
	// This would need more sophisticated implementation for real use
	return nil
}
