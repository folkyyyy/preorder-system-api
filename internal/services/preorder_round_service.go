package services

import (
	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
	"time"
)

type PreorderRoundService interface {
	CreatePreorderRound(round *models.PreorderRound, menus []models.PreorderMenu) error
	GetRoundByID(id uint) (*models.PreorderRound, error)
	UpdateRoundWithMenus(round *models.PreorderRound, newMenus []models.PreorderMenu) error
	DeleteRound(id uint) error
	GetRoundsByDateRange(startDate, endDate time.Time) ([]models.PreorderRound, error)
	ChangeRoundStatus(id uint, newStatus string) error
}

type preorderRoundService struct {
	repo repositories.PreorderRoundRepository
}

func NewPreorderRoundService(repo repositories.PreorderRoundRepository) PreorderRoundService {
	return &preorderRoundService{repo}
}

func (s *preorderRoundService) CreatePreorderRound(round *models.PreorderRound, menus []models.PreorderMenu) error {
	if round.DeliveryDate.Before(time.Now()) {
		return apperrors.ErrPastDate
	}

	if len(menus) == 0 {
		return apperrors.ErrInvalidMenu
	}

	for _, menu := range menus {
		if menu.Quota <= 0 {
			return apperrors.ErrInvalidQuota
		}
	}

	exists, err := s.repo.CheckRoundExistsByDate(round.DeliveryDate)
	if err != nil {
		return err // ถ้า Database มีปัญหา ให้โยน error ออกไป
	}

	if exists {
		// ถ้ามีแล้ว ให้เบรกการทำงาน และส่ง Error แจ้งเตือนกลับไป
		return apperrors.ErrRoundConflict
	}

	return s.repo.CreateRoundWithMenus(round, menus)
}

// Get round by ID
func (s *preorderRoundService) GetRoundByID(id uint) (*models.PreorderRound, error) {
	return s.repo.GetRoundByID(id)
}

// Update round with new menus
func (s *preorderRoundService) UpdateRoundWithMenus(round *models.PreorderRound, newMenus []models.PreorderMenu) error {
	if round.DeliveryDate.Before(time.Now()) {
		return apperrors.ErrPastDate
	}

	exists, err := s.repo.CheckRoundExistsByDate(round.DeliveryDate)
	if err != nil {
		return err // ถ้า Database มีปัญหา ให้โยน error ออกไป
	}

	if exists {
		// ถ้ามีแล้ว ให้เบรกการทำงาน และส่ง Error แจ้งเตือนกลับไป
		return apperrors.ErrRoundConflict
	}

	if len(newMenus) == 0 {
		return apperrors.ErrInvalidMenu
	}

	for _, menu := range newMenus {
		if menu.Quota <= 0 {
			return apperrors.ErrInvalidQuota
		}
	}

	return s.repo.UpdateRoundWithMenus(round, newMenus)
}

// Delete round by ID
func (s *preorderRoundService) DeleteRound(id uint) error {
	return s.repo.DeleteRound(id)
}

// Get rounds by date range
func (s *preorderRoundService) GetRoundsByDateRange(startDate, endDate time.Time) ([]models.PreorderRound, error) {
	return s.repo.GetRoundsByDateRange(startDate, endDate)
}

// change round status
func (s *preorderRoundService) ChangeRoundStatus(id uint, newStatus string) error {
	if(newStatus != "open" && newStatus != "closed") {
		return apperrors.ErrInvalidStatus
	}
	return s.repo.ChangeRoundStatus(id, newStatus)
}
