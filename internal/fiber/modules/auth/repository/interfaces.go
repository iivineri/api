package repository

import (
	"context"
	"iivineri/internal/fiber/modules/auth/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByNickname(ctx context.Context, nickname string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	SoftDelete(ctx context.Context, id int) error
	EmailExists(ctx context.Context, email string) (bool, error)
	NicknameExists(ctx context.Context, nickname string) (bool, error)
}

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id int) (*models.Session, error)
	GetActiveByUserID(ctx context.Context, userID int) ([]*models.Session, error)
	Delete(ctx context.Context, id int) error
	DeleteAllByUserID(ctx context.Context, userID int) error
	CleanupExpired(ctx context.Context) error
}

type ResetPasswordRepositoryInterface interface {
	Create(ctx context.Context, reset *models.ResetPassword) error
	GetByID(ctx context.Context, id string) (*models.ResetPassword, error)
	GetByEmail(ctx context.Context, email string) (*models.ResetPassword, error)
	Delete(ctx context.Context, id string) error
	DeleteByEmail(ctx context.Context, email string) error
	CleanupExpired(ctx context.Context) error
}

type User2FASecretRepositoryInterface interface {
	Create(ctx context.Context, secret *models.User2FASecret) error
	GetActiveByUserID(ctx context.Context, userID int) ([]*models.User2FASecret, error)
	DeleteByUserID(ctx context.Context, userID int) error
	MarkUsed(ctx context.Context, id int) error
}

type BanRepositoryInterface interface {
	Create(ctx context.Context, ban *models.Ban) error
	GetActiveByUserID(ctx context.Context, userID int) (*models.Ban, error)
	IsUserBanned(ctx context.Context, userID int) (bool, error)
}