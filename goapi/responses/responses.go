// Package responses provides response handling functionality for GoAPI
package responses

import (
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Detail interface{} `json:"detail"`
	Type   string      `json:"type,omitempty"`
}

// ValidationErrorResponse represents validation errors
type ValidationErrorResponse struct {
	Detail []ResponseValidationError `json:"detail"`
	Type   string                    `json:"type"`
}

// ResponseValidationError represents a single validation error
type ResponseValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ResponseBuilder helps build standardized responses
type ResponseBuilder struct {
	statusCode int
	data       interface{}
	message    string
	errors     interface{}
}

// NewResponse creates a new response builder
func NewResponse() *ResponseBuilder {
	return &ResponseBuilder{
		statusCode: http.StatusOK,
	}
}

// WithStatus sets the HTTP status code
func (rb *ResponseBuilder) WithStatus(code int) *ResponseBuilder {
	rb.statusCode = code
	return rb
}

// WithData sets the response data
func (rb *ResponseBuilder) WithData(data interface{}) *ResponseBuilder {
	rb.data = data
	return rb
}

// WithMessage sets the response message
func (rb *ResponseBuilder) WithMessage(message string) *ResponseBuilder {
	rb.message = message
	return rb
}

// WithErrors sets the response errors
func (rb *ResponseBuilder) WithErrors(errors interface{}) *ResponseBuilder {
	rb.errors = errors
	return rb
}

// Send sends the response
func (rb *ResponseBuilder) Send(c *gin.Context) {
	response := Response{
		Data:    rb.data,
		Message: rb.message,
		Success: rb.statusCode >= 200 && rb.statusCode < 300,
		Errors:  rb.errors,
	}
	
	c.JSON(rb.statusCode, response)
}

// Success response helpers
func Success(c *gin.Context, data interface{}) {
	NewResponse().WithData(data).Send(c)
}

func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	NewResponse().WithData(data).WithMessage(message).Send(c)
}

func Created(c *gin.Context, data interface{}) {
	NewResponse().WithStatus(http.StatusCreated).WithData(data).Send(c)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error response helpers
func BadRequest(c *gin.Context, detail interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Detail: detail,
		Type:   "bad_request",
	})
}

func Unauthorized(c *gin.Context, detail interface{}) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Detail: detail,
		Type:   "unauthorized",
	})
}

func Forbidden(c *gin.Context, detail interface{}) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Detail: detail,
		Type:   "forbidden",
	})
}

func NotFound(c *gin.Context, detail interface{}) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Detail: detail,
		Type:   "not_found",
	})
}

func InternalServerError(c *gin.Context, detail interface{}) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Detail: detail,
		Type:   "internal_server_error",
	})
}

func ValidationError(c *gin.Context, errors []ResponseValidationError) {
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Detail: errors,
		Type:   "validation_error",
	})
}

// Paginated response helper
func Paginated(c *gin.Context, items interface{}, total, page, pageSize int) {
	totalPages := (total + pageSize - 1) / pageSize
	
	response := PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	
	Success(c, response)
}

// ResponseSchema represents a response schema for documentation
type ResponseSchema struct {
	StatusCode  int         `json:"status_code"`
	Description string      `json:"description"`
	Schema      interface{} `json:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// ResponseSchemas holds multiple response schemas
type ResponseSchemas map[int]ResponseSchema

// NewResponseSchemas creates a new response schemas collection
func NewResponseSchemas() ResponseSchemas {
	return make(ResponseSchemas)
}

// Add adds a response schema
func (rs ResponseSchemas) Add(statusCode int, description string, schema interface{}, example interface{}) ResponseSchemas {
	rs[statusCode] = ResponseSchema{
		StatusCode:  statusCode,
		Description: description,
		Schema:      schema,
		Example:     example,
	}
	return rs
}

// AddSuccess adds a success response schema
func (rs ResponseSchemas) AddSuccess(description string, schema interface{}, example interface{}) ResponseSchemas {
	return rs.Add(http.StatusOK, description, schema, example)
}

// AddCreated adds a created response schema
func (rs ResponseSchemas) AddCreated(description string, schema interface{}, example interface{}) ResponseSchemas {
	return rs.Add(http.StatusCreated, description, schema, example)
}

// AddBadRequest adds a bad request response schema
func (rs ResponseSchemas) AddBadRequest(description string) ResponseSchemas {
	return rs.Add(http.StatusBadRequest, description, ErrorResponse{}, ErrorResponse{
		Detail: "Validation failed",
		Type:   "validation_error",
	})
}

// AddUnauthorized adds an unauthorized response schema
func (rs ResponseSchemas) AddUnauthorized(description string) ResponseSchemas {
	return rs.Add(http.StatusUnauthorized, description, ErrorResponse{}, ErrorResponse{
		Detail: "Authentication required",
		Type:   "unauthorized",
	})
}

// AddNotFound adds a not found response schema
func (rs ResponseSchemas) AddNotFound(description string) ResponseSchemas {
	return rs.Add(http.StatusNotFound, description, ErrorResponse{}, ErrorResponse{
		Detail: "Resource not found",
		Type:   "not_found",
	})
}

// JSONResponse sends a JSON response with the specified status code
func JSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// XMLResponse sends an XML response with the specified status code
func XMLResponse(c *gin.Context, statusCode int, data interface{}) {
	c.XML(statusCode, data)
}

// HTMLResponse sends an HTML response
func HTMLResponse(c *gin.Context, statusCode int, template string, data interface{}) {
	c.HTML(statusCode, template, data)
}

// FileResponse sends a file response
func FileResponse(c *gin.Context, filepath string) {
	c.File(filepath)
}

// RedirectResponse sends a redirect response
func RedirectResponse(c *gin.Context, statusCode int, location string) {
	c.Redirect(statusCode, location)
}

// StreamResponse sends a streaming response
func StreamResponse(c *gin.Context, step func(w io.Writer) bool) {
	c.Stream(step)
}

// ResponseModel represents a model for response documentation
type ResponseModel struct {
	Type        reflect.Type
	Description string
	Example     interface{}
}

// NewResponseModel creates a new response model
func NewResponseModel(model interface{}, description string, example interface{}) ResponseModel {
	return ResponseModel{
		Type:        reflect.TypeOf(model),
		Description: description,
		Example:     example,
	}
}

// GetTypeName returns the type name of the model
func (rm ResponseModel) GetTypeName() string {
	if rm.Type.Kind() == reflect.Ptr {
		return rm.Type.Elem().Name()
	}
	return rm.Type.Name()
}

// Common response models
var (
	// StandardResponse is the standard response model
	StandardResponse = NewResponseModel(Response{}, "Standard API response", Response{
		Data:    "example data",
		Message: "Operation successful",
		Success: true,
	})
	
	// ErrorResponseModel is the error response model
	ErrorResponseModel = NewResponseModel(ErrorResponse{}, "Error response", ErrorResponse{
		Detail: "Error description",
		Type:   "error_type",
	})
	
	// ValidationErrorResponseModel is the validation error response model
	ValidationErrorResponseModel = NewResponseModel(ValidationErrorResponse{}, "Validation error response", ValidationErrorResponse{
		Detail: []ResponseValidationError{
			{
				Field:   "field_name",
				Message: "Field is required",
				Value:   "invalid_value",
			},
		},
		Type: "validation_error",
	})
	
	// PaginatedResponseModel is the paginated response model
	PaginatedResponseModel = NewResponseModel(PaginatedResponse{}, "Paginated response", PaginatedResponse{
		Items:      []interface{}{"item1", "item2"},
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	})
)
