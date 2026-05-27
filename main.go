package main

import (
	"context"
	"log"
	"strings"

	"flamespot-api/src/config"
	"flamespot-api/src/database"
	"flamespot-api/src/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg := config.LoadConfig()

	client, db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting database: %v", err)
		}
	}()

	app := fiber.New(fiber.Config{
		BodyLimit:    config.MaxBodySize,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: strings.Split(config.CORSMethods, ","),
		AllowHeaders: strings.Split(config.CORSHeaders, ","),
	}))

	router.SetupRoutes(app, db, cfg)

	log.Printf("Server started successfully and is listening on port %s", cfg.Port)

	if err := app.Listen(":"+cfg.Port, fiber.ListenConfig{
		EnablePrefork:         false,
		DisableStartupMessage: true,
	}); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
