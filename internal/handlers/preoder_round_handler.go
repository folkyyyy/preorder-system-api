package handlers

import (
	"time"

	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PreorderRoundHandler struct {
	service services.PreorderRoundService
}

func NewPreorderRoundHandler(service services.PreorderRoundService) *PreorderRoundHandler {
	return &PreorderRoundHandler{service}
}

// --- สร้าง DTO สำหรับรับข้อมูล JSON จากหน้าเว็บ ---
type PreorderMenuInput struct {
	MenuID uint `json:"menuId"`
	Quota  int  `json:"quota"`
}

type CreateRoundInput struct {
	Title        string              `json:"title"`
	DeliveryDate string              `json:"deliveryDate"` // รับเป็น string เช่น "2026-04-15"
	Menus        []PreorderMenuInput `json:"menus"`         // รับเมนูมาเป็น Array
}


func (h *PreorderRoundHandler) CreateRound(c *fiber.Ctx) error {
	var input CreateRoundInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
	}

	// 1. แปลงวันที่จาก String เป็น time.Time ของ Go
	deliveryDate, err := time.Parse("2006-01-02", input.DeliveryDate) // "2006-01-02" คือ format มาตรฐานของ Go
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบวันที่ไม่ถูกต้อง ต้องเป็น YYYY-MM-DD"})
	}

	// 2. ประกอบร่าง Model PreorderRound
	round := models.PreorderRound{
		Title:        input.Title,
		DeliveryDate: deliveryDate,
		Status:       "open",
	}

	// 3. ประกอบร่าง Model PreorderMenu (วนลูปตามที่ส่งมา)
	var preorderMenus []models.PreorderMenu
	for _, m := range input.Menus {
		preorderMenus = append(preorderMenus, models.PreorderMenu{
			MenuID: m.MenuID,
			Quota:  m.Quota,
		})
	}

	// 4. ส่งให้ Service จัดการ
	if err := h.service.CreatePreorderRound(&round, preorderMenus); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เปิดรอบพรีออเดอร์สำเร็จ",
	})
}
