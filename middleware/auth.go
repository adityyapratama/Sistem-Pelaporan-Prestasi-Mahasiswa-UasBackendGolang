package middleware

import (
	"strings"
	"uas-pelaporan-prestasi-mahasiswa/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Token wajib ada"})
		}

		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Format token salah"})
		}

		tokenString := parts[1]

		
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Token tidak valid atau kadaluwarsa"})
		}

		c.Locals("user_id", claims.UserID.String())
		c.Locals("role", claims.RoleName)

		return c.Next()
	}
}

func VerifyRole(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        
        userRole := c.Locals("role")

        
        if userRole == nil {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Role tidak ditemukan"})
        }

        roleStr := userRole.(string)

        
        isAllowed := false
        for _, role := range allowedRoles {
        
            if strings.EqualFold(roleStr, role) {
                isAllowed = true
                break
            }
        }

        if !isAllowed {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden: Anda tidak memiliki akses untuk fitur ini",
            })
        }

        return c.Next()
    }
}