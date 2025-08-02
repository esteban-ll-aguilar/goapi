package advanced

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi"
	"github.com/esteban-ll-aguilar/goapi/goapi/middleware"
	"github.com/esteban-ll-aguilar/goapi/goapi/responses"
	"github.com/esteban-ll-aguilar/goapi/goapi/validation"
)

// User representa un usuario del sistema
type User struct {
	ID       int    `json:"id" validate:"required,min=1" example:"1"`
	Name     string `json:"name" validate:"required,min=2,max=50" example:"Juan Pérez"`
	Email    string `json:"email" validate:"required,email" example:"juan@example.com"`
	Age      int    `json:"age" validate:"required,min=18,max=120" example:"25"`
	IsActive bool   `json:"is_active" example:"true"`
}

// CreateUserRequest representa la petición para crear un usuario
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50" example:"Juan Pérez"`
	Email string `json:"email" validate:"required,email" example:"juan@example.com"`
	Age   int    `json:"age" validate:"required,min=18,max=120" example:"25"`
}

// UpdateUserRequest representa la petición para actualizar un usuario
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=50" example:"Juan Pérez"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"juan@example.com"`
	Age      *int    `json:"age,omitempty" validate:"omitempty,min=18,max=120" example:"25"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

// UserService simula un servicio de usuarios
type UserService struct {
	users  []User
	nextID int
}

// NewUserService crea un nuevo servicio de usuarios
func NewUserService() *UserService {
	return &UserService{
		users: []User{
			{ID: 1, Name: "Juan Pérez", Email: "juan@example.com", Age: 25, IsActive: true},
			{ID: 2, Name: "María García", Email: "maria@example.com", Age: 30, IsActive: true},
		},
		nextID: 3,
	}
}

// GetAll devuelve todos los usuarios
func (s *UserService) GetAll() []User {
	return s.users
}

// GetByID devuelve un usuario por ID
func (s *UserService) GetByID(id int) (*User, bool) {
	for _, user := range s.users {
		if user.ID == id {
			return &user, true
		}
	}
	return nil, false
}

// Create crea un nuevo usuario
func (s *UserService) Create(req CreateUserRequest) User {
	user := User{
		ID:       s.nextID,
		Name:     req.Name,
		Email:    req.Email,
		Age:      req.Age,
		IsActive: true,
	}
	s.users = append(s.users, user)
	s.nextID++
	return user
}

// Update actualiza un usuario existente
func (s *UserService) Update(id int, req UpdateUserRequest) (*User, bool) {
	for i, user := range s.users {
		if user.ID == id {
			if req.Name != nil {
				s.users[i].Name = *req.Name
			}
			if req.Email != nil {
				s.users[i].Email = *req.Email
			}
			if req.Age != nil {
				s.users[i].Age = *req.Age
			}
			if req.IsActive != nil {
				s.users[i].IsActive = *req.IsActive
			}
			return &s.users[i], true
		}
	}
	return nil, false
}

// Delete elimina un usuario
func (s *UserService) Delete(id int) bool {
	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return true
		}
	}
	return false
}

// UserHandlers contiene los handlers para usuarios
type UserHandlers struct {
	service *UserService
}

// NewUserHandlers crea nuevos handlers de usuario
func NewUserHandlers(service *UserService) *UserHandlers {
	return &UserHandlers{service: service}
}

// GetUsers obtiene todos los usuarios con paginación
// @Summary      Obtener usuarios
// @Description  Obtiene una lista paginada de usuarios
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Número de página"  default(1)
// @Param        page_size query     int     false  "Tamaño de página"  default(10)
// @Param        active    query     bool    false  "Filtrar por usuarios activos"
// @Success      200       {object}  responses.PaginatedResponse
// @Failure      400       {object}  responses.ErrorResponse
// @Router       /api/v1/users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
	// Parsear parámetros de consulta
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	activeStr := c.Query("active")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		responses.BadRequest(c, "Parámetro 'page' inválido")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		responses.BadRequest(c, "Parámetro 'page_size' inválido (1-100)")
		return
	}

	users := h.service.GetAll()

	// Filtrar por activos si se especifica
	if activeStr != "" {
		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			responses.BadRequest(c, "Parámetro 'active' inválido")
			return
		}
		
		var filteredUsers []User
		for _, user := range users {
			if user.IsActive == active {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	// Aplicar paginación
	total := len(users)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		users = []User{}
	} else {
		if end > total {
			end = total
		}
		users = users[start:end]
	}

	responses.Paginated(c, users, total, page, pageSize)
}

// GetUser obtiene un usuario por ID
// @Summary      Obtener usuario por ID
// @Description  Obtiene un usuario específico por su ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del usuario"
// @Success      200  {object}  User
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (h *UserHandlers) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return
	}

	user, found := h.service.GetByID(id)
	if !found {
		responses.NotFound(c, "Usuario no encontrado")
		return
	}

	responses.Success(c, user)
}

// CreateUser crea un nuevo usuario
// @Summary      Crear usuario
// @Description  Crea un nuevo usuario con los datos proporcionados
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "Datos del usuario"
// @Success      201   {object}  User
// @Failure      400   {object}  responses.ValidationErrorResponse
// @Router       /api/v1/users [post]
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Formato de datos inválido")
		return
	}

	// Validar datos
	validator := validation.NewValidator()
	if err := validator.ValidateStruct(req); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		var responseErrors []responses.ResponseValidationError
		for _, ve := range validationErrors {
			responseErrors = append(responseErrors, responses.ResponseValidationError{
				Field:   ve.Field,
				Message: ve.Message,
				Value:   ve.Value,
			})
		}
		responses.ValidationError(c, responseErrors)
		return
	}

	user := h.service.Create(req)
	responses.Created(c, user)
}

// UpdateUser actualiza un usuario existente
// @Summary      Actualizar usuario
// @Description  Actualiza un usuario existente con los datos proporcionados
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "ID del usuario"
// @Param        user  body      UpdateUserRequest  true  "Datos a actualizar"
// @Success      200   {object}  User
// @Failure      400   {object}  responses.ErrorResponse
// @Failure      404   {object}  responses.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Formato de datos inválido")
		return
	}

	// Validar datos
	validator := validation.NewValidator()
	if err := validator.ValidateStruct(req); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		var responseErrors []responses.ResponseValidationError
		for _, ve := range validationErrors {
			responseErrors = append(responseErrors, responses.ResponseValidationError{
				Field:   ve.Field,
				Message: ve.Message,
				Value:   ve.Value,
			})
		}
		responses.ValidationError(c, responseErrors)
		return
	}

	user, found := h.service.Update(id, req)
	if !found {
		responses.NotFound(c, "Usuario no encontrado")
		return
	}

	responses.SuccessWithMessage(c, user, "Usuario actualizado correctamente")
}

// DeleteUser elimina un usuario
// @Summary      Eliminar usuario
// @Description  Elimina un usuario por su ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del usuario"
// @Success      204  "Usuario eliminado correctamente"
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return
	}

	if !h.service.Delete(id) {
		responses.NotFound(c, "Usuario no encontrado")
		return
	}

	responses.NoContent(c)
}

// @title           GoAPI Advanced Example
// @version         1.0
// @description     Ejemplo avanzado de GoAPI con todas las funcionalidades de FastAPI
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

func main() {
	// Crear configuración de la API
	config := goapi.DefaultConfig()
	config.Title = "GoAPI Advanced Example"
	config.Description = "Ejemplo avanzado de GoAPI con funcionalidades similares a FastAPI"
	config.Version = "1.0.0"
	config.BasePath = "/api/v1"

	// Crear instancia de GoAPI
	api := goapi.New(config)

	// Configurar middleware adicional
	api.AddRateLimit(middleware.RateLimitConfig{
		RequestsPerMinute: 100,
		BurstSize:         10,
	})

	// Registrar dependencias
	userService := NewUserService()
	api.RegisterSingletonDependency(func(c *gin.Context) (interface{}, error) {
		return userService, nil
	}, (*UserService)(nil))

	// Crear handlers
	userHandlers := NewUserHandlers(userService)

	// Definir rutas de la API
	v1 := api.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userHandlers.GetUsers)
			users.GET("/:id", userHandlers.GetUser)
			users.POST("", userHandlers.CreateUser)
			users.PUT("/:id", userHandlers.UpdateUser)
			users.DELETE("/:id", userHandlers.DeleteUser)
		}
	}

	// Ruta de ejemplo con dependency injection
	api.GET("/api/v1/stats", func(c *gin.Context) {
		// Resolver dependencia
		var service *UserService
		if err := api.GetDependencyContainer().Resolve(c, &service); err != nil {
			responses.InternalServerError(c, "Error resolviendo dependencias")
			return
		}

		users := service.GetAll()
		activeUsers := 0
		for _, user := range users {
			if user.IsActive {
				activeUsers++
			}
		}

		stats := gin.H{
			"total_users":  len(users),
			"active_users": activeUsers,
			"inactive_users": len(users) - activeUsers,
		}

		responses.Success(c, stats)
	}, goapi.WithTags("stats"), goapi.WithSummary("Obtener estadísticas"))

	// Ejecutar la API
	log.Println("Iniciando GoAPI Advanced Example...")
	log.Println("Funcionalidades implementadas:")
	log.Println("✓ Validación automática de datos (como Pydantic)")
	log.Println("✓ Dependency Injection")
	log.Println("✓ Middleware personalizable")
	log.Println("✓ Respuestas estandarizadas")
	log.Println("✓ Paginación automática")
	log.Println("✓ Documentación Swagger/OpenAPI")
	log.Println("✓ Rate limiting")
	log.Println("✓ CORS configuración")
	log.Println("✓ Manejo de errores centralizado")
	
	if err := api.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar la API:", err)
	}
}
