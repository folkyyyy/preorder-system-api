package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, handler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")

	authGroup.Post("/register", handler.Register) // POST /api/auth/register
	authGroup.Post("/login", handler.Login) // POST /api/auth/login
}