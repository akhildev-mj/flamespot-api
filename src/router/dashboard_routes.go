package router

import (
	"flamespot-api/src/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupDashboardRoutes(api fiber.Router, dashboardController *controller.DashboardController) {
	dashboard := api.Group("/dashboard")
	dashboard.Get("/summary", dashboardController.GetSummary)

	export := api.Group("/export")
	export.Get("/dashboard", dashboardController.ExportDashboard)
}
