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

type MenuHandler struct {
	db *mongo.Database
}

// Unified: This struct perfectly matches Android and MongoDB now
type MenuItem struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Price    float64       `json:"price" bson:"price"`
	URL      string        `json:"url" bson:"url"` // 1:1 mapping
	Category string        `json:"category" bson:"category"`
}

var menuCacheStore = memory.New()

func RegisterMenuRoutes(router fiber.Router, db *mongo.Database) {
	handler := &MenuHandler{
		db: db,
	}

	// Unified endpoint
	menu := router.Group("/menu")

	menuCache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             menuCacheStore,
	})

	menu.Get("/", menuCache, handler.getMenu)
}

func (h *MenuHandler) getMenu(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unified collection
	collection := h.db.Collection("menu")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch menu from the database",
		})
	}
	defer cursor.Close(ctx)

	var menuItems []MenuItem
	if err := cursor.All(ctx, &menuItems); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse menu items",
		})
	}

	if menuItems == nil {
		menuItems = []MenuItem{}
	}

	return c.JSON(menuItems)
}
