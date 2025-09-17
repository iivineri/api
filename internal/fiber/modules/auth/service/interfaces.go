package service

import (
	"context"
	"iivineri/internal/fiber/modules/auth/models"
)

type AuthService interface {
	// Authentication
	Register(ctx context.Context, dto *models.RegisterRequest) (*models.UserPublic, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID int) error
	LogoutAll(ctx context.Context, userID int) error

	// User management
	GetProfile(ctx context.Context, userID int) (*models.UserPublic, error)
	ChangePassword(ctx context.Context, userID int, req *models.ChangePasswordRequest) error

	// Password reset
	RequestPasswordReset(ctx context.Context, req *models.PasswordResetRequest) error
	ConfirmPasswordReset(ctx context.Context, req *models.PasswordResetConfirmRequest) error

	// 2FA
	Enable2FA(ctx context.Context, userID int, req *models.Enable2FARequest) (*models.Enable2FAResponse, error)
	Confirm2FA(ctx context.Context, userID int, req *models.Confirm2FARequest) error
	Disable2FA(ctx context.Context, userID int, req *models.Disable2FARequest) error

	// Token management
	GenerateTokens(ctx context.Context, user *models.User) (accessToken, refreshToken string, expiresAt int64, err error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)

	// Validation
	ValidateUser(ctx context.Context, userID int) (*models.User, error)
}
