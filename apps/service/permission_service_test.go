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

// TestCreatePermission - Test untuk POST /permissions
func TestCreatePermission(t *testing.T) {
	mockRepo := mocks.NewManualMockPermissionRepo()
	permissionService := service.NewPermissionService(mockRepo)
	app := fiber.New()
	app.Post("/permissions", permissionService.Create)

	tests := []struct {
		name           string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name: "Create Permission Sukses",
			inputBody: map[string]string{
				"name":        "manage_users",
				"resource":    "users",
				"action":      "create",
				"description": "Can manage users",
			},
			expectedStatus: 201,
		},
		{
			name: "Nama Kosong",
			inputBody: map[string]string{
				"name":        "",
				"resource":    "users",
				"action":      "create",
				"description": "Description",
			},
			expectedStatus: 400,
		},
		{
			name: "Resource Kosong",
			inputBody: map[string]string{
				"name":        "manage_users",
				"resource":    "",
				"action":      "create",
				"description": "Description",
			},
			expectedStatus: 400,
		},
		{
			name: "Action Kosong",
			inputBody: map[string]string{
				"name":        "manage_users",
				"resource":    "users",
				"action":      "",
				"description": "Description",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/permissions", bytes.NewReader(bodyBytes))
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

// TestGetAllPermissions - Test untuk GET /permissions
func TestGetAllPermissions(t *testing.T) {
	mockRepo := mocks.NewManualMockPermissionRepo()
	permissionService := service.NewPermissionService(mockRepo)
	app := fiber.New()
	app.Get("/permissions", permissionService.GetAll)

	// Seed data
	mockRepo.Create(nil, &models.Permission{
		ID:          uuid.New(),
		Name:        "manage_users",
		Resource:    "users",
		Action:      "create",
		Description: "Can manage users",
	})

	// Test
	req := httptest.NewRequest("GET", "/permissions", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Errorf("Error request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Status salah! Dapat %d, Harapan 200", resp.StatusCode)
	}
}

// TestAssignPermissionToRole - Test untuk POST /permissions/assign
func TestAssignPermissionToRole(t *testing.T) {
	mockRepo := mocks.NewManualMockPermissionRepo()
	permissionService := service.NewPermissionService(mockRepo)
	app := fiber.New()
	app.Post("/permissions/assign", permissionService.AssignToRole)

	// Seed permission data
	permissionID := uuid.New()
	mockRepo.Create(nil, &models.Permission{
		ID:          permissionID,
		Name:        "manage_users",
		Resource:    "users",
		Action:      "create",
		Description: "Can manage users",
	})

	tests := []struct {
		name           string
		inputBody      map[string]string
		expectedStatus int
	}{
		{
			name: "Assign Permission Sukses",
			inputBody: map[string]string{
				"role_id":       uuid.New().String(),
				"permission_id": permissionID.String(),
			},
			expectedStatus: 200,
		},
		{
			name: "Role ID Format Salah",
			inputBody: map[string]string{
				"role_id":       "invalid-uuid",
				"permission_id": permissionID.String(),
			},
			expectedStatus: 400,
		},
		{
			name: "Permission ID Format Salah",
			inputBody: map[string]string{
				"role_id":       uuid.New().String(),
				"permission_id": "invalid-uuid",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest("POST", "/permissions/assign", bytes.NewReader(bodyBytes))
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
