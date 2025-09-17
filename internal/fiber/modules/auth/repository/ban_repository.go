package repository

import (
	"context"
	"database/sql"
	"fmt"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/logger"
)

type BanRepository struct {
	db     database.DatabaseInterface
	logger *logger.Logger
}

func NewBanRepository(db database.DatabaseInterface, logger *logger.Logger) BanRepositoryInterface {
	return &BanRepository{
		db:     db,
		logger: logger,
	}
}

func (r *BanRepository) Create(ctx context.Context, ban *models.Ban) error {
	query := `
		INSERT INTO bans (user_id, banned_by, reason)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, ban.UserID, ban.BannedBy, ban.Reason).Scan(&ban.ID, &ban.CreatedAt)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create ban")
		return fmt.Errorf("failed to create ban: %w", err)
	}

	return nil
}

func (r *BanRepository) GetActiveByUserID(ctx context.Context, userID int) (*models.Ban, error) {
	query := `
		SELECT id, user_id, banned_by, reason, created_at
		FROM bans 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	ban := &models.Ban{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&ban.ID,
		&ban.UserID,
		&ban.BannedBy,
		&ban.Reason,
		&ban.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no ban found")
		}
		r.logger.WithError(err).Error("Failed to get ban by user ID")
		return nil, fmt.Errorf("failed to get ban: %w", err)
	}

	return ban, nil
}

func (r *BanRepository) IsUserBanned(ctx context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM bans WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		r.logger.WithError(err).Error("Failed to check if user is banned")
		return false, fmt.Errorf("failed to check ban status: %w", err)
	}

	return exists, nil
}