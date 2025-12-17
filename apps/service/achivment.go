package service

import (
	
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

func NewAchievementService(aRepo repository.AchievementRepository, sRepo repository.StudentsRepository) *AchievementService {
	return &AchievementService{
		achievementRepo: aRepo,
		studentRepo:     sRepo,
	}
}

func (s *AchievementService) Create(c *fiber.Ctx) error {
	
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Judul prestasi wajib diisi"})
	}

	
	userVal := c.Locals("user_id")
	if userVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, err := uuid.Parse(userVal.(string))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	
	ctx := c.Context()
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profil mahasiswa tidak ditemukan. Silakan lengkapi biodata terlebih dahulu."})
	}

	
	newAchievement := &models.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: "dummy-mongo-id-" + uuid.NewString(), // Placeholder Mongo
		Status:             "submitted",
		SubmittedAt:        time.Now(),
	}

	
	if err := s.achievementRepo.CreateReference(ctx, newAchievement); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Prestasi berhasil dilaporkan",
		"data":    newAchievement,
	})
}