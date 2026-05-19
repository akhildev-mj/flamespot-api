package router

import (
	"flamespot-api/config"
	"flamespot-api/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

func SetupMenuRoutes(api fiber.Router, menuController *controller.MenuController) {
	menuCacheStore := memory.New()

	menuCache := cache.New(cache.Config{
		Expiration:          config.CacheExpiration,
		DisableCacheControl: false,
		Storage:             menuCacheStore,
	})

	menu := api.Group("/menu")
	menu.Get("/", menuCache, menuController.GetMenu)
}
