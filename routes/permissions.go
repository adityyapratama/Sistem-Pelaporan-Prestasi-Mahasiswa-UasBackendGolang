package routes


import (
	"uas-pelaporan-prestasi-mahasiswa/apps/service"

	"github.com/gofiber/fiber/v2"
)


func PermissionRoutes(router fiber.Router, permService *service.PermissionService) {
	
	
	permission := router.Group("/permissions")
	

	permission.Post("/", permService.Create)
	permission.Get("/", permService.GetAll)
	

	

	


}