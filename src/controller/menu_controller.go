package controller

import (
	"strconv"

	"flamespot-api/src/model"
	"flamespot-api/src/response"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
)

type MenuController struct {
	menuService service.MenuService
	cache       *memory.Storage
}

func NewMenuController(menuService service.MenuService, cache *memory.Storage) *MenuController {
	return &MenuController{
		menuService: menuService,
		cache:       cache,
	}
}

func (c *MenuController) GetMenu(ctx fiber.Ctx) error {
	defaultActive := true
	isActiveFilter := &defaultActive

	if activeStr := ctx.Query("isActive"); activeStr != "" {
		if val, err := strconv.ParseBool(activeStr); err == nil {
			isActiveFilter = &val
		}
	}

	categoryIdFilter := ctx.Query("categoryId")

	sortFilter := ctx.Query("sort")
	if sortFilter == "" {
		sortFilter = string(model.SortPrice)
	} else if !model.Sort(sortFilter).IsValidForMenu() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid sort parameter. Allowed values: RECENT, NAME, PRICE",
		})
	}

	menuItems, err := c.menuService.GetMenu(isActiveFilter, categoryIdFilter, sortFilter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch menu",
		})
	}
	return ctx.JSON(menuItems)
}

func (c *MenuController) CreateMenu(ctx fiber.Ctx) error {
	var item model.MenuItem

	if err := ctx.Bind().JSON(&item); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	item.IsActive = true

	_, err := c.menuService.CreateMenuItem(item)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to create menu item",
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusCreated).Send(nil)
}

func (c *MenuController) UpdateMenu(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var updates map[string]interface{}

	if err := ctx.Bind().JSON(&updates); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	_, err := c.menuService.UpdateMenuItem(id, updates)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to update menu item",
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
