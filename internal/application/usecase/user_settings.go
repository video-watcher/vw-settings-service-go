package usecase

import (
	"context"
	"errors"
	"settings-service-go/internal/application/dto"
	"settings-service-go/internal/domain/entity"
	"settings-service-go/internal/domain/repository"
)

// UserSettingsUseCase handles user settings business logic
type UserSettingsUseCase struct {
	repo repository.UserSettingsRepository
}

// NewUserSettingsUseCase creates a new UserSettingsUseCase
func NewUserSettingsUseCase(repo repository.UserSettingsRepository) *UserSettingsUseCase {
	return &UserSettingsUseCase{repo: repo}
}

// GetSettings retrieves user settings by user ID
func (uc *UserSettingsUseCase) GetSettings(ctx context.Context, userID string) (*dto.UserSettingsResponse, error) {
	settings, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create default settings if not found
	if settings == nil {
		settings = entity.NewUserSettings(userID)
		if err := uc.repo.Create(ctx, settings); err != nil {
			return nil, err
		}
	}

	return &dto.UserSettingsResponse{
		ID:                settings.ID.String(),
		UserID:            settings.UserID,
		NotificationEmail: settings.NotificationEmail,
		NotificationPush:  settings.NotificationPush,
		Theme:             settings.Theme,
		Language:          settings.Language,
		CreatedAt:         settings.CreatedAt,
		UpdatedAt:         settings.UpdatedAt,
	}, nil
}

// UpdateSettings updates user settings
func (uc *UserSettingsUseCase) UpdateSettings(ctx context.Context, userID string, req *dto.UserSettingsUpdateRequest) (*dto.UserSettingsResponse, error) {
	settings, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, errors.New("user settings not found")
	}

	// Update fields
	settings.Update(req.NotificationEmail, req.NotificationPush, req.Theme, req.Language)

	// Validate
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, settings); err != nil {
		return nil, err
	}

	return &dto.UserSettingsResponse{
		ID:                settings.ID.String(),
		UserID:            settings.UserID,
		NotificationEmail: settings.NotificationEmail,
		NotificationPush:  settings.NotificationPush,
		Theme:             settings.Theme,
		Language:          settings.Language,
		CreatedAt:         settings.CreatedAt,
		UpdatedAt:         settings.UpdatedAt,
	}, nil
}
