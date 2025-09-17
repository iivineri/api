package middleware

import (
	"iivineri/internal/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		entry := log.WithFields(map[string]interface{}{
			"method":      c.Method(),
			"path":        c.Path(),
			"status":      c.Response().StatusCode(),
			"duration_ms": duration.Milliseconds(),
			"ip":          c.IP(),
			"user_agent":  c.Get("User-Agent"),
			"request_id":  c.Locals("requestid"),
		})

		if err != nil {
			entry.WithError(err).Error("Request failed")
		} else {
			entry.Info("Request completed")
		}

		return err
	}
}
