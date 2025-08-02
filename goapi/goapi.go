// Package goapi provides a Go library for creating APIs in FastAPI style
package goapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"github.com/esteban-ll-aguilar/goapi/goapi/core"
	"github.com/esteban-ll-aguilar/goapi/goapi/dependencies"
	"github.com/esteban-ll-aguilar/goapi/goapi/middleware"
	"github.com/esteban-ll-aguilar/goapi/goapi/router"
	"github.com/esteban-ll-aguilar/goapi/goapi/validation"
)

// APIConfig contains the API configuration
type APIConfig struct {
	Title       string
	Description string
	Version     string
	BasePath    string
	Host        string
	Schemes     []string
	Contact     Contact
	License     License
	Debug       bool
}

// Contact contains contact information for the API
type Contact struct {
	Name  string
	URL   string
	Email string
}

// License contains license information for the API
type License struct {
	Name string
	URL  string
}

// DefaultConfig returns a default API configuration with sensible defaults
// This configuration can be used as a starting point and customized as needed
func DefaultConfig() APIConfig {
	return APIConfig{
		Title:       "GoAPI",
		Description: "API created with GoAPI framework",
		Version:     "1.0.0",
		BasePath:    "",
		Host:        "localhost:8080",
		Schemes:     []string{"http"},
		Debug:       true,
		Contact: Contact{
			Name:  "API Support",
			URL:   "https://github.com/esteban-ll-aguilar/goapi",
			Email: "support@example.com",
		},
		License: License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
	}
}

// GoAPI is the main structure that encapsulates all API functionality
// It provides a FastAPI-like interface for building REST APIs in Go
type GoAPI struct {
	config       APIConfig                         // API configuration settings
	router       *gin.Engine                       // Underlying Gin router instance
	routes       []router.Route                    // Collection of registered API routes
	endpoints    map[string]interface{}            // Map of endpoint configurations
	dependencies *dependencies.DependencyContainer // Dependency injection container
	validator    *validation.Validator             // Request validation handler
	middlewares  []gin.HandlerFunc                 // Collection of registered middlewares
}

// New creates and initializes a new GoAPI instance with the provided configuration
// It sets up the Gin router, initializes all components, and configures default middleware
// Parameters:
//   - configuration: APIConfig containing the API settings and metadata
//
// Returns:
//   - *GoAPI: Pointer to the newly created and configured GoAPI instance
func New(configuration APIConfig) *GoAPI {
	// Configure Gin mode based on debug setting
	if configuration.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create new Gin router instance
	ginRouterInstance := gin.New()

	// Initialize GoAPI instance with all components
	apiInstance := &GoAPI{
		config:       configuration,
		router:       ginRouterInstance,
		routes:       make([]router.Route, 0),
		endpoints:    make(map[string]interface{}),
		dependencies: dependencies.NewDependencyContainer(),
		validator:    validation.NewValidator(),
		middlewares:  make([]gin.HandlerFunc, 0),
	}

	// Setup default middleware stack
	apiInstance.setupDefaultMiddleware()

	return apiInstance
}

// AddRoute adds a new route to the API with the specified method, path, and handler
// It applies any provided route options to configure the route's metadata
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: URL path for the route
//   - handler: Gin handler function to process requests
//   - opts: Optional route configuration options
func (apiInstance *GoAPI) AddRoute(method, path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	newRoute := router.Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	// Apply all provided options to configure the route
	for _, routeOption := range opts {
		routeOption(&newRoute)
	}

	apiInstance.routes = append(apiInstance.routes, newRoute)
}

// WithTags adds tags to a route for API documentation grouping
// Tags are used to group related endpoints in Swagger UI
func WithTags(tags ...string) router.RouteOption {
	return router.WithTags(tags...)
}

// WithSummary adds a summary to a route for API documentation
// The summary provides a brief description of what the endpoint does
func WithSummary(summary string) router.RouteOption {
	return router.WithSummary(summary)
}

// WithDescription adds a detailed description to a route for API documentation
// The description provides comprehensive information about the endpoint's functionality
func WithDescription(description string) router.RouteOption {
	return router.WithDescription(description)
}

// WithPathParameter adds a path parameter configuration to a route
// Path parameters are part of the URL path (e.g., /users/{id})
func WithPathParameter(name, paramType, description string) router.RouteOption {
	return router.WithPathParameter(name, paramType, description)
}

