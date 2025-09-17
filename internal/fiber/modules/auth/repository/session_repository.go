package repository

import (
	"context"
	"database/sql"
	"fmt"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/logger"
)

type SessionRepository struct {
	db     database.DatabaseInterface
	logger *logger.Logger
}

func NewSessionRepository(db database.DatabaseInterface, logger *logger.Logger) SessionRepositoryInterface {
	return &SessionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, user_agent, ip_address, mime_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		session.UserID,
		session.UserAgent,
		session.IPAddress,
		session.MimeType,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create session")
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id int) (*models.Session, error) {
	query := `
		SELECT id, user_id, user_agent, ip_address, mime_type, created_at, deleted_at
		FROM sessions
		WHERE id = $1`

	session := &models.Session{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.UserAgent,
		&session.IPAddress,
		&session.MimeType,
		&session.CreatedAt,
		&session.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		r.logger.WithError(err).Error("Failed to get session by ID")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) GetActiveByUserID(ctx context.Context, userID int) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, user_agent, ip_address, mime_type, created_at, deleted_at
		FROM sessions
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get active sessions by user ID")
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.UserAgent,
			&session.IPAddress,
			&session.MimeType,
			&session.CreatedAt,
			&session.DeletedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan session")
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE sessions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete session")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *SessionRepository) DeleteAllByUserID(ctx context.Context, userID int) error {
	query := `UPDATE sessions SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete all sessions for user")
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	return nil
}

func (r *SessionRepository) CleanupExpired(ctx context.Context) error {
	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE deleted_at IS NULL
		AND created_at < NOW() - INTERVAL '30 days'`

	err := r.db.Exec(ctx, query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to cleanup expired sessions")
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	return nil
}
