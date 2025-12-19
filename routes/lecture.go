package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)

func LectureRoutes(router fiber.Router, lectService *service.LectureService) {

	lecture := router.Group("/lectures", middleware.AuthProtected())
	lecture.Post("/", lectService.Create)
	lecture.Get("/", lectService.GetAll)
	lecture.Get("/current", lectService.GetCurrentLecture)
	lecture.Get("/:id", lectService.GetByID)
	lecture.Put("/:id", lectService.Update)
	lecture.Delete("/:id", middleware.VerifyRole("Admin"), lectService.Delete)

}
