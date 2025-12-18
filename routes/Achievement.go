package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(router fiber.Router, Achievservice *service.AchievementService) {
	
	achievements := router.Group("/achievements", middleware.AuthProtected())

	achievements.Post("/upload", middleware.VerifyRole("Mahasiswa"), Achievservice.UploadAttachment)
	
	
	achievements.Get("/raw-mongo", Achievservice.GetAllMongoData) 
	achievements.Post("/", middleware.VerifyRole("Mahasiswa"), Achievservice.Create)
	achievements.Get("/", Achievservice.GetAll)
	achievements.Get("/:id", Achievservice.GetByID)
	achievements.Put("/:id", middleware.VerifyRole("Mahasiswa"), Achievservice.Update)
	achievements.Patch("/:id/verify", middleware.VerifyRole("Dosen"), Achievservice.Verify)
}