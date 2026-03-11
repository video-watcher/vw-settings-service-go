package dto

import "time"

// UserSettingsResponse represents user settings response
type UserSettingsResponse struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	NotificationEmail bool      `json:"notification_email"`
	NotificationPush  bool      `json:"notification_push"`
	Theme             string    `json:"theme"`
	Language          string    `json:"language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserSettingsUpdateRequest represents user settings update request
type UserSettingsUpdateRequest struct {
	NotificationEmail *bool   `json:"notification_email"`
	NotificationPush  *bool   `json:"notification_push"`
	Theme             *string `json:"theme" validate:"omitempty,oneof=light dark"`
	Language          *string `json:"language" validate:"omitempty,len=2"`
}

// SystemSettingResponse represents system setting response
type SystemSettingResponse struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Value       map[string]interface{} `json:"value"`
	Description *string                `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// SystemSettingCreateRequest represents system setting create/update request
type SystemSettingCreateRequest struct {
	Value       map[string]interface{} `json:"value" validate:"required"`
	Description *string                `json:"description"`
}
