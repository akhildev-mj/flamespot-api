package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CategoryHandler struct {
	db *mongo.Database
}

var categoryCacheStore = memory.New()

func RegisterCategoryRoutes(router fiber.Router, db *mongo.Database) {
	handler := &CategoryHandler{
		db: db,
	}

	categories := router.Group("/categories")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             categoryCacheStore,
	})

	categories.Get("/", cache, handler.getAllCategories)
}

func (h *CategoryHandler) getAllCategories(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all categories"})
}
