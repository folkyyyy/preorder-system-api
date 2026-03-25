package services

import (
	"errors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
)

type MenuService interface {
	CreateMenu(menu *models.Menu) error
	GetAllMenus() ([]models.Menu, error)
}

type menuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) MenuService {
	return &menuService{repo}
}

func (s *menuService) CreateMenu(menu *models.Menu) error {
	if menu.Price < 0 {
		return errors.New("ราคาอาหารต้องไม่ติดลบ")
	}

	return s.repo.CreateMenu(menu)
}

func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	return s.repo.GetAllMenus()
}