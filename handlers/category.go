package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

var categoryCacheStore = memory.New()

func RegisterCategoryRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             categoryCacheStore,
	})

	categories.Get("/", cache, getAllCategories)
}

func getAllCategories(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all categories"})
}
