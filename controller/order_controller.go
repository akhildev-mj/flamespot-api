package controller

import (
	"strconv"

	"flamespot-api/config"
	"flamespot-api/model"
	"flamespot-api/response"
	"flamespot-api/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
)

type OrderController struct {
	orderService service.OrderService
	cache        *memory.Storage
}

func NewOrderController(orderService service.OrderService, cache *memory.Storage) *OrderController {
	return &OrderController{
		orderService: orderService,
		cache:        cache,
	}
}

func (c *OrderController) CreateOrder(ctx fiber.Ctx) error {
	var order model.Order

	if err := ctx.Bind().JSON(&order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := c.orderService.CreateOrder(order); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to save/update order",
		})
	}

	_ = c.cache.Reset()

	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *OrderController) GetAllOrders(ctx fiber.Ctx) error {
	limitStr := ctx.Query("limit", strconv.Itoa(config.LimitDefault))
	cursorStr := ctx.Query("cursor", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = config.LimitDefault
	}

	cursor, err := strconv.ParseInt(cursorStr, 10, 64)
	if err != nil {
		cursor = 0
	}

	orders, err := c.orderService.GetOrders(limit, cursor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch orders",
		})
	}

	return ctx.JSON(orders)
}
