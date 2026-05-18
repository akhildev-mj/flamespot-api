package main

import (
	"flamespot-api/handlers"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:    1 * 1024 * 1024,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "healthy",
		})
	})

	v1 := app.Group("/api/v1")
	handlers.RegisterItemRoutes(v1)
	handlers.RegisterCategoryRoutes(v1)
	handlers.RegisterOrderRoutes(v1)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Flame Spot API server booting on port %s...", port)

	err := app.Listen(":"+port, fiber.ListenConfig{
		EnablePrefork:         false,
		DisableStartupMessage: true,
	})

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
