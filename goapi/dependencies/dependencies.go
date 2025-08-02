// Package dependencies provides dependency injection functionality for GoAPI
package dependencies

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

// DependencyProvider is a function that provides a dependency
type DependencyProvider func(c *gin.Context) (interface{}, error)

// DependencyContainer manages dependencies
type DependencyContainer struct {
	providers map[reflect.Type]DependencyProvider
	instances map[reflect.Type]interface{}
	mutex     sync.RWMutex
}

// NewDependencyContainer creates a new dependency container
func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		providers: make(map[reflect.Type]DependencyProvider),
		instances: make(map[reflect.Type]interface{}),
	}
}

// Register registers a dependency provider
func (dc *DependencyContainer) Register(provider DependencyProvider, target interface{}) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	
	dc.providers[targetType] = provider
}

// RegisterSingleton registers a singleton dependency
func (dc *DependencyContainer) RegisterSingleton(provider DependencyProvider, target interface{}) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	
	dc.providers[targetType] = func(c *gin.Context) (interface{}, error) {
		dc.mutex.RLock()
		if instance, exists := dc.instances[targetType]; exists {
			dc.mutex.RUnlock()
			return instance, nil
		}
		dc.mutex.RUnlock()
		
		dc.mutex.Lock()
		defer dc.mutex.Unlock()
		
		// Double-check pattern
		if instance, exists := dc.instances[targetType]; exists {
			return instance, nil
		}
		
		instance, err := provider(c)
		if err != nil {
			return nil, err
		}
		
		dc.instances[targetType] = instance
		return instance, nil
	}
}

// Resolve resolves a dependency
func (dc *DependencyContainer) Resolve(c *gin.Context, target interface{}) error {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()
	
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	
	elementType := targetType.Elem()
	provider, exists := dc.providers[elementType]
	if !exists {
		return fmt.Errorf("no provider registered for type %s", elementType.String())
	}
	
	instance, err := provider(c)
	if err != nil {
		return fmt.Errorf("error resolving dependency: %w", err)
	}
	
	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(instance))
	return nil
}

// Dependency represents a dependency that can be injected
type Dependency interface {
	GetType() reflect.Type
}

// CommonDependencies provides common dependencies used in FastAPI-style applications
type CommonDependencies struct {
	container *DependencyContainer
}

// NewCommonDependencies creates common dependencies
func NewCommonDependencies() *CommonDependencies {
	container := NewDependencyContainer()
	
	// Register common dependencies
	container.Register(func(c *gin.Context) (interface{}, error) {
		return c, nil
	}, (*gin.Context)(nil))
	
	return &CommonDependencies{
		container: container,
	}
}

// GetContainer returns the dependency container
func (cd *CommonDependencies) GetContainer() *DependencyContainer {
	return cd.container
}

// Database dependency example
type Database struct {
	ConnectionString string
	Connected        bool
}

// Connect simulates database connection
func (db *Database) Connect() error {
	db.Connected = true
	return nil
}

// Close simulates database disconnection
func (db *Database) Close() error {
	db.Connected = false
	return nil
}

// DatabaseProvider provides a database dependency
func DatabaseProvider(connectionString string) DependencyProvider {
	return func(c *gin.Context) (interface{}, error) {
		db := &Database{
			ConnectionString: connectionString,
		}
		if err := db.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return db, nil
	}
}

// CurrentUser represents the current authenticated user
type CurrentUser struct {
	ID       string
	Username string
	Email    string
	Roles    []string
}

// CurrentUserProvider provides the current user from the request context
func CurrentUserProvider() DependencyProvider {
	return func(c *gin.Context) (interface{}, error) {
		// This would typically extract user information from JWT token or session
		// For now, we'll return a mock user
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			return nil, fmt.Errorf("user not authenticated")
		}
		
		return &CurrentUser{
			ID:       userID,
			Username: c.GetHeader("X-Username"),
			Email:    c.GetHeader("X-User-Email"),
			Roles:    []string{"user"}, // This would come from the token/session
		}, nil
	}
}

// Settings represents application settings
type Settings struct {
	AppName     string
	Version     string
	Environment string
	Debug       bool
}

// SettingsProvider provides application settings
func SettingsProvider(settings *Settings) DependencyProvider {
	return func(c *gin.Context) (interface{}, error) {
		return settings, nil
	}
}

// Logger represents a logger dependency
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// SimpleLogger is a simple logger implementation
type SimpleLogger struct {
	prefix string
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger(prefix string) *SimpleLogger {
	return &SimpleLogger{prefix: prefix}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	fmt.Printf("[INFO] %s: %s\n", l.prefix, fmt.Sprintf(msg, fields...))
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	fmt.Printf("[ERROR] %s: %s\n", l.prefix, fmt.Sprintf(msg, fields...))
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	fmt.Printf("[DEBUG] %s: %s\n", l.prefix, fmt.Sprintf(msg, fields...))
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	fmt.Printf("[WARN] %s: %s\n", l.prefix, fmt.Sprintf(msg, fields...))
}

// LoggerProvider provides a logger dependency
func LoggerProvider(prefix string) DependencyProvider {
	return func(c *gin.Context) (interface{}, error) {
		return NewSimpleLogger(prefix), nil
	}
}
