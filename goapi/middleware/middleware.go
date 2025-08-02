// Package middleware provides middleware functionality for GoAPI
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi/validation"
)

// MiddlewareFunc represents a middleware function
type MiddlewareFunc func() gin.HandlerFunc

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// CORS returns a CORS middleware
func CORS(config ...CORSConfig) gin.HandlerFunc {
	var cfg CORSConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		// Set other CORS headers
		if len(cfg.AllowMethods) > 0 {
			methods := ""
			for i, method := range cfg.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Header("Access-Control-Allow-Methods", methods)
		}

		if len(cfg.AllowHeaders) > 0 {
			headers := ""
			for i, header := range cfg.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Allow-Headers", headers)
		}

		if len(cfg.ExposeHeaders) > 0 {
			headers := ""
			for i, header := range cfg.ExposeHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Expose-Headers", headers)
		}

		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if cfg.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%.0f", cfg.MaxAge.Seconds()))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// ErrorHandler handles errors in a FastAPI-like manner
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Check if it's a validation error
			if validationErrors, ok := err.Err.(validation.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"detail": validationErrors,
					"type":   "validation_error",
				})
				return
			}

			// Handle other types of errors
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"detail": "Invalid request format",
					"type":   "bind_error",
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, gin.H{
					"detail": err.Error(),
					"type":   "public_error",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"detail": "Internal server error",
					"type":   "internal_error",
				})
			}
		}
	}
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// RateLimit provides basic rate limiting (simplified implementation)
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	// This is a simplified rate limiter
	// In production, you'd want to use a more sophisticated implementation
	// with Redis or similar for distributed rate limiting

	requestCounts := make(map[string]int)
	lastReset := time.Now()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Reset counts every minute
		if time.Since(lastReset) > time.Minute {
			requestCounts = make(map[string]int)
			lastReset = time.Now()
		}

		// Check current request count
		currentCount := requestCounts[clientIP]
		if currentCount >= config.RequestsPerMinute {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"detail": "Rate limit exceeded",
				"type":   "rate_limit_error",
			})
			c.Abort()
			return
		}

		// Increment request count
		requestCounts[clientIP] = currentCount + 1

		c.Next()
	}
}

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// CSP más permisivo para permitir que Swagger UI y ReDoc funcionen
		csp := "default-src 'self'; " +
			"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; " +
			"worker-src 'self' blob:; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self'"

		c.Header("Content-Security-Policy", csp)
		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Timeout middleware adds request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set a timeout for the request context
		ctx := c.Request.Context()
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		c.Request = c.Request.WithContext(timeoutCtx)
		c.Next()
	}
}

// Recovery middleware with custom error handling
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": fmt.Sprintf("Internal server error: %s", err),
				"type":   "panic_error",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "Internal server error",
				"type":   "panic_error",
			})
		}
		c.Abort()
	})
}

// Authentication middleware (basic implementation)
func Authentication(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Authorization header required",
				"type":   "authentication_error",
			})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// In a real implementation, you would validate the JWT token here
		// For now, we'll just check if it matches a simple secret
		if token != secretKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid token",
				"type":   "authentication_error",
			})
			c.Abort()
			return
		}

		// Set user information in context (mock data)
		c.Set("user_id", "user123")
		c.Set("username", "testuser")
		c.Next()
	}
}

// Compression middleware
func Compression() gin.HandlerFunc {
	// This would typically use gzip compression
	// For now, we'll return a placeholder
	return func(c *gin.Context) {
		// In a real implementation, you'd use gin-contrib/gzip
		c.Next()
	}
}
