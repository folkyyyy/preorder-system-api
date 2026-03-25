package handlers

import (
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service}
}

// ฟังก์ชันสำหรับรับ POST Request เพื่อสร้างเมนู
func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	var menu models.Menu

	// แปลงข้อมูล JSON ที่ส่งมาให้อยู่ในรูปแบบ struct
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ส่งข้อมูลไม่ถูกต้อง"})
	}

	// เรียกใช้ Service
	if err := h.service.CreateMenu(&menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "สร้างเมนูสำเร็จ",
		"data":    menu,
	})
}

// ฟังก์ชันสำหรับรับ GET Request เพื่อดึงเมนูทั้งหมด
func (h *MenuHandler) GetAllMenus(c *fiber.Ctx) error {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": menus})
}
