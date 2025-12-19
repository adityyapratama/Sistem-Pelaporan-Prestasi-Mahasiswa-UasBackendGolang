package mocks

import (
	"context"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

type ManualMockLectureRepo struct {
	lectures map[uuid.UUID]*models.Lecture
	byUserID map[uuid.UUID]*models.Lecture
}

// NewManualMockLectureRepo - Constructor
func NewManualMockLectureRepo() *ManualMockLectureRepo {
	return &ManualMockLectureRepo{
		lectures: make(map[uuid.UUID]*models.Lecture),
		byUserID: make(map[uuid.UUID]*models.Lecture),
	}
}


func (m *ManualMockLectureRepo) Create(ctx context.Context, lecture *models.Lecture) error {
	if lecture.LecturerID == "" {
		return errors.New("lecturer_id (NID) empty")
	}

	// Cek duplikat user_id
	if _, exists := m.byUserID[lecture.UserID]; exists {
		return errors.New("lecture profile already exists for this user")
	}

	// Cek duplikat NID
	for _, l := range m.lectures {
		if l.LecturerID == lecture.LecturerID {
			return errors.New("NID already registered")
		}
	}

	if lecture.ID == uuid.Nil {
		lecture.ID = uuid.New()
	}

	m.lectures[lecture.ID] = lecture
	m.byUserID[lecture.UserID] = lecture
	return nil
}


func (m *ManualMockLectureRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecture, error) {
	if lecture, exists := m.byUserID[userID]; exists {
		return lecture, nil
	}
	return nil, nil
}


func (m *ManualMockLectureRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Lecture, error) {
	if lecture, exists := m.lectures[id]; exists {
		return lecture, nil
	}
	return nil, nil
}


func (m *ManualMockLectureRepo) GetAll(ctx context.Context) ([]models.Lecture, error) {
	var result []models.Lecture
	for _, l := range m.lectures {
		result = append(result, *l)
	}
	return result, nil
}


func (m *ManualMockLectureRepo) Update(ctx context.Context, lecture *models.Lecture) error {
	if _, exists := m.lectures[lecture.ID]; !exists {
		return errors.New("lecture not found")
	}

	m.lectures[lecture.ID] = lecture
	m.byUserID[lecture.UserID] = lecture
	return nil
}


func (m *ManualMockLectureRepo) Delete(ctx context.Context, id uuid.UUID) error {
	lecture, exists := m.lectures[id]
	if !exists {
		return errors.New("lecture not found")
	}

	delete(m.byUserID, lecture.UserID)
	delete(m.lectures, id)
	return nil
}
