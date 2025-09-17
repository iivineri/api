package repository

import (
	"context"
	"database/sql"
	"fmt"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/logger"
)

type UserRepository struct {
	db     database.DatabaseInterface
	logger *logger.Logger
}

func NewUserRepository(db database.DatabaseInterface, logger *logger.Logger) UserRepositoryInterface {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (nickname, email, password, date_of_birth)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		user.Nickname,
		user.Email,
		user.Password,
		user.DateOfBirth,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, nickname, email, password, enabled_2fa, secret_2fa, date_of_birth,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Enabled2FA,
		&user.Secret2FA,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, nickname, email, password, enabled_2fa, secret_2fa, date_of_birth,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Enabled2FA,
		&user.Secret2FA,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	query := `
		SELECT id, nickname, email, password, enabled_2fa, secret_2fa, date_of_birth,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE nickname = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, nickname).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Enabled2FA,
		&user.Secret2FA,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by nickname")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET nickname = $2, email = $3, password = $4, enabled_2fa = $5,
		    secret_2fa = $6, date_of_birth = $7, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query,
		user.ID,
		user.Nickname,
		user.Email,
		user.Password,
		user.Enabled2FA,
		user.Secret2FA,
		user.DateOfBirth,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to soft delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.WithError(err).Error("Failed to check if email exists")
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) NicknameExists(ctx context.Context, nickname string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nickname = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRow(ctx, query, nickname).Scan(&exists)
	if err != nil {
		r.logger.WithError(err).Error("Failed to check if nickname exists")
		return false, fmt.Errorf("failed to check nickname existence: %w", err)
	}

	return exists, nil
}
