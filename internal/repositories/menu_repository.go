package repositories

import (
	"github/folkyyyy/preorder-api/internal/models"
	"gorm.io/gorm"
)

// 1. สร้าง Interface กำหนดว่า Repository นี้ทำอะไรได้บ้าง
type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
	GetAllMenus() ([]models.Menu, error)
	GetMenuByID(id uint) (*models.Menu, error)
	UpdateMenu(menu *models.Menu) error
	DeleteMenu(id uint) error
}

type menuRepository struct {
	db *gorm.DB
}

// 2. ฟังก์ชันสำหรับสร้าง Repository Instance
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

// 3. เขียนฟังก์ชันสร้างเมนูใหม่ (Insert ลง DB)
func (r *menuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// 4. เขียนฟังก์ชันดึงเมนูทั้งหมด (Select จาก DB)
func (r *menuRepository) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Find(&menus).Error
	return menus, err
}

// Get Menu by ID (Select จาก DB โดยใช้ ID)
func (r *menuRepository) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update Menu
func (r *menuRepository) UpdateMenu(menu *models.Menu) error {
	result := r.db.Model(&models.Menu{}).Where("id = ?", menu.ID).Updates(menu)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete Menu
func (r *menuRepository) DeleteMenu(id uint) error {
	result := r.db.Delete(&models.Menu{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
