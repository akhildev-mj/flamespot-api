package controller

import (
	"flamespot-api/src/response"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
)

type LookupController struct {
	lookupService service.LookupService
}

func NewLookupController(lookupService service.LookupService) *LookupController {
	return &LookupController{lookupService: lookupService}
}

func (c *LookupController) GetActiveCategories(ctx fiber.Ctx) error {
	cats, err := c.lookupService.GetActiveCategories()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch category lookups",
		})
	}
	return ctx.JSON(cats)
}

func (c *LookupController) GetConfigEnums(ctx fiber.Ctx) error {
	configData, err := c.lookupService.GetConfigEnums()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch config enums",
		})
	}
	return ctx.JSON(configData)
}
