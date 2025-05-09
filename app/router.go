package app

import (
	"github.com/MRegterschot/trackmania-server-fm/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/UserData/*", handlers.HandleListFiles)
	app.Post("/upload", handlers.HandleUploadFiles)
	app.Delete("/delete", handlers.HandleDeleteFiles)

	app.Get("/scripts", handlers.HandleListScripts)
}
