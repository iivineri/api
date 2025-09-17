package serve

import (
	"context"
	"fmt"
	"iivineri/internal/container"
	"iivineri/internal/fiber/modules/auth"
	"iivineri/internal/wire"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "iivineri/runtime/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// configureSwagger updates the Swagger configuration with values from environment
func configureSwagger(container *container.Container) {
	docs.SwaggerInfo.Host = container.Config.App.SwaggerHost()
	docs.SwaggerInfo.BasePath = container.Config.App.SwaggerBasePath()
	docs.SwaggerInfo.Schemes = container.Config.App.SwaggerSchemes()
}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Initialize container with Wire
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}

		// Configure Swagger with environment values
		configureSwagger(container)

		// Ensure cleanup on exit
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		// Connect to database
		if err := container.ConnectDatabase(ctx); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		// Create a context that cancels on interrupt signals
		ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		defer cancel()

		// Health check before starting server
		if err := container.HealthCheck(ctx); err != nil {
			return fmt.Errorf("health check failed: %w", err)
		}

		// Start system metrics collector
		go container.SystemCollector.Start(ctx)

		// Setup routes
		container.Server.RegisterRoutes(func(app *fiber.App) {

			fmt.Println(container.Config.App.Environment())

			if container.Config.App.SwaggerEnabled() {
				app.Get("/swagger/*", fiberSwagger.WrapHandler)
			}

			auth.RegisterRoutes(app, container.AuthHandler, container.AuthMiddleware)

			// API routes
			api := app.Group("/")
			api.Get("/", GetAPIInfo(container))
			api.Get("/health", GetHealthCheck(container))
		})

		// Start Fiber server in a goroutine
		serverErr := make(chan error, 1)
		go func() {
			if err := container.Server.Start(); err != nil {
				serverErr <- fmt.Errorf("failed to start server: %w", err)
			}
		}()

		container.Logger.Info("Server started successfully")

		// Wait for either interrupt signal or server error
		select {
		case <-ctx.Done():
			container.Logger.Info("Received shutdown signal, shutting down...")
		case err := <-serverErr:
			return err
		}

		return nil
	},
}

// GetAPIInfo godoc
// @Summary Get API information
// @Description Returns basic information about the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API information"
// @Router / [get]
func GetAPIInfo(container *container.Container) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record business metric
		container.Metrics.IncrementUserAction("api_request")

		return c.JSON(fiber.Map{
			"message": "Iivineri API v1",
			"version": "1.0.0",
		})
	}
}

// GetHealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health check passed"
// @Failure 503 {object} map[string]interface{} "Service unavailable"
// @Router /health [get]
func GetHealthCheck(container *container.Container) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check database health
		if err := container.HealthCheck(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":    "unhealthy",
				"error":     "database connection failed",
				"timestamp": time.Now().UTC(),
			})
		}

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
			"services": fiber.Map{
				"database": "healthy",
			},
		})
	}
}
