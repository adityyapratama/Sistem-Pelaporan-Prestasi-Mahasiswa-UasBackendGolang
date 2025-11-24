package mocks

import (
	"context"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models" // Pastikan path ini sesuai go.mod kamu

	"github.com/google/uuid"
)


type ManualMockUserRepo struct {
	users map[string]*models.User
}


func NewManualMockUserRepo() *ManualMockUserRepo {
	return &ManualMockUserRepo{
		users: make(map[string]*models.User),
	}
}


func (m *ManualMockUserRepo) Create(ctx context.Context, user *models.User) error {
	if user.Username == "" {
		return errors.New("username empty")
	}
	
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	
	
	m.users[user.Username] = user
	m.users[user.Email] = user
	return nil
}


func (m *ManualMockUserRepo) GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error) {
	if user, exists := m.users[login]; exists {
		return user, nil
	}
	return nil, errors.New("user not found") // Error mirip sql.ErrNoRows
}


func (m *ManualMockUserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	
	uniqueUsers := make(map[uuid.UUID]models.User)
	
	for _, u := range m.users {
		uniqueUsers[u.ID] = *u
	}

	var result []models.User
	for _, u := range uniqueUsers {
		result = append(result, u)
	}
	return result, nil
}


func (m *ManualMockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}


func (m *ManualMockUserRepo) Update(ctx context.Context, user *models.User) error {

	var existingUser *models.User
	for _, u := range m.users {
		if u.ID == user.ID {
			existingUser = u
			break
		}
	}

	if existingUser == nil {
		return errors.New("user not found")
	}


	existingUser.FullName = user.FullName
	existingUser.Username = user.Username
	existingUser.Email = user.Email



	
	return nil
}


func (m *ManualMockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	var keysToDelete []string


	for key, u := range m.users {
		if u.ID == id {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if len(keysToDelete) == 0 {
		return errors.New("user not found")
	}


	for _, key := range keysToDelete {
		delete(m.users, key)
	}

	return nil
}


func (m *ManualMockUserRepo) UpdateRole(ctx context.Context, uid uuid.UUID, rid uuid.UUID) error {
	for _, u := range m.users {
		if u.ID == uid {
			u.RoleID = rid
			return nil
		}
	}
	return errors.New("user not found")
}