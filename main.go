package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi"
	"github.com/esteban-ll-aguilar/goapi/goapi/middleware"
	"github.com/esteban-ll-aguilar/goapi/goapi/responses"
	"github.com/esteban-ll-aguilar/goapi/goapi/validation"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id" validate:"required,min=1" example:"1"`
	Name     string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Age      int    `json:"age" validate:"required,min=18,max=120" example:"25"`
	IsActive bool   `json:"is_active" example:"true"`
}

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Age   int    `json:"age" validate:"required,min=18,max=120" example:"25"`
}

// UpdateUserRequest represents the request payload for updating an existing user
// All fields are optional to allow partial updates
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=2,max=50" example:"John Doe Updated"`
	Email    string `json:"email,omitempty" validate:"omitempty,email" example:"john.updated@example.com"`
	Age      int    `json:"age,omitempty" validate:"omitempty,min=18,max=120" example:"26"`
	IsActive bool   `json:"is_active,omitempty" example:"true"`
}

// UserService provides business logic operations for user management
// It acts as a service layer between handlers and data storage
type UserService struct {
	users        []User // Collection of users stored in memory
	nextUserID   int    // Counter for generating unique user IDs
}

// NewUserService creates and initializes a new UserService instance
// It pre-populates the service with sample users for demonstration purposes
func NewUserService() *UserService {
	return &UserService{
		users: []User{
			{ID: 1, Name: "John Doe", Email: "john.doe@example.com", Age: 25, IsActive: true},
			{ID: 2, Name: "Jane Smith", Email: "jane.smith@example.com", Age: 30, IsActive: true},
			{ID: 3, Name: "Alice Johnson", Email: "alice.johnson@example.com", Age: 28, IsActive: true},
			{ID: 4, Name: "Bob Brown", Email: "bob.brown@example.com", Age: 35, IsActive: true},
			{ID: 5, Name: "Charlie White", Email: "charlie.white@example.com", Age: 40, IsActive: true},
			{ID: 6, Name: "Diana Green", Email: "diana.green@example.com", Age: 32, IsActive: true},
			{ID: 7, Name: "Ethan Blue", Email: "ethan.blue@example.com", Age: 29, IsActive: true},
			{ID: 8, Name: "Fiona Black", Email: "fiona.black@example.com", Age: 26, IsActive: true},
		},
		nextUserID: 9,
	}
}

// GetAll retrieves all users from the service
// Returns a slice containing all users currently stored in the service
func (userService *UserService) GetAll() []User {
	return userService.users
}

// GetByID retrieves a specific user by their unique identifier
// Parameters:
//   - userID: The unique identifier of the user to retrieve
// Returns:
//   - *User: Pointer to the user if found, nil otherwise
//   - bool: true if user was found, false otherwise
func (userService *UserService) GetByID(userID int) (*User, bool) {
	for _, currentUser := range userService.users {
		if currentUser.ID == userID {
			return &currentUser, true
		}
	}
	return nil, false
}

// Create adds a new user to the service with the provided data
// Parameters:
//   - request: CreateUserRequest containing the user data to create
// Returns:
//   - User: The newly created user with assigned ID and default values
func (userService *UserService) Create(request CreateUserRequest) User {
	newUser := User{
		ID:       userService.nextUserID,
		Name:     request.Name,
		Email:    request.Email,
		Age:      request.Age,
		IsActive: true, // New users are active by default
	}
	userService.users = append(userService.users, newUser)
	userService.nextUserID++
	return newUser
}

// Update modifies an existing user with the provided data
// Only non-empty/non-zero fields in the request will be updated
// Parameters:
//   - userID: The unique identifier of the user to update
//   - request: UpdateUserRequest containing the fields to update
// Returns:
//   - *User: Pointer to the updated user if found, nil otherwise
//   - bool: true if user was found and updated, false otherwise
func (userService *UserService) Update(userID int, request UpdateUserRequest) (*User, bool) {
	for userIndex, currentUser := range userService.users {
		if currentUser.ID == userID {
			// Update only non-empty fields to allow partial updates
			if request.Name != "" {
				userService.users[userIndex].Name = request.Name
			}
			if request.Email != "" {
				userService.users[userIndex].Email = request.Email
			}
			if request.Age > 0 {
				userService.users[userIndex].Age = request.Age
			}
			// For boolean fields, we assume false means "don't update" and true means "update to true"
			// In a production implementation, you might want a more sophisticated approach using pointers
			if request.IsActive {
				userService.users[userIndex].IsActive = request.IsActive
			}
			return &userService.users[userIndex], true
		}
	}
	return nil, false
}

