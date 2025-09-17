package models

// Generic response models for documentation
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"The request contains invalid data"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Tag     string `json:"tag" example:"required"`
	Value   string `json:"value" example:""`
	Message string `json:"message" example:"This field is required"`
}

type ValidationErrorResponse struct {
	Error   string                  `json:"error" example:"Validation Failed"`
	Message string                  `json:"message" example:"The request contains invalid data"`
	Details []ValidationErrorDetail `json:"details"`
}

type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// Auth-specific response models
type RegisterSuccessResponse struct {
	Message string      `json:"message" example:"User registered successfully"`
	Data    *UserPublic `json:"data"`
}

type LoginSuccessResponse struct {
	Message string         `json:"message" example:"Login successful"`
	Data    *LoginResponse `json:"data"`
}

type ProfileSuccessResponse struct {
	Message string      `json:"message" example:"Profile retrieved successfully"`
	Data    *UserPublic `json:"data"`
}

type LogoutSuccessResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

type Enable2FASuccessResponse struct {
	Message string             `json:"message" example:"2FA setup initiated"`
	Data    *Enable2FAResponse `json:"data"`
}

type PasswordResetSuccessResponse struct {
	Message string `json:"message" example:"Password reset email sent"`
}

type PasswordResetConfirmSuccessResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
}

// Error response examples
type AuthError struct {
	Error   string `json:"error" example:"Authentication Failed"`
	Message string `json:"message" example:"Invalid email or password"`
}

type UnauthorizedError struct {
	Error   string `json:"error" example:"Unauthorized"`
	Message string `json:"message" example:"Authentication required"`
}

type ForbiddenError struct {
	Error   string `json:"error" example:"Forbidden"`
	Message string `json:"message" example:"Insufficient permissions"`
}

type NotFoundError struct {
	Error   string `json:"error" example:"Not Found"`
	Message string `json:"message" example:"Resource not found"`
}

type ConflictError struct {
	Error   string `json:"error" example:"Conflict"`
	Message string `json:"message" example:"Email already exists"`
}

type TooManyRequestsError struct {
	Error   string `json:"error" example:"Too Many Requests"`
	Message string `json:"message" example:"Rate limit exceeded"`
}

type InternalServerError struct {
	Error   string `json:"error" example:"Internal Server Error"`
	Message string `json:"message" example:"An unexpected error occurred"`
}
