// Package router provides functionality for managing routes in GoAPI
package router

import (
	"github.com/gin-gonic/gin"
)

// APIProvider defines the interface that the API must implement
type APIProvider interface {
	AddRoute(method, path string, handler gin.HandlerFunc, opts ...RouteOption)
}

// Route represents a route in the API
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Tags        []string
	Summary     string
	Description string
	Responses   map[int]string
	Parameters  []Parameter
}

// Parameter represents a parameter in the API
type Parameter struct {
	Name        string
	In          string // "path", "query", "header", "body"
	Type        string // "string", "integer", "boolean", etc.
	Format      string // "int64", "date-time", etc.
	Required    bool
	Description string
	Schema      interface{} // For body parameters
}

// Schema represents a request/response schema
type Schema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
	Example    interface{}            `json:"example,omitempty"`
}

// RouteOption is a function that modifies a route
type RouteOption func(*Route)

// WithTags adds tags to a route for API documentation grouping
// Tags are used to organize related endpoints in Swagger UI
func WithTags(tags ...string) RouteOption {
	return func(route *Route) {
		route.Tags = append(route.Tags, tags...)
	}
}

// WithSummary adds a summary to a route for API documentation
// The summary provides a brief description of what the endpoint does
func WithSummary(summary string) RouteOption {
	return func(route *Route) {
		route.Summary = summary
	}
}

// WithDescription adds a detailed description to a route for API documentation
// The description provides comprehensive information about the endpoint's functionality
func WithDescription(description string) RouteOption {
	return func(route *Route) {
		route.Description = description
	}
}

// WithResponse adds an expected response configuration to a route
// This defines the possible HTTP status codes and their descriptions
func WithResponse(statusCode int, description string) RouteOption {
	return func(route *Route) {
		if route.Responses == nil {
			route.Responses = make(map[int]string)
		}
		route.Responses[statusCode] = description
	}
}

// WithParameter adds a parameter configuration to a route
// This allows for flexible parameter definitions with custom locations and types
func WithParameter(name, in, paramType, description string, required bool) RouteOption {
	return func(route *Route) {
		newParameter := Parameter{
			Name:        name,
			In:          in,
			Type:        paramType,
			Required:    required,
			Description: description,
		}
		route.Parameters = append(route.Parameters, newParameter)
	}
}

// WithPathParameter adds a path parameter to a route
func WithPathParameter(name, paramType, description string) RouteOption {
	return WithParameter(name, "path", paramType, description, true)
}

// WithQueryParameter adds a query parameter to a route
func WithQueryParameter(name, paramType, description string, required bool) RouteOption {
	return WithParameter(name, "query", paramType, description, required)
}

// WithRequestBody adds a request body schema configuration to a route
// This defines the expected structure and format of the request payload
func WithRequestBody(schema interface{}, description string) RouteOption {
	return func(route *Route) {
		bodyParameter := Parameter{
			Name:        "body",
			In:          "body",
			Required:    true,
			Description: description,
			Schema:      schema,
		}
		route.Parameters = append(route.Parameters, bodyParameter)
	}
}

// WithJSONSchema creates a JSON schema configuration from a struct example
// This automatically generates OpenAPI schema from Go struct definitions
func WithJSONSchema(example interface{}, description string) RouteOption {
	return WithRequestBody(example, description)
}

// RouterGroup represents a group of routes with a common path prefix
// It allows for organizing related routes and applying common middleware
type RouterGroup struct {
	apiProvider APIProvider // The API instance that will handle route registration
	pathPrefix  string      // Common path prefix for all routes in this group
}

// NewRouterGroup creates and initializes a new route group
// Parameters:
//   - apiProvider: The API instance that will handle route registration
//   - pathPrefix: Common path prefix for all routes in this group
// Returns:
//   - *RouterGroup: Pointer to the newly created route group
func NewRouterGroup(apiProvider APIProvider, pathPrefix string) *RouterGroup {
	return &RouterGroup{
		apiProvider: apiProvider,
		pathPrefix:  pathPrefix,
	}
}

// GET registers a new GET route in the group with the specified path and handler
// The final route path will be the group prefix combined with the provided path
func (routerGroup *RouterGroup) GET(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	routerGroup.apiProvider.AddRoute("GET", routerGroup.pathPrefix+path, handler, opts...)
}

// POST registers a new POST route in the group with the specified path and handler
// The final route path will be the group prefix combined with the provided path
func (routerGroup *RouterGroup) POST(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	routerGroup.apiProvider.AddRoute("POST", routerGroup.pathPrefix+path, handler, opts...)
}

// PUT registers a new PUT route in the group with the specified path and handler
// The final route path will be the group prefix combined with the provided path
func (routerGroup *RouterGroup) PUT(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	routerGroup.apiProvider.AddRoute("PUT", routerGroup.pathPrefix+path, handler, opts...)
}

// DELETE registers a new DELETE route in the group with the specified path and handler
// The final route path will be the group prefix combined with the provided path
func (routerGroup *RouterGroup) DELETE(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	routerGroup.apiProvider.AddRoute("DELETE", routerGroup.pathPrefix+path, handler, opts...)
}

// PATCH registers a new PATCH route in the group with the specified path and handler
// The final route path will be the group prefix combined with the provided path
func (routerGroup *RouterGroup) PATCH(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	routerGroup.apiProvider.AddRoute("PATCH", routerGroup.pathPrefix+path, handler, opts...)
}

// Group creates a new route subgroup with an additional path prefix
// This allows for nested route organization and hierarchical path structures
func (routerGroup *RouterGroup) Group(path string) *RouterGroup {
	return NewRouterGroup(routerGroup.apiProvider, routerGroup.pathPrefix+path)
}