// Delete removes a user from the service by their unique identifier
// Parameters:
//   - userID: The unique identifier of the user to delete
// Returns:
//   - bool: true if user was found and deleted, false otherwise
func (userService *UserService) Delete(userID int) bool {
	for userIndex, currentUser := range userService.users {
		if currentUser.ID == userID {
			// Remove user from slice by combining elements before and after the target index
			userService.users = append(userService.users[:userIndex], userService.users[userIndex+1:]...)
			return true
		}
	}
	return false
}

// UserHandlers contains HTTP handlers for user-related operations
// It acts as the presentation layer, handling HTTP requests and responses
type UserHandlers struct {
	userService *UserService // Service layer for user business logic
}

// NewUserHandlers creates and initializes a new UserHandlers instance
// Parameters:
//   - userService: The user service instance to handle business logic
// Returns:
//   - *UserHandlers: Pointer to the newly created handlers instance
func NewUserHandlers(userService *UserService) *UserHandlers {
	return &UserHandlers{userService: userService}
}

// GetUsers retrieves a paginated list of users with optional filtering
// @Summary      Get users
// @Description  Retrieves a paginated list of users with optional active status filtering
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"  default(10)
// @Param        active    query     bool    false  "Filter by active users"
// @Success      200       {object}  responses.PaginatedResponse
// @Failure      400       {object}  responses.ErrorResponse
// @Router       /api/v1/users [get]
func (handlers *UserHandlers) GetUsers(context *gin.Context) {
	// Parse query parameters with default values
	pageString := context.DefaultQuery("page", "1")
	pageSizeString := context.DefaultQuery("page_size", "10")
	activeString := context.Query("active")

	// Validate and convert page parameter
	pageNumber, parseError := strconv.Atoi(pageString)
	if parseError != nil || pageNumber < 1 {
		responses.BadRequest(context, "Invalid 'page' parameter")
		return
	}

	// Validate and convert page size parameter
	pageSize, parseError := strconv.Atoi(pageSizeString)
	if parseError != nil || pageSize < 1 || pageSize > 100 {
		responses.BadRequest(context, "Invalid 'page_size' parameter (must be between 1-100)")
		return
	}

	// Get all users from service
	allUsers := handlers.userService.GetAll()

	// Apply active status filter if specified
	if activeString != "" {
		isActiveFilter, parseError := strconv.ParseBool(activeString)
		if parseError != nil {
			responses.BadRequest(context, "Invalid 'active' parameter")
			return
		}
		
		var filteredUsers []User
		for _, currentUser := range allUsers {
			if currentUser.IsActive == isActiveFilter {
				filteredUsers = append(filteredUsers, currentUser)
			}
		}
		allUsers = filteredUsers
	}

	// Apply pagination logic
	totalUsers := len(allUsers)
	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	// Handle edge cases for pagination
	if startIndex >= totalUsers {
		allUsers = []User{}
	} else {
		if endIndex > totalUsers {
			endIndex = totalUsers
		}
		allUsers = allUsers[startIndex:endIndex]
	}

	// Return paginated response
	responses.Paginated(context, allUsers, totalUsers, pageNumber, pageSize)
}

// GetUser retrieves a specific user by their unique identifier
// @Summary      Get user by ID
// @Description  Retrieves a specific user by their unique identifier
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (handlers *UserHandlers) GetUser(context *gin.Context) {
	// Extract and validate user ID from URL parameter
	userIDString := context.Param("id")
	userID, parseError := strconv.Atoi(userIDString)
	if parseError != nil {
		responses.BadRequest(context, "Invalid user ID")
		return
	}

	// Retrieve user from service
	foundUser, userExists := handlers.userService.GetByID(userID)
	if !userExists {
		responses.NotFound(context, "User not found")
		return
	}

	// Return successful response with user data
	responses.Success(context, foundUser)
}

