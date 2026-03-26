package main

import (
	"github/folkyyyy/preorder-api/config"
	"github/folkyyyy/preorder-api/internal/handlers"
	"github/folkyyyy/preorder-api/internal/repositories"
	"github/folkyyyy/preorder-api/internal/routes"
	"github/folkyyyy/preorder-api/internal/services"
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
	// log.Println("Running Auto Migration...")
	// err = config.DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Menu{},
	// 	&models.PreorderRound{},
	// 	&models.PreorderMenu{},
	// 	&models.Order{},
	// 	&models.OrderItem{},
	// )
	// if err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	// 3. สร้างแอป Fiber
	app := fiber.New()

	// ---- Dependency Injection ----
	menuRepo := repositories.NewMenuRepository(config.DB)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	// ---- Setup Routes ----
	api := app.Group("/api")

	routes.SetupMenuRoutes(api, menuHandler)

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
