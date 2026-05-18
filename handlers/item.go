package handlers

import "github.com/gofiber/fiber/v3"

func RegisterItemRoutes(router fiber.Router) {
	items := router.Group("/items")

	items.Get("/", getAllItems)
}

func getAllItems(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all items"})
}
