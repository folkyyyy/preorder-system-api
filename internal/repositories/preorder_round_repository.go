package repositories

import (
	"github/folkyyyy/preorder-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type PreorderRoundRepository interface {
	CreateRoundWithMenus(round *models.PreorderRound, menus []models.PreorderMenu) error
	CheckRoundExistsByDate(date time.Time) (bool, error)
	GetRoundByID(id uint) (*models.PreorderRound, error)
	UpdateRoundWithMenus(round *models.PreorderRound, newMenus []models.PreorderMenu) error
	DeleteRound(id uint) error
	GetRoundsByDateRange(startDate, endDate time.Time) ([]models.PreorderRound, error)
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

// CheckRoundExistsByDate ตรวจสอบว่ามีรอบพรีออเดอร์เปิดรับในวันที่นี้แล้วหรือยัง
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

// Get preorder round by ID
func (r *preorderRoundRepository) GetRoundByID(id uint) (*models.PreorderRound, error) {
	var round models.PreorderRound
	err := r.db.Preload("PreorderMenus").Preload("PreorderMenus.Menu").First(&round, id).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// Update preorder round with new menus
func (r *preorderRoundRepository) UpdateRoundWithMenus(round *models.PreorderRound, newMenus []models.PreorderMenu) error {
	tx := r.db.Begin()

	// 1. อัปเดตข้อมูลรายละเอียดรอบ (Title, DeliveryDate, Status)
	result := tx.Model(&models.PreorderRound{}).Where("id = ?", round.ID).Updates(round);
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}


	// ใช้ Unscoped() เพื่อลบขาดจากตารางไปเลย (Hard Delete)
	if err := tx.Unscoped().Where("preorder_round_id = ?", round.ID).Delete(&models.PreorderMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. เอา ID ของรอบ มาแปะให้เมนูชุดใหม่ แล้ว Insert เข้าไปใหม่ทั้งหมด
	if len(newMenus) > 0 {
		for i := range newMenus {
			newMenus[i].PreorderRoundID = round.ID
		}

		if err := tx.Create(&newMenus).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if tx.Commit().RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return tx.Commit().Error
	 
}

// Delete preorder round
func (r *preorderRoundRepository) DeleteRound(id uint) error {
	result := r.db.Delete(&models.PreorderRound{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *preorderRoundRepository) GetRoundsByDateRange(startDate, endDate time.Time) ([]models.PreorderRound, error) {
	var rounds []models.PreorderRound
	// ตัดเวลาเป็น 00:00:00 เพื่อให้การค้นหาเทียบแค่วันที่ (ไม่เอาเวลา)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
	// ใช้ <= สำหรับ endDate เพื่อให้รวมข้อมูลของวันสุดท้ายเข้าไปด้วย
	err := r.db.
		Preload("PreorderMenus").
		Preload("PreorderMenus.Menu"). // ดึงข้อมูลอาหารมาให้ครบ
		Where("delivery_date >= ? AND delivery_date <= ?", startDate, endDate).
		Find(&rounds).Error
	
	return rounds, err
}
