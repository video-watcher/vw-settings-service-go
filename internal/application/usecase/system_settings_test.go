package usecase

import (
	"context"
	"settings-service-go/internal/application/dto"
	"settings-service-go/internal/domain/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSystemSettingsRepository struct {
	mock.Mock
}

func (m *MockSystemSettingsRepository) GetAll(ctx context.Context, skip, limit int) ([]*entity.SystemSettings, error) {
	args := m.Called(ctx, skip, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.SystemSettings), args.Error(1)
}

func (m *MockSystemSettingsRepository) GetByKey(ctx context.Context, key string) (*entity.SystemSettings, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SystemSettings), args.Error(1)
}

func (m *MockSystemSettingsRepository) Set(ctx context.Context, settings *entity.SystemSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockSystemSettingsRepository) DeleteByKey(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockSystemSettingsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetAllSettings_Success(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	desc := "Test setting"
	expectedSettings := []*entity.SystemSettings{
		{
			ID:          uuid.New(),
			Key:         "test_key",
			Value:       map[string]interface{}{"enabled": true},
			Description: &desc,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetAll", mock.Anything, 0, 10).Return(expectedSettings, nil)

	result, err := useCase.GetAllSettings(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test_key", result[0].Key)
	mockRepo.AssertExpectations(t)
}

func TestGetSetting_Success(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	desc := "Test setting"
	expectedSetting := &entity.SystemSettings{
		ID:          uuid.New(),
		Key:         "test_key",
		Value:       map[string]interface{}{"enabled": true},
		Description: &desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByKey", mock.Anything, "test_key").Return(expectedSetting, nil)

	result, err := useCase.GetSetting(context.Background(), "test_key")

	assert.NoError(t, err)
	assert.Equal(t, "test_key", result.Key)
	assert.Equal(t, true, result.Value["enabled"])
	mockRepo.AssertExpectations(t)
}

func TestGetSetting_NotFound(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	mockRepo.On("GetByKey", mock.Anything, "nonexistent").Return(nil, nil)

	result, err := useCase.GetSetting(context.Background(), "nonexistent")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestSetSetting_Create(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	mockRepo.On("GetByKey", mock.Anything, "new_key").Return(nil, nil)
	mockRepo.On("Set", mock.Anything, mock.AnythingOfType("*entity.SystemSettings")).Return(nil)

	desc := "New setting"
	req := &dto.SystemSettingCreateRequest{
		Value:       map[string]interface{}{"enabled": true},
		Description: &desc,
	}

	result, err := useCase.SetSetting(context.Background(), "new_key", req)

	assert.NoError(t, err)
	assert.Equal(t, "new_key", result.Key)
	assert.Equal(t, true, result.Value["enabled"])
	mockRepo.AssertExpectations(t)
}

func TestSetSetting_Update(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	desc := "Existing setting"
	existingSetting := &entity.SystemSettings{
		ID:          uuid.New(),
		Key:         "existing_key",
		Value:       map[string]interface{}{"enabled": false},
		Description: &desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByKey", mock.Anything, "existing_key").Return(existingSetting, nil)
	mockRepo.On("Set", mock.Anything, mock.AnythingOfType("*entity.SystemSettings")).Return(nil)

	newDesc := "Updated setting"
	req := &dto.SystemSettingCreateRequest{
		Value:       map[string]interface{}{"enabled": true},
		Description: &newDesc,
	}

	result, err := useCase.SetSetting(context.Background(), "existing_key", req)

	assert.NoError(t, err)
	assert.Equal(t, "existing_key", result.Key)
	assert.Equal(t, true, result.Value["enabled"])
	mockRepo.AssertExpectations(t)
}

func TestDeleteSetting_Success(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	desc := "Test setting"
	existingSetting := &entity.SystemSettings{
		ID:          uuid.New(),
		Key:         "test_key",
		Value:       map[string]interface{}{"enabled": true},
		Description: &desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByKey", mock.Anything, "test_key").Return(existingSetting, nil)
	mockRepo.On("DeleteByKey", mock.Anything, "test_key").Return(nil)

	err := useCase.DeleteSetting(context.Background(), "test_key")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSetting_NotFound(t *testing.T) {
	mockRepo := new(MockSystemSettingsRepository)
	useCase := NewSystemSettingsUseCase(mockRepo)

	mockRepo.On("GetByKey", mock.Anything, "nonexistent").Return(nil, nil)

	err := useCase.DeleteSetting(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setting not found")
	mockRepo.AssertExpectations(t)
}
