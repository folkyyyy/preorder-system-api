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
	tokenUserId := c.Locals("user_id").(uint)

	var finalUserId *uint

	switch tokenRole {
	case "user":
		finalUserId = &tokenUserId
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