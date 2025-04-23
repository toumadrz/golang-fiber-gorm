package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	// โหลด .env มาใช้
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// ผูก .env
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//เชื่อมต่อ
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	fmt.Println("Successfully connected!")

	app := fiber.New()
	app.Get("/users", getUsers)
	app.Get("/user/:id", getUser)
	app.Post("/user", createUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)
	app.Listen(":8000")
}

// Fiber get user ทั้งหมด
func getUsers(c *fiber.Ctx) error {
	user, err := GetUsers()
	if err != nil {
		log.Println("GetUsers error:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)
}

// Fiber get user จาก Id
func getUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user, err := GetUserByID(Id)
	if err != nil {
		log.Println("GetUsers error:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)
}

// Fiber create user
func createUser(c *fiber.Ctx) error {
	u := new(User)
	if err := c.BodyParser(u); err != nil {
		log.Println("BodyParser error:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := CreateUser(*u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(user)
}

// Fiber update user
func updateUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	u := new(User)
	if err := c.BodyParser(u); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := UpdateUser(Id, u)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)

}

// Fiber delete user
func deleteUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = DeleteUser(Id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendString("Delete user successfully.")
}
