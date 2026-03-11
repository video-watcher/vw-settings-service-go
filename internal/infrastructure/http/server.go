package http

import (
	"context"
	"fmt"
	"settings-service-go/internal/application/usecase"
	"settings-service-go/internal/config"
	"settings-service-go/internal/infrastructure/http/handler"
	"settings-service-go/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

// Server represents the HTTP server
type Server struct {
	app                   *fiber.App
	config                *config.Config
	userSettingsHandler   *handler.UserSettingsHandler
	systemSettingsHandler *handler.SystemSettingsHandler
}

// NewServer creates a new HTTP server
func NewServer(
	cfg *config.Config,
	userSettingsUseCase *usecase.UserSettingsUseCase,
	systemSettingsUseCase *usecase.SystemSettingsUseCase,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "Settings Service",
		ServerHeader: "Settings Service",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(middleware.Logger())

	// Handlers
	userSettingsHandler := handler.NewUserSettingsHandler(userSettingsUseCase)
	systemSettingsHandler := handler.NewSystemSettingsHandler(systemSettingsUseCase)

	server := &Server{
		app:                   app,
		config:                cfg,
		userSettingsHandler:   userSettingsHandler,
		systemSettingsHandler: systemSettingsHandler,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Health check
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "settings-service",
		})
	})

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "settings-service",
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})

	// API routes
	api := s.app.Group("/api")
	settings := api.Group("/settings")

	// JWT middleware for authenticated routes
	jwtMiddleware := middleware.JWTAuth(s.config.JWT.Secret)

	// User settings routes (require authentication)
	settings.Get("/me", jwtMiddleware, s.userSettingsHandler.GetUserSettings)
	settings.Patch("/me", jwtMiddleware, s.userSettingsHandler.UpdateUserSettings)

	// System settings routes (require admin)
	adminMiddleware := middleware.RequireAdmin()
	system := settings.Group("/system", jwtMiddleware, adminMiddleware)
	system.Get("/", s.systemSettingsHandler.ListSystemSettings)
	system.Get("/:key", s.systemSettingsHandler.GetSystemSetting)
	system.Put("/:key", s.systemSettingsHandler.SetSystemSetting)
	system.Delete("/:key", s.systemSettingsHandler.DeleteSystemSetting)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	log.Info().Msgf("Starting server on %s", addr)
	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server...")
	return s.app.ShutdownWithContext(ctx)
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Error().Err(err).Int("status", code).Msg("HTTP error")

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
