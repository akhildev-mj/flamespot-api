package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OrderHandler struct {
	db *mongo.Database
}

type CartItem struct {
	MenuItem MenuItem `json:"menuItem" bson:"menuItem"`
	Quantity int      `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID        string     `json:"id" bson:"_id"`
	TableName string     `json:"tableName" bson:"tableName"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"total"`
	Status    string     `json:"status" bson:"status"`
	Timestamp int64      `json:"timestamp" bson:"timestamp"`
}

var orderCacheStore = memory.New()

func RegisterOrderRoutes(router fiber.Router, db *mongo.Database) {
	handler := &OrderHandler{
		db: db,
	}

	orders := router.Group("/orders")

	orderCache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             orderCacheStore,
	})

	orders.Get("/", orderCache, handler.getAllOrders)
	orders.Post("/", handler.createOrder)
}

func (h *OrderHandler) createOrder(c fiber.Ctx) error {
	var order Order

	if err := c.Bind().JSON(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.Collection("orders")

	opts := options.Replace().SetUpsert(true)

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save/update order",
		})
	}

	_ = orderCacheStore.Reset()

	return c.SendStatus(fiber.StatusCreated)
}

func (h *OrderHandler) getAllOrders(c fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	cursorStr := c.Query("cursor", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}

	cursor, err := strconv.ParseInt(cursorStr, 10, 64)
	if err != nil {
		cursor = 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.Collection("orders")

	filter := bson.M{}
	if cursor > 0 {
		filter = bson.M{"timestamp": bson.M{"$lt": cursor}}
	}

	opts := options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(limit)

	cursorResult, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}
	defer cursorResult.Close(ctx)

	var orders []Order
	if err := cursorResult.All(ctx, &orders); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse orders",
		})
	}

	if orders == nil {
		orders = []Order{}
	}

	return c.JSON(orders)
}
