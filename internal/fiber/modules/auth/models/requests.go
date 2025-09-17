package models

import (
	"fmt"
	"time"
)

type LoginRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	TOTPCode *string `json:"totp_code,omitempty" validate:"omitempty,len=6,numeric"`
}

type RegisterRequest struct {
	Nickname    string `json:"nickname" validate:"required,min=3,max=32,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

// GetDateOfBirth parses the date string and returns a time.Time
func (r *RegisterRequest) GetDateOfBirth() (time.Time, error) {
	// Try multiple date formats
	formats := []string{
		"2006-01-02",           // "2020-12-12"
		"2006-01-02T15:04:05Z", // "2020-12-12T00:00:00Z"
		time.RFC3339,           // "2020-12-12T00:00:00Z07:00"
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, r.DateOfBirth); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("invalid date format: %s, expected format: YYYY-MM-DD", r.DateOfBirth)
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token" validate:"required,uuid"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type Enable2FARequest struct {
	Password string `json:"password" validate:"required"`
}

type Confirm2FARequest struct {
	Secret   string `json:"secret" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required,len=6,numeric"`
}

type Disable2FARequest struct {
	Password string  `json:"password" validate:"required"`
	TOTPCode *string `json:"totp_code,omitempty" validate:"omitempty,len=6,numeric"`
}

type ChangePasswordRequest struct {
	CurrentPassword string  `json:"current_password" validate:"required"`
	NewPassword     string  `json:"new_password" validate:"required,min=8"`
	TOTPCode        *string `json:"totp_code,omitempty" validate:"omitempty,len=6,numeric"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      *UserPublic `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
	Requires2FA bool      `json:"requires_2fa,omitempty"`
}

type Enable2FAResponse struct {
	Secret     string   `json:"secret"`
	QRCodeURL  string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}