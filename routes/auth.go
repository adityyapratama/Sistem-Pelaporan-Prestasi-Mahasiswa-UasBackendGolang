package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"github.com/gofiber/fiber/v2"
)


func RegisterAuthRoutes(router fiber.Router, authService *service.AuthService) {
	
	
	auth := router.Group("/auth")

	auth.Post("Register",authService.Register)
	auth.Post("Login",authService.Login)
}