package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(router fiber.Router, reportService *service.ReportService) {
	reports := router.Group("/reports", middleware.AuthProtected())

	reports.Get("/statistics", middleware.VerifyRole("Admin", "Dosen"), reportService.GetStatistics)
	reports.Get("/student/:id", middleware.VerifyRole("Admin", "Dosen"), reportService.GetStudentReport)
}
