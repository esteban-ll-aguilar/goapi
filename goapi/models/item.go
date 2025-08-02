package models

import (
	"errors"
)

// Item representa un elemento de datos
type Item struct {
	BaseModel
	ID     string  `json:"id" example:"1"`
	Name   string  `json:"name" example:"Test Item"`
	Price  float64 `json:"price" example:"10.5"`
	IsUsed bool    `json:"is_used" example:"false"`
}

// Validate valida los campos del item
func (i *Item) Validate() error {
	if i.ID == "" {
		return errors.New("el ID no puede estar vacío")
	}
	if i.Name == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if i.Price < 0 {
		return errors.New("el precio no puede ser negativo")
	}
	return nil
}

// ItemStore almacena los items en memoria para este ejemplo
type ItemStore struct {
	Items []Item
}

// NewItemStore crea un nuevo almacén de items con datos de ejemplo
func NewItemStore() *ItemStore {
	return &ItemStore{
		Items: []Item{
			{ID: "1", Name: "Item 1", Price: 10.5, IsUsed: false},
			{ID: "2", Name: "Item 2", Price: 20.0, IsUsed: true},
		},
	}
}

// GetAll devuelve todos los items
func (s *ItemStore) GetAll() []Item {
	return s.Items
}

// GetByID devuelve un item por su ID
func (s *ItemStore) GetByID(id string) (Item, bool) {
	for _, item := range s.Items {
		if item.ID == id {
			return item, true
		}
	}
	return Item{}, false
}

// Create añade un nuevo item al almacén
func (s *ItemStore) Create(item Item) {
	s.Items = append(s.Items, item)
}
