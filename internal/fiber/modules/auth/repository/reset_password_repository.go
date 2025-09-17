package repository

import (
	"context"
	"database/sql"
	"fmt"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/logger"
	"time"

	"github.com/google/uuid"
)

type ResetPasswordRepository struct {
	db     database.DatabaseInterface
	logger *logger.Logger
}

func NewResetPasswordRepository(db database.DatabaseInterface, logger *logger.Logger) ResetPasswordRepositoryInterface {
	return &ResetPasswordRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ResetPasswordRepository) Create(ctx context.Context, reset *models.ResetPassword) error {
	reset.ID = uuid.New().String()
	reset.CreatedAt = time.Now()
	reset.ExpiredAt = time.Now().Add(24 * time.Hour)

	query := `
		INSERT INTO reset_passwords (id, email, created_at, expired_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) 
		DO UPDATE SET id = $1, created_at = $3, expired_at = $4`

	err := r.db.Exec(ctx, query, reset.ID, reset.Email, reset.CreatedAt, reset.ExpiredAt)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create reset password token")
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	return nil
}

func (r *ResetPasswordRepository) GetByID(ctx context.Context, id string) (*models.ResetPassword, error) {
	query := `
		SELECT id, email, created_at, expired_at
		FROM reset_passwords 
		WHERE id = $1 AND expired_at > NOW()`

	reset := &models.ResetPassword{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&reset.ID,
		&reset.Email,
		&reset.CreatedAt,
		&reset.ExpiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reset token not found or expired")
		}
		r.logger.WithError(err).Error("Failed to get reset password token")
		return nil, fmt.Errorf("failed to get reset token: %w", err)
	}

	return reset, nil
}

func (r *ResetPasswordRepository) GetByEmail(ctx context.Context, email string) (*models.ResetPassword, error) {
	query := `
		SELECT id, email, created_at, expired_at
		FROM reset_passwords 
		WHERE email = $1 AND expired_at > NOW()`

	reset := &models.ResetPassword{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&reset.ID,
		&reset.Email,
		&reset.CreatedAt,
		&reset.ExpiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reset token not found or expired")
		}
		r.logger.WithError(err).Error("Failed to get reset password token by email")
		return nil, fmt.Errorf("failed to get reset token: %w", err)
	}

	return reset, nil
}

func (r *ResetPasswordRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reset_passwords WHERE id = $1`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete reset password token")
		return fmt.Errorf("failed to delete reset token: %w", err)
	}

	return nil
}

func (r *ResetPasswordRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM reset_passwords WHERE email = $1`

	err := r.db.Exec(ctx, query, email)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete reset password tokens by email")
		return fmt.Errorf("failed to delete reset tokens: %w", err)
	}

	return nil
}

func (r *ResetPasswordRepository) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM reset_passwords WHERE expired_at <= NOW()`

	err := r.db.Exec(ctx, query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to cleanup expired reset tokens")
		return fmt.Errorf("failed to cleanup reset tokens: %w", err)
	}

	return nil
}