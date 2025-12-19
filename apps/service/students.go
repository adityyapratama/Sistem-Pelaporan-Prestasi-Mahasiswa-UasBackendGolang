package service

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService struct {
	studentRepo repository.StudentsRepository
}

func NewStudentService(studentRepo repository.StudentsRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo}
}

// Create godoc
// @Summary      Buat profil mahasiswa
// @Description  Membuat profil mahasiswa baru untuk pengguna yang sudah terdaftar
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateStudentRequest true "Data profil mahasiswa"
// @Success      201  {object}  map[string]interface{} "Profil mahasiswa berhasil dibuat"
// @Failure      400  {object}  map[string]interface{} "Request tidak valid"
// @Failure      409  {object}  map[string]interface{} "Data mahasiswa sudah ada"
// @Failure      500  {object}  map[string]interface{} "Gagal menyimpan data"
// @Router       /students [post]
func (s *StudentService) Create(c *fiber.Ctx) error {
	var req models.Students
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "request body gk valid"})
	}

	if req.StudentID == "" || req.ProgramStudy == "" {
		return c.Status(400).JSON(fiber.Map{"error": "NIM dan Program Studi wajib diisi"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Context()
	existing, _ := s.studentRepo.GetByUserID(ctx, userID)
	if existing != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Data mahasiswa untuk user ini sudah ada"})
	}

	student := &models.Students{
		UserID:       userID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
	}

	if err := s.studentRepo.Create(ctx, student); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data mahasiswa"})
	}

	return c.Status(201).JSON(fiber.Map{
		"messege": "profile mahasiswa berhasil di buat",
		"data":    student,
	})

}

// GetCurrentStudent godoc
// @Summary      Dapatkan profil mahasiswa saat ini
// @Description  Mengambil profil mahasiswa yang sedang login berdasarkan token JWT
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Data profil mahasiswa"
// @Failure      404  {object}  map[string]interface{} "Profil belum dibuat"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /students/current [get]
func (s *StudentService) GetCurrentStudent(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Context()
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	if student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profil mahasiswa belum dibuat. Silakan lengkapi data."})
	}

	return c.JSON(fiber.Map{
		"data": student,
	})

}

// GetByID godoc
// @Summary      Dapatkan mahasiswa berdasarkan ID
// @Description  Mengambil data mahasiswa berdasarkan ID mahasiswa
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID mahasiswa (UUID)"
// @Success      200  {object}  map[string]interface{} "Data mahasiswa"
// @Failure      404  {object}  map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /students/{id} [get]
func (s *StudentService) GetByID(c *fiber.Ctx) error {
	studentIDStr := c.Params("id")
	studentID, _ := uuid.Parse(studentIDStr)

	ctx := c.Context()
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	if student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "mahasiswa tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"data": student,
	})

}

// Update godoc
// @Summary      Update profil mahasiswa
// @Description  Memperbarui data profil mahasiswa berdasarkan ID
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID mahasiswa (UUID)"
// @Param        request body models.CreateStudentRequest true "Data yang akan diupdate"
// @Success      200  {object}  map[string]interface{} "Profil berhasil diperbarui"
// @Failure      400  {object}  map[string]interface{} "ID atau body tidak valid"
// @Failure      404  {object}  map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal update data"
// @Router       /students/{id} [put]
func (s *StudentService) Update(c *fiber.Ctx) error {
	studentIDStr := c.Params("id")
	studentID, err := uuid.Parse(studentIDStr)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id mahasiswa tidak valid"})
	}

	var req models.Students
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	ctx := c.Context()
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa tidak ditemukan"})
	}

	if req.StudentID != "" {
		student.StudentID = req.StudentID
	}
	if req.ProgramStudy != "" {
		student.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		student.AcademicYear = req.AcademicYear
	}

	if err := s.studentRepo.Update(ctx, student); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update data mahasiswa"})
	}

	return c.JSON(fiber.Map{
		"message": "Data mahasiswa berhasil diupdate",
		"data":    student,
	})
}

// AssignAdvisor godoc
// @Summary      Tetapkan dosen pembimbing
// @Description  Menetapkan dosen pembimbing untuk mahasiswa berdasarkan ID (khusus Admin)
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID mahasiswa (UUID)"
// @Param        request body models.SetAdvisorRequest true "ID dosen pembimbing"
// @Success      200  {object}  map[string]interface{} "Dosen pembimbing berhasil ditetapkan"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      500  {object}  map[string]interface{} "Gagal menetapkan dosen wali"
// @Router       /students/{id}/advisor [put]
func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
	studentIDStr := c.Params("id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id mahasiswa tidak valid"})
	}

	var req models.SetAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	advisorUUID, err := uuid.Parse(req.AdvisorID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID Dosen Wali tidak valid"})
	}

	ctx := c.Context()
	err = s.studentRepo.AssignAdvisor(ctx, studentID, advisorUUID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menetapkan dosen wali: " + err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":    "Berhasil! Dosen wali sudah ditetapkan.",
		"student_id": studentID,
		"advisor_id": advisorUUID,
	})

}

// GetAll godoc
// @Summary      Dapatkan semua mahasiswa
// @Description  Mengambil daftar seluruh mahasiswa (khusus Dosen/Admin)
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Daftar semua mahasiswa"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /students [get]
func (s *StudentService) GetAll(c *fiber.Ctx) error {
	ctx := c.Context()
	students, err := s.studentRepo.GetAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data students"})
	}
	return c.JSON(fiber.Map{"data": students})
}

// GetByAdvisorID godoc
// @Summary      Dapatkan mahasiswa berdasarkan dosen pembimbing
// @Description  Mengambil daftar mahasiswa berdasarkan ID dosen pembimbing
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID dosen pembimbing (UUID)"
// @Success      200  {object}  map[string]interface{} "Daftar mahasiswa"
// @Failure      404  {object}  map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /students/advisor/{id} [get]
func (s *StudentService) GetByAdvisorID(c *fiber.Ctx) error {
	advisorIDStr := c.Params("id")
	advisorID, _ := uuid.Parse(advisorIDStr)

	ctx := c.Context()
	student, err := s.studentRepo.GetByAdvisorID(ctx, advisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	if student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "mahasiswa tidak ditemukan"})
	}
	return c.JSON(fiber.Map{
		"data": student,
	})

}

// Delete godoc
// @Summary      Hapus mahasiswa
// @Description  Menghapus data mahasiswa berdasarkan ID (khusus Admin)
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID mahasiswa (UUID)"
// @Success      200  {object}  map[string]interface{} "Data mahasiswa berhasil dihapus"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      500  {object}  map[string]interface{} "Gagal menghapus data"
// @Router       /students/{id} [delete]
func (s *StudentService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	studentUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx := c.Context()
	if err := s.studentRepo.Delete(ctx, studentUUID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data mahasiswa"})
	}

	return c.JSON(fiber.Map{"message": "Data mahasiswa berhasil dihapus"})
}
