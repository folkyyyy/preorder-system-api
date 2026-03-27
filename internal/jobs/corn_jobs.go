package jobs

import (
	"log"
	"time"

	"github/folkyyyy/preorder-api/internal/models" // เปลี่ยน path ให้ตรงกับโปรเจกต์คุณ

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// StartAutoCloseJob ฟังก์ชันสำหรับเปิดการทำงานของ Cron
func StartAutoCloseJob(db *gorm.DB) {
	// 1. ตั้งค่า Timezone ให้เป็นเวลาไทย (สำคัญมาก ไม่งั้นเที่ยงคืนจะกลายเป็นเวลาประเทศอื่น)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("โหลด Timezone ไม่ได้: %v", err)
	}

	// 2. สร้างตัวจัดการ Cron Job
	c := cron.New(cron.WithLocation(loc))

	// 3. ตั้งเวลาการทำงาน: "0 0 * * *" แปลว่า "นาทีที่ 0, ชั่วโมงที่ 0 (เที่ยงคืน), ของทุกวัน"
	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println(" [Cron Job] กำลังรันระบบปิดรอบพรีออเดอร์อัตโนมัติ...")

		// หาวันที่ของวันนี้ (เอาเฉพาะ วัน-เดือน-ปี ไม่เอาเวลา)
		today := time.Now().In(loc)
		todayMidnight := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)

		// สั่งให้ GORM อัปเดตข้อมูล:
		// "ปรับ status เป็น closed ถ้า status เดิมเป็น open และ delivery_date น้อยกว่าวันนี้"
		result := db.Model(&models.PreorderRound{}).
			Where("status = ? AND delivery_date < ?", "open", todayMidnight).
			Update("status", "closed")

		if result.Error != nil {
			log.Printf(" [Cron Job] เกิดข้อผิดพลาดในการปิดรอบ: %v\n", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf(" [Cron Job] ปิดรอบอัตโนมัติสำเร็จ จำนวน %d รอบ\n", result.RowsAffected)
		} else {
			log.Println(" [Cron Job] ไม่มีรอบพรีออเดอร์ที่หมดอายุในวันนี้")
		}
	})

	if err != nil {
		log.Fatalf("สร้าง Cron Job ไม่สำเร็จ: %v", err)
	}

	// 4. สั่งให้ Cron เริ่มทำงานแบบ Background
	c.Start()
	log.Println(" ระบบตั้งเวลาปิดรอบพรีออเดอร์ (Cron Job) เริ่มทำงานแล้ว!")
}
