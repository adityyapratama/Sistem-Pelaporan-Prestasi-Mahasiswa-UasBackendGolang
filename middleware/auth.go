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
        // 1. Ambil role user yang sudah disimpan oleh AuthProtected tadi
        userRole := c.Locals("role")

        // Jaga-jaga kalau AuthProtected belum dijalankan / token error
        if userRole == nil {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Role tidak ditemukan"})
        }

        roleStr := userRole.(string)

        // 2. Cek apakah role user ada di daftar role yang diizinkan
        isAllowed := false
        for _, role := range allowedRoles {
            // Kita pakai EqualFold biar tidak sensitif huruf besar/kecil (Admin == admin)
            if strings.EqualFold(roleStr, role) {
                isAllowed = true
                break
            }
        }

        // 3. Kalau role tidak cocok, tendang!
        if !isAllowed {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden: Anda tidak memiliki akses untuk fitur ini",
            })
        }

        return c.Next()
    }
}