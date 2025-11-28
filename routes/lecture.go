package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)





func LectureRoutes(router fiber.Router, lectService *service.LectureService) {
	
	
	lecture := router.Group("/lectures",middleware.AuthProtected())
	lecture.Post("/", lectService.Create)
	lecture.Get("/", lectService.GetAll)
	
}