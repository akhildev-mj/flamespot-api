package controller

import (
	"flamespot-api/response"
	"flamespot-api/service"

	"github.com/gofiber/fiber/v3"
)

type MenuController struct {
	menuService service.MenuService
}

func NewMenuController(menuService service.MenuService) *MenuController {
	return &MenuController{menuService: menuService}
}

func (c *MenuController) GetMenu(ctx fiber.Ctx) error {
	menuItems, err := c.menuService.GetMenu()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch menu from the database",
		})
	}

	return ctx.JSON(menuItems)
}
