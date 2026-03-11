package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserSettings represents user preferences and notification settings
type UserSettings struct {
	ID                uuid.UUID `json:"id"`
	UserID            string    `json:"user_id"`
	NotificationEmail bool      `json:"notification_email"`
	NotificationPush  bool      `json:"notification_push"`
	Theme             string    `json:"theme"`
	Language          string    `json:"language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewUserSettings creates a new UserSettings with default values
func NewUserSettings(userID string) *UserSettings {
	now := time.Now()
	return &UserSettings{
		ID:                uuid.New(),
		UserID:            userID,
		NotificationEmail: true,
		NotificationPush:  true,
		Theme:             "light",
		Language:          "en",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// Validate validates the UserSettings entity
func (u *UserSettings) Validate() error {
	if u.UserID == "" {
		return errors.New("user_id is required")
	}

	if u.Theme != "light" && u.Theme != "dark" {
		return errors.New("theme must be 'light' or 'dark'")
	}

	if len(u.Language) != 2 {
		return errors.New("language must be 2-character code (e.g., 'en', 'vi')")
	}

	return nil
}

// Update updates user settings with provided values
func (u *UserSettings) Update(notificationEmail, notificationPush *bool, theme, language *string) {
	if notificationEmail != nil {
		u.NotificationEmail = *notificationEmail
	}
	if notificationPush != nil {
		u.NotificationPush = *notificationPush
	}
	if theme != nil {
		u.Theme = *theme
	}
	if language != nil {
		u.Language = *language
	}
	u.UpdatedAt = time.Now()
}
