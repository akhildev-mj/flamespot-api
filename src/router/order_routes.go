package router

import (
	"flamespot-api/src/config"
	"flamespot-api/src/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

func SetupOrderRoutes(api fiber.Router, orderController *controller.OrderController, cacheStore *memory.Storage) {
	orderCache := cache.New(cache.Config{
		Expiration:          config.CacheExpiration,
		DisableCacheControl: false,
		Storage:             cacheStore,
	})

	orders := api.Group("/orders")
	orders.Get("/", orderCache, orderController.GetAllOrders)
	orders.Post("/", orderController.CreateOrder)
	orders.Patch("/:id", orderController.UpdateOrder)
}