// WithQueryParameter adds a query parameter configuration to a route
// Query parameters are passed in the URL query string (e.g., ?page=1&size=10)
func WithQueryParameter(name, paramType, description string, required bool) router.RouteOption {
	return router.WithQueryParameter(name, paramType, description, required)
}

// WithParameter adds a custom parameter configuration to a route
// This allows for flexible parameter definitions with custom locations and types
func WithParameter(name, in, paramType, description string, required bool) router.RouteOption {
	return router.WithParameter(name, in, paramType, description, required)
}

// WithRequestBody adds a request body schema configuration to a route
// This defines the expected structure and format of the request payload
func WithRequestBody(schema interface{}, description string) router.RouteOption {
	return router.WithRequestBody(schema, description)
}

// WithJSONSchema creates a JSON schema configuration from a struct example
// This automatically generates OpenAPI schema from Go struct definitions
func WithJSONSchema(example interface{}, description string) router.RouteOption {
	return router.WithJSONSchema(example, description)
}

// GET registers a new GET route with the specified path and handler
// GET routes are typically used for retrieving data without side effects
func (apiInstance *GoAPI) GET(path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	apiInstance.AddRoute(http.MethodGet, path, handler, opts...)
}

// POST registers a new POST route with the specified path and handler
// POST routes are typically used for creating new resources
func (apiInstance *GoAPI) POST(path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	apiInstance.AddRoute(http.MethodPost, path, handler, opts...)
}

// PUT registers a new PUT route with the specified path and handler
// PUT routes are typically used for updating existing resources completely
func (apiInstance *GoAPI) PUT(path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	apiInstance.AddRoute(http.MethodPut, path, handler, opts...)
}

// DELETE registers a new DELETE route with the specified path and handler
// DELETE routes are used for removing existing resources
func (apiInstance *GoAPI) DELETE(path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	apiInstance.AddRoute(http.MethodDelete, path, handler, opts...)
}

// PATCH registers a new PATCH route with the specified path and handler
// PATCH routes are typically used for partial updates to existing resources
func (apiInstance *GoAPI) PATCH(path string, handler gin.HandlerFunc, opts ...router.RouteOption) {
	apiInstance.AddRoute(http.MethodPatch, path, handler, opts...)
}

// Group creates a new route group with the specified path prefix
// Route groups allow for organizing related routes and applying common middleware
func (apiInstance *GoAPI) Group(path string) *router.RouterGroup {
	return router.NewRouterGroup(apiInstance, path)
}

// SetupRoutes configures and registers all defined routes with the underlying router
// This method should be called before starting the server to ensure all routes are available
func (apiInstance *GoAPI) SetupRoutes() {
	// Configure API documentation routes
	apiInstance.setupDocs()

	// Register all defined API routes with the Gin router
	for _, currentRoute := range apiInstance.routes {
		apiInstance.router.Handle(currentRoute.Method, currentRoute.Path, currentRoute.Handler)
	}
}

// setupDocs configures documentation routes
func (a *GoAPI) setupDocs() {
	// Generar documentaciรณn automรกticamente basรกndose en las rutas
	a.generateSwaggerSpec()

	// Escribir el archivo swagger.json dinรกmicamente ANTES del wildcard
	a.writeSwaggerFile()

	// Main route in FastAPI style
	a.router.GET("/", core.IndexHandler(a.config, a.routes))

	// Documentation routes
	a.router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	a.router.GET("/redoc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/redoc/index.html")
	})

	// Servir archivos estรกticos de documentaciรณn
	a.router.Static("/docs-static", "./goapi/docs")

	// Swagger documentation con URL personalizada
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/openapi.json")))

	// ReDoc documentation
	a.router.GET("/redoc/index.html", core.RedocHandler())
}

// writeSwaggerFile escribe el archivo swagger.json dinรกmicamente
func (a *GoAPI) writeSwaggerFile() {
	// Generar el contenido del archivo swagger.json
	swaggerContent := a.getSwaggerJSON()

	// Servir dinรกmicamente en una ruta que no conflicte con el wildcard
	a.router.GET("/openapi.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, swaggerContent)
	})
}

// generateSwaggerSpec genera automรกticamente la especificaciรณn Swagger
func (a *GoAPI) generateSwaggerSpec() {
	// Crear la especificaciรณn Swagger dinรกmicamente
	spec := &swag.Spec{
		Version:          a.config.Version,
		Host:             a.config.Host,
		BasePath:         a.config.BasePath,
		Schemes:          a.config.Schemes,
		Title:            a.config.Title,
		Description:      a.config.Description,
		InfoInstanceName: "swagger",
		SwaggerTemplate:  a.getSwaggerJSON(),
		LeftDelim:        "",
		RightDelim:       "",
	}

	// Registrar la especificaciรณn
	swag.Register(spec.InstanceName(), spec)
}

