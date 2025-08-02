package basic

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/esteban-ll-aguilar/goapi/goapi"
	"github.com/esteban-ll-aguilar/goapi/goapi/models"
)

// @title           GoAPI - FastAPI Style for Go
// @version         1.0
// @description     A FastAPI-like API built with Go
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http
// @securityDefinitions.basic  BasicAuth

// @tag.name  items
// @tag.description  Operaciones con items

// ItemHandler maneja las peticiones relacionadas con items
type ItemHandler struct {
	store *models.ItemStore
}

// NewItemHandler crea un nuevo manejador de items
func NewItemHandler(store *models.ItemStore) *ItemHandler {
	return &ItemHandler{
		store: store,
	}
}

// GetItems devuelve todos los items
// @Summary      Obtener todos los items
// @Description  Devuelve una lista de todos los items
// @Tags         items
// @Produce      json
// @Success      200  {array}   models.Item
// @Router       /api/v1/items [get]
func (h *ItemHandler) GetItems(c *gin.Context) {
	items := h.store.GetAll()
	c.JSON(http.StatusOK, items)
}

// GetItemByID devuelve un item por su ID
// @Summary      Obtener un item por ID
// @Description  Devuelve un item específico por su ID
// @Tags         items
// @Produce      json
// @Param        id   path      string  true  "ID del item"
// @Success      200  {object}  models.Item
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/items/{id} [get]
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	id := c.Param("id")

	item, found := h.store.GetByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item no encontrado"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem crea un nuevo item
// @Summary      Crear un nuevo item
// @Description  Crea un nuevo item con los datos proporcionados
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      models.Item  true  "Datos del nuevo item"
// @Success      201   {object}  models.Item
// @Router       /api/v1/items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var newItem models.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.store.Create(newItem)
	c.JSON(http.StatusCreated, newItem)
}

func main() {
	// Crear configuración de la API
	config := goapi.DefaultConfig()
	config.Title = "GoAPI Demo"
	config.Description = "Demostración de GoAPI - Un framework al estilo FastAPI para Go"
	config.BasePath = "/api/v1"

	// Crear instancia de GoAPI
	api := goapi.New(config)

	// Crear almacén de datos y handlers
	store := models.NewItemStore()
	itemHandler := NewItemHandler(store)

	// Definir rutas de la API
	v1 := api.Group("/api/v1")
	{
		items := v1.Group("/items")
		{
			items.GET("", itemHandler.GetItems)
			items.GET("/:id", itemHandler.GetItemByID)
			items.POST("", itemHandler.CreateItem)
		}
	}

	// Ejecutar la API
	log.Println("Iniciando GoAPI...")
	if err := api.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar la API:", err)
	}
}
