package service
import (
	"net/http"
	"time"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentsRepository
}

func NewAchivmentService(achievementRepo repository.AchievementRepository, studentRepo repository.StudentsRepository) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
	}
}

func (s *AchievementService) Create(c *fiber.Ctx) error {
	// A. Validasi Input Body
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Judul prestasi wajib diisi"})
	}

	// B. Ambil User ID dari Token (Siapa yang login?)
	userVal := c.Locals("user_id")
	if userVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, _ := uuid.Parse(userVal.(string))

	// C. PENTING: Cari Student ID berdasarkan User ID
	// Karena tabel achievement butuh student_id, bukan user_id.
	ctx := c.Context()
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profil mahasiswa belum dibuat. Silakan lengkapi profil dulu."})
	}

	// D. Simpan ke Database (PostgreSQL)
	// Catatan: Jika kamu pakai MongoDB juga, simpan ke Mongo dulu, ambil ID-nya,
	// lalu masukkan ke field MongoAchievementID. Disini kita simulasi string dulu.
	newAchievement := &models.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: "dummy-mongo-id-" + uuid.NewString(), // Nanti diganti logic Mongo asli
		Status:             "submitted", // Status awal selalu submitted/draft
		SubmittedAt:        time.Now()
	}

	if err := s.achievementRepo.Create(ctx, newAchievement); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan prestasi"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Prestasi berhasil dilaporkan",
		"data":    newAchievement,
	})
}
