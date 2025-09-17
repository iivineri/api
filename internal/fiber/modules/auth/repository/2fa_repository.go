package repository

import (
	"context"
	"fmt"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/logger"
)

type User2FASecretRepository struct {
	db     database.DatabaseInterface
	logger *logger.Logger
}

func New2FARepository(db database.DatabaseInterface, logger *logger.Logger) User2FASecretRepositoryInterface {
	return &User2FASecretRepository{
		db:     db,
		logger: logger,
	}
}

func (r *User2FASecretRepository) Create(ctx context.Context, secret *models.User2FASecret) error {
	query := `
		INSERT INTO user_2fa_secrets (user_id, hash)
		VALUES ($1, $2)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, secret.UserID, secret.Hash).Scan(&secret.ID, &secret.CreatedAt)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create 2FA secret")
		return fmt.Errorf("failed to create 2FA secret: %w", err)
	}

	return nil
}

func (r *User2FASecretRepository) GetActiveByUserID(ctx context.Context, userID int) ([]*models.User2FASecret, error) {
	query := `
		SELECT id, user_id, hash, created_at, deleted_at
		FROM user_2fa_secrets 
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get active 2FA secrets by user ID")
		return nil, fmt.Errorf("failed to get 2FA secrets: %w", err)
	}
	defer rows.Close()

	var secrets []*models.User2FASecret
	for rows.Next() {
		secret := &models.User2FASecret{}
		err := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.Hash,
			&secret.CreatedAt,
			&secret.DeletedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan 2FA secret")
			return nil, fmt.Errorf("failed to scan 2FA secret: %w", err)
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (r *User2FASecretRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `UPDATE user_2fa_secrets SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete 2FA secrets for user")
		return fmt.Errorf("failed to delete 2FA secrets: %w", err)
	}

	return nil
}

func (r *User2FASecretRepository) MarkUsed(ctx context.Context, id int) error {
	query := `UPDATE user_2fa_secrets SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to mark 2FA secret as used")
		return fmt.Errorf("failed to mark 2FA secret as used: %w", err)
	}

	return nil
}