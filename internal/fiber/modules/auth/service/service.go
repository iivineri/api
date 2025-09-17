package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"iivineri/internal/config"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/fiber/modules/auth/repository"
	"iivineri/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo          repository.UserRepositoryInterface
	sessionRepo       repository.SessionRepositoryInterface
	resetPasswordRepo repository.ResetPasswordRepositoryInterface
	user2FASecretRepo repository.User2FASecretRepositoryInterface
	banRepo           repository.BanRepositoryInterface
	config            *config.Config
	logger            *logger.Logger
	jwtSecret         []byte
}

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	sessionRepo repository.SessionRepositoryInterface,
	resetPasswordRepo repository.ResetPasswordRepositoryInterface,
	user2FASecretRepo repository.User2FASecretRepositoryInterface,
	banRepo repository.BanRepositoryInterface,
	config *config.Config,
	logger *logger.Logger,
) AuthService {
	return &AuthServiceImpl{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		resetPasswordRepo: resetPasswordRepo,
		user2FASecretRepo: user2FASecretRepo,
		banRepo:           banRepo,
		config:            config,
		logger:            logger,
		jwtSecret:         []byte(config.App.JWTSecret()),
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *models.RegisterRequest) (*models.UserPublic, error) {
	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if nickname already exists
	nicknameExists, err := s.userRepo.NicknameExists(ctx, req.Nickname)
	if err != nil {
		return nil, fmt.Errorf("failed to check nickname existence: %w", err)
	}
	if nicknameExists {
		return nil, fmt.Errorf("nickname already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse date of birth
	dateOfBirth, err := req.GetDateOfBirth()
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth format: %w", err)
	}

	// Create user
	user := &models.User{
		Nickname:    req.Nickname,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Enabled2FA:  false,
		Secret2FA:   "",
		DateOfBirth: dateOfBirth,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Infof("User registered successfully: %s", user.Email)
	return user.ToPublic(), nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is banned
	isBanned, err := s.banRepo.IsUserBanned(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ban status: %w", err)
	}
	if isBanned {
		return nil, fmt.Errorf("user is banned")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check 2FA if enabled
	if user.Is2FAEnabled() {
		if req.TOTPCode == nil || *req.TOTPCode == "" {
			return &models.LoginResponse{
				Requires2FA: true,
			}, nil
		}

		valid := totp.Validate(*req.TOTPCode, user.Secret2FA)
		if !valid {
			return nil, fmt.Errorf("invalid 2FA code")
		}
	}

	// Generate tokens
	accessToken, _, expiresAt, err := s.GenerateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Infof("User logged in successfully: %s", user.Email)

	return &models.LoginResponse{
		Token:       accessToken,
		User:        user.ToPublic(),
		ExpiresAt:   time.Unix(expiresAt, 0),
		Requires2FA: false,
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, userID int) error {
	// For now, we'll just log the logout
	// In a real implementation, you might want to blacklist the token
	s.logger.Infof("User logged out: %d", userID)
	return nil
}

func (s *AuthServiceImpl) LogoutAll(ctx context.Context, userID int) error {
	err := s.sessionRepo.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to logout all sessions: %w", err)
	}

	s.logger.Infof("User logged out from all devices: %d", userID)
	return nil
}

func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID int) (*models.UserPublic, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user.ToPublic(), nil
}

func (s *AuthServiceImpl) ChangePassword(ctx context.Context, userID int, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Check 2FA if enabled
	if user.Is2FAEnabled() {
		if req.TOTPCode == nil || *req.TOTPCode == "" {
			return fmt.Errorf("2FA code required")
		}

		valid := totp.Validate(*req.TOTPCode, user.Secret2FA)
		if !valid {
			return fmt.Errorf("invalid 2FA code")
		}
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Infof("Password changed for user: %d", userID)
	return nil
}

func (s *AuthServiceImpl) RequestPasswordReset(ctx context.Context, req *models.PasswordResetRequest) error {
	// Check if user exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		s.logger.Infof("Password reset requested for non-existent email: %s", req.Email)
		return nil
	}

	// Create reset token
	reset := &models.ResetPassword{
		Email: req.Email,
	}

	err = s.resetPasswordRepo.Create(ctx, reset)
	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send email with reset token
	s.logger.Infof("Password reset token created for: %s (token: %s)", req.Email, reset.ID)

	return nil
}

func (s *AuthServiceImpl) ConfirmPasswordReset(ctx context.Context, req *models.PasswordResetConfirmRequest) error {
	// Get reset token
	reset, err := s.resetPasswordRepo.GetByID(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Get user
	user, err := s.userRepo.GetByEmail(ctx, reset.Email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete reset token
	err = s.resetPasswordRepo.Delete(ctx, req.Token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete reset token")
	}

	// Logout all sessions
	err = s.sessionRepo.DeleteAllByUserID(ctx, user.ID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to logout all sessions after password reset")
	}

	s.logger.Infof("Password reset completed for user: %s", user.Email)
	return nil
}

func (s *AuthServiceImpl) Enable2FA(ctx context.Context, userID int, req *models.Enable2FARequest) (*models.Enable2FAResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate 2FA secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Iivineri",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate 2FA secret: %w", err)
	}

	// Generate backup codes
	backupCodes := make([]string, 10)
	for i := range backupCodes {
		code := make([]byte, 10)
		rand.Read(code)
		backupCodes[i] = base32.StdEncoding.EncodeToString(code)[:8]
	}

	// Store backup codes
	for _, code := range backupCodes {
		hashedCode, _ := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		secret := &models.User2FASecret{
			UserID: userID,
			Hash:   string(hashedCode),
		}
		s.user2FASecretRepo.Create(ctx, secret)
	}

	return &models.Enable2FAResponse{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

func (s *AuthServiceImpl) Confirm2FA(ctx context.Context, userID int, req *models.Confirm2FARequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Validate TOTP code
	valid := totp.Validate(req.TOTPCode, req.Secret)
	if !valid {
		return fmt.Errorf("invalid 2FA code")
	}

	// Enable 2FA
	user.Enabled2FA = true
	user.Secret2FA = req.Secret

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	s.logger.Infof("2FA enabled for user: %d", userID)
	return nil
}

func (s *AuthServiceImpl) Disable2FA(ctx context.Context, userID int, req *models.Disable2FARequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	// Check 2FA code if provided
	if req.TOTPCode != nil && *req.TOTPCode != "" {
		valid := totp.Validate(*req.TOTPCode, user.Secret2FA)
		if !valid {
			return fmt.Errorf("invalid 2FA code")
		}
	}

	// Disable 2FA
	user.Enabled2FA = false
	user.Secret2FA = ""

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	// Delete backup codes
	err = s.user2FASecretRepo.DeleteByUserID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete 2FA backup codes")
	}

	s.logger.Infof("2FA disabled for user: %d", userID)
	return nil
}

func (s *AuthServiceImpl) GenerateTokens(ctx context.Context, user *models.User) (accessToken, refreshToken string, expiresAt int64, err error) {
	now := time.Now()
	expiresAt = now.Add(24 * time.Hour).Unix()

	// Create access token
	claims := &JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "iivineri",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign access token: %w", err)
	}

	// For now, refresh token is the same as access token with longer expiry
	refreshClaims := &JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "iivineri-refresh",
		},
	}

	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenJWT.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, expiresAt, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.LoginResponse, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(req.RefreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if it's a refresh token
	if claims.Issuer != "iivineri-refresh" {
		return nil, fmt.Errorf("invalid refresh token type")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is banned
	isBanned, err := s.banRepo.IsUserBanned(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ban status: %w", err)
	}
	if isBanned {
		return nil, fmt.Errorf("user is banned")
	}

	// Generate new tokens
	accessToken, _, expiresAt, err := s.GenerateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.LoginResponse{
		Token:     accessToken,
		User:      user.ToPublic(),
		ExpiresAt: time.Unix(expiresAt, 0),
	}, nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if it's an access token
	if claims.Issuer != "iivineri" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *AuthServiceImpl) ValidateUser(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is banned
	isBanned, err := s.banRepo.IsUserBanned(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ban status: %w", err)
	}
	if isBanned {
		return nil, fmt.Errorf("user is banned")
	}

	return user, nil
}
