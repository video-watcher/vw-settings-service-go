package repository

import (
	"context"
	"settings-service-go/internal/domain/entity"

	"github.com/google/uuid"
)

// SystemSettingsRepository defines the interface for system settings data access
type SystemSettingsRepository interface {
	// GetAll retrieves all system settings with pagination
	GetAll(ctx context.Context, skip, limit int) ([]*entity.SystemSettings, error)

	// GetByKey retrieves system setting by key
	GetByKey(ctx context.Context, key string) (*entity.SystemSettings, error)

	// Set creates or updates a system setting
	Set(ctx context.Context, settings *entity.SystemSettings) error

	// DeleteByKey deletes system setting by key
	DeleteByKey(ctx context.Context, key string) error

	// Delete deletes system setting by ID
	Delete(ctx context.Context, id uuid.UUID) error
}
