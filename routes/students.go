package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)


func StudentRoutes(router fiber.Router, studentService *service.StudentService) {
	
	
	studens := router.Group("/students",middleware.AuthProtected())
	

	studens.Post("/", studentService.Create)
	studens.Get("/current", studentService.GetCurrentStudent)
	

	

	


}