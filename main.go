package main

import (
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

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

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
