package handlers

import (
	"errors"
	"strconv"
	"time"

	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	Menus        []PreorderMenuInput `json:"menus"`        // รับเมนูมาเป็น Array
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
		switch {
		case errors.Is(err, apperrors.ErrPastDate):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		case errors.Is(err, apperrors.ErrRoundConflict):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})

		default:
			// ถ้าหลุดจากข้างบนมา แปลว่าเป็น Error ที่เราไม่ได้คาดคิด (เช่น DB ล่ม)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เปิดรอบพรีออเดอร์สำเร็จ",
	})
}

// get round by ID
func (h *PreorderRoundHandler) GetRoundByID(c *fiber.Ctx) error {
	id := c.Params("id")

	roundid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	round, err := h.service.GetRoundByID(uint(roundid))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลรอบพรีออเดอร์ที่ต้องการ", // หารอบพรีออเดอร์นี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": round})
}

// update round with new menus
func (h *PreorderRoundHandler) UpdateRound(c *fiber.Ctx) error {
	var input CreateRoundInput
	// แปลงข้อมูล JSON ที่ส่งมาให้อยู่ในรูปแบบ struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ส่งข้อมูลไม่ถูกต้อง"})
	}
	id := c.Params("id")
	roundid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}
	// แปลงวันที่จาก String เป็น time.Time ของ Go
	deliveryDate, err := time.Parse("2006-01-02", input.DeliveryDate) // "2006-01-02" คือ format มาตรฐานของ Go
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบวันที่ไม่ถูกต้อง ต้องเป็น YYYY-MM-DD"})
	}
	// ประกอบร่าง Model PreorderRound
	round := models.PreorderRound{
		Title:        input.Title,
		DeliveryDate: deliveryDate,
	}
	round.ID = uint(roundid)
	// ประกอบร่าง Model PreorderMenu (วนลูปตามที่ส่งมา)
	var preorderMenus []models.PreorderMenu
	for _, m := range input.Menus {
		preorderMenus = append(preorderMenus, models.PreorderMenu{
			MenuID: m.MenuID,
			Quota:  m.Quota,
		})
	}
	// ส่งให้ Service จัดการ
	if err := h.service.UpdateRoundWithMenus(&round, preorderMenus); err != nil {

		// ดักจับและแปลงเป็น Status Code
		switch {

		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลรอบพรีออเดอร์ที่ต้องการ", // หารอบพรีออเดอร์นี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		case errors.Is(err, apperrors.ErrPastDate):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		case errors.Is(err, apperrors.ErrRoundConflict):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})

		default:
			// ถ้าหลุดจากข้างบนมา แปลว่าเป็น Error ที่เราไม่ได้คาดคิด (เช่น DB ล่ม)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัปเดตรอบพรีออเดอร์สำเร็จ",
	})
}

// delete round by ID
func (h *PreorderRoundHandler) DeleteRound(c *fiber.Ctx) error {
	id := c.Params("id")
	roundid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}
	if err := h.service.DeleteRound(uint(roundid)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลรอบพรีออเดอร์ที่ต้องการ", // หารอบพรีออเดอร์นี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// get rounds by date range
func (h *PreorderRoundHandler) GetRoundsByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบวันที่เริ่มต้นไม่ถูกต้อง ต้องเป็น YYYY-MM-DD"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบวันที่สิ้นสุดไม่ถูกต้อง ต้องเป็น YYYY-MM-DD"})
	}
	rounds, err := h.service.GetRoundsByDateRange(startDate, endDate)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่มีข้อมูลรอบพรีออเดอร์ที่ต้องการ", // หารอบพรีออเดอร์นี้ไม่เจอ! (ส่ง 404 กลับไป)
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": rounds})
}

type ChangeStatusInput struct {
	Status string `json:"status"`
}

// change status
func (h *PreorderRoundHandler) ChangeRoundStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	roundid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}
	var input ChangeStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
	}
	if err := h.service.ChangeRoundStatus(uint(roundid), input.Status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่มีข้อมูลรอบพรีออเดอร์ที่ต้องการ", // หารอบพรีออเดอร์นี้ไม่เจอ! (ส่ง 404 กลับไป)
			})
		}
		if errors.Is(err,apperrors.ErrInvalidStatus) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	statusText := ""

	if input.Status == "closed" {
		statusText = "ปิดรับพรีออเดอร์สำเร็จ"
	}
	if input.Status == "open" {
		statusText = "เปิดรับพรีออเดอร์สำเร็จ"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": statusText,
	})
}
