package auth

import (
	"iivineri/internal/config"
	"iivineri/internal/database"
	"iivineri/internal/fiber/modules/auth/handler"
	"iivineri/internal/fiber/modules/auth/repository"
	"iivineri/internal/fiber/modules/auth/service"
	"iivineri/internal/fiber/shared/middleware"
	"iivineri/internal/logger"

	"github.com/google/wire"
)

// Repository providers
func ProvideUserRepository(db database.DatabaseInterface, logger *logger.Logger) repository.UserRepositoryInterface {
	return repository.NewUserRepository(db, logger)
}

func ProvideSessionRepository(db database.DatabaseInterface, logger *logger.Logger) repository.SessionRepositoryInterface {
	return repository.NewSessionRepository(db, logger)
}

func ProvideResetPasswordRepository(db database.DatabaseInterface, logger *logger.Logger) repository.ResetPasswordRepositoryInterface {
	return repository.NewResetPasswordRepository(db, logger)
}

func ProvideUser2FASecretRepository(db database.DatabaseInterface, logger *logger.Logger) repository.User2FASecretRepositoryInterface {
	return repository.New2FARepository(db, logger)
}

func ProvideBanRepository(db database.DatabaseInterface, logger *logger.Logger) repository.BanRepositoryInterface {
	return repository.NewBanRepository(db, logger)
}

// Service providers
func ProvideAuthService(
	userRepo repository.UserRepositoryInterface,
	sessionRepo repository.SessionRepositoryInterface,
	resetPasswordRepo repository.ResetPasswordRepositoryInterface,
	user2FASecretRepo repository.User2FASecretRepositoryInterface,
	banRepo repository.BanRepositoryInterface,
	config *config.Config,
	logger *logger.Logger,
) service.AuthService {
	return service.NewAuthService(
		userRepo,
		sessionRepo,
		resetPasswordRepo,
		user2FASecretRepo,
		banRepo,
		config,
		logger,
	)
}

// Handler providers
func ProvideAuthHandler(authService service.AuthService, logger *logger.Logger) *handler.AuthHandler {
	return handler.NewAuthHandler(authService, logger)
}

// Middleware providers
func ProvideAuthMiddleware(authService service.AuthService) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(authService)
}

// Wire provider set for auth module
var AuthProviderSet = wire.NewSet(
	// Repositories
	ProvideUserRepository,
	ProvideSessionRepository,
	ProvideResetPasswordRepository,
	ProvideUser2FASecretRepository,
	ProvideBanRepository,
	
	// Services
	ProvideAuthService,
	
	// Handlers
	ProvideAuthHandler,
	
	// Middleware
	ProvideAuthMiddleware,
)