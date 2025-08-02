// Package core provides core functionality for GoAPI
package core

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi/router"
)

// IndexHandler generates a handler for the main page in FastAPI style
func IndexHandler(config interface{}, routes []router.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, generateIndexHTML(config, routes))
	}
}

// RedocHandler generates a handler for ReDoc documentation
func RedocHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, redocHTML)
	}
}

// generateIndexHTML generates HTML for the main page
func generateIndexHTML(config interface{}, routes []router.Route) string {
	// Use reflection to access fields of the config structure
	v := reflect.ValueOf(config)
	title := ""
	description := ""
	basePath := ""

	// If it's a pointer, get the value it points to
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Only proceed if it's a struct
	if v.Kind() == reflect.Struct {
		if titleField := v.FieldByName("Title"); titleField.IsValid() && titleField.Kind() == reflect.String {
			title = titleField.String()
		}
		if descField := v.FieldByName("Description"); descField.IsValid() && descField.Kind() == reflect.String {
			description = descField.String()
		}
		if basePathField := v.FieldByName("BasePath"); basePathField.IsValid() && basePathField.Kind() == reflect.String {
			basePath = basePathField.String()
		}
	}

	return `
<!DOCTYPE html>
<html>
<head>
    <title>` + title + `</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">
    <style>
        body { margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { background-color: #1a1a1a; color: white; padding: 20px; text-align: center; }
        .header h1 { margin: 0; }
        .main { background-color: white; border-radius: 5px; padding: 20px; margin-top: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .endpoint { border: 1px solid #eee; border-radius: 5px; margin-bottom: 10px; padding: 10px; }
        .method { display: inline-block; padding: 3px 8px; border-radius: 3px; font-weight: bold; margin-right: 10px; }
        .get { background-color: #61affe; color: white; }
        .post { background-color: #49cc90; color: white; }
        .put { background-color: #fca130; color: white; }
        .delete { background-color: #f93e3e; color: white; }
        .patch { background-color: #50e3c2; color: white; }
        .docs-link { margin-top: 20px; text-align: center; }
        .docs-button { display: inline-block; background-color: #1a1a1a; color: white; padding: 10px 20px; 
                       border-radius: 5px; text-decoration: none; font-weight: bold; margin: 0 10px; transition: background-color 0.2s; }
        .docs-button:hover { background-color: #444; color: white; text-decoration: none; }
    </style>
</head>
<body>
    <div class="header">
        <h1>` + title + `</h1>
        <p>` + description + `</p>
    </div>
    <div class="container">
        <div class="main">
            <h2>API Endpoints</h2>
            ` + generateEndpointsHTML(routes) + `
        </div>
        <div class="docs-link">
            <a href="/docs" class="docs-button">Swagger UI</a>
            <a href="/redoc" class="docs-button">ReDoc</a>
            <a href="` + basePath + `" class="docs-button">API Root</a>
        </div>
    </div>
</body>
</html>
`
}

// generateEndpointsHTML generates HTML for the endpoints
func generateEndpointsHTML(routes []router.Route) string {
	html := ""
	for _, route := range routes {
		// Only show routes with a path that are not documentation
		if route.Path != "/" && route.Path != "/docs" && route.Path != "/redoc" &&
			route.Path != "/swagger/*any" && route.Path != "/redoc/index.html" {

			methodClass := "get"
			switch route.Method {
			case "POST":
				methodClass = "post"
			case "PUT":
				methodClass = "put"
			case "DELETE":
				methodClass = "delete"
			case "PATCH":
				methodClass = "patch"
			}

			description := route.Description
			if description == "" {
				description = route.Summary
			}

			html += `
            <div class="endpoint">
                <span class="method ` + methodClass + `">` + route.Method + `</span>
                <span class="path">` + route.Path + `</span>
                <p>` + description + `</p>
            </div>`
		}
	}
	return html
}

// HTML page for ReDoc
const redocHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>GoAPI - ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body { margin: 0; padding: 0; }
        #redoc-container { min-height: 100vh; }
    </style>
</head>
<body>
    <div id="redoc-container"></div>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
    <script>
        Redoc.init('/openapi.json', {
            scrollYOffset: 50,
            hideDownloadButton: false,
            expandResponses: '200,201',
            theme: {
                colors: { primary: { main: '#1a1a1a' } },
                sidebar: { width: '300px' }
            }
        }, document.getElementById('redoc-container'))
    </script>
</body>
</html>
`
