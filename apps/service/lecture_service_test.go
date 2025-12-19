package service_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository/mocks"
	"uas-pelaporan-prestasi-mahasiswa/apps/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TestCreateLecture_TableDriven - Test untuk POST /lectures
func TestCreateLecture_TableDriven(t *testing.T) {
	// Setup
	mockRepo := mocks.NewManualMockLectureRepo()
	lectureService := service.NewLectureService(mockRepo)
	app := fiber.New()

	// Simulasi middleware auth dengan user_id dan role di Locals
	app.Post("/lectures", func(c *fiber.Ctx) error {
		c.Locals("user_id", "550e8400-e29b-41d4-a716-446655440001")
		c.Locals("role", "Dosen") // Tambahkan role
		return c.Next()
	}, lectureService.Create)

	tests := []struct {
		name           string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name: "Create Lecture Sukses",
			inputBody: map[string]string{
				"lecturer_id": "D001",
				"department":  "Teknik Informatika",
			},
			expectedStatus: 201,
		},
		{
			name: "NID Kosong",
			inputBody: map[string]string{
				"lecturer_id": "",
				"department":  "Teknik Informatika",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/lectures", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("Error request: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

// TestGetAllLectures - Test untuk GET /lectures
func TestGetAllLectures(t *testing.T) {
	// Setup
	mockRepo := mocks.NewManualMockLectureRepo()
	lectureService := service.NewLectureService(mockRepo)
	app := fiber.New()
	app.Get("/lectures", lectureService.GetAll)

	// Seed data
	mockRepo.Create(nil, &models.Lecture{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		LecturerID: "D001",
		Department: "Teknik Informatika",
	})

	// Test
	req := httptest.NewRequest("GET", "/lectures", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Errorf("Error request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Status salah! Dapat %d, Harapan 200", resp.StatusCode)
	}
}

// TestGetLectureByID_TableDriven - Test untuk GET /lectures/:id
func TestGetLectureByID_TableDriven(t *testing.T) {
	// Setup
	mockRepo := mocks.NewManualMockLectureRepo()
	lectureService := service.NewLectureService(mockRepo)
	app := fiber.New()
	app.Get("/lectures/:id", lectureService.GetByID)

	// Seed data
	existingID := uuid.New()
	mockRepo.Create(nil, &models.Lecture{
		ID:         existingID,
		UserID:     uuid.New(),
		LecturerID: "D001",
		Department: "Teknik Informatika",
	})

	tests := []struct {
		name           string
		lectureID      string
		expectedStatus int
	}{
		{
			name:           "Lecture Ditemukan",
			lectureID:      existingID.String(),
			expectedStatus: 200,
		},
		{
			name:           "Lecture Tidak Ditemukan",
			lectureID:      uuid.New().String(),
			expectedStatus: 404,
		},
		{
			name:           "ID Tidak Valid",
			lectureID:      "invalid-uuid",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/lectures/"+tt.lectureID, nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("Error request: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

// TestUpdateLecture_TableDriven - Test untuk PUT /lectures/:id
func TestUpdateLecture_TableDriven(t *testing.T) {
	// Setup
	mockRepo := mocks.NewManualMockLectureRepo()
	lectureService := service.NewLectureService(mockRepo)
	app := fiber.New()
	app.Put("/lectures/:id", lectureService.Update)

	// Seed data
	existingID := uuid.New()
	mockRepo.Create(nil, &models.Lecture{
		ID:         existingID,
		UserID:     uuid.New(),
		LecturerID: "D001",
		Department: "Teknik Informatika",
	})

	tests := []struct {
		name           string
		lectureID      string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name:      "Update Sukses",
			lectureID: existingID.String(),
			inputBody: map[string]string{
				"department": "Sistem Informasi",
			},
			expectedStatus: 200,
		},
		{
			name:      "Lecture Tidak Ditemukan",
			lectureID: uuid.New().String(),
			inputBody: map[string]string{
				"department": "Teknik Informatika",
			},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("PUT", "/lectures/"+tt.lectureID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("Error request: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

// TestDeleteLecture_TableDriven - Test untuk DELETE /lectures/:id
func TestDeleteLecture_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(*mocks.ManualMockLectureRepo) string
		expectedStatus int
	}{
		{
			name: "Delete Sukses",
			setupFunc: func(repo *mocks.ManualMockLectureRepo) string {
				id := uuid.New()
				repo.Create(nil, &models.Lecture{
					ID:         id,
					UserID:     uuid.New(),
					LecturerID: "D001",
					Department: "Teknik Informatika",
				})
				return id.String()
			},
			expectedStatus: 200,
		},
		{
			name: "Lecture Tidak Ditemukan",
			setupFunc: func(repo *mocks.ManualMockLectureRepo) string {
				return uuid.New().String()
			},
			expectedStatus: 500, // Service return 500 saat delete gagal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewManualMockLectureRepo()
			lectureService := service.NewLectureService(mockRepo)
			app := fiber.New()
			app.Delete("/lectures/:id", lectureService.Delete)

			lectureID := tt.setupFunc(mockRepo)
			req := httptest.NewRequest("DELETE", "/lectures/"+lectureID, nil)

			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("Error request: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}
