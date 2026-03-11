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

type MockUserSettingsRepository struct {
	mock.Mock
}

func (m *MockUserSettingsRepository) GetByUserID(ctx context.Context, userID string) (*entity.UserSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSettings), args.Error(1)
}

func (m *MockUserSettingsRepository) Create(ctx context.Context, settings *entity.UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockUserSettingsRepository) Update(ctx context.Context, settings *entity.UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockUserSettingsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetSettings_ExistingUser(t *testing.T) {
	mockRepo := new(MockUserSettingsRepository)
	useCase := NewUserSettingsUseCase(mockRepo)

	expectedSettings := &entity.UserSettings{
		ID:                uuid.New(),
		UserID:            "user-123",
		NotificationEmail: true,
		NotificationPush:  false,
		Theme:             "dark",
		Language:          "vi",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mockRepo.On("GetByUserID", mock.Anything, "user-123").Return(expectedSettings, nil)

	result, err := useCase.GetSettings(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedSettings.UserID, result.UserID)
	assert.Equal(t, expectedSettings.Theme, result.Theme)
	assert.Equal(t, expectedSettings.Language, result.Language)
	assert.False(t, result.NotificationPush)
	assert.True(t, result.NotificationEmail)
	mockRepo.AssertExpectations(t)
}

func TestGetSettings_NewUser(t *testing.T) {
	mockRepo := new(MockUserSettingsRepository)
	useCase := NewUserSettingsUseCase(mockRepo)

	// Repository returns nil, nil when user not found (not an error)
	mockRepo.On("GetByUserID", mock.Anything, "new-user").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSettings")).Return(nil)

	result, err := useCase.GetSettings(context.Background(), "new-user")

	assert.NoError(t, err)
	assert.Equal(t, "new-user", result.UserID)
	assert.Equal(t, "light", result.Theme)
	assert.Equal(t, "en", result.Language)
	assert.True(t, result.NotificationEmail)
	assert.True(t, result.NotificationPush)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSettings_Success(t *testing.T) {
	mockRepo := new(MockUserSettingsRepository)
	useCase := NewUserSettingsUseCase(mockRepo)

	existingSettings := &entity.UserSettings{
		ID:                uuid.New(),
		UserID:            "user-123",
		NotificationEmail: true,
		NotificationPush:  true,
		Theme:             "light",
		Language:          "en",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mockRepo.On("GetByUserID", mock.Anything, "user-123").Return(existingSettings, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.UserSettings")).Return(nil)

	theme := "dark"
	language := "vi"
	notificationPush := false
	req := &dto.UserSettingsUpdateRequest{
		Theme:            &theme,
		Language:         &language,
		NotificationPush: &notificationPush,
	}

	result, err := useCase.UpdateSettings(context.Background(), "user-123", req)

	assert.NoError(t, err)
	assert.Equal(t, "dark", result.Theme)
	assert.Equal(t, "vi", result.Language)
	assert.False(t, result.NotificationPush)
	assert.True(t, result.NotificationEmail)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSettings_InvalidTheme(t *testing.T) {
	mockRepo := new(MockUserSettingsRepository)
	useCase := NewUserSettingsUseCase(mockRepo)

	existingSettings := &entity.UserSettings{
		ID:                uuid.New(),
		UserID:            "user-123",
		NotificationEmail: true,
		NotificationPush:  true,
		Theme:             "light",
		Language:          "en",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mockRepo.On("GetByUserID", mock.Anything, "user-123").Return(existingSettings, nil)

	invalidTheme := "invalid-theme"
	req := &dto.UserSettingsUpdateRequest{
		Theme: &invalidTheme,
	}

	_, err := useCase.UpdateSettings(context.Background(), "user-123", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "theme must be 'light' or 'dark'")
	mockRepo.AssertExpectations(t)
}

func TestUpdateSettings_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserSettingsRepository)
	useCase := NewUserSettingsUseCase(mockRepo)

	mockRepo.On("GetByUserID", mock.Anything, "nonexistent").Return(nil, nil)

	theme := "dark"
	req := &dto.UserSettingsUpdateRequest{
		Theme: &theme,
	}

	_, err := useCase.UpdateSettings(context.Background(), "nonexistent", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user settings not found")
	mockRepo.AssertExpectations(t)
}
