package handlers

import "github.com/gofiber/fiber/v3"

func RegisterOrderRoutes(router fiber.Router) {
	orders := router.Group("/orders")

	orders.Get("/", getAllOrders)
}

func getAllOrders(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get all orders"})
}
