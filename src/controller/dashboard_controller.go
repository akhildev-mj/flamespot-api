package controller

import (
	"flamespot-api/src/response"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
)

type DashboardController struct {
	dashboardService service.DashboardService
}

func NewDashboardController(dashboardService service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

func (c *DashboardController) GetSummary(ctx fiber.Ctx) error {
	summary, err := c.dashboardService.GetSummary()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to generate dashboard summary",
		})
	}
	return ctx.JSON(summary)
}

func (c *DashboardController) ExportDashboard(ctx fiber.Ctx) error {
	csvData, err := c.dashboardService.ExportDashboard()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to export dashboard data",
		})
	}

	ctx.Set("Content-Type", "text/csv")
	ctx.Set("Content-Disposition", "attachment; filename=\"dashboard_export.csv\"")
	return ctx.Send(csvData)
}
