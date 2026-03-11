package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserSettings(t *testing.T) {
	userID := "user-123"
	settings := NewUserSettings(userID)

	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	assert.True(t, settings.NotificationEmail)
	assert.True(t, settings.NotificationPush)
	assert.Equal(t, "light", settings.Theme)
	assert.Equal(t, "en", settings.Language)
	assert.NotEqual(t, "", settings.ID.String())
}

func TestUserSettings_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *UserSettings
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid settings",
			setup: func() *UserSettings {
				return NewUserSettings("user-123")
			},
			wantErr: false,
		},
		{
			name: "empty user_id",
			setup: func() *UserSettings {
				s := NewUserSettings("")
				return s
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "invalid theme",
			setup: func() *UserSettings {
				s := NewUserSettings("user-123")
				s.Theme = "invalid"
				return s
			},
			wantErr: true,
			errMsg:  "theme must be 'light' or 'dark'",
		},
		{
			name: "invalid language length",
			setup: func() *UserSettings {
				s := NewUserSettings("user-123")
				s.Language = "eng"
				return s
			},
			wantErr: true,
			errMsg:  "language must be 2-character code",
		},
		{
			name: "valid dark theme",
			setup: func() *UserSettings {
				s := NewUserSettings("user-123")
				s.Theme = "dark"
				return s
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := tt.setup()
			err := settings.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSettings_Update(t *testing.T) {
	settings := NewUserSettings("user-123")
	originalUpdatedAt := settings.UpdatedAt

	// Update notification settings
	notifEmail := false
	notifPush := false
	settings.Update(&notifEmail, &notifPush, nil, nil)

	assert.False(t, settings.NotificationEmail)
	assert.False(t, settings.NotificationPush)
	assert.Equal(t, "light", settings.Theme) // unchanged
	assert.Equal(t, "en", settings.Language) // unchanged
	assert.False(t, settings.UpdatedAt.Before(originalUpdatedAt))

	// Update theme and language
	theme := "dark"
	language := "vi"
	settings.Update(nil, nil, &theme, &language)

	assert.False(t, settings.NotificationEmail) // unchanged
	assert.False(t, settings.NotificationPush)  // unchanged
	assert.Equal(t, "dark", settings.Theme)
	assert.Equal(t, "vi", settings.Language)
}
