package service

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
)

type PermissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) *PermissionService {
	return &PermissionService{permissionRepo: permissionRepo}
}

func ( s *PermissionService) Create( c *fiber.Ctx) error{
	var req models.Permission

	if err := c.BodyParser(&req) ; err !=nil{
		return c.Status(400).JSON(fiber.Map{"error" :"request body gk valid"})
	}

	if req.Name == ""|| req.Resource == "" || req.Action == "" {
		return c.Status(400).JSON(fiber.Map{"error" : " isi form wajib di isi semua"})
	}

	ctx := c.Context()
	if err := s.permissionRepo.Create(ctx, &req) ; err != nil{
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat permission"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message" :"data permission berhasil di buat",
		"data" : req,
	})
}

func (s *PermissionService)GetAll(c *fiber.Ctx) error {
	ctx := c.Context()
	permissions, err := s.permissionRepo.GetAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data permission"})
	}

	return c.JSON(fiber.Map{
		"message": "List Permission",
		"data":    permissions,
	})
}