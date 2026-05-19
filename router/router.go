package router

import (
	"flamespot-api/config"
	"flamespot-api/controller"
	"flamespot-api/middleware"
	"flamespot-api/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database, cfg config.Config) {
	orderCacheStore := memory.New()

	// Services
	authService := service.NewAuthService(db, cfg)
	menuService := service.NewMenuService(db)
	orderService := service.NewOrderService(db)

	// Controllers
	authController := controller.NewAuthController(authService)
	menuController := controller.NewMenuController(menuService)
	orderController := controller.NewOrderController(orderService, orderCacheStore)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	v1 := app.Group("/api/v1")

	v1.Post("/login", authController.Login)

	v1.Use(middleware.Protected(cfg.JWTSecret))

	SetupMenuRoutes(v1, menuController)
	SetupOrderRoutes(v1, orderController, orderCacheStore)
}
