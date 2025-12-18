package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)


func StudentRoutes(router fiber.Router, studentService *service.StudentService) {
	
	
	students := router.Group("/students",middleware.AuthProtected())
	
	students.Get("/", middleware.VerifyRole("Admin"),studentService.GetAll)
	students.Get("/advisor/:id", middleware.VerifyRole("Admin"),studentService.GetByAdvisorID)
	students.Post("/", studentService.Create)
	students.Get("/current", studentService.GetCurrentStudent)
	students.Get("/:id", studentService.GetByID)
	students.Put("/:id/advisor", middleware.VerifyRole("Admin"), studentService.AssignAdvisor)
	

	

	


}