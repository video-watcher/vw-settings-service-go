package usecase

import (
	"context"
	"errors"
	"settings-service-go/internal/application/dto"
	"settings-service-go/internal/domain/entity"
	"settings-service-go/internal/domain/repository"
)

// SystemSettingsUseCase handles system settings business logic
type SystemSettingsUseCase struct {
	repo repository.SystemSettingsRepository
}

// NewSystemSettingsUseCase creates a new SystemSettingsUseCase
func NewSystemSettingsUseCase(repo repository.SystemSettingsRepository) *SystemSettingsUseCase {
	return &SystemSettingsUseCase{repo: repo}
}

// GetAllSettings retrieves all system settings with pagination
func (uc *SystemSettingsUseCase) GetAllSettings(ctx context.Context, skip, limit int) ([]*dto.SystemSettingResponse, error) {
	settings, err := uc.repo.GetAll(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.SystemSettingResponse, len(settings))
	for i, s := range settings {
		responses[i] = &dto.SystemSettingResponse{
			ID:          s.ID.String(),
			Key:         s.Key,
			Value:       s.Value,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
	}

	return responses, nil
}

// GetSetting retrieves system setting by key
func (uc *SystemSettingsUseCase) GetSetting(ctx context.Context, key string) (*dto.SystemSettingResponse, error) {
	setting, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, nil
	}

	return &dto.SystemSettingResponse{
		ID:          setting.ID.String(),
		Key:         setting.Key,
		Value:       setting.Value,
		Description: setting.Description,
		CreatedAt:   setting.CreatedAt,
		UpdatedAt:   setting.UpdatedAt,
	}, nil
}

// SetSetting creates or updates a system setting
func (uc *SystemSettingsUseCase) SetSetting(ctx context.Context, key string, req *dto.SystemSettingCreateRequest) (*dto.SystemSettingResponse, error) {
	// Check if exists
	existing, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	var setting *entity.SystemSettings
	if existing != nil {
		// Update existing
		existing.Update(req.Value, req.Description)
		setting = existing
	} else {
		// Create new
		setting = entity.NewSystemSettings(key, req.Value, req.Description)
	}

	// Validate
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Set(ctx, setting); err != nil {
		return nil, err
	}

	return &dto.SystemSettingResponse{
		ID:          setting.ID.String(),
		Key:         setting.Key,
		Value:       setting.Value,
		Description: setting.Description,
		CreatedAt:   setting.CreatedAt,
		UpdatedAt:   setting.UpdatedAt,
	}, nil
}

// DeleteSetting deletes system setting by key
func (uc *SystemSettingsUseCase) DeleteSetting(ctx context.Context, key string) error {
	setting, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return err
	}

	if setting == nil {
		return errors.New("setting not found")
	}

	return uc.repo.DeleteByKey(ctx, key)
}
