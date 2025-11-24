package service

import (
    "uas-pelaporan-prestasi-mahasiswa/apps/models"
    "uas-pelaporan-prestasi-mahasiswa/apps/repository"
    "uas-pelaporan-prestasi-mahasiswa/utils"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type AuthService struct{
    userRepo repository.UserRepository
}


func NewAuthService(userRepo repository.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

func(s *AuthService) Register(c *fiber.Ctx) error {
    var req models.RegisterRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error" : "request body gk valid"})
    }

    ctx := c.Context()
    existUser, _ := s.userRepo.GetByUsernameOrEmail(ctx, req.Username)
    
    if existUser != nil{
        return c.Status(400).JSON(fiber.Map{"error" : "username atau email udah ada / terdaftar"})
    }

    hashedPassword , err := utils.HashPassword(req.Password)
    if err !=nil{
        return c.Status(500).JSON(fiber.Map{"error": "password gagal di hash"})
    }

    roleUUID, err := uuid.Parse(req.RoleID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Role ID tidak valid"})
    }

    newUser := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        FullName:     req.FullName,
        RoleID:       roleUUID,
    }
    
    if err := s.userRepo.Create(ctx, newUser); err !=nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user ke database"})
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "User berhasil dibuat",
        "data":    newUser,
    })
}

func(s *AuthService) Login(c *fiber.Ctx) error {
    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error" : "request body gk valid"})
    }
    ctx := c.Context()
    
    
    user, err := s.userRepo.GetByUsernameOrEmail(ctx, req.Username)
    
    
    if err != nil || user == nil {
        return c.Status(401).JSON(fiber.Map{"error" : "username atau password salah"})
    }

    
    if !utils.CheckPassword(req.Password, user.PasswordHash){
        return c.Status(401).JSON(fiber.Map{"error" : "username atau password salah"})
    }

    
    if !user.IsActive {
        return c.Status(403).JSON(fiber.Map{"error" : "akun anda mati"})
    }

    // Ambil Role Name
    roleName := "Unknown"
    if user.Role != nil{
        roleName = user.Role.Name
    }

    
    token, err := utils.GenerateToken(user.ID, roleName)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
    }

    return c.JSON(fiber.Map{
        "message": "Login berhasil",
        "data": fiber.Map{
            "token": token,
            "user":  user,
        },
    })
}