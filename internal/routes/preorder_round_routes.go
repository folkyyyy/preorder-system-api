package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"
	"github/folkyyyy/preorder-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupPreorderRoundRoutes รับ Router group และ Handler เข้ามาเพื่อจัดการเส้นทาง
func SetupPreorderRoundRoutes(router fiber.Router, handler *handlers.PreorderRoundHandler) {
	roundGroup := router.Group("/preorder")
	roundGroup.Post("/", middlewares.Protected(), middlewares.AdminOnly(), handler.CreateRound)
	roundGroup.Get("/date-range", handler.GetRoundsByDateRange)

	roundGroup.Get("/:id", middlewares.Protected(), handler.GetRoundByID)
	roundGroup.Put("/:id", middlewares.Protected(), middlewares.AdminOnly(), handler.UpdateRound)
	roundGroup.Delete("/:id", middlewares.Protected(), middlewares.AdminOnly(), handler.DeleteRound)

	roundGroup.Patch("/:id/status", middlewares.Protected(), middlewares.AdminOnly(), handler.ChangeRoundStatus)
}
