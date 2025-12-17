package config

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/routes"

	"github.com/gofiber/fiber/v2"
)

func NewApp(
	authService *service.AuthService,
	permService *service.PermissionService,
	studentService *service.StudentService,
	lectureService *service.LectureService,
	achievmentService *service.AchievementService,
	) *fiber.App{
	app := fiber.New()

	routes.SetupRoutes(app,authService,permService,studentService,lectureService,achievmentService)
	
	

	return app
}