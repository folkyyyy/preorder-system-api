package handlers

import (
	"errors"
	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		switch {
		case errors.Is(err, apperrors.ErrPriceNegative):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		default:
			// ถ้าหลุดจากข้างบนมา แปลว่าเป็น Error ที่เราไม่ได้คาดคิด (เช่น DB ล่ม)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
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

func (h *MenuHandler) GetMenuByID(c *fiber.Ctx) error {
	id := c.Params("id")

	menuid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	menu, err := h.service.GetMenuByID(uint(menuid))
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลเมนูที่ต้องการ", // หาเมนูนี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": menu})
}

func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	// แปลงข้อมูล JSON ที่ส่งมาให้อยู่ในรูปแบบ struct
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ส่งข้อมูลไม่ถูกต้อง"})
	}
	id := c.Params("id")
	Menuid, _ := strconv.ParseUint(id, 10, 32)
	menu.ID = uint(Menuid)
	// เรียกใช้ Service
	if err := h.service.UpdateMenu(&menu); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPriceNegative):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลเมนูที่ต้องการ", // หาเมนูนี้ไม่เจอ! (ส่ง 404 กลับไป)
			})

		default:
			// ถ้าหลุดจากข้างบนมา แปลว่าเป็น Error ที่เราไม่ได้คาดคิด (เช่น DB ล่ม)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "แก้ไขเมนูสำเร็จ",
		"data":    menu,
	})
}

func (h *MenuHandler) DeleteMenu(c *fiber.Ctx) error {
	id := c.Params("id")
	menuid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	if err := h.service.DeleteMenu(uint(menuid)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลเมนูที่ต้องการ", // หาเมนูนี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