// generateSwaggerTemplate genera el template de Swagger basรกndose en las rutas
func (a *GoAPI) generateSwaggerTemplate() string {
	paths := make(map[string]interface{})

	// Generar paths basรกndose en las rutas registradas
	for _, route := range a.routes {
		if route.Path == "/" || route.Path == "/docs" || route.Path == "/redoc" ||
			route.Path == "/swagger/*any" || route.Path == "/redoc/index.html" {
			continue // Skip documentation routes
		}

		pathItem := make(map[string]interface{})
		operation := map[string]interface{}{
			"summary":     route.Summary,
			"description": route.Description,
			"tags":        route.Tags,
			"parameters":  a.getRouteParameters(route),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
					"schema": map[string]interface{}{
						"type": "object",
					},
				},
			},
		}

		if route.Summary == "" {
			operation["summary"] = "API endpoint"
		}
		if route.Description == "" {
			operation["description"] = "API endpoint description"
		}
		if len(route.Tags) == 0 {
			operation["tags"] = []string{"default"}
		}

		methodLower := strings.ToLower(route.Method)
		pathItem[methodLower] = operation
		paths[route.Path] = pathItem
	}

	// Crear el template completo
	template := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "{{.Title}}",
			"description": "{{escape .Description}}",
			"version":     "{{.Version}}",
			"contact": map[string]interface{}{
				"name":  a.config.Contact.Name,
				"url":   a.config.Contact.URL,
				"email": a.config.Contact.Email,
			},
			"license": map[string]interface{}{
				"name": a.config.License.Name,
				"url":  a.config.License.URL,
			},
		},
		"host":     "{{.Host}}",
		"basePath": "{{.BasePath}}",
		"schemes":  []string{"{{range .Schemes}}{{.}}{{end}}"},
		"paths":    paths,
	}

	// Convertir a JSON string
	templateBytes, _ := json.MarshalIndent(template, "", "  ")
	return string(templateBytes)
}

// getSwaggerJSON devuelve el JSON de Swagger
func (a *GoAPI) getSwaggerJSON() string {
	paths := make(map[string]interface{})

	// Generar paths basรกndose en las rutas registradas
	for _, route := range a.routes {
		if route.Path == "/" || route.Path == "/docs" || route.Path == "/redoc" ||
			route.Path == "/swagger/*any" || route.Path == "/redoc/index.html" ||
			route.Path == "/openapi.json" || route.Path == "/docs-static/*filepath" {
			continue // Skip documentation routes
		}

		// Convertir ruta de Gin (:id) a formato OpenAPI ({id})
		openAPIPath := a.convertToOpenAPIPath(route.Path)

		// Obtener o crear el pathItem para esta ruta
		pathItem, exists := paths[openAPIPath]
		if !exists {
			pathItem = make(map[string]interface{})
		}

		operation := map[string]interface{}{
			"summary":     route.Summary,
			"description": route.Description,
			"tags":        route.Tags,
			"parameters":  a.getRouteParameters(route),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
					"schema": map[string]interface{}{
						"type": "object",
					},
				},
			},
		}

		if route.Summary == "" {
			operation["summary"] = "API endpoint"
		}
		if route.Description == "" {
			operation["description"] = "API endpoint description"
		}
		if len(route.Tags) == 0 {
			operation["tags"] = []string{"default"}
		}

		methodLower := strings.ToLower(route.Method)
		pathItem.(map[string]interface{})[methodLower] = operation
		paths[openAPIPath] = pathItem
	}

	// Crear la especificaciรณn completa
	spec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       a.config.Title,
			"description": a.config.Description,
			"version":     a.config.Version,
			"contact": map[string]interface{}{
				"name":  a.config.Contact.Name,
				"url":   a.config.Contact.URL,
				"email": a.config.Contact.Email,
			},
			"license": map[string]interface{}{
				"name": a.config.License.Name,
				"url":  a.config.License.URL,
			},
		},
		"host":     a.config.Host,
		"basePath": a.config.BasePath,
		"schemes":  a.config.Schemes,
		"paths":    paths,
	}

	// Convertir a JSON string
	specBytes, _ := json.MarshalIndent(spec, "", "  ")
	return string(specBytes)
}

