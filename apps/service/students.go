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

func NewStudentService(studentRepo repository.StudentsRepository)*StudentService {
	return &StudentService{studentRepo: studentRepo}
}


func( s *StudentService) Create(c *fiber.Ctx) error {
	var req models.Students
	if err := c.BodyParser(&req) ; err !=nil{
		return c.Status(400).JSON(fiber.Map{"error" :"request body gk valid"})
	}

	if req.StudentID =="" || req.ProgramStudy== ""  {
		return c.Status(400).JSON(fiber.Map{"error": "NIM dan Program Studi wajib diisi"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID,_ := uuid.Parse(userIDStr)

	ctx := c.Context()
	existing, _ := s.studentRepo.GetByUserID(ctx, userID)
	if existing != nil {
    return c.Status(409).JSON(fiber.Map{"error": "Data mahasiswa untuk user ini sudah ada"})
}

	student := &models.Students{
		UserID: userID,
		StudentID: req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
	}

	if err := s.studentRepo.Create(ctx, student); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data mahasiswa"})
    }

	return c.Status(201).JSON(fiber.Map{
		"messege" :"profile mahasiswa berhasil di buat",
		"data" : student,
	})

}

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
		"data" :student,
	})

}

func (s *StudentService) SetAdvisor(c *fiber.Ctx)error {
	studentIDStr := c.Params("id")
	studentID, err := uuid.Parse(studentIDStr)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id mahasiswa tidak valid"})
	}

	var req models.SetAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}


	advisorID, err := uuid.Parse(req.AdvisorID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID Dosen Wali tidak valid"})
	}

	ctx := c.Context()
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa tidak ditemukan"})
	}

	student.AdvisorID= &advisorID

	if err := s.studentRepo.Update(ctx, student); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update dosen wali",
        "detail": err.Error(),})
	}

	return c.JSON(fiber.Map{
		"message": "Dosen wali berhasil ditetapkan",
		"data": student,
	})

}
