package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"
	"github/folkyyyy/preorder-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupPreorderRoundRoutes รับ Router group และ Handler เข้ามาเพื่อจัดการเส้นทาง
func SetupPreorderRoundRoutes(router fiber.Router, handler *handlers.PreorderRoundHandler) {
	roundGroup := router.Group("/preorder", middlewares.Protected())
	roundGroup.Post("/", middlewares.AdminOnly(), handler.CreateRound)
}
