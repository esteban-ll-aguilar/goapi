# GoAPI Examples

This folder contains complete examples of how to use GoAPI in different scenarios.

## 📁 Examples Structure

```
examples/
├── README.md           # This file
├── basic/             # Basic example
│   ├── main.go        # Basic example code
│   └── go.mod         # Example dependencies
└── advanced/          # Advanced example
    ├── main.go        # Advanced example code
    └── go.mod         # Example dependencies
```

## 🚀 Basic Example

The basic example demonstrates:
- Initial GoAPI configuration
- Simple route definitions
- Basic response handling
- Automatic documentation

### Run the basic example:

```bash
cd basic
go mod tidy
go run main.go
```

**Available endpoints:**
- `GET /api/v1/items` - Get all items
- `GET /api/v1/items/{id}` - Get an item by ID
- `POST /api/v1/items` - Create a new item

**Documentation:**
- Swagger UI: http://localhost:8080/docs
- ReDoc: http://localhost:8080/redoc
- Main page: http://localhost:8080/

## 🔥 Advanced Example

The advanced example demonstrates all GoAPI features:
- ✅ Automatic data validation
- ✅ Dependency Injection
- ✅ Customizable middleware
- ✅ Standardized responses
- ✅ Automatic pagination
- ✅ Rate limiting
- ✅ CORS configuration
- ✅ Centralized error handling

### Run the advanced example:

```bash
cd advanced
go mod tidy
go run main.go
```

**Available endpoints:**
- `GET /api/v1/users` - Get users (with pagination and filters)
- `GET /api/v1/users/{id}` - Get a user by ID
- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/{id}` - Update a user
- `DELETE /api/v1/users/{id}` - Delete a user
- `GET /api/v1/stats` - Get statistics (with dependency injection)

**Demonstrated features:**

### Automatic Validation
```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=18,max=120"`
}
```

### Pagination
```bash
GET /api/v1/users?page=1&page_size=10&active=true
```

### Rate Limiting
- 100 requests per minute
- Burst size of 10

### Standardized Responses
- Success responses with consistent format
- Centralized error handling
- Data validation with descriptive messages

## 🧪 Test the Examples

### Using curl

**Get users:**
```bash
curl http://localhost:8080/api/v1/users
```

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25
  }'
```

**Get user by ID:**
```bash
curl http://localhost:8080/api/v1/users/1
```

**Update user:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "age": 26
  }'
```

**Delete user:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**Get statistics:**
```bash
curl http://localhost:8080/api/v1/stats
```

### Using Postman

1. Import the collection from Swagger documentation
2. Visit http://localhost:8080/docs
3. Download the OpenAPI/Swagger file
4. Import into Postman

## 📊 Example Responses

### Success response:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25,
    "is_active": true
  },
  "message": "User created successfully"
}
```

### Paginated response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 25,
      "is_active": true
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

### Validation error response:
```json
{
  "success": false,
  "error": "Validation failed",
  "details": [
    {
      "field": "email",
      "message": "must be a valid email address",
      "value": "invalid-email"
    },
    {
      "field": "age",
      "message": "must be at least 18",
      "value": 15
    }
  ]
}
```

## 🔧 Customization

### Modify Configuration

You can modify the configuration in any example:

```go
config := goapi.DefaultConfig()
config.Title = "My Custom API"
config.Description = "Custom description"
config.Version = "1.0.0"
config.BasePath = "/api/v1"
config.Host = "localhost:8080"
config.Debug = true
```

### Add Middleware

```go
// Custom CORS
api.AddCORS(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
})

// Custom rate limiting
api.AddRateLimit(middleware.RateLimitConfig{
    RequestsPerMinute: 60,
    BurstSize: 10,
})

// Custom middleware
api.AddMiddleware(func(c *gin.Context) {
    // Your logic here
    c.Next()
})
```

### Add New Routes

```go
// Simple route
api.GET("/health", func(c *gin.Context) {
    responses.Success(c, gin.H{"status": "healthy"})
})

// Route group
v1 := api.Group("/api/v1")
{
    products := v1.Group("/products")
    {
        products.GET("", GetProducts)
        products.POST("", CreateProduct)
    }
}
```

## 📚 Additional Resources

- [Complete documentation](../README.md)
- [Installation guide](../INSTRUCTIONS.md)
- [Gin documentation](https://gin-gonic.com/)
- [Validator documentation](https://github.com/go-playground/validator)
- [Swagger documentation](https://swagger.io/)

## 🤝 Contributing

Have ideas for new examples? Contribute!

1. Fork the project
2. Create a new directory in `examples/`
3. Add your example with documentation
4. Submit a Pull Request

## 📞 Support

If you have problems with the examples:

1. Verify you have Go 1.21+ installed
2. Run `go mod tidy` in the example directory
3. Check that port 8080 is available
4. Review the [documentation](../README.md)
<!-- 5. Open an [issue](https://github.com/esteban-ll-aguilar/goapi/issues) -->

Enjoy exploring GoAPI! 🚀
