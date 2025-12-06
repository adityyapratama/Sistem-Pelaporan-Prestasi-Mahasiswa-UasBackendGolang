package service

import "C"
import (
	"net/http"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"

	"github.com/gofiber/fiber/v2"
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

func (s *AchievementService) Create(fiber fiber.Ctx) {
	var req models.Achievement
	if err := c.BodyParser(fiber, &req); err != nil {
		fiber.Status(http.StatusBadRequest)

	}
}
