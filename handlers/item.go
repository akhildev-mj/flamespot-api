package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/storage/memory/v2"
)

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Image    string  `json:"image"`
	Category string  `json:"category"`
}

var itemsData = []Item{
	{
		Name:     "Shawai",
		Price:    550.0,
		Image:    "https://ik.imagekit.io/akhildev/flamespot/shawai.jpeg",
		Category: "Mains",
	},
	{
		Name:     "Alfaham",
		Price:    480.0,
		Image:    "https://ik.imagekit.io/akhildev/flamespot/alfaham.jpeg",
		Category: "Mains",
	},
	{
		Name:     "Kuboos",
		Price:    15.0,
		Image:    "https://ik.imagekit.io/akhildev/flamespot/kuboos.jpeg",
		Category: "Breads",
	},
	{
		Name:     "Parotta",
		Price:    20.0,
		Image:    "https://ik.imagekit.io/akhildev/flamespot/parotta.jpeg",
		Category: "Breads",
	},
	{
		Name:     "Lime Juice",
		Price:    25.0,
		Image:    "https://ik.imagekit.io/akhildev/flamespot/lime.jpeg",
		Category: "Drinks",
	},
}

var itemCacheStore = memory.New()

func RegisterItemRoutes(router fiber.Router) {
	items := router.Group("/items")

	cache := cache.New(cache.Config{
		Expiration:          30 * 24 * time.Hour,
		DisableCacheControl: false,
		Storage:             itemCacheStore,
	})

	items.Get("/", cache, getAllItems)
}

func getAllItems(c fiber.Ctx) error {
	return c.JSON(itemsData)
}
