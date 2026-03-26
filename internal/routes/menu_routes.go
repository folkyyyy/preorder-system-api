package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupMenuRoutes รับ Router group และ Handler เข้ามาเพื่อจัดการเส้นทาง
func SetupMenuRoutes(router fiber.Router, handler *handlers.MenuHandler) {
	menuGroup := router.Group("/menu")

	menuGroup.Post("/", handler.CreateMenu) // POST /api/menu
	menuGroup.Get("/", handler.GetAllMenus) // GET /api/menu
	menuGroup.Get("/:id", handler.GetMenuByID) // GET /api/menu/:id
	menuGroup.Put("/:id", handler.UpdateMenu) // PUT /api/menu/:id
	menuGroup.Delete("/:id", handler.DeleteMenu) // DELETE /api/menu/:id
}
