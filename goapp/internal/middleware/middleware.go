package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/yash/user-dob-api/internal/logger"
)

const RequestIDHeader = "X-Request-Id"

// RequestID injects a unique request ID into every response header.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(RequestIDHeader, id)
		c.Locals("requestId", id)
		return c.Next()
	}
}

// Logger logs request method, path, status code, and duration using Uber Zap.
func Logger() fiber.Handler {
	log := logger.Get()
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("request_id", c.Locals("requestId").(string)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		}

		if status >= 500 {
			log.Error("request completed with server error", fields...)
		} else if status >= 400 {
			log.Warn("request completed with client error", fields...)
		} else {
			log.Info("request completed", fields...)
		}

		return err
	}
}

// Recover catches panics and returns a 500 response instead of crashing.
func Recover() fiber.Handler {
	log := logger.Get()
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Path()),
				)
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}()
		return c.Next()
	}
}
