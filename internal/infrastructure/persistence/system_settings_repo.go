package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"settings-service-go/internal/domain/entity"
	"settings-service-go/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SystemSettingsRepositoryImpl implements SystemSettingsRepository
type SystemSettingsRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewSystemSettingsRepository creates a new SystemSettingsRepositoryImpl
func NewSystemSettingsRepository(db *pgxpool.Pool) repository.SystemSettingsRepository {
	return &SystemSettingsRepositoryImpl{db: db}
}

// GetAll retrieves all system settings with pagination
func (r *SystemSettingsRepositoryImpl) GetAll(ctx context.Context, skip, limit int) ([]*entity.SystemSettings, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM system_settings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*entity.SystemSettings
	for rows.Next() {
		var s entity.SystemSettings
		var valueJSON []byte

		err := rows.Scan(
			&s.ID,
			&s.Key,
			&valueJSON,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSONB value
		if err := json.Unmarshal(valueJSON, &s.Value); err != nil {
			return nil, err
		}

		settings = append(settings, &s)
	}

	return settings, rows.Err()
}

// GetByKey retrieves system setting by key
func (r *SystemSettingsRepositoryImpl) GetByKey(ctx context.Context, key string) (*entity.SystemSettings, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM system_settings
		WHERE key = $1
	`

	var s entity.SystemSettings
	var valueJSON []byte

	err := r.db.QueryRow(ctx, query, key).Scan(
		&s.ID,
		&s.Key,
		&valueJSON,
		&s.Description,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSONB value
	if err := json.Unmarshal(valueJSON, &s.Value); err != nil {
		return nil, err
	}

	return &s, nil
}

// Set creates or updates a system setting
func (r *SystemSettingsRepositoryImpl) Set(ctx context.Context, settings *entity.SystemSettings) error {
	// Marshal value to JSON
	valueJSON, err := json.Marshal(settings.Value)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO system_settings (id, key, value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.Exec(ctx, query,
		settings.ID,
		settings.Key,
		valueJSON,
		settings.Description,
		settings.CreatedAt,
		settings.UpdatedAt,
	)

	return err
}

// DeleteByKey deletes system setting by key
func (r *SystemSettingsRepositoryImpl) DeleteByKey(ctx context.Context, key string) error {
	query := `DELETE FROM system_settings WHERE key = $1`
	result, err := r.db.Exec(ctx, query, key)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("setting not found")
	}

	return nil
}

// Delete deletes system setting by ID
func (r *SystemSettingsRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM system_settings WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("setting not found")
	}

	return nil
}
