package handler

import (
	"settings-service-go/internal/application/dto"
	"settings-service-go/internal/application/usecase"
	"settings-service-go/internal/infrastructure/http/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserSettingsHandler handles user settings HTTP requests
type UserSettingsHandler struct {
	useCase  *usecase.UserSettingsUseCase
	validate *validator.Validate
}

// NewUserSettingsHandler creates a new UserSettingsHandler
func NewUserSettingsHandler(useCase *usecase.UserSettingsUseCase) *UserSettingsHandler {
	return &UserSettingsHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// GetUserSettings handles GET /api/settings/me
func (h *UserSettingsHandler) GetUserSettings(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user_id not found in token",
		})
	}

	settings, err := h.useCase.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// UpdateUserSettings handles PATCH /api/settings/me
func (h *UserSettingsHandler) UpdateUserSettings(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user_id not found in token",
		})
	}

	var req dto.UserSettingsUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	settings, err := h.useCase.UpdateSettings(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}
