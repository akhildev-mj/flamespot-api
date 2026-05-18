package handlers

import "github.com/gofiber/fiber/v3"

func RegisterCategoryRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	categories.Get("/", getAllCategories)
}

func getAllCategories(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all categories"})
}
