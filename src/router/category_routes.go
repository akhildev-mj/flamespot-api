package router

import (
	"flamespot-api/src/config"
	"flamespot-api/src/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

func SetupCategoryRoutes(api fiber.Router, ctrl *controller.CategoryController, cacheStore *memory.Storage) {
	catCache := cache.New(cache.Config{
		Expiration:          config.CacheExpiration,
		DisableCacheControl: false,
		Storage:             cacheStore,
	})

	catGroup := api.Group("/categories")
	catGroup.Get("/", catCache, ctrl.GetCategories)
	catGroup.Post("/", ctrl.CreateCategory)
	catGroup.Patch("/:id", ctrl.UpdateCategory)
}
