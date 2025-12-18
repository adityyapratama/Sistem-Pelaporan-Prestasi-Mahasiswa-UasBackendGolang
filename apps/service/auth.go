package service

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"
	"uas-pelaporan-prestasi-mahasiswa/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo repository.UserRepository
	permissionRepo repository.PermissionRepository
}

func NewAuthService(userRepo repository.UserRepository, permissionRepo repository.PermissionRepository) *AuthService {
	return &AuthService{userRepo: userRepo, permissionRepo: permissionRepo}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "request body gk valid"})
	}

	ctx := c.Context()
	existUser, _ := s.userRepo.GetByUsernameOrEmail(ctx, req.Username)

	if existUser != nil {
		return c.Status(400).JSON(fiber.Map{"error": "username atau email udah ada / terdaftar"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
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

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user ke database"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User berhasil dibuat",
		"data":    newUser,
	})
}

func (s *AuthService) Login(c *fiber.Ctx) error {
    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    
    ctx := c.Context()
    user, err := s.userRepo.GetByUsernameOrEmail(ctx, req.Username)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    
    permissions, err := s.permissionRepo.GetByRoleID(ctx, user.RoleID)
    
    
    var permissionList []string
    if err == nil {
        for _, p := range permissions {
            permissionList = append(permissionList, p.Name)
        }
    }
    

    
    token, refreshToken, err := utils.GenerateToken(user.ID, user.Role.Name)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
    }
    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "token":        token,
            "refreshToken": refreshToken,
            "user": fiber.Map{
                "id":       user.ID,
                "username": user.Username,
                "fullName": user.FullName,
                "role":     user.Role.Name,
                
                "permissions": permissionList, 
            },
        },
    })
}
func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "user id tidak valid",
		})
	}

	ctx := c.Context()
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user tidak di temukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil data profile user",
		"data":    user,
	})

}

func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "request body tidak valid",
		})
	}

	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "refresh token tidak valid atau expired",
		})
	}

	ctx := c.Context()
	user, err := s.userRepo.GetByID(ctx, claims.UserID)

	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	roleName := "unknow"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, roleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token baru"})
	}

	return c.JSON(fiber.Map{
		"message":      "Token berhasil diperbarui",
		"access_token": newAccessToken,
	})
}


func (s *AuthService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userUUID, err := uuid.Parse(id)
	if err != nil{
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req models.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	req.ID = userUUID
	ctx := c.Context()
	if err := s.userRepo.Update(ctx, &req); err != nil{
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update user"})
	
	}

	return  c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		})

}


func (s *AuthService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx := c.Context()
	if err := s.userRepo.Delete(ctx, userUUID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}


func (s *AuthService) GetAllUser(c *fiber.Ctx) error {
	ctx := c.Context()
	users , err := s.userRepo.GetAll(ctx)
	if err!= nil{
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}
	return c.JSON(fiber.Map{
		"message" : "data user berhasil di hapus",
		"data": users,
	})
}

func (s *AuthService) UpdateRoleUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	roleUUID, err := uuid.Parse(req.RoleID)
	if err != nil{
		return c.Status(400).JSON(fiber.Map{"error": "Role ID tidak valid"})
	}

	ctx := c.Context()
	if err := s.userRepo.UpdateRole(ctx, roleUUID,userUUID) ; err != nil{
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update role"})
	}
	return c.JSON(fiber.Map{"message": "Role user berhasil diupdate"})
	
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
    userIdStr := c.Locals("user_id")
    if userIdStr == nil {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "Berhasil logout",
    })
}