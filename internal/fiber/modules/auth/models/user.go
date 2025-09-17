package models

import (
	"time"
)

type User struct {
	ID          int        `json:"id" db:"id"`
	Nickname    string     `json:"nickname" db:"nickname"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password"`
	Enabled2FA  bool       `json:"enabled_2fa" db:"enabled_2fa"`
	Secret2FA   string     `json:"-" db:"secret_2fa"`
	DateOfBirth time.Time  `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Session struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	UserAgent *string    `json:"user_agent,omitempty" db:"user_agent"`
	IPAddress *string    `json:"ip_address,omitempty" db:"ip_address"`
	MimeType  *string    `json:"mime_type,omitempty" db:"mime_type"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ResetPassword struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
}

type User2FASecret struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Hash      string     `json:"-" db:"hash"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Ban struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	BannedBy  int       `json:"banned_by" db:"banned_by"`
	Reason    string    `json:"reason" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) Is2FAEnabled() bool {
	return u.Enabled2FA && u.Secret2FA != ""
}

func (s *Session) IsActive() bool {
	return s.DeletedAt == nil
}

func (u *User) ToPublic() *UserPublic {
	return &UserPublic{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Email:       u.Email,
		Enabled2FA:  u.Enabled2FA,
		DateOfBirth: u.DateOfBirth,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type UserPublic struct {
	ID          int        `json:"id"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Enabled2FA  bool       `json:"enabled_2fa"`
	DateOfBirth time.Time  `json:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}