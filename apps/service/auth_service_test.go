package service_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository/mocks"
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func TestLogin_TableDriven(t *testing.T) {
	
	// 1. ARRANGE (Siapkan Mock Manual & Data Awal)
	mockRepo := mocks.NewManualMockUserRepo() // Pakai Mock Manual buatan sendiri
	authService := service.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	// Kita isi "Database Palsu" kita dengan 1 user
	passHash, _ := utils.HashPassword("rahasia123")
	existingUser := &models.User{
		ID:           uuid.New(),
		Username:     "aditya",
		PasswordHash: passHash,
		IsActive:     true,
		Role:         &models.Role{Name: "Mahasiswa"},
	}
	// Simpan user ke mock manual (tanpa lewat service)
	mockRepo.Create(nil, existingUser)

	// 2. Definisi Tabel Test (Table-Driven) [cite: 2074-2085]
	tests := []struct {
		name           string // Nama skenario
		inputUsername  string
		inputPassword  string
		expectedStatus int    // Harapan status code (200, 401, dll)
		wantErr        bool   // Apakah expect error?
	}{
		{
			name:           "Login Sukses",
			inputUsername:  "aditya",
			inputPassword:  "rahasia123",
			expectedStatus: 200,
			wantErr:        false,
		},
		{
			name:           "Password Salah",
			inputUsername:  "aditya",
			inputPassword:  "salahbanget",
			expectedStatus: 401,
			wantErr:        true,
		},
		{
			name:           "User Tidak Ditemukan",
			inputUsername:  "hantu",
			inputPassword:  "rahasia123",
			expectedStatus: 401,
			wantErr:        true,
		},
	}

	// 3. Loop Eksekusi Test [cite: 2136-2140]
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Siapkan Request Body
			reqBody := map[string]string{
				"username": tt.inputUsername,
				"password": tt.inputPassword,
			}
			bodyBytes, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// ACT (Jalankan)
			resp, err := app.Test(req)

			// ASSERT (Cek Hasil)
			// [cite: 2141] Cek error sistem
			if err != nil {
				t.Errorf("Error sistem tidak diharapkan: %v", err)
			}

			// Cek Status Code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status code salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}