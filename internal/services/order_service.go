package services

import (
	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
)

type OrderService interface {
	CreateOrder(order *models.Order, items []models.OrderItem) error
	GetOrdersByRoundID(roundID uint) ([]models.Order, error)
	UpdateOrderStatus(orderID uint, newStatus string) error
	GetKitchenSummary(roundID uint) ([]models.KitchenSummary, error)
	GetOrderById(orderID uint) (*models.Order, error)
	UpdateOrder(orderID uint, updateData *models.Order, newItems []models.OrderItem) error
}

type orderService struct {
	repo      repositories.OrderRepository
	roundRepo repositories.PreorderRoundRepository
}

func NewOrderService(repo repositories.OrderRepository, roundRepo repositories.PreorderRoundRepository) OrderService {
	return &orderService{
		repo:      repo,
		roundRepo: roundRepo,
	}
}

func (s *orderService) CreateOrder(order *models.Order, items []models.OrderItem) error {
	// 1. ดักจับ Logic พื้นฐาน: ถ้าแอดมินส่งบิลมาแต่ไม่มีเมนูอะไรเลย ให้เตะออก
	if len(items) == 0 {
		return apperrors.ErrEmptyOrderItems
	}

	// ถ้าเป็น "closed" แล้วก็ return error กลับไปว่า "รอบนี้ปิดรับไปแล้ว"
	if order.PreorderRound.Status == "closed" {
		return apperrors.ErrRoundClosed
	}

	round, err := s.roundRepo.GetRoundByID(order.PreorderRoundID)
	if err != nil {
		return err // ถ้ารอบไม่มีจริง เตะออก
	}

	// เช็คว่าเมนูที่สั่งมาทั้งหมดมีอยู่ในรอบพรีออเดอร์นี้ไหม? ถ้าไม่เจอเมนูไหนเลย ก็เตะออก
	for _, item := range items {
		found := false
		for _, menu := range round.PreorderMenus {
			if item.PreorderMenuID == menu.ID {
				found = true
				break
			}
		}
		if !found {
			return apperrors.ErrMenuNotFound
		}
	}

	// 2. ถ้าผ่านหมด ส่งให้ Repository จัดการต่อ (หักโควต้า + บันทึกลง DB)
	return s.repo.CreateOrder(order, items)
}

func (s *orderService) GetOrdersByRoundID(roundID uint) ([]models.Order, error) {
	return s.repo.GetOrdersByRoundID(roundID)
}

func (s *orderService) UpdateOrderStatus(orderID uint, newStatus string) error {
	// 1. ดักตรวจสอบความถูกต้องของ Status ก่อน
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"cancelled": true,
	}

	if !validStatuses[newStatus] {
		return apperrors.ErrInvalidOrderStatus
	}

	// 2. ถ้าสถานะถูกต้อง ส่งต่อให้ Repository จัดการ Database และเรื่องโควต้า
	return s.repo.UpdateOrderStatus(orderID, newStatus)
}


func (s *orderService) GetKitchenSummary(roundID uint) ([]models.KitchenSummary, error) {
	return s.repo.GetKitchenSummary(roundID)
}

func (s *orderService) GetOrderById(orderID uint) (*models.Order, error) {
	return s.repo.GetOrderById(orderID)
}

func (s *orderService) UpdateOrder(orderID uint, updateData *models.Order, newItems []models.OrderItem) error {
	// ดักทาง: บิลนึงควรต้องมีอาหารอย่างน้อย 1 อย่าง ไม่งั้นก็ควรใช้วิธียกเลิกบิลแทน
	if len(newItems) == 0 {
		return apperrors.ErrEmptyOrderItems
	}


	// ส่งต่อให้ Repository ทำการ Wipe & Recreate สุดโหดของเรา
	return s.repo.UpdateOrder(orderID, updateData, newItems)
}