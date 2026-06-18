package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver registered as "pgx"

	"github.com/yash/user-dob-api/config"
	"github.com/yash/user-dob-api/internal/handler"
	"github.com/yash/user-dob-api/internal/logger"
	"github.com/yash/user-dob-api/internal/middleware"
	"github.com/yash/user-dob-api/internal/repository"
	"github.com/yash/user-dob-api/internal/routes"
	"github.com/yash/user-dob-api/internal/service"
)

func main() {
	// ── Logger ──────────────────────────────────────────────────────────────
	isProd := os.Getenv("APP_ENV") == "production"
	if err := logger.Init(isProd); err != nil {
		log.Fatalf("failed to initialise logger: %v", err)
	}
	defer logger.Sync()

	// ── Config ───────────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Database ─────────────────────────────────────────────────────────────
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// ── Wire dependencies ────────────────────────────────────────────────────
	repo    := repository.New(db)
	svc     := service.New(repo)
	handler := handler.New(svc)

	// ── Fiber app ────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "user-dob-api",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())

	// Routes
	routes.Register(app, handler)

	// Health-check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// customErrorHandler returns a JSON error for unhandled Fiber errors.
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var fe *fiber.Error
	if ok := fiber.As(err, &fe); ok {
		code = fe.Code
	}
	return c.Status(code).JSON(fiber.Map{"message": err.Error()})
}
