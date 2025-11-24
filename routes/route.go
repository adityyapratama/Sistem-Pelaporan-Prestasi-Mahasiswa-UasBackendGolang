package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "go.mongodb.org/mongo-driver/internal/logger"
)



func SetupRoutes(app *fiber.App,
authService *service.AuthService){
	// app.Use(logger.new())
	app.Use(cors.New())

	api := app.Group("/api/v1")
	RegisterAuthRoutes(api, authService)


	
}