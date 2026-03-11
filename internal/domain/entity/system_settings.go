package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SystemSettings represents system-wide configuration stored as key-value pairs
type SystemSettings struct {
	ID          uuid.UUID              `json:"id"`
	Key         string                 `json:"key"`
	Value       map[string]interface{} `json:"value"` // JSONB in PostgreSQL
	Description *string                `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NewSystemSettings creates a new SystemSettings
func NewSystemSettings(key string, value map[string]interface{}, description *string) *SystemSettings {
	now := time.Now()
	return &SystemSettings{
		ID:          uuid.New(),
		Key:         key,
		Value:       value,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate validates the SystemSettings entity
func (s *SystemSettings) Validate() error {
	if s.Key == "" {
		return errors.New("key is required")
	}

	if len(s.Key) > 100 {
		return errors.New("key too long (max 100 characters)")
	}

	if s.Value == nil {
		return errors.New("value is required")
	}

	return nil
}

// Update updates system setting value and description
func (s *SystemSettings) Update(value map[string]interface{}, description *string) {
	if value != nil {
		s.Value = value
	}
	if description != nil {
		s.Description = description
	}
	s.UpdatedAt = time.Now()
}
