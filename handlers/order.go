package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderHandler struct {
	db *mongo.Database
}

var orderCacheStore = memory.New()

func RegisterOrderRoutes(router fiber.Router, db *mongo.Database) {
	handler := &OrderHandler{
		db: db,
	}

	orders := router.Group("/orders")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             orderCacheStore,
	})

	orders.Get("/", cache, handler.getAllOrders)
}

func (h *OrderHandler) getAllOrders(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all orders"})
}
