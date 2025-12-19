package service

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LectureService struct {
	lectureRepo repository.LectureRepository
}

func NewLectureService(lectureRepo repository.LectureRepository) *LectureService {
	return &LectureService{lectureRepo: lectureRepo}
}

// Create godoc
// @Summary      Buat profil dosen
// @Description  Membuat profil dosen baru. Admin dapat membuat untuk user lain, dosen hanya untuk diri sendiri
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateLectureRequest true "Data profil dosen"
// @Success      201  {object}  map[string]interface{} "Profil dosen berhasil dibuat"
// @Failure      400  {object}  map[string]interface{} "Request tidak valid"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      409  {object}  map[string]interface{} "Profil sudah ada"
// @Failure      500  {object}  map[string]interface{} "Gagal menyimpan data"
// @Router       /lectures [post]
func (s *LectureService) Create(c *fiber.Ctx) error {
	var req models.CreateLectureRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.LecturerID == "" || req.Department == "" {
		return c.Status(400).JSON(fiber.Map{"error": "NIP dan Departemen wajib diisi"})
	}

	userVal := c.Locals("user_id")
	roleVal := c.Locals("role")
	if userVal == nil || roleVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Silakan login ulang"})
	}

	userIDLogin := userVal.(string)
	roleLogin := roleVal.(string)

	var targetUserID uuid.UUID
	var err error

	if roleLogin == "Admin" || roleLogin == "admin" {

		if req.UserID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Admin wajib mengisi field user_id"})
		}
		targetUserID, err = uuid.Parse(req.UserID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format user_id tidak valid"})
		}
	} else {

		targetUserID, _ = uuid.Parse(userIDLogin)
	}

	ctx := c.Context()

	existing, _ := s.lectureRepo.GetByUserID(ctx, targetUserID)
	if existing != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Profil dosen untuk user ini sudah ada"})
	}

	lecture := &models.Lecture{
		UserID:     targetUserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	if err := s.lectureRepo.Create(ctx, lecture); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data dosen"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Profil dosen berhasil dibuat",
		"data":    lecture,
	})
}

// GetAll godoc
// @Summary      Dapatkan semua dosen
// @Description  Mengambil daftar seluruh dosen dalam sistem
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Daftar semua dosen"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /lectures [get]
func (s *LectureService) GetAll(c *fiber.Ctx) error {
	ctx := c.Context()
	lecturers, err := s.lectureRepo.GetAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data dosen"})
	}
	return c.JSON(fiber.Map{"data": lecturers})
}

// GetByID godoc
// @Summary      Dapatkan dosen berdasarkan ID
// @Description  Mengambil data dosen berdasarkan ID
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID dosen (UUID)"
// @Success      200  {object}  map[string]interface{} "Data dosen"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      404  {object}  map[string]interface{} "Dosen tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /lectures/{id} [get]
func (s *LectureService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	lectureUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx := c.Context()
	lecture, err := s.lectureRepo.GetByID(ctx, lectureUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if lecture == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dosen tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"data": lecture})
}

// GetCurrentLecture godoc
// @Summary      Dapatkan profil dosen saat ini
// @Description  Mengambil profil dosen yang sedang login berdasarkan token JWT
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Data profil dosen"
// @Failure      404  {object}  map[string]interface{} "Profil belum dibuat"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /lectures/current [get]
func (s *LectureService) GetCurrentLecture(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Context()
	lecture, err := s.lectureRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	if lecture == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profil dosen belum dibuat. Silakan lengkapi data."})
	}

	return c.JSON(fiber.Map{"data": lecture})
}

// Update godoc
// @Summary      Update profil dosen
// @Description  Memperbarui data profil dosen berdasarkan ID
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID dosen (UUID)"
// @Param        request body models.CreateLectureRequest true "Data yang akan diupdate"
// @Success      200  {object}  map[string]interface{} "Profil berhasil diperbarui"
// @Failure      400  {object}  map[string]interface{} "ID atau body tidak valid"
// @Failure      404  {object}  map[string]interface{} "Dosen tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal update data"
// @Router       /lectures/{id} [put]
func (s *LectureService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	lectureUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req models.CreateLectureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	ctx := c.Context()
	existing, err := s.lectureRepo.GetByID(ctx, lectureUUID)
	if err != nil || existing == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dosen tidak ditemukan"})
	}

	existing.LecturerID = req.LecturerID
	existing.Department = req.Department

	if err := s.lectureRepo.Update(ctx, existing); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update data dosen"})
	}

	return c.JSON(fiber.Map{
		"message": "Data dosen berhasil diupdate",
		"data":    existing,
	})
}

// Delete godoc
// @Summary      Hapus dosen
// @Description  Menghapus data dosen berdasarkan ID (khusus Admin)
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID dosen (UUID)"
// @Success      200  {object}  map[string]interface{} "Data dosen berhasil dihapus"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      500  {object}  map[string]interface{} "Gagal menghapus data"
// @Router       /lectures/{id} [delete]
func (s *LectureService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	lectureUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx := c.Context()
	if err := s.lectureRepo.Delete(ctx, lectureUUID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data dosen"})
	}

	return c.JSON(fiber.Map{"message": "Data dosen berhasil dihapus"})
}