// CreateUser creates a new user with the provided data
// @Summary      Create user
// @Description  Creates a new user with the provided data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      201   {object}  User
// @Failure      400   {object}  responses.ValidationErrorResponse
// @Router       /api/v1/users [post]
func (handlers *UserHandlers) CreateUser(context *gin.Context) {
	var createRequest CreateUserRequest
	
	// Parse and bind JSON request body
	if bindError := context.ShouldBindJSON(&createRequest); bindError != nil {
		responses.BadRequest(context, "Invalid data format")
		return
	}

	// Validate request data using validator
	requestValidator := validation.NewValidator()
	if validationError := requestValidator.ValidateStruct(createRequest); validationError != nil {
		validationErrors := validation.FormatValidationErrors(validationError)
		var responseErrors []responses.ResponseValidationError
		
		for _, validationError := range validationErrors {
			responseErrors = append(responseErrors, responses.ResponseValidationError{
				Field:   validationError.Field,
				Message: validationError.Message,
				Value:   validationError.Value,
			})
		}
		responses.ValidationError(context, responseErrors)
		return
	}

	// Create user through service layer
	createdUser := handlers.userService.Create(createRequest)
	responses.Created(context, createdUser)
}

// UpdateUser updates an existing user with the provided data
// @Summary      Update user
// @Description  Updates an existing user with the provided data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "User ID"
// @Param        user  body      UpdateUserRequest  true  "Data to update"
// @Success      200   {object}  User
// @Failure      400   {object}  responses.ErrorResponse
// @Failure      404   {object}  responses.ErrorResponse
// @Router       /api/v1/users/{id} [put]
func (handlers *UserHandlers) UpdateUser(context *gin.Context) {
	// Extract and validate user ID from URL parameter
	userIDString := context.Param("id")
	userID, parseError := strconv.Atoi(userIDString)
	if parseError != nil {
		responses.BadRequest(context, "Invalid user ID")
		return
	}

	var updateRequest UpdateUserRequest
	
	// Parse and bind JSON request body
	if bindError := context.ShouldBindJSON(&updateRequest); bindError != nil {
		responses.BadRequest(context, "Invalid data format")
		return
	}

	// Validate request data using validator
	requestValidator := validation.NewValidator()
	if validationError := requestValidator.ValidateStruct(updateRequest); validationError != nil {
		validationErrors := validation.FormatValidationErrors(validationError)
		var responseErrors []responses.ResponseValidationError
		
		for _, validationError := range validationErrors {
			responseErrors = append(responseErrors, responses.ResponseValidationError{
				Field:   validationError.Field,
				Message: validationError.Message,
				Value:   validationError.Value,
			})
		}
		responses.ValidationError(context, responseErrors)
		return
	}

	// Update user through service layer
	updatedUser, userExists := handlers.userService.Update(userID, updateRequest)
	if !userExists {
		responses.NotFound(context, "User not found")
		return
	}

	// Return successful response with updated user data
	responses.SuccessWithMessage(context, updatedUser, "User updated successfully")
}

// DeleteUser removes a user by their unique identifier
// @Summary      Delete user
// @Description  Removes a user by their unique identifier
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "User deleted successfully"
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /api/v1/users/{id} [delete]
func (handlers *UserHandlers) DeleteUser(context *gin.Context) {
	// Extract and validate user ID from URL parameter
	userIDString := context.Param("id")
	userID, parseError := strconv.Atoi(userIDString)
	if parseError != nil {
		responses.BadRequest(context, "Invalid user ID")
		return
	}

	// Delete user through service layer
	if !handlers.userService.Delete(userID) {
		responses.NotFound(context, "User not found")
		return
	}

	// Return successful no-content response
	responses.NoContent(context)
}

// @title           GoAPI Advanced Example
// @version         1.0
// @description     Advanced example of GoAPI with all FastAPI-like functionalities
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https
// @securityDefinitions.basic  BasicAuth

