// Package models provides data models for GoAPI
package models

import (
	"encoding/json"
	"fmt"
)

// Model is an interface that all models must implement
type Model interface {
	Validate() error
}

// BaseModel is a base model that provides common functionality
type BaseModel struct{}

// Validate implements the Model interface
func (m *BaseModel) Validate() error {
	return nil
}

// ToJSON converts a model to JSON
func ToJSON(model interface{}) (string, error) {
	bytes, err := json.Marshal(model)
	if err != nil {
		return "", fmt.Errorf("error converting model to JSON: %w", err)
	}
	return string(bytes), nil
}

// FromJSON converts JSON to a model
func FromJSON(data string, model interface{}) error {
	if err := json.Unmarshal([]byte(data), model); err != nil {
		return fmt.Errorf("error converting JSON to model: %w", err)
	}
	return nil
}
