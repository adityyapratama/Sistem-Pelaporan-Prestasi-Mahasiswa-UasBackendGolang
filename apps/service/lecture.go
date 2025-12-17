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


// [FIX] Ganti 'unc' jadi 'func'
func (s *LectureService) Create(c *fiber.Ctx) error {
    var req models.CreateLectureRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
    }

    if req.LecturerID == "" || req.Department == "" {
        return c.Status(400).JSON(fiber.Map{"error": "NIP dan Departemen wajib diisi"})
    }

    // [FIX] Defensive Coding: Cek dulu apakah data dari token ada?
    userVal := c.Locals("user_id")
    roleVal := c.Locals("role")
    
    // Kalau tidak ada data di locals, berarti belum login / middleware auth tidak jalan
    if userVal == nil || roleVal == nil {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Silakan login ulang"})
    }

    // Aman untuk di-convert ke string
    userIDLogin := userVal.(string)
    roleLogin := roleVal.(string)

    var targetUserID uuid.UUID
    var err error
    
    // Cek Role (Case Insensitive biar aman)
    if roleLogin == "Admin" || roleLogin == "admin" {
        // SKENARIO 1: ADMIN YANG INPUT
        if req.UserID == "" {
            return c.Status(400).JSON(fiber.Map{"error": "Admin wajib mengisi field user_id"})
        }
        targetUserID, err = uuid.Parse(req.UserID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Format user_id tidak valid"})
        }
    } else {
        // SKENARIO 2: DOSEN INPUT SENDIRI
        targetUserID, _ = uuid.Parse(userIDLogin)
    }

    ctx := c.Context()

    // [Catatan] Pastikan nama field di struct service kamu 'lectureRepo' atau 'lecturerRepo'
    existing, _ := s.lectureRepo.GetByUserID(ctx, targetUserID)
    if existing != nil {
        return c.Status(409).JSON(fiber.Map{"error": "Profil dosen untuk user ini sudah ada"})
    }

    // [FIX] Pastikan nama struct modelnya benar. Biasanya 'Lecturer'
    lecture := &models.Lecture{ 
        UserID:     targetUserID,
        LecturerID: req.LecturerID,
        Department: req.Department,
    }

    if err := s.lectureRepo.Create(ctx, lecture); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data dosen"})
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "Profil dosen berhasil dibuat", // [FIX] Pesan disesuaikan
        "data":    lecture,
    })
}




func (s *LectureService) GetAll(c *fiber.Ctx) error {
	ctx := c.Context()
	lecturers, err := s.lectureRepo.GetAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data dosen"})
	}
	return c.JSON(fiber.Map{"data": lecturers})
}