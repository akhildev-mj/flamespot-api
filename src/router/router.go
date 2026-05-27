package router

import (
	"flamespot-api/src/config"
	"flamespot-api/src/controller"
	"flamespot-api/src/middleware"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database, cfg config.Config) {
	cacheStoreCats := memory.New()
	cacheStoreMenu := memory.New()
	cacheStoreOrders := memory.New()

	authService := service.NewAuthService(db, cfg)
	catService := service.NewCategoryService(db)
	menuService := service.NewMenuService(db)
	orderService := service.NewOrderService(db)
	dashboardService := service.NewDashboardService(db)
	lookupService := service.NewLookupService(db)

	authController := controller.NewAuthController(authService)
	catController := controller.NewCategoryController(catService, cacheStoreCats)
	menuController := controller.NewMenuController(menuService, cacheStoreMenu)
	orderController := controller.NewOrderController(orderService, cacheStoreOrders)
	dashboardController := controller.NewDashboardController(dashboardService)
	lookupController := controller.NewLookupController(lookupService)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/api/v1")
	v1.Post("/login", authController.Login)

	v1.Use(middleware.Protected(cfg.JWTSecret))

	SetupCategoryRoutes(v1, catController, cacheStoreCats)
	SetupMenuRoutes(v1, menuController, cacheStoreMenu)
	SetupOrderRoutes(v1, orderController, cacheStoreOrders)
	SetupDashboardRoutes(v1, dashboardController)
	SetupLookupRoutes(v1, lookupController)
}
