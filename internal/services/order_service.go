package services

import (
	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
)

type OrderService interface {
	CreateOrder(order *models.Order, items []models.OrderItem) error
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
