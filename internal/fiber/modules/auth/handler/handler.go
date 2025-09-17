package handler

import (
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/fiber/modules/auth/service"
	"iivineri/internal/fiber/shared/middleware"
	"iivineri/internal/fiber/shared/utils"
	"iivineri/internal/logger"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register godoc
// @Summary Register a new user account
// @Description Create a new user account with validation and return user profile
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration data (date_of_birth format: YYYY-MM-DD)"
// @Success 201 {object} models.RegisterSuccessResponse "User registered successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Validation error - invalid input data"
// @Failure 409 {object} models.ConflictError "Conflict - email or nickname already exists"
// @Failure 500 {object} models.InternalServerError "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	dto := c.Locals("validatedRequest").(*models.RegisterRequest)

	user, err := h.authService.Register(c.Context(), dto)
	if err != nil {
		h.logger.WithError(err).Error("Registration failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User registered successfully", user)
}

// Login godoc
// @Summary User authentication
// @Description Authenticate user with email and password, optionally with 2FA code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "User login credentials (include totp_code if 2FA enabled)"
// @Success 200 {object} models.LoginSuccessResponse "Login successful - returns JWT token and user profile"
// @Failure 400 {object} models.ValidationErrorResponse "Validation error - invalid input data"
// @Failure 401 {object} models.AuthError "Authentication failed - invalid credentials or 2FA code"
// @Failure 403 {object} models.ForbiddenError "Account banned or disabled"
// @Failure 429 {object} models.TooManyRequestsError "Too many login attempts"
// @Failure 500 {object} models.InternalServerError "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Login failed")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Login successful", response)
}

// Logout godoc
// @Summary Logout current session
// @Description Logout current user session and invalidate the current session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.LogoutSuccessResponse "Session logged out successfully"
// @Failure 401 {object} models.UnauthorizedError "Authentication required - invalid or missing token"
// @Failure 500 {object} models.InternalServerError "Internal server error"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err := h.authService.Logout(c.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Logout failed")
		return utils.InternalErrorResponse(c, "Failed to logout")
	}

	return utils.SuccessResponse(c, "Logged out successfully", nil)
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err := h.authService.LogoutAll(c.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Logout all failed")
		return utils.InternalErrorResponse(c, "Failed to logout from all devices")
	}

	return utils.SuccessResponse(c, "Logged out from all devices successfully", nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve current authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ProfileSuccessResponse "Profile retrieved successfully"
// @Failure 401 {object} models.UnauthorizedError "Authentication required - invalid or missing token"
// @Failure 404 {object} models.NotFoundError "User not found"
// @Failure 500 {object} models.InternalServerError "Internal server error"
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	profile, err := h.authService.GetProfile(c.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get profile")
		return utils.InternalErrorResponse(c, "Failed to get profile")
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req models.ChangePasswordRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Password change failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Password changed successfully", nil)
}

func (h *AuthHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req models.PasswordResetRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.RequestPasswordReset(c.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Password reset request failed")
		return utils.InternalErrorResponse(c, "Failed to process password reset request")
	}

	return utils.SuccessResponse(c, "Password reset email sent if account exists", nil)
}

func (h *AuthHandler) ConfirmPasswordReset(c *fiber.Ctx) error {
	var req models.PasswordResetConfirmRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.ConfirmPasswordReset(c.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Password reset confirmation failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Password reset successfully", nil)
}

func (h *AuthHandler) Enable2FA(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req models.Enable2FARequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Enable2FA(c.Context(), userID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Enable 2FA failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "2FA setup initiated", response)
}

func (h *AuthHandler) Confirm2FA(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req models.Confirm2FARequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.Confirm2FA(c.Context(), userID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Confirm 2FA failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "2FA enabled successfully", nil)
}

func (h *AuthHandler) Disable2FA(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req models.Disable2FARequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.Disable2FA(c.Context(), userID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Disable 2FA failed")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "2FA disabled successfully", nil)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} models.LoginSuccessResponse "Token refreshed successfully - returns new access token"
// @Failure 400 {object} models.ValidationErrorResponse "Validation error - invalid refresh token format"
// @Failure 401 {object} models.AuthError "Authentication failed - invalid or expired refresh token"
// @Failure 500 {object} models.InternalServerError "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.RefreshToken(c.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Token refresh failed")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
}
