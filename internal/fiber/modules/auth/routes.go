package auth

import (
	"iivineri/internal/fiber/modules/auth/handler"
	authMiddleware "iivineri/internal/fiber/modules/auth/middleware"
	"iivineri/internal/fiber/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, authHandler *handler.AuthHandler, sharedAuthMiddleware *middleware.AuthMiddleware) {
	// Create auth validation middleware
	validationMiddleware := authMiddleware.NewAuthValidationMiddleware()
	
	// Public routes (no authentication required)
	auth := app.Group("/api/v1/auth")

	// Authentication endpoints with validation
	auth.Post("/register", validationMiddleware.ValidateRegisterRequest(), authHandler.Register)
	auth.Post("/login", validationMiddleware.ValidateLoginRequest(), authHandler.Login)
	auth.Post("/refresh", validationMiddleware.ValidateRefreshTokenRequest(), authHandler.RefreshToken)

	// Password reset endpoints with validation
	auth.Post("/password/reset", validationMiddleware.ValidatePasswordResetRequest(), authHandler.RequestPasswordReset)
	auth.Post("/password/reset/confirm", validationMiddleware.ValidatePasswordResetConfirmRequest(), authHandler.ConfirmPasswordReset)

	// Protected routes (authentication required)
	protected := auth.Use(sharedAuthMiddleware.RequireAuth())

	// User profile
	protected.Get("/profile", authHandler.GetProfile)
	protected.Post("/logout", authHandler.Logout)
	protected.Post("/logout/all", authHandler.LogoutAll)

	// Password management with validation
	protected.Post("/password/change", validationMiddleware.ValidateChangePasswordRequest(), authHandler.ChangePassword)

	// 2FA management with validation
	protected.Post("/2fa/enable", validationMiddleware.ValidateEnable2FARequest(), authHandler.Enable2FA)
	protected.Post("/2fa/confirm", validationMiddleware.ValidateConfirm2FARequest(), authHandler.Confirm2FA)
	protected.Post("/2fa/disable", validationMiddleware.ValidateDisable2FARequest(), authHandler.Disable2FA)
}