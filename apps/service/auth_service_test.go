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
	
	
	mockRepo := mocks.NewManualMockUserRepo() 
	authService := service.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	
	passHash, _ := utils.HashPassword("rahasia123")
	existingUser := &models.User{
		ID:           uuid.New(),
		Username:     "aditya",
		PasswordHash: passHash,
		IsActive:     true,
		Role:         &models.Role{Name: "Mahasiswa"},
	}
	
	mockRepo.Create(nil, existingUser)

	
	tests := []struct {
		name           string 
		inputUsername  string
		inputPassword  string
		expectedStatus int    
		wantErr        bool   
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

	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
	
			reqBody := map[string]string{
				"username": tt.inputUsername,
				"password": tt.inputPassword,
			}
			bodyBytes, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

	
			resp, err := app.Test(req)

	
	
			if err != nil {
				t.Errorf("Error sistem tidak diharapkan: %v", err)
			}

	
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status code salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}


func TestRegister_TableDriven(t *testing.T) {
	
	mockRepo := mocks.NewManualMockUserRepo()
	authService := service.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/register", authService.Register)

	
	existingUser := &models.User{
		ID:       uuid.New(),
		Username: "sudahada",
		Email:    "ada@gmail.com",
	}
	mockRepo.Create(nil, existingUser)

	
	tests := []struct {
		name           string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name: "Register Sukses",
			inputBody: map[string]string{
				"username":  "maba_baru",
				"email":     "maba@univ.ac.id",
				"password":  "rahasia123",
				"full_name": "Maba Univ",
				"role_id":   uuid.New().String(), 
			},
			expectedStatus: 201, 
		},
		{
			name: "Username Sudah Dipakai",
			inputBody: map[string]string{
				"username":  "sudahada", 
				"email":     "baru@gmail.com",
				"password":  "123456",
				"full_name": "Orang Lama",
				"role_id":   uuid.New().String(),
			},
			expectedStatus: 400, 
		},
		{
			name: "Input Tidak Lengkap",
			inputBody: map[string]string{
				"username": "kosong",
				
			},
			expectedStatus: 400, 
		},
	}

	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
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




// func TestRefreshToken_TableDriven(t *testing.T) {
// 	// 1. Setup
// 	mockRepo := mocks.NewManualMockUserRepo()
// 	authService := service.NewAuthService(mockRepo)
// 	app := fiber.New()
// 	app.Post("/refresh", authService.RefreshToken)

// 	// 2. Siapkan User Valid di DB Mock
// 	userID := uuid.New()
// 	validUser := &models.User{
// 		ID:       userID,
// 		Username: "user_refresh",
// 		IsActive: true,
// 		Role:     &models.Role{Name: "Mahasiswa"},
// 	}
// 	mockRepo.Create(nil, validUser)

// 	// 3. Generate Refresh Token Asli (Valid)
// 	validToken, _ := utils.GenerateRefreshToken(userID)

// 	// 4. Tabel Skenario
// 	tests := []struct {
// 		name           string
// 		tokenInput     string
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "Refresh Sukses",
// 			tokenInput:     validToken,
// 			expectedStatus: 200, // OK -> Dapat Access Token baru
// 		},
// 		{
// 			name:           "Token Tidak Valid/Rusak",
// 			tokenInput:     "token.ngawur.palsu",
// 			expectedStatus: 401, // Unauthorized
// 		},
// 		{
// 			name:           "Token Kosong",
// 			tokenInput:     "",
// 			expectedStatus: 400, // Bad Request / Unauthorized
// 		},
// 	}

// 	// 5. Eksekusi Loop
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Body request: { "refresh_token": "..." }
// 			reqBody := map[string]string{
// 				"refresh_token": tt.tokenInput,
// 			}
// 			bodyBytes, _ := json.Marshal(reqBody)
			
// 			req := httptest.NewRequest("POST", "/refresh", bytes.NewReader(bodyBytes))
// 			req.Header.Set("Content-Type", "application/json")

// 			resp, err := app.Test(req)

// 			if err != nil {
// 				t.Errorf("Error request: %v", err)
// 			}
// 			if resp.StatusCode != tt.expectedStatus {
// 				t.Errorf("Status salah! Dapat %d, Harapan %d", resp.StatusCode, tt.expectedStatus)
// 			}
// 		})
// 	}
// }