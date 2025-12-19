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

// TestCreateStudent - Test untuk POST /students
func TestCreateStudent(t *testing.T) {
	mockRepo := mocks.NewManualMockStudentRepo()
	studentService := service.NewStudentService(mockRepo)
	app := fiber.New()

	app.Post("/students", func(c *fiber.Ctx) error {
		c.Locals("user_id", "550e8400-e29b-41d4-a716-446655440001")
		return c.Next()
	}, studentService.Create)

	tests := []struct {
		name           string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name: "Create Student Sukses",
			inputBody: map[string]string{
				"student_id":    "123456",
				"name":          "Test User",
				"program_study": "Informatika",
				"academic_year": "2024",
			},
			expectedStatus: 201,
		},
		{
			name: "Student ID Kosong",
			inputBody: map[string]string{
				"student_id":    "",
				"name":          "Test User",
				"program_study": "Informatika",
				"academic_year": "2024",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/students", bytes.NewReader(bodyBytes))
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

// TestGetAllStudents - Test untuk GET /students
func TestGetAllStudents(t *testing.T) {
	mockRepo := mocks.NewManualMockStudentRepo()
	studentService := service.NewStudentService(mockRepo)
	app := fiber.New()
	app.Get("/students", studentService.GetAll)

	// Seed data
	mockRepo.Create(nil, &models.Students{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		StudentID:    "123456",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	})

	req := httptest.NewRequest("GET", "/students", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Errorf("Error request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Status salah! Dapat %d, Harapan 200", resp.StatusCode)
	}
}

// TestGetStudentByID - Test untuk GET /students/:id
func TestGetStudentByID_TableDriven(t *testing.T) {
	mockRepo := mocks.NewManualMockStudentRepo()
	studentService := service.NewStudentService(mockRepo)
	app := fiber.New()
	app.Get("/students/:id", studentService.GetByID)

	existingID := uuid.New()
	mockRepo.Create(nil, &models.Students{
		ID:           existingID,
		UserID:       uuid.New(),
		StudentID:    "123456",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	})

	tests := []struct {
		name           string
		studentID      string
		expectedStatus int
	}{
		{
			name:           "Student Ditemukan",
			studentID:      existingID.String(),
			expectedStatus: 200,
		},
		{
			name:           "Student Tidak Ditemukan",
			studentID:      uuid.New().String(),
			expectedStatus: 404,
		},
		{
			name:           "ID Tidak Valid",
			studentID:      "invalid-uuid",
			expectedStatus: 404, // Service returns 404 for invalid UUID (not 400)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/students/"+tt.studentID, nil)
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

// TestUpdateStudent - Test untuk PUT /students/:id
func TestUpdateStudent_TableDriven(t *testing.T) {
	mockRepo := mocks.NewManualMockStudentRepo()
	studentService := service.NewStudentService(mockRepo)
	app := fiber.New()
	app.Put("/students/:id", studentService.Update)

	existingID := uuid.New()
	mockRepo.Create(nil, &models.Students{
		ID:           existingID,
		UserID:       uuid.New(),
		StudentID:    "123456",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	})

	tests := []struct {
		name           string
		studentID      string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name:      "Update Sukses",
			studentID: existingID.String(),
			inputBody: map[string]string{
				"program_study": "Sistem Informasi",
			},
			expectedStatus: 200,
		},
		{
			name:      "Student Tidak Ditemukan",
			studentID: uuid.New().String(),
			inputBody: map[string]string{
				"program_study": "Sistem Informasi",
			},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("PUT", "/students/"+tt.studentID, bytes.NewReader(bodyBytes))
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

// TestDeleteStudent - Test untuk DELETE /students/:id
func TestDeleteStudent_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(*mocks.ManualMockStudentRepo) string
		expectedStatus int
	}{
		{
			name: "Delete Sukses",
			setupFunc: func(repo *mocks.ManualMockStudentRepo) string {
				id := uuid.New()
				repo.Create(nil, &models.Students{
					ID:           id,
					UserID:       uuid.New(),
					StudentID:    "123456",
					ProgramStudy: "Informatika",
					AcademicYear: "2024",
				})
				return id.String()
			},
			expectedStatus: 200,
		},
		{
			name: "Student Tidak Ditemukan",
			setupFunc: func(repo *mocks.ManualMockStudentRepo) string {
				return uuid.New().String()
			},
			expectedStatus: 500, // Service return 500 saat delete gagal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewManualMockStudentRepo()
			studentService := service.NewStudentService(mockRepo)
			app := fiber.New()
			app.Delete("/students/:id", studentService.Delete)

			studentID := tt.setupFunc(mockRepo)
			req := httptest.NewRequest("DELETE", "/students/"+studentID, nil)

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

// TestAssignAdvisor - Test untuk POST /students/:id/advisor
func TestAssignAdvisor_TableDriven(t *testing.T) {
	mockRepo := mocks.NewManualMockStudentRepo()
	studentService := service.NewStudentService(mockRepo)
	app := fiber.New()
	app.Post("/students/:id/advisor", studentService.AssignAdvisor)

	existingID := uuid.New()
	mockRepo.Create(nil, &models.Students{
		ID:           existingID,
		UserID:       uuid.New(),
		StudentID:    "123456",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	})

	tests := []struct {
		name           string
		studentID      string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name:      "Assign Advisor Sukses",
			studentID: existingID.String(),
			inputBody: map[string]string{
				"advisor_id": uuid.New().String(),
			},
			expectedStatus: 200,
		},
		{
			name:      "Student Tidak Ditemukan",
			studentID: uuid.New().String(),
			inputBody: map[string]string{
				"advisor_id": uuid.New().String(),
			},
			expectedStatus: 500, // Service return 500 saat assign gagal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/students/"+tt.studentID+"/advisor", bytes.NewReader(bodyBytes))
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
