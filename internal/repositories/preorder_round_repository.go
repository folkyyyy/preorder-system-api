package repositories

import (
	"github/folkyyyy/preorder-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type PreorderRoundRepository interface {
	CreateRoundWithMenus(round *models.PreorderRound, menus []models.PreorderMenu) error
	CheckRoundExistsByDate(date time.Time) (bool, error)
}

type preorderRoundRepository struct {
	db *gorm.DB
}

func NewPreorderRoundRepository(db *gorm.DB) PreorderRoundRepository {
	return &preorderRoundRepository{db}
}

// CreateRoundWithMenus บันทึกรอบและเมนูพร้อมกัน
func (r *preorderRoundRepository) CreateRoundWithMenus(round *models.PreorderRound, menus []models.PreorderMenu) error {
	// 1. เริ่ม Transaction (ถ้ามีอะไรพังหลังจากนี้ เราจะสั่ง Rollback ย้อนกลับได้)
	tx := r.db.Begin()

	// 2. บันทึกข้อมูลลงตาราง preorder_rounds ก่อน เพื่อให้ได้ ID ของรอบนี้มา
	if err := tx.Create(round).Error; err != nil {
		tx.Rollback() // พังปุ๊บ ยกเลิกการบันทึกทันที
		return err
	}

	// 3. เอา ID ของรอบที่เพิ่งสร้างเสร็จ ไปใส่ให้กับเมนูแต่ละตัว แล้วบันทึกลง preorder_menus
	for i := range menus {
		menus[i].PreorderRoundID = round.ID
	}

	if err := tx.Create(&menus).Error; err != nil {
		tx.Rollback() // ถ้าบันทึกเมนูพัง ก็ยกเลิกการสร้างรอบในข้อ 2 ด้วย (ข้อมูลจะได้ไม่แหว่ง)
		return err
	}

	// 4. ถ้าผ่านฉลุยทุกขั้นตอน สั่งยืนยันการบันทึก (Commit)
	return tx.Commit().Error
}

func (r *preorderRoundRepository) CheckRoundExistsByDate(date time.Time) (bool, error) {
	var count int64
	// ค้นหาในตาราง preorder_rounds โดยเทียบเฉพาะส่วนของวันที่ (ไม่เอาเวลา)
	// ใช้ DATE() เพื่อให้ชัวร์ว่าเทียบแค่วัน เดือน ปี
	err := r.db.Model(&models.PreorderRound{}).Where("DATE(delivery_date) = DATE(?)", date).Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	// ถ้า count มากกว่า 0 แปลว่ามีวันที่นี้อยู่แล้ว (return true)
	return count > 0, nil
}