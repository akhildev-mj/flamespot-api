package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ItemHandler struct {
	db *mongo.Database
}

type Item struct {
	ID       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Price    float64       `json:"price" bson:"price"`
	Image    string        `json:"image" bson:"image"`
	Category string        `json:"category" bson:"category"`
}

var itemCacheStore = memory.New()

func RegisterItemRoutes(router fiber.Router, db *mongo.Database) {
	handler := &ItemHandler{
		db: db,
	}

	items := router.Group("/items")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             itemCacheStore,
	})

	items.Get("/", cache, handler.getAllItems)
}

func (h *ItemHandler) getAllItems(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := h.db.Collection("items")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch items from the database",
		})
	}
	defer cursor.Close(ctx)

	var items []Item
	if err := cursor.All(ctx, &items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse items",
		})
	}

	if items == nil {
		items = []Item{}
	}

	return c.JSON(items)
}