// Run runs the server on the specified port
func (a *GoAPI) Run(addr ...string) error {
	// Configure routes
	a.SetupRoutes()

	// Show server information
	serverAddr := ":8080"
	if len(addr) > 0 {
		serverAddr = addr[0]
	}

	log.Println("Server started at http://localhost" + serverAddr)
	log.Println("Documentation available at:")
	log.Println("- Swagger UI: http://localhost" + serverAddr + "/docs")
	log.Println("- ReDoc: http://localhost" + serverAddr + "/redoc")

	// Ejecutar servidor
	return a.router.Run(addr...)
}

// setupDefaultMiddleware configura middleware por defecto
func (a *GoAPI) setupDefaultMiddleware() {
	// Recovery middleware
	a.router.Use(middleware.Recovery())

	// Request logger
	if a.config.Debug {
		a.router.Use(middleware.RequestLogger())
	}

	// Error handler
	a.router.Use(middleware.ErrorHandler())

	// Security headers
	a.router.Use(middleware.SecurityHeaders())

	// Request ID
	a.router.Use(middleware.RequestID())

	// CORS con configuraciรณn por defecto
	a.router.Use(middleware.CORS())
}

// AddMiddleware agrega middleware personalizado
func (a *GoAPI) AddMiddleware(middlewareFunc gin.HandlerFunc) {
	a.middlewares = append(a.middlewares, middlewareFunc)
	a.router.Use(middlewareFunc)
}

// AddCORS configura CORS con configuraciรณn personalizada
func (a *GoAPI) AddCORS(config middleware.CORSConfig) {
	a.router.Use(middleware.CORS(config))
}

// AddRateLimit agrega rate limiting
func (a *GoAPI) AddRateLimit(config middleware.RateLimitConfig) {
	a.router.Use(middleware.RateLimit(config))
}

// AddAuthentication agrega autenticaciรณn
func (a *GoAPI) AddAuthentication(secretKey string) {
	a.router.Use(middleware.Authentication(secretKey))
}

// RegisterDependency registra una dependencia
func (a *GoAPI) RegisterDependency(provider dependencies.DependencyProvider, target interface{}) {
	a.dependencies.Register(provider, target)
}

// RegisterSingletonDependency registra una dependencia singleton
func (a *GoAPI) RegisterSingletonDependency(provider dependencies.DependencyProvider, target interface{}) {
	a.dependencies.RegisterSingleton(provider, target)
}

// GetDependencyContainer devuelve el contenedor de dependencias
func (a *GoAPI) GetDependencyContainer() *dependencies.DependencyContainer {
	return a.dependencies
}

// GetValidator devuelve el validador
func (a *GoAPI) GetValidator() *validation.Validator {
	return a.validator
}

// getRouteParameters obtiene los parรกmetros de una ruta, priorizando los configurados por el usuario
func (a *GoAPI) getRouteParameters(route router.Route) []map[string]interface{} {
	var parameters []map[string]interface{}

	// Primero, usar parรกmetros configurados por el usuario
	for _, param := range route.Parameters {
		parameter := map[string]interface{}{
			"name":        param.Name,
			"in":          param.In,
			"required":    param.Required,
			"description": param.Description,
		}

		// Manejar parรกmetros de body con schema
		if param.In == "body" && param.Schema != nil {
			parameter["schema"] = a.generateSchemaFromStruct(param.Schema)
		} else {
			parameter["type"] = param.Type
			if param.Format != "" {
				parameter["format"] = param.Format
			}
		}

		parameters = append(parameters, parameter)
	}

	// Si no hay parรกmetros configurados, usar detecciรณn automรกtica para parรกmetros de ruta
	if len(route.Parameters) == 0 {
		parameters = a.extractParameters(route.Path)
	}

	return parameters
}

// convertToOpenAPIPath convierte rutas de Gin (:id) a formato OpenAPI ({id})
func (a *GoAPI) convertToOpenAPIPath(path string) string {
	// Reemplazar :param con {param}
	segments := strings.Split(path, "/")

	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			paramName := segment[1:] // Remover el ":"
			segments[i] = "{" + paramName + "}"
		}
	}

	return strings.Join(segments, "/")
}

