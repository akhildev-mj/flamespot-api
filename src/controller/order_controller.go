package controller

import (
	"strconv"

	"flamespot-api/src/config"
	"flamespot-api/src/model"
	"flamespot-api/src/response"
	"flamespot-api/src/service"

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

	statusFilter := ctx.Query("status", string(model.StatusAll))
	if !model.Status(statusFilter).IsValidFilter() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid status filter",
		})
	}

	typeFilter := ctx.Query("type", string(model.TypeAll))
	if !model.Type(typeFilter).IsValidFilter() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid type filter",
		})
	}

	sortFilter := ctx.Query("sort")
	if sortFilter == "" {
		sortFilter = string(model.SortRecent)
	} else if !model.Sort(sortFilter).IsValidForOrder() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid sort parameter. Allowed values: RECENT, PRICE",
		})
	}

	orders, err := c.orderService.GetOrders(limit, cursor, statusFilter, typeFilter, sortFilter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch orders",
		})
	}

	return ctx.JSON(orders)
}

func (c *OrderController) CreateOrder(ctx fiber.Ctx) error {
	var order model.Order

	if err := ctx.Bind().JSON(&order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if order.ID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Order ID is required. The frontend must generate and provide the ID.",
		})
	}

	if !order.Type.IsValid() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid type. Must be TAKEAWAY or DINE_IN",
		})
	}

	if order.Status != "" && !order.Status.IsValid() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid order status value",
		})
	}

	_, err := c.orderService.CreateOrder(order)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to save order",
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusCreated).Send(nil)
}

func (c *OrderController) UpdateOrder(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var updates map[string]interface{}

	if err := ctx.Bind().JSON(&updates); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if statusVal, ok := updates["status"].(string); ok {
		if !model.Status(statusVal).IsValid() {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: "Invalid status value",
			})
		}
	}

	if typeVal, ok := updates["type"].(string); ok {
		if !model.Type(typeVal).IsValid() {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: "Invalid type value",
			})
		}
	}

	_, err := c.orderService.UpdateOrder(id, updates)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: err.Error(),
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
