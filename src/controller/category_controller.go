package controller

import (
	"strconv"

	"flamespot-api/src/model"
	"flamespot-api/src/response"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
)

type CategoryController struct {
	catService service.CategoryService
	cache      *memory.Storage
}

func NewCategoryController(catService service.CategoryService, cache *memory.Storage) *CategoryController {
	return &CategoryController{
		catService: catService,
		cache:      cache,
	}
}

func (c *CategoryController) GetCategories(ctx fiber.Ctx) error {
	defaultActive := true
	isActiveFilter := &defaultActive

	if activeStr := ctx.Query("isActive"); activeStr != "" {
		if val, err := strconv.ParseBool(activeStr); err == nil {
			isActiveFilter = &val
		}
	}

	sortFilter := ctx.Query("sort")
	if sortFilter == "" {
		sortFilter = string(model.SortName)
	} else if !model.Sort(sortFilter).IsValidForCategory() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid sort parameter. Allowed values: RECENT, NAME",
		})
	}

	cats, err := c.catService.GetCategories(isActiveFilter, sortFilter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to fetch categories",
		})
	}
	return ctx.JSON(cats)
}

func (c *CategoryController) CreateCategory(ctx fiber.Ctx) error {
	var item model.Category

	if err := ctx.Bind().JSON(&item); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	item.IsActive = true

	_, err := c.catService.CreateCategory(item)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to create category",
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusCreated).Send(nil)
}

func (c *CategoryController) UpdateCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var updates map[string]interface{}

	if err := ctx.Bind().JSON(&updates); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	_, err := c.catService.UpdateCategory(id, updates)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to update category",
		})
	}

	_ = c.cache.Reset()
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
