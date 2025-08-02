# Instructions for Running GoAPI

## 🚀 Steps to run the project

### 1. Install dependencies
```bash
go mod tidy
```

### 2. Install swag (tool for generating documentation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 3. Generate Swagger documentation
```bash
swag init -g examples/advanced_example.go
```

### 4. Run the advanced example
```bash
go run examples/advanced_example.go
```

## 📚 Access documentation

Once the server is running, you can access:

- **Main page**: http://localhost:8080/
- **Swagger UI**: http://localhost:8080/docs or http://localhost:8080/swagger/index.html
- **ReDoc**: http://localhost:8080/redoc

## 🧪 Test the API

### Get all users
```bash
curl http://localhost:8080/api/v1/users
```

### Get users with pagination
```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=5"
```

### Create a new user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "age": 28
  }'
```

### Get a user by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update a user
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Carlos Smith",
    "age": 26
  }'
```

### Delete a user
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Get statistics (dependency injection example)
```bash
curl http://localhost:8080/api/v1/stats
```

## 🔧 Troubleshooting

### If Swagger UI doesn't show content:

1. Make sure you have run `swag init -g examples/advanced_example.go`
2. Verify that files have been generated in the `docs/` directory
3. Restart the server after generating documentation

### If there are compilation errors:

1. Run `go mod tidy` to update dependencies
2. Verify that all imports are correct
3. Make sure you have Go 1.19 or higher

## 📁 Generated file structure

After running `swag init`, you should see:

```
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

These files contain the automatically generated OpenAPI documentation.

## ✨ Demonstrated features

- ✅ Automatic data validation
- ✅ Dependency Injection
- ✅ Middleware (CORS, Rate Limiting, etc.)
- ✅ Standardized responses
- ✅ Automatic pagination
- ✅ Swagger/OpenAPI documentation
- ✅ Centralized error handling
- ✅ Declarative syntax similar to FastAPI

Enjoy exploring GoAPI! 🎉
