package services

import (
	"errors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
	"time"
)

type PreorderRoundService interface {
	CreatePreorderRound(round *models.PreorderRound, menus []models.PreorderMenu) error
}

type preorderRoundService struct {
	repo repositories.PreorderRoundRepository
}

func NewPreorderRoundService(repo repositories.PreorderRoundRepository) PreorderRoundService {
	return &preorderRoundService{repo}
}

func (s *preorderRoundService) CreatePreorderRound(round *models.PreorderRound, menus []models.PreorderMenu) error {
	if round.DeliveryDate.Before(time.Now()) {
		return errors.New("วันที่จัดส่งต้องไม่ใช่วันที่ผ่านไปแล้ว")
	}

	exists, err := s.repo.CheckRoundExistsByDate(round.DeliveryDate)
	if err != nil {
		return err // ถ้า Database มีปัญหา ให้โยน error ออกไป
	}
	
	if exists {
		// ถ้ามีแล้ว ให้เบรกการทำงาน และส่ง Error แจ้งเตือนกลับไป
		return errors.New("มีรอบพรีออเดอร์เปิดรับในวันที่นี้แล้ว ไม่สามารถเปิดซ้ำได้")
	}

	
	return s.repo.CreateRoundWithMenus(round, menus)
}