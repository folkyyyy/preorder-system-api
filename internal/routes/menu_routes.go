package routes

import (
	"github/folkyyyy/preorder-api/internal/handlers"
	"github/folkyyyy/preorder-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupMenuRoutes รับ Router group และ Handler เข้ามาเพื่อจัดการเส้นทาง
func SetupMenuRoutes(router fiber.Router, handler *handlers.MenuHandler) {
	menuGroup := router.Group("/menu", middlewares.Protected())

	menuGroup.Post("/", middlewares.AdminOnly() , handler.CreateMenu)      // POST /api/menu
	menuGroup.Get("/", handler.GetAllMenus)      // GET /api/menu
	menuGroup.Get("/:id", handler.GetMenuByID)   // GET /api/menu/:id
	menuGroup.Put("/:id", middlewares.AdminOnly() , handler.UpdateMenu)    // PUT /api/menu/:id
	menuGroup.Delete("/:id", middlewares.AdminOnly() , handler.DeleteMenu) // DELETE /api/menu/:id
}
