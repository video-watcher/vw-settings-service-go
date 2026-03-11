package persistence

import (
	"context"
	"settings-service-go/internal/domain/entity"
	"settings-service-go/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserSettingsRepositoryImpl implements UserSettingsRepository
type UserSettingsRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewUserSettingsRepository creates a new UserSettingsRepositoryImpl
func NewUserSettingsRepository(db *pgxpool.Pool) repository.UserSettingsRepository {
	return &UserSettingsRepositoryImpl{db: db}
}

// GetByUserID retrieves user settings by user ID
func (r *UserSettingsRepositoryImpl) GetByUserID(ctx context.Context, userID string) (*entity.UserSettings, error) {
	query := `
		SELECT id, user_id, notification_email, notification_push, theme, language, created_at, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	var s entity.UserSettings
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.NotificationEmail,
		&s.NotificationPush,
		&s.Theme,
		&s.Language,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Create creates new user settings
func (r *UserSettingsRepositoryImpl) Create(ctx context.Context, settings *entity.UserSettings) error {
	query := `
		INSERT INTO user_settings (id, user_id, notification_email, notification_push, theme, language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		settings.ID,
		settings.UserID,
		settings.NotificationEmail,
		settings.NotificationPush,
		settings.Theme,
		settings.Language,
		settings.CreatedAt,
		settings.UpdatedAt,
	)

	return err
}

// Update updates existing user settings
func (r *UserSettingsRepositoryImpl) Update(ctx context.Context, settings *entity.UserSettings) error {
	query := `
		UPDATE user_settings
		SET notification_email = $2,
		    notification_push = $3,
		    theme = $4,
		    language = $5,
		    updated_at = $6
		WHERE user_id = $1
	`

	result, err := r.db.Exec(ctx, query,
		settings.UserID,
		settings.NotificationEmail,
		settings.NotificationPush,
		settings.Theme,
		settings.Language,
		settings.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Delete deletes user settings by ID
func (r *UserSettingsRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_settings WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
