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

// Create godoc
// @Summary      Buat prestasi baru
// @Description  Melaporkan prestasi baru ke sistem (khusus Mahasiswa). Data disimpan di MongoDB dan PostgreSQL
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateAchievementRequest true "Data prestasi (type, title, description, details, tags, attachments)"
// @Success      201  {object}  map[string]interface{} "Prestasi berhasil dilaporkan"
// @Failure      400  {object}  map[string]interface{} "Request tidak valid"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      404  {object}  map[string]interface{} "Profil mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal menyimpan data"
// @Router       /achievements [post]
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

// UploadAttachment godoc
// @Summary      Upload lampiran prestasi
// @Description  Mengunggah file lampiran untuk prestasi (jpg, png, pdf). Maks 10MB
// @Tags         Achievements
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file formData file true "File yang akan di-upload (Maks 10MB, Tipe: jpg/png/pdf)"
// @Success      200  {object}  map[string]interface{} "File berhasil diupload"
// @Failure      400  {object}  map[string]interface{} "Format file tidak diizinkan"
// @Failure      500  {object}  map[string]interface{} "Gagal menyimpan file"
// @Router       /achievements/upload [post]
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal mengambil file. Pastikan key-nya 'file'"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
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

// GetAll godoc
// @Summary      Dapatkan semua prestasi
// @Description  Mengambil daftar seluruh prestasi dengan filter status opsional (draft, submitted, verified, rejected)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status query string false "Filter berdasarkan status (draft/submitted/verified/rejected)"
// @Success      200  {object}  map[string]interface{} "Daftar prestasi"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /achievements [get]
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

// GetByID godoc
// @Summary      Dapatkan detail prestasi
// @Description  Mengambil detail prestasi berdasarkan ID termasuk data dari PostgreSQL dan MongoDB
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID prestasi (UUID)"
// @Success      200  {object}  map[string]interface{} "Detail prestasi (info dari PostgreSQL, detail dari MongoDB)"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid"
// @Failure      404  {object}  map[string]interface{} "Prestasi tidak ditemukan"
// @Router       /achievements/{id} [get]
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

// Verify godoc
// @Summary      Verifikasi prestasi
// @Description  Memverifikasi atau menolak prestasi mahasiswa (khusus Dosen). Wajib sertakan alasan jika ditolak
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID prestasi (UUID)"
// @Param        request body models.VerifyAchievementRequest true "Status verifikasi (verified/rejected) dan catatan"
// @Success      200  {object}  map[string]interface{} "Status berhasil diperbarui"
// @Failure      400  {object}  map[string]interface{} "ID atau status tidak valid"
// @Failure      500  {object}  map[string]interface{} "Gagal update status"
// @Router       /achievements/{id}/verify [patch]
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

// Update godoc
// @Summary      Update prestasi
// @Description  Memperbarui data prestasi (khusus pemilik, status belum verified)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID prestasi (UUID)"
// @Param        request body models.CreateAchievementRequest true "Data prestasi yang akan diupdate"
// @Success      200  {object}  map[string]interface{} "Data berhasil diperbarui"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid atau prestasi sudah verified"
// @Failure      403  {object}  map[string]interface{} "Tidak berhak mengedit"
// @Failure      404  {object}  map[string]interface{} "Prestasi tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal update data"
// @Router       /achievements/{id} [put]
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

// GetAllMongoData godoc
// @Summary      Dapatkan semua data dari MongoDB
// @Description  Mengambil semua raw data prestasi dari MongoDB (untuk debugging)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Data dari MongoDB"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /achievements/mongo [get]
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

// GetMyAchievment godoc
// @Summary      Dapatkan prestasi saya
// @Description  Mengambil daftar prestasi milik mahasiswa yang sedang login
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Daftar prestasi Anda"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      404  {object}  map[string]interface{} "Profil mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal mengambil data"
// @Router       /achievements/achiev [get]
func (s *AchievementService) GetMyAchievment(c *fiber.Ctx) error {
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
		return c.Status(404).JSON(fiber.Map{"error": "Profil mahasiswa tidak ditemukan"})
	}

	refs, err := s.achievementRepo.GetAllByStudentID(ctx, student.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data prestasi: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil daftar prestasi Anda",
		"data":    refs,
	})
}

// Delete godoc
// @Summary      Hapus prestasi (soft delete)
// @Description  Menghapus prestasi dengan soft delete (khusus pemilik, status belum verified)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID prestasi (UUID)"
// @Success      200  {object}  map[string]interface{} "Prestasi berhasil dihapus"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid atau prestasi sudah verified"
// @Failure      403  {object}  map[string]interface{} "Tidak berhak menghapus"
// @Failure      404  {object}  map[string]interface{} "Prestasi tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal menghapus prestasi"
// @Router       /achievements/{id} [delete]
func (s *AchievementService) Delete(c *fiber.Ctx) error {
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
		return c.Status(403).JSON(fiber.Map{"error": "Anda tidak berhak menghapus prestasi ini"})
	}

	if ref.Status == "verified" {
		return c.Status(400).JSON(fiber.Map{"error": "Prestasi yang sudah diverifikasi tidak bisa dihapus"})
	}

	err = s.achievementRepo.SoftDelete(ctx, achievUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus prestasi: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Prestasi berhasil dihapus",
	})
}

// Submit godoc
// @Summary      Submit prestasi untuk verifikasi
// @Description  Mengubah status prestasi dari draft menjadi submitted (khusus pemilik)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID prestasi (UUID)"
// @Success      200  {object}  map[string]interface{} "Prestasi berhasil disubmit"
// @Failure      400  {object}  map[string]interface{} "ID tidak valid atau status bukan draft"
// @Failure      403  {object}  map[string]interface{} "Tidak berhak submit"
// @Failure      404  {object}  map[string]interface{} "Prestasi tidak ditemukan"
// @Failure      500  {object}  map[string]interface{} "Gagal submit prestasi"
// @Router       /achievements/{id}/submit [patch]
func (s *AchievementService) Submit(c *fiber.Ctx) error {
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
		return c.Status(403).JSON(fiber.Map{"error": "Anda tidak berhak submit prestasi ini"})
	}

	if ref.Status != "Draft" && ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "Hanya prestasi dengan status draft yang bisa disubmit"})
	}

	err = s.achievementRepo.Submit(ctx, achievUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal submit prestasi: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Prestasi berhasil disubmit untuk verifikasi",
	})
}
