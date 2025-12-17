package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)


func RegisterAuthRoutes(router fiber.Router, authService *service.AuthService) {
	
	
	auth := router.Group("/auth")

	
	auth.Post("Register",authService.Register)
	auth.Post("Login",authService.Login)
	auth.Put("/:id",authService.UpdateUser)
	auth.Get("/",authService.GetAllUser)

	// update role
	auth.Put("/:id",authService.UpdateRoleUser)

	auth.Get("profiles",middleware.AuthProtected(),authService.GetProfile)
	auth.Post("/refresh", authService.RefreshToken)


}