package middleware

import (
	"encoding/json"
	"strings"

	"iivineri/internal/fiber/modules/auth/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthValidationMiddleware struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details"`
}

func NewAuthValidationMiddleware() *AuthValidationMiddleware {
	return &AuthValidationMiddleware{
		validator: validator.New(),
	}
}

func (vm *AuthValidationMiddleware) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + err.Param() + " characters"
	case "max":
		return "Must be at most " + err.Param() + " characters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "len":
		return "Must be exactly " + err.Param() + " characters"
	case "numeric":
		return "Must contain only numbers"
	case "uuid":
		return "Must be a valid UUID"
	default:
		return "Validation failed for tag: " + err.Tag()
	}
}

func (vm *AuthValidationMiddleware) handleValidationErrors(c *fiber.Ctx, err error) error {
	var validationErrors []ValidationError
	
	for _, err := range err.(validator.ValidationErrors) {
		validationError := ValidationError{
			Field:   strings.ToLower(err.Field()),
			Tag:     err.Tag(),
			Value:   err.Param(),
			Message: vm.getValidationMessage(err),
		}
		validationErrors = append(validationErrors, validationError)
	}

	return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
		Error:   "Validation Failed",
		Message: "The request contains invalid data",
		Details: validationErrors,
	})
}

// ValidateRegisterRequest validates registration requests
func (vm *AuthValidationMiddleware) ValidateRegisterRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		
		// Parse JSON body with better error handling
		if err := c.BodyParser(&req); err != nil {
			// Check if it's a JSON parsing error
			if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid JSON format",
					"message": "Invalid value for field '" + jsonErr.Field + "'. Expected " + jsonErr.Type.String() + " but got " + jsonErr.Value,
					"field":   jsonErr.Field,
				})
			}
			
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format: " + err.Error(),
			})
		}

		// Validate the request
		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		// Additional validation for date format
		if _, err := req.GetDateOfBirth(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
				Error:   "Validation Failed",
				Message: "The request contains invalid data",
				Details: []ValidationError{
					{
						Field:   "date_of_birth",
						Tag:     "format",
						Message: err.Error(),
					},
				},
			})
		}

		// Store validated request in context
		c.Locals("validatedRequest", &req)
		
		return c.Next()
	}
}

// ValidateLoginRequest validates login requests
func (vm *AuthValidationMiddleware) ValidateLoginRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.LoginRequest
		
		// Parse JSON body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		// Validate the request
		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		// Store validated request in context
		c.Locals("validatedRequest", &req)
		
		return c.Next()
	}
}

// ValidatePasswordResetRequest validates password reset requests
func (vm *AuthValidationMiddleware) ValidatePasswordResetRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.PasswordResetRequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidatePasswordResetConfirmRequest validates password reset confirmation requests
func (vm *AuthValidationMiddleware) ValidatePasswordResetConfirmRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.PasswordResetConfirmRequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidateChangePasswordRequest validates change password requests
func (vm *AuthValidationMiddleware) ValidateChangePasswordRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ChangePasswordRequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidateEnable2FARequest validates enable 2FA requests
func (vm *AuthValidationMiddleware) ValidateEnable2FARequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Enable2FARequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidateConfirm2FARequest validates confirm 2FA requests
func (vm *AuthValidationMiddleware) ValidateConfirm2FARequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Confirm2FARequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidateDisable2FARequest validates disable 2FA requests
func (vm *AuthValidationMiddleware) ValidateDisable2FARequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Disable2FARequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}

// ValidateRefreshTokenRequest validates refresh token requests
func (vm *AuthValidationMiddleware) ValidateRefreshTokenRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RefreshTokenRequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Invalid JSON format",
			})
		}

		if err := vm.validator.Struct(&req); err != nil {
			return vm.handleValidationErrors(c, err)
		}

		c.Locals("validatedRequest", &req)
		return c.Next()
	}
}