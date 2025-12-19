package routes

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	_ "uas-pelaporan-prestasi-mahasiswa/docs" // Import docs yang di-generate swag

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App,
	authService *service.AuthService,
	permService *service.PermissionService,
	studentService *service.StudentService,
	lectureService *service.LectureService,
	achievService *service.AchievementService,
	reportService *service.ReportService,
) {
	// app.Use(logger.new())
	app.Use(cors.New())

	api := app.Group("/api/v1")
	RegisterAuthRoutes(api, authService)
	PermissionRoutes(api, permService)
	StudentRoutes(api, studentService)
	LectureRoutes(api, lectureService)
	AchievementRoutes(api, achievService)
	ReportRoutes(api, reportService)

	app.Get("/swagger/*", swagger.HandlerDefault)

}
