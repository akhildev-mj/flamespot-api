package controller

import (
	"flamespot-api/src/response"
	"flamespot-api/src/service"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	token, user, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"token":    token,
		"username": user.Username,
		"nickName": user.NickName,
	})
}
