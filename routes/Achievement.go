package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)





func AchievementRoutes(router fiber.Router, Achievservice *service.AchievementService) {
	achievements := router.Group("/achievements", middleware.AuthProtected())

	achievements.Post("/", middleware.VerifyRole("Mahasiswa"), Achievservice.Create)
	
}