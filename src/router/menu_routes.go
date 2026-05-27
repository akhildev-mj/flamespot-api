package router

import (
	"flamespot-api/src/config"
	"flamespot-api/src/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

func SetupMenuRoutes(api fiber.Router, menuController *controller.MenuController, cacheStore *memory.Storage) {
	menuCache := cache.New(cache.Config{
		Expiration:          config.CacheExpiration,
		DisableCacheControl: false,
		Storage:             cacheStore,
	})

	menu := api.Group("/menu")
	menu.Get("/", menuCache, menuController.GetMenu)
	menu.Post("/", menuController.CreateMenu)
	menu.Patch("/:id", menuController.UpdateMenu)
}
