package app

import (
	"strconv"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// Sets up the application and runs the HTTP server
func SetupAndRunApp() error {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		return err
	}

	// Setup logger
	config.SetupLogger()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1 GB limit
	})

	// Attach middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// Setup routes
	SetupRoutes(app)

	// Get the port from the config and start the server
	port := config.AppEnv.Port
	zap.L().Info("Starting server", zap.Int("port", port))
	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		zap.L().Error("Failed to start server", zap.Error(err))
		return err
	}

	return nil
}
