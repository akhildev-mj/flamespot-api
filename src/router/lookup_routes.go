package router

import (
	"flamespot-api/src/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupLookupRoutes(api fiber.Router, lookupController *controller.LookupController) {
	lookups := api.Group("/lookups")
	lookups.Get("/categories", lookupController.GetActiveCategories)

	config := api.Group("/config")
	config.Get("/enums", lookupController.GetConfigEnums)
}
