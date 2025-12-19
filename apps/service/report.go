package service

import (
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportService struct {
	reportRepo      *repository.ReportRepo
	studentRepo     repository.StudentsRepository
	achievementRepo repository.AchievementRepository
}

func NewReportService(
	reportRepo *repository.ReportRepo,
	studentRepo repository.StudentsRepository,
	achievementRepo repository.AchievementRepository,
) *ReportService {
	return &ReportService{
		reportRepo:      reportRepo,
		studentRepo:     studentRepo,
		achievementRepo: achievementRepo,
	}
}

// GetStatistics godoc
// @Summary      Dapatkan statistik sistem
// @Description  Mengambil statistik umum sistem termasuk jumlah mahasiswa, dosen, dan prestasi (khusus Admin/Dosen)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Statistik sistem"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil statistik"
// @Router       /reports/statistics [get]
func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := s.reportRepo.GetStatistics(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil statistik: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Statistik berhasil diambil",
		"data":    stats,
	})
}

// GetStudentReport godoc
// @Summary      Dapatkan laporan mahasiswa
// @Description  Mengambil laporan prestasi lengkap per mahasiswa berdasarkan ID (khusus Admin/Dosen)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID mahasiswa (UUID)"
// @Success      200  {object}  map[string]interface{} "Laporan prestasi mahasiswa"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      404  {object}  map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /reports/student/{id} [get]
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	id := c.Params("id")
	studentUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx := c.Context()

	// Get student data
	student, err := s.studentRepo.GetByID(ctx, studentUUID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa tidak ditemukan"})
	}

	// Get achievement stats for student
	achievementStats, total, err := s.reportRepo.GetAchievementStatsByStudentID(ctx, studentUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil statistik prestasi: " + err.Error()})
	}

	// Get all achievements for student
	achievements, err := s.achievementRepo.GetAllByStudentID(ctx, studentUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data prestasi: " + err.Error()})
	}

	response := models.StudentReportResponse{
		Student:           *student,
		TotalAchievements: total,
		AchievementStats:  *achievementStats,
		Achievements:      achievements,
	}

	return c.JSON(fiber.Map{
		"message": "Laporan mahasiswa berhasil diambil",
		"data":    response,
	})
}