// main is the entry point of the application
// It initializes the API server with all necessary configurations, middleware, and routes
func main() {
	// Create API configuration with default settings
	apiConfiguration := goapi.DefaultConfig()
	apiConfiguration.Title = "GoAPI Advanced Example"
	apiConfiguration.Description = "Advanced example of GoAPI with FastAPI-like functionalities"
	apiConfiguration.Version = "1.0.0"

	// Create new GoAPI instance with the configuration
	apiInstance := goapi.New(apiConfiguration)

	// Configure additional middleware for enhanced functionality
	apiInstance.AddRateLimit(middleware.RateLimitConfig{
		RequestsPerMinute: 100, // Allow 100 requests per minute per client
		BurstSize:         10,  // Allow burst of 10 requests
	})

	// Register dependencies for dependency injection
	userServiceInstance := NewUserService()
	apiInstance.RegisterSingletonDependency(func(context *gin.Context) (interface{}, error) {
		return userServiceInstance, nil
	}, (*UserService)(nil))

	// Create handlers with injected dependencies
	userHandlersInstance := NewUserHandlers(userServiceInstance)

	// Define API routes with proper grouping and documentation
	apiV1Group := apiInstance.Group("/api/v1")
	{
		usersGroup := apiV1Group.Group("/users")
		{
			// GET /api/v1/users - Retrieve paginated list of users
			usersGroup.GET("", userHandlersInstance.GetUsers,
				goapi.WithSummary("Get users"),
				goapi.WithDescription("Retrieves a paginated list of users with optional filtering"),
				goapi.WithTags("users"),
				goapi.WithQueryParameter("page", "integer", "Page number", false),
				goapi.WithQueryParameter("page_size", "integer", "Page size", false),
				goapi.WithQueryParameter("active", "boolean", "Filter by active users", false))
			
			// GET /api/v1/users/:id - Retrieve specific user by ID
			usersGroup.GET("/:id", userHandlersInstance.GetUser,
				goapi.WithSummary("Get user by ID"),
				goapi.WithDescription("Retrieves a specific user by their unique identifier"),
				goapi.WithTags("users"),
				goapi.WithPathParameter("id", "integer", "User ID"))
			
			// POST /api/v1/users - Create new user
			usersGroup.POST("", userHandlersInstance.CreateUser,
				goapi.WithSummary("Create user"),
				goapi.WithDescription("Creates a new user with the provided data"),
				goapi.WithTags("users"),
				goapi.WithRequestBody(CreateUserRequest{
					Name:  "John Doe",
					Email: "john.doe@example.com",
					Age:   25,
				}, "User data for creation"))
			
			// PUT /api/v1/users/:id - Update existing user
			usersGroup.PUT("/:id", userHandlersInstance.UpdateUser,
				goapi.WithSummary("Update user"),
				goapi.WithDescription("Updates an existing user with the provided data"),
				goapi.WithTags("users"),
				goapi.WithPathParameter("id", "integer", "User ID"),
				goapi.WithRequestBody(UpdateUserRequest{
					Name:     "John Doe Updated",
					Email:    "john.updated@example.com",
					Age:      26,
					IsActive: true,
				}, "User data to update"))
			
			// DELETE /api/v1/users/:id - Delete user
			usersGroup.DELETE("/:id", userHandlersInstance.DeleteUser,
				goapi.WithSummary("Delete user"),
				goapi.WithDescription("Removes a user by their unique identifier"),
				goapi.WithTags("users"),
				goapi.WithPathParameter("id", "integer", "User ID"))
		}
	}

	// Example route with dependency injection for statistics
	apiInstance.GET("/api/v1/stats", func(context *gin.Context) {
		// Resolve dependency from container
		var resolvedService *UserService
		if resolutionError := apiInstance.GetDependencyContainer().Resolve(context, &resolvedService); resolutionError != nil {
			responses.InternalServerError(context, "Error resolving dependencies")
			return
		}

		// Calculate user statistics
		allUsers := resolvedService.GetAll()
		activeUsersCount := 0
		for _, currentUser := range allUsers {
			if currentUser.IsActive {
				activeUsersCount++
			}
		}

		// Prepare statistics response
		statisticsData := gin.H{
			"total_users":    len(allUsers),
			"active_users":   activeUsersCount,
			"inactive_users": len(allUsers) - activeUsersCount,
		}

		responses.Success(context, statisticsData)
	}, goapi.WithTags("stats"), goapi.WithSummary("Get statistics"))

	// Start the API server
	log.Println("Starting GoAPI Advanced Example...")
	log.Println("Implemented features:")
	log.Println("✓ Automatic data validation (like Pydantic)")
	log.Println("✓ Dependency Injection")
	log.Println("✓ Customizable middleware")
	log.Println("✓ Standardized responses")
	log.Println("✓ Automatic pagination")
	log.Println("✓ Swagger/OpenAPI documentation")
	log.Println("✓ Rate limiting")
	log.Println("✓ CORS configuration")
	log.Println("✓ Centralized error handling")
	
	if serverError := apiInstance.Run(":8080"); serverError != nil {
		log.Fatal("Error starting the API server:", serverError)
	}
}
