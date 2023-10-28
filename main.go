package main

import (
	"auth-service/controllers"
	"auth-service/middlewares"
	"auth-service/models"
	"auth-service/routes"
	"auth-service/utils"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	server          *fiber.App
	authController  controllers.AuthController
	routeController routes.RouteController
	authMiddleware  middlewares.AuthMiddleware
	env             *utils.Environment
)

func init() {
	envFileName := ".env"
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("Error loading %v file", envFileName)
		panic(err)
	}

	dsn := os.Getenv("DB_DSN")
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Failed to connect to the database")
		panic(err)
	}
	DB.AutoMigrate(&models.User{}, &models.Profile{})

	fmt.Println("Connected to the database successfully")
	authController = controllers.NewAuthCotroller(DB)
	authMiddleware = middlewares.NewAuthMiddleware(DB)
	routeController = routes.NewRouteController(authController, authMiddleware)
	server = fiber.New()
}

func main() {
	api := server.Group("/api")
	api.Use(logger.New(logger.Config{
		Output:        os.Stdout,
		Format:        "[${time}] ${status} - ${latency} ${method} ${path}\n",
		DisableColors: false,
	}))

	api.Get("/healthcheck", func(ctx *fiber.Ctx) error {
		message := "Auth service is healthy"
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
	})

	v1 := api.Group("/v1")
	routeController.AuthRoute(v1)
	log.Fatal(server.Listen(":3008"))

	// svc := NewUserService("http://localhost:3005/user")
	// svc = NewLoggingService(svc)

	// apiServer := NewApiServer(svc)
	// log.Fatal(apiServer.Start(":3006"))

	// user, err := svc.GetUser(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
