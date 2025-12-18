package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

    newDetail := &models.AchievementDetail{
        StudentID:       student.ID.String(),
        AchievementType: req.Type,
        Title:           req.Title,
        Description:     req.Description,
        Details:         req.Details,
        Attachments:     req.Attachments,
        Tags:            req.Tags,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    if err := s.achievementRepo.CreateDetail(ctx, newDetail); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal ke MongoDB: " + err.Error()})
    }
    
    
    mongoIDString := newDetail.ID.Hex()
    newRef := &models.AchievementReference{
        StudentID:          student.ID,
        MongoAchievementID: mongoIDString, 
        Status:             "Draft",
        SubmittedAt:        time.Now(),
    }

    
    if err := s.achievementRepo.CreateReference(ctx, newRef); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal ke PostgreSQL: " + err.Error()})
    }


 return c.Status(201).JSON(fiber.Map{
        "message": "Prestasi berhasil dilaporkan",
        "data": fiber.Map{
            "id":                   newRef.ID,
            "mongo_achievement_id": mongoIDString,
            "status":               newRef.Status,
            "detail":               newDetail, 
        },
    })
}

func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err!= nil{
		return c.Status(400).JSON(fiber.Map{"error": "Gagal mengambil file. Pastikan key-nya 'file'"})
	}

	ext:=strings.ToLower(filepath.Ext(file.Filename))
	allowedExts :=map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !allowedExts[ext] {
        return c.Status(400).JSON(fiber.Map{"error": "Format file tidak diizinkan (hanya jpg, png, pdf)"})
    }

	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
    savePath := fmt.Sprintf("./uploads/%s", filename)

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
    os.Mkdir("./uploads", 0755)
    }
	if err := c.SaveFile(file, savePath); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file ke server"})
    }
    
    fileURL := fmt.Sprintf("http://localhost:3000/uploads/%s", filename)

    return c.Status(200).JSON(fiber.Map{
        "message": "File berhasil diupload",
        "data": fiber.Map{
            "fileName": file.Filename,
            "fileUrl":  fileURL,
            "fileType": file.Header["Content-Type"][0],
        },
    })
}



func (s *AchievementService) GetAll(c *fiber.Ctx) error {
    status := c.Query("status")
    
    refs, err := s.achievementRepo.GetAll(c.Context(), status)
    if err != nil {
         return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "message": "Berhasil mengambil daftar prestasi",
        "data":    refs,
    })
}


func (s *AchievementService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	achievementUUID, err := uuid.Parse(id)
	if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

	ctx := c.Context()

	ref, err := s.achievementRepo.GetReferenceByID(ctx, achievementUUID)
    if err != nil || ref == nil {
        return c.Status(404).JSON(fiber.Map{"error": err.Error()})
    }

	mongoID := ref.MongoAchievementID
	detail, err := s.achievementRepo.GetDetailByID(ctx, mongoID)
	var detailData interface{} = detail
    if err != nil {
        detailData = "Detail data di MongoDB tidak ditemukan atau rusak."
    }

    
    return c.JSON(fiber.Map{
        "message": "Detail prestasi ditemukan",
        "data": fiber.Map{
            "info":   ref,    
            "detail": detailData, 
        },
    })



}

func (s *AchievementService) Verify(c *fiber.Ctx) error {
    id := c.Params("id")
    achievUUID, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }


    dosenIDStr := c.Locals("user_id").(string)
    dosenUUID, _ := uuid.Parse(dosenIDStr)

    
    var req models.VerifyAchievementRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    
    if req.Status != "verified" && req.Status != "rejected" {
        return c.Status(400).JSON(fiber.Map{"error": "Status harus sesuai dengan format 'verified' atau 'rejected'"})
    }

    if req.Status == "rejected" && req.Notes == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Wajib sertakan alasan penolakan"})
    }

    ctx := c.Context()
    now := time.Now() 

    updateData := &models.AchievementReference{
        ID:            achievUUID,
        Status:        req.Status,
        VerifiedBy:    &dosenUUID, 
        VerifiedAt:    &now,       
        RejectionNote: &req.Notes,
    }

    
    err = s.achievementRepo.UpdateStatus(ctx, updateData)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal update status: " + err.Error()})
    }

    return c.JSON(fiber.Map{
        "message": "Status prestasi berhasil diperbarui",
        "status":  req.Status,
    })
}



func (s *AchievementService) Update(c *fiber.Ctx) error {

    id := c.Params("id")
    achievUUID, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }


    userIDStr := c.Locals("user_id").(string)
    userUUID, _ := uuid.Parse(userIDStr) 

    ctx := c.Context()

    
    ref, err := s.achievementRepo.GetReferenceByID(ctx, achievUUID)
    if err != nil || ref == nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data prestasi tidak ditemukan"})
    }

    
    student, err := s.studentRepo.GetByUserID(ctx, userUUID)
    if err != nil || student == nil {
        return c.Status(403).JSON(fiber.Map{"error": "Data mahasiswa tidak ditemukan"})
    }

    
    if ref.StudentID != student.ID {
        return c.Status(403).JSON(fiber.Map{"error": "Anda tidak berhak mengedit prestasi ini"})
    }

    
    if ref.Status == "verified" {
        return c.Status(400).JSON(fiber.Map{"error": "Prestasi yang sudah diverifikasi tidak bisa diedit"})
    }

    
    var req models.CreateAchievementRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    
    updateDetail := &models.AchievementDetail{
        Title:           req.Title,
        AchievementType: req.Type,
        Description:     req.Description,
        Details:         req.Details,
        Attachments:     req.Attachments,
        Tags:            req.Tags,
        
    }

    
    err = s.achievementRepo.UpdateDetail(ctx, ref.MongoAchievementID, updateDetail)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal update data: " + err.Error()})
    }

    return c.JSON(fiber.Map{
        "message": "Data prestasi berhasil diperbarui",
    })
}


func (s *AchievementService) GetAllMongoData(c *fiber.Ctx) error {
    details, err := s.achievementRepo.GetAllDetailsFromMongo(c.Context())
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Gagal mengambil data dari MongoDB: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Berhasil mengambil semua raw data dari MongoDB",
        "total":   len(details),
        "data":    details,
    })
}