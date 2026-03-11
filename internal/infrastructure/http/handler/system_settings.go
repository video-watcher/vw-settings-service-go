package handler

import (
	"settings-service-go/internal/application/dto"
	"settings-service-go/internal/application/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// SystemSettingsHandler handles system settings HTTP requests
type SystemSettingsHandler struct {
	useCase  *usecase.SystemSettingsUseCase
	validate *validator.Validate
}

// NewSystemSettingsHandler creates a new SystemSettingsHandler
func NewSystemSettingsHandler(useCase *usecase.SystemSettingsUseCase) *SystemSettingsHandler {
	return &SystemSettingsHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// ListSystemSettings handles GET /api/settings/system
func (h *SystemSettingsHandler) ListSystemSettings(c *fiber.Ctx) error {
	skip, _ := strconv.Atoi(c.Query("skip", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "100"))

	if limit > 100 {
		limit = 100
	}

	settings, err := h.useCase.GetAllSettings(c.Context(), skip, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// GetSystemSetting handles GET /api/settings/system/:key
func (h *SystemSettingsHandler) GetSystemSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "key is required",
		})
	}

	setting, err := h.useCase.GetSetting(c.Context(), key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if setting == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "setting not found",
		})
	}

	return c.JSON(setting)
}

// SetSystemSetting handles PUT /api/settings/system/:key
func (h *SystemSettingsHandler) SetSystemSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "key is required",
		})
	}

	var req dto.SystemSettingCreateRequest
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

	setting, err := h.useCase.SetSetting(c.Context(), key, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(setting)
}

// DeleteSystemSetting handles DELETE /api/settings/system/:key
func (h *SystemSettingsHandler) DeleteSystemSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "key is required",
		})
	}

	if err := h.useCase.DeleteSetting(c.Context(), key); err != nil {
		if err.Error() == "setting not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
