package middleware

import (
	"context"
	"iivineri/internal/fiber/modules/auth/models"
	"iivineri/internal/fiber/modules/auth/service"
	"iivineri/internal/fiber/shared/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Authorization header required")
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.UnauthorizedResponse(c, "Invalid authorization format")
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return utils.UnauthorizedResponse(c, "Token required")
		}

		// Validate token
		user, err := m.authService.ValidateToken(context.Background(), token)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid or expired token")
		}

		// Validate user (check if banned, etc.)
		user, err = m.authService.ValidateUser(context.Background(), user.ID)
		if err != nil {
			return utils.UnauthorizedResponse(c, err.Error())
		}

		// Store user in context
		c.Locals("user", user)
		c.Locals("user_id", user.ID)

		return c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next()
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Next()
		}

		// Validate token
		user, err := m.authService.ValidateToken(context.Background(), token)
		if err != nil {
			return c.Next()
		}

		// Validate user (check if banned, etc.)
		user, err = m.authService.ValidateUser(context.Background(), user.ID)
		if err != nil {
			return c.Next()
		}

		// Store user in context
		c.Locals("user", user)
		c.Locals("user_id", user.ID)

		return c.Next()
	}
}

// Helper function to get user from context
func GetUserFromContext(c *fiber.Ctx) (*models.User, bool) {
	user, ok := c.Locals("user").(*models.User)
	return user, ok
}

// Helper function to get user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (int, bool) {
	userID, ok := c.Locals("user_id").(int)
	return userID, ok
}