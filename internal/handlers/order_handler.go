package handlers

import (
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"

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

	// 2. ประกอบร่าง Model ตาราง Order
	order := models.Order{
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