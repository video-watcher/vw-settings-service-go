package repository

import (
	"context"
	"settings-service-go/internal/domain/entity"

	"github.com/google/uuid"
)

// UserSettingsRepository defines the interface for user settings data access
type UserSettingsRepository interface {
	// GetByUserID retrieves user settings by user ID
	GetByUserID(ctx context.Context, userID string) (*entity.UserSettings, error)

	// Create creates new user settings
	Create(ctx context.Context, settings *entity.UserSettings) error

	// Update updates existing user settings
	Update(ctx context.Context, settings *entity.UserSettings) error

	// Delete deletes user settings by ID
	Delete(ctx context.Context, id uuid.UUID) error
}