// extractParameters extrae parรกmetros de una ruta para la especificaciรณn OpenAPI
func (a *GoAPI) extractParameters(path string) []map[string]interface{} {
	var parameters []map[string]interface{}

	// Dividir la ruta en segmentos
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		// Buscar parรกmetros de ruta (formato :param)
		if strings.HasPrefix(segment, ":") {
			paramName := segment[1:] // Remover el ":"

			parameter := map[string]interface{}{
				"name":        paramName,
				"in":          "path",
				"required":    true,
				"type":        "string",
				"description": fmt.Sprintf("ID del %s", paramName),
			}

			// Personalizar descripciรณn segรบn el nombre del parรกmetro
			switch paramName {
			case "id":
				parameter["description"] = "ID del recurso"
				parameter["type"] = "integer"
				parameter["format"] = "int64"
			case "userId", "user_id":
				parameter["description"] = "ID del usuario"
				parameter["type"] = "integer"
				parameter["format"] = "int64"
			default:
				parameter["description"] = fmt.Sprintf("Parรกmetro %s", paramName)
			}

			parameters = append(parameters, parameter)
		}
	}

	return parameters
}

// generateSchemaFromStruct genera un schema OpenAPI desde un struct de Go
func (a *GoAPI) generateSchemaFromStruct(example interface{}) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	// Si el ejemplo es directamente un valor, usarlo como ejemplo
	if example != nil {
		schema["example"] = example

		// Usar reflection para generar propiedades del schema
		v := reflect.ValueOf(example)
		t := reflect.TypeOf(example)

		// Si es un puntero, obtener el valor al que apunta
		if v.Kind() == reflect.Ptr {
			if !v.IsNil() {
				v = v.Elem()
				t = t.Elem()
			}
		}

		// Solo procesar structs
		if v.Kind() == reflect.Struct {
			properties := make(map[string]interface{})
			var required []string

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				fieldValue := v.Field(i)

				// Obtener el nombre del campo JSON
				jsonTag := field.Tag.Get("json")
				fieldName := field.Name
				if jsonTag != "" && jsonTag != "-" {
					// Usar el nombre del tag JSON
					parts := strings.Split(jsonTag, ",")
					if parts[0] != "" {
						fieldName = parts[0]
					}

					// Verificar si es omitempty
					isOptional := false
					for _, part := range parts[1:] {
						if part == "omitempty" {
							isOptional = true
							break
						}
					}

					if !isOptional {
						required = append(required, fieldName)
					}
				} else {
					// Si no hay tag JSON, el campo es requerido por defecto
					required = append(required, fieldName)
				}

				// Generar el tipo del campo
				fieldSchema := a.getFieldSchema(fieldValue, field)
				properties[fieldName] = fieldSchema
			}

			schema["properties"] = properties
			if len(required) > 0 {
				schema["required"] = required
			}
		}
	}

	return schema
}

// getFieldSchema obtiene el schema de un campo especรญfico
func (a *GoAPI) getFieldSchema(fieldValue reflect.Value, field reflect.StructField) map[string]interface{} {
	fieldSchema := make(map[string]interface{})

	// Obtener el ejemplo del tag
	if example := field.Tag.Get("example"); example != "" {
		fieldSchema["example"] = example
	}

	// Determinar el tipo basรกndose en el tipo de Go
	switch fieldValue.Kind() {
	case reflect.String:
		fieldSchema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldSchema["type"] = "integer"
		fieldSchema["format"] = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldSchema["type"] = "integer"
		fieldSchema["format"] = "int64"
	case reflect.Float32, reflect.Float64:
		fieldSchema["type"] = "number"
		if fieldValue.Kind() == reflect.Float32 {
			fieldSchema["format"] = "float"
		} else {
			fieldSchema["format"] = "double"
		}
	case reflect.Bool:
		fieldSchema["type"] = "boolean"
	case reflect.Slice, reflect.Array:
		fieldSchema["type"] = "array"
		// Para arrays/slices, podrรญamos analizar el tipo del elemento
		fieldSchema["items"] = map[string]interface{}{
			"type": "string", // Por defecto
		}
	case reflect.Ptr:
		// Para punteros, analizar el tipo al que apuntan
		if !fieldValue.IsNil() {
			return a.getFieldSchema(fieldValue.Elem(), field)
		}
		fieldSchema["type"] = "string" // Por defecto para punteros nulos
	default:
		fieldSchema["type"] = "string" // Por defecto
	}

	return fieldSchema
}

// Router devuelve el router Gin subyacente
func (a *GoAPI) Router() *gin.Engine {
	return a.router
}
