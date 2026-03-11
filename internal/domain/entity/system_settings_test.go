package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSystemSettings(t *testing.T) {
	key := "max_upload_size"
	value := map[string]interface{}{"size": 10485760}
	desc := "Maximum upload size in bytes"

	settings := NewSystemSettings(key, value, &desc)

	assert.NotNil(t, settings)
	assert.Equal(t, key, settings.Key)
	assert.Equal(t, value, settings.Value)
	assert.Equal(t, &desc, settings.Description)
	assert.NotEqual(t, "", settings.ID.String())
}

func TestSystemSettings_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *SystemSettings
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid settings",
			setup: func() *SystemSettings {
				value := map[string]interface{}{"enabled": true}
				return NewSystemSettings("feature_flag", value, nil)
			},
			wantErr: false,
		},
		{
			name: "empty key",
			setup: func() *SystemSettings {
				value := map[string]interface{}{"test": "value"}
				return NewSystemSettings("", value, nil)
			},
			wantErr: true,
			errMsg:  "key is required",
		},
		{
			name: "key too long",
			setup: func() *SystemSettings {
				longKey := string(make([]byte, 101))
				for i := range longKey {
					longKey = longKey[:i] + "a" + longKey[i+1:]
				}
				value := map[string]interface{}{"test": "value"}
				return NewSystemSettings(longKey, value, nil)
			},
			wantErr: true,
			errMsg:  "key too long",
		},
		{
			name: "nil value",
			setup: func() *SystemSettings {
				s := NewSystemSettings("test_key", nil, nil)
				return s
			},
			wantErr: true,
			errMsg:  "value is required",
		},
		{
			name: "valid with description",
			setup: func() *SystemSettings {
				value := map[string]interface{}{"count": 100}
				desc := "Test description"
				return NewSystemSettings("test_key", value, &desc)
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

func TestSystemSettings_Update(t *testing.T) {
	value := map[string]interface{}{"enabled": true}
	settings := NewSystemSettings("feature_flag", value, nil)
	originalUpdatedAt := settings.UpdatedAt

	// Update value
	newValue := map[string]interface{}{"enabled": false, "reason": "maintenance"}
	settings.Update(newValue, nil)

	assert.Equal(t, newValue, settings.Value)
	assert.Nil(t, settings.Description)
	assert.False(t, settings.UpdatedAt.Before(originalUpdatedAt))

	// Update description
	newDesc := "Updated description"
	settings.Update(nil, &newDesc)

	assert.Equal(t, newValue, settings.Value) // unchanged
	assert.Equal(t, &newDesc, settings.Description)
}
