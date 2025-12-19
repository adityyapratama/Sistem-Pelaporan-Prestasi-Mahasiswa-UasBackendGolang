package mocks

import (
	"context"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

// ManualMockPermissionRepo - Mock untuk PermissionRepository
type ManualMockPermissionRepo struct {
	permissions     map[string]*models.Permission
	rolePermissions map[uuid.UUID][]uuid.UUID // roleID -> []permissionID
}

// NewManualMockPermissionRepo - Constructor
func NewManualMockPermissionRepo() *ManualMockPermissionRepo {
	return &ManualMockPermissionRepo{
		permissions:     make(map[string]*models.Permission),
		rolePermissions: make(map[uuid.UUID][]uuid.UUID),
	}
}


func (m *ManualMockPermissionRepo) Create(ctx context.Context, perm *models.Permission) error {
	if perm.Name == "" {
		return errors.New("permission name empty")
	}

	if perm.ID == uuid.Nil {
		perm.ID = uuid.New()
	}

	// Cek duplikat
	if _, exists := m.permissions[perm.Name]; exists {
		return errors.New("permission already exists")
	}

	m.permissions[perm.Name] = perm
	return nil
}


func (m *ManualMockPermissionRepo) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	if perm, exists := m.permissions[name]; exists {
		return perm, nil
	}
	return nil, errors.New("permission not found")
}


func (m *ManualMockPermissionRepo) GetAll(ctx context.Context) ([]models.Permission, error) {
	var result []models.Permission
	for _, p := range m.permissions {
		result = append(result, *p)
	}
	return result, nil
}


func (m *ManualMockPermissionRepo) AssignToRole(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error {
	// Cek apakah sudah di-assign
	for _, pid := range m.rolePermissions[roleID] {
		if pid == permID {
			return errors.New("permission already assigned to role")
		}
	}

	m.rolePermissions[roleID] = append(m.rolePermissions[roleID], permID)
	return nil
}


func (m *ManualMockPermissionRepo) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	var result []models.Permission
	permIDs := m.rolePermissions[roleID]

	for _, pid := range permIDs {
		for _, p := range m.permissions {
			if p.ID == pid {
				result = append(result, *p)
			}
		}
	}
	return result, nil
}
