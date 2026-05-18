package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"flamespot-api/database"
	"flamespot-api/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    1 * 1024 * 1024,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("Fatal Error: MONGODB_URI environment variable is not set!")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Successfully connected to MongoDB Atlas!")

	db := client.Database("flamespot_db")
	database.RunMigrationsAndSeed(db)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "healthy",
		})
	})

	v1 := app.Group("/api/v1")

	handlers.RegisterMenuRoutes(v1, db)
	handlers.RegisterOrderRoutes(v1, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Flame Spot API server booting on port %s...", port)

	if err := app.Listen(":"+port, fiber.ListenConfig{
		EnablePrefork:         false,
		DisableStartupMessage: true,
	}); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
