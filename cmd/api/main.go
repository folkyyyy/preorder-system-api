package main

import (
	"github/folkyyyy/preorder-api/config"
	"github/folkyyyy/preorder-api/internal/models"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. โหลดไฟล์ .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// 2. เชื่อมต่อ Database
	config.ConnectDB()

	// 2.5 ทำ AutoMigrate เพื่อสร้างตาราง
	log.Println("Running Auto Migration...")
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.Menu{},
		&models.PreorderRound{},   
		&models.PreorderMenu{},    
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 3. สร้างแอป Fiber
	app := fiber.New()

	// 4. สร้าง Route ง่ายๆ เพื่อทดสอบ
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Pre-order API!",
			"status":  "success",
		})
	})

	// 5. รันเซิร์ฟเวอร์
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(app.Listen(port))
}
