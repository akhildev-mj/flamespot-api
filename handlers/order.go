package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

var orderCacheStore = memory.New()

func RegisterOrderRoutes(router fiber.Router) {
	orders := router.Group("/orders")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             orderCacheStore,
	})

	orders.Get("/", cache, getAllOrders)
}

func getAllOrders(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all orders"})
}
