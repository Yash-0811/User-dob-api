package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/yash/user-dob-api/internal/logger"
	"github.com/yash/user-dob-api/internal/models"
	"github.com/yash/user-dob-api/internal/service"
)

// UserHandler holds dependencies for the user HTTP handlers.
type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
	log      *zap.Logger
}

// New returns a new UserHandler.
func New(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: validator.New(),
		log:      logger.Get(),
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp, err := h.svc.Create(c.Context(), req)
	if err != nil {
		h.log.Error("failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "failed to create user",
		})
	}

	h.log.Info("user created", zap.Int32("id", resp.ID), zap.String("name", resp.Name))
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "invalid user id"})
	}

	resp, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "user not found"})
		}
		h.log.Error("failed to get user", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "failed to get user"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "invalid user id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Message: err.Error()})
	}

	resp, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "user not found"})
		}
		h.log.Error("failed to update user", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "failed to update user"})
	}

	h.log.Info("user updated", zap.Int32("id", resp.ID))
	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "invalid user id"})
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "user not found"})
		}
		h.log.Error("failed to delete user", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "failed to delete user"})
	}

	h.log.Info("user deleted", zap.Int32("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers handles GET /users with optional pagination query params.
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	resp, err := h.svc.List(c.Context(), page, pageSize)
	if err != nil {
		h.log.Error("failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "failed to list users"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// parseID extracts and validates the :id URL parameter.
func parseID(c *fiber.Ctx) (int32, error) {
	raw := c.Params("id")
	n, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || n <= 0 {
		return 0, errors.New("invalid id")
	}
	return int32(n), nil
}
