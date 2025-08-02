package core

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi/models"
)

// Handler is the base interface for all handlers
type Handler interface {
	Register(api interface{})
}

// ResponseError represents an error in the API response
type ResponseError struct {
	Error string `json:"error"`
}

// SendOK sends a successful response
func SendOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// SendCreated sends a successful creation response
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// SendError sends an error response
func SendError(c *gin.Context, status int, err error) {
	c.JSON(status, ResponseError{Error: err.Error()})
}

// ValidateJSON validates a JSON model
func ValidateJSON(c *gin.Context, model models.Model) bool {
	if err := c.ShouldBindJSON(model); err != nil {
		SendError(c, http.StatusBadRequest, err)
		return false
	}

	if err := model.Validate(); err != nil {
		SendError(c, http.StatusBadRequest, err)
		return false
	}

	return true
}
