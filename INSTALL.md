# Instalación y Configuración de GoAPI

Esta guía te ayudará a instalar y configurar GoAPI en tu proyecto.

## 📋 Requisitos Previos

- Go 1.21 o superior
- Git (para clonar el repositorio)

## 🚀 Instalación

### Opción 1: Usando go get (Recomendado)

```bash
# Crear un nuevo proyecto
mkdir mi-proyecto-api
cd mi-proyecto-api

# Inicializar módulo Go
go mod init mi-proyecto-api

# Instalar GoAPI
go get github.com/esteban-ll-aguilar/goapi
```

### Opción 2: Clonando el repositorio

```bash
# Clonar el repositorio
git clone https://github.com/esteban-ll-aguilar/goapi.git
cd goapi

# Instalar dependencias
go mod tidy
```

## 🏃‍♂️ Primer Proyecto

### 1. Crear main.go

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/esteban-ll-aguilar/goapi/goapi"
    "github.com/esteban-ll-aguilar/goapi/goapi/responses"
)

func main() {
    // Crear configuración
    config := goapi.DefaultConfig()
    config.Title = "Mi Primera API"
    config.Description = "Una API creada con GoAPI"
    
    // Crear instancia de GoAPI
    api := goapi.New(config)
    
    // Definir una ruta simple
    api.GET("/hello", func(c *gin.Context) {
        responses.Success(c, gin.H{
            "message": "¡Hola desde GoAPI!",
            "status":  "success",
        })
    })
    
    // Ejecutar servidor
    api.Run(":8080")
}
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Ejecutar la aplicación

```bash
go run main.go
```

### 4. Probar la API

Abre tu navegador y visita:
- **API**: http://localhost:8080/hello
- **Documentación**: http://localhost:8080/docs
- **Página principal**: http://localhost:8080/

## 📚 Ejemplos Incluidos

El proyecto incluye ejemplos completos que puedes ejecutar:

### Ejemplo Básico

```bash
cd examples/basic
go mod tidy
go run main.go
```

### Ejemplo Avanzado

```bash
cd examples/advanced
go mod tidy
go run main.go
```

## 🔧 Configuración Avanzada

### Estructura de Proyecto Recomendada

```
mi-proyecto/
├── main.go
├── go.mod
├── go.sum
├── handlers/
│   ├── users.go
│   ├── auth.go
│   └── products.go
├── models/
│   ├── user.go
│   └── product.go
├── services/
│   ├── user_service.go
│   └── product_service.go
└── config/
    └── config.go
```

### Ejemplo de Configuración Personalizada

```go
package main

import (
    "github.com/esteban-ll-aguilar/goapi/goapi"
    "github.com/esteban-ll-aguilar/goapi/goapi/middleware"
)

func main() {
    // Configuración personalizada
    config := goapi.APIConfig{
        Title:       "Mi API Empresarial",
        Description: "API para gestión empresarial",
        Version:     "2.0.0",
        BasePath:    "/api/v2",
        Host:        "api.miempresa.com",
        Schemes:     []string{"https"},
        Debug:       false,
        Contact: goapi.Contact{
            Name:  "Soporte Técnico",
            URL:   "https://miempresa.com/soporte",
            Email: "soporte@miempresa.com",
        },
        License: goapi.License{
            Name: "MIT",
            URL:  "https://opensource.org/licenses/MIT",
        },
    }
    
    api := goapi.New(config)
    
    // Configurar middleware
    api.AddCORS(middleware.CORSConfig{
        AllowOrigins: []string{"https://miempresa.com"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    })
    
    api.AddRateLimit(middleware.RateLimitConfig{
        RequestsPerMinute: 1000,
        BurstSize:         50,
    })
    
    // Tus rutas aquí...
    
    api.Run(":8080")
}
```

## 🛠️ Desarrollo

### Ejecutar en Modo Debug

```go
config := goapi.DefaultConfig()
config.Debug = true  // Habilita logs detallados
```

### Configurar CORS para Desarrollo

```go
api.AddCORS(middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"*"},
    AllowCredentials: true,
})
```

### Hot Reload (Opcional)

Instala Air para hot reload durante el desarrollo:

```bash
# Instalar Air
go install github.com/cosmtrek/air@latest

# Crear archivo .air.toml
air init

# Ejecutar con hot reload
air
```

## 🚀 Despliegue

### Variables de Entorno

```bash
export API_PORT=8080
export API_DEBUG=false
export API_HOST=api.midominio.com
```

### Docker

Crear `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
```

Construir y ejecutar:

```bash
docker build -t mi-api .
docker run -p 8080:8080 mi-api
```

## 🔍 Solución de Problemas

### Error: "could not import github.com/esteban-ll-aguilar/goapi"

```bash
# Limpiar caché de módulos
go clean -modcache

# Reinstalar dependencias
go mod tidy
```

### Error: "port already in use"

```bash
# Cambiar puerto en main.go
api.Run(":8081")  // o cualquier otro puerto disponible
```

### Error de CORS

```go
// Agregar configuración CORS permisiva para desarrollo
api.AddCORS(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"*"},
    AllowHeaders: []string{"*"},
})
```

## 📞 Soporte

Si encuentras problemas:

1. Revisa la [documentación](README.md)
2. Consulta los [ejemplos](examples/)
3. Abre un [issue](https://github.com/esteban-ll-aguilar/goapi/issues)

## 🎯 Próximos Pasos

1. Lee la [documentación completa](README.md)
2. Explora los [ejemplos](examples/)
3. Únete a la comunidad
4. Contribuye al proyecto

¡Feliz desarrollo con GoAPI! 🚀
