package handlers

import (
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service}
}

// --- DTO สำหรับรับข้อมูล JSON ---
type OrderItemInput struct {
	PreorderMenuID uint   `json:"preorderMenuId"`
	Quantity       int    `json:"quantity"`
	Note           string `json:"note"`
}

type CreateOrderInput struct {
	PreorderRoundID  uint             `json:"preorderRoundId"`
	Name             string           `json:"name"`
	DeliveryLocation string           `json:"deliveryLocation"`
	Items            []OrderItemInput `json:"items"`
}

// ------------------------------

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var input CreateOrderInput

	// 1. แปลง JSON เข้า Struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
	}

	tokenRole := c.Locals("role").(string)
	var tokenUserID uint
	if idVal, ok := c.Locals("userID").(float64); ok {
		tokenUserID = uint(idVal)
	}

	var finalUserId *uint

	switch tokenRole {
	case "user":
		finalUserId = &tokenUserID
	case "admin":
		finalUserId = nil
	}

	// 2. ประกอบร่าง Model ตาราง Order
	order := models.Order{
		UserID:           finalUserId,
		PreorderRoundID:  input.PreorderRoundID,
		Name:             input.Name,
		DeliveryLocation: input.DeliveryLocation,
		Status:           "pending", // บิลใหม่ สถานะรอโอนเงินเสมอ
	}

	// 3. ประกอบร่าง Model ตาราง OrderItem
	var orderItems []models.OrderItem
	for _, item := range input.Items {
		// ป้องกันแอดมินพิมพ์จำนวนติดลบ หรือ 0
		if item.Quantity <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "จำนวนอาหารต้องมากกว่า 0"})
		}

		orderItems = append(orderItems, models.OrderItem{
			PreorderMenuID: item.PreorderMenuID,
			Quantity:       item.Quantity,
			Note:           item.Note,
		})
	}

	// 4. ส่งให้ Service ทำงาน
	if err := h.service.CreateOrder(&order, orderItems); err != nil {
		// ถ้ามี Error (เช่น โควต้าไม่พอ, เมนูไม่มีจริง) ส่ง 400 Bad Request กลับไป
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 5. ตอบกลับเมื่อสำเร็จ
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "บันทึกรายการสั่งซื้อสำเร็จ",
	})
}

func (h *OrderHandler) GetOrdersByRound(c *fiber.Ctx) error {
	idParam := c.Params("roundId")
	roundID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ ID รอบพรีออเดอร์ไม่ถูกต้อง",
		})
	}

	// 2. ส่งให้ Service ดึงข้อมูล
	orders, err := h.service.GetOrdersByRoundID(uint(roundID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลออเดอร์ได้",
		})
	}

	// 3. ถ้าไม่มีออเดอร์เลย ส่งเป็น Array ว่างกลับไป (Frontend จะได้ไม่พัง)
	if len(orders) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"data":    []models.Order{},
			"message": "ยังไม่มีรายการสั่งซื้อในรอบนี้",
		})
	}

	// 4. ส่งข้อมูลกลับไป
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": orders,
	})
}

type UpdateOrderStatusInput struct {
	Status string `json:"status"`
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ ID ออเดอร์ไม่ถูกต้อง",
		})
	}

	// 2. รับข้อมูล Status จาก Body JSON
	var input UpdateOrderStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	// 3. ทำความสะอาดข้อมูล (ตัดช่องว่างหน้าหลัง และแปลงเป็นตัวพิมพ์เล็ก)
	cleanStatus := strings.ToLower(strings.TrimSpace(input.Status))

	// 4. ส่งให้ Service ทำงาน
	if err := h.service.UpdateOrderStatus(uint(orderID), cleanStatus); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 5. ตอบกลับเมื่อสำเร็จ
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัปเดตสถานะออเดอร์สำเร็จ",
	})
}

func (h *OrderHandler) GetKitchenSummary(c *fiber.Ctx) error {
	idParam := c.Params("roundId")
	roundID, err := strconv.ParseUint(strings.TrimSpace(idParam), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ ID รอบพรีออเดอร์ไม่ถูกต้อง",
		})
	}

	summary, err := h.service.GetKitchenSummary(uint(roundID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลสรุปยอดพรีออเดอร์รอบนี้ได้",
		})
	}

	// ถ้ายังไม่มีออเดอร์เลย
	if len(summary) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":    []models.KitchenSummary{},
			"message": "ยังไม่มียอดสั่งอาหารสำหรับรอบนี้",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": summary,
	})
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ ID ออเดอร์ไม่ถูกต้อง",
		})
	}
	order, err := h.service.GetOrderById(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลออเดอร์ได้",
		})
	}
	if order == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบออเดอร์ที่ต้องการ",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": order,
	})
}

type UpdateOrderInput struct {
	Name             string           `json:"name"`
	DeliveryLocation string           `json:"deliveryLocation"`
	Items            []OrderItemInput `json:"items"` // ใช้ Struct เดิมที่คุณเคยสร้างไว้ตอน Create ได้เลย
}

func (h *OrderHandler) UpdateOrderDetails(c *fiber.Ctx) error {
	// 1. รับ ID ของบิลจาก URL
	idParam := strings.TrimSpace(c.Params("id"))
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ ID ออเดอร์ไม่ถูกต้อง",
		})
	}

	// 2. รับข้อมูล JSON จาก Request Body
	var input UpdateOrderInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	tokenRole := c.Locals("role").(string)
	var tokenUserID uint
	if idVal, ok := c.Locals("userID").(float64); ok {
		tokenUserID = uint(idVal)
	}

	var finalUserId *uint

	switch tokenRole {
	case "user":
		finalUserId = &tokenUserID
	case "admin":
		finalUserId = nil
	}

	// 3. ประกอบร่างข้อมูลพื้นฐานของบิล (ข้อมูลที่อาจจะเปลี่ยน)
	updateData := &models.Order{
		UserID:           finalUserId,
		Name:             input.Name,
		DeliveryLocation: input.DeliveryLocation,
	}

	// 4. แปลงรายการอาหาร (Items) จาก Input ให้กลายเป็น Model ของ GORM
	var newItems []models.OrderItem
	for _, item := range input.Items {
		newItems = append(newItems, models.OrderItem{
			PreorderMenuID: item.PreorderMenuID,
			Quantity:       item.Quantity,
			// Note: item.Note, // ถ้าใน OrderItemInput ของคุณมีช่อง Note (หมายเหตุ) ให้เปิดคอมเมนต์บรรทัดนี้ด้วยครับ
		})
	}

	// 5. ส่งให้ Service ทำงาน
	if err := h.service.UpdateOrder(uint(orderID), updateData, newItems); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 6. ส่งผลลัพธ์กลับเมื่อสำเร็จ
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "แก้ไขข้อมูลออเดอร์และอัปเดตโควต้าสำเร็จ",
	})
}
