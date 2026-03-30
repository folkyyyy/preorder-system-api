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
	// orderGroup.Get("/", middlewares.Protected(), handler.ListOrders)
	// orderGroup.Get("/:id", middlewares.Protected(), handler.GetOrderByID)
}

