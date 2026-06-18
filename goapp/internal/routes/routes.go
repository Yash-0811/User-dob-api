package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yash/user-dob-api/internal/handler"
)

// Register mounts all application routes onto the Fiber app.
func Register(app *fiber.App, h *handler.UserHandler) {
	users := app.Group("/users")

	users.Post("/", h.CreateUser)
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}
