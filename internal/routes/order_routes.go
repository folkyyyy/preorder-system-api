package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"
	"github/folkyyyy/preorder-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupOrderRoutes รับ Router group และ Handler เข้ามาเพื่อจัดการเส้นทาง
func SetupOrderRoutes(router fiber.Router, handler *handlers.OrderHandler) {
	orderGroup := router.Group("/orders")
	orderGroup.Post("/", middlewares.Protected(), handler.CreateOrder)
	orderGroup.Get("/round/:roundId", middlewares.Protected(), handler.GetOrdersByRound)
	orderGroup.Patch("/:id/status", middlewares.Protected(),middlewares.AdminOnly(), handler.UpdateOrderStatus)
	orderGroup.Get("/kitchen-summary/round/:roundId", middlewares.Protected(),middlewares.AdminOnly(), handler.GetKitchenSummary)
	orderGroup.Get("/:id", middlewares.Protected(), handler.GetOrderByID)
	orderGroup.Put("/:id", middlewares.Protected(), handler.UpdateOrderDetails)
}

