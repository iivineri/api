package fiber

import (
	"context"
	"fmt"
	"iivineri/internal/config"
	"iivineri/internal/fiber/shared/middleware"
	"iivineri/internal/logger"
	"iivineri/internal/metrics"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type ServerInterface interface {
	Start() error
	Shutdown(ctx context.Context) error
	RegisterRoutes(fn func(app *fiber.App))
	GetApp() *fiber.App
}

type Server struct {
	app     *fiber.App
	config  *config.AppConfig
	logger  *logger.Logger
	metrics metrics.MetricsInterface
}

func NewServer(
	cfg *config.AppConfig,
	log *logger.Logger,
	metrics metrics.MetricsInterface,
) ServerInterface {
	// Fiber config
	fiberConfig := fiber.Config{
		AppName:               "Iivineri API",
		ServerHeader:          "Iivineri",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		Prefork:               cfg.Prefork(),
		ErrorHandler:          errorHandler(log),
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           60 * time.Second,
	}

	app := fiber.New(fiberConfig)

	server := &Server{
		app:     app,
		config:  cfg,
		logger:  log,
		metrics: metrics,
	}

	// Setup middleware
	server.setupMiddleware()

	return server
}

func (s *Server) setupMiddleware() {
	// Request ID middleware
	s.app.Use(requestid.New())

	// Metrics middleware (before logger to capture all requests)
	s.app.Use(metrics.MetricsMiddleware(s.metrics))
	s.app.Use(metrics.SystemMetricsMiddleware(s.metrics))
	s.app.Use(metrics.HealthMetricsMiddleware(s.metrics))

	// Custom logger middleware
	s.app.Use(middleware.LoggerMiddleware(s.logger))

	// Recovery middleware
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: s.config.IsDevelopment(),
	}))

	// Security middleware
	s.app.Use(helmet.New())

	// CORS middleware
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Compression middleware
	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Rate limiting (only in production)
	if s.config.IsProduction() {
		s.app.Use(limiter.New(limiter.Config{
			Max:        100,
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
		}))
	}

	// Metrics endpoint
	s.app.Get("/metrics", metrics.PrometheusHandler())

	// Health check endpoint
	s.app.Get("/health", s.healthCheck)
}

func (s *Server) RegisterRoutes(fn func(app *fiber.App)) {
	fn(s.app)
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port())

	s.logger.Tracef("Starting Fiber server on %s", addr)
	s.logger.Tracef("Environment: %s", s.config.Environment())
	s.logger.Tracef("Prefork: %v", s.config.Prefork())

	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		done <- s.app.Shutdown()
	}()

	select {
	case err := <-done:
		if err != nil {
			s.logger.WithError(err).Error("Error during server shutdown")
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":      "ok",
		"timestamp":   time.Now().Unix(),
		"environment": s.config.Environment(),
	})
}

func errorHandler(log *logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		log.WithError(err).
			WithField("method", c.Method()).
			WithField("path", c.Path()).
			WithField("ip", c.IP()).
			Error("HTTP error occurred")

		return c.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": message,
		})
	}
}
