package mocks

import (
	"context"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

// ManualMockStudentRepo - Mock untuk StudentsRepository
type ManualMockStudentRepo struct {
	students map[uuid.UUID]*models.Students
	byUserID map[uuid.UUID]*models.Students
}

// NewManualMockStudentRepo - Constructor
func NewManualMockStudentRepo() *ManualMockStudentRepo {
	return &ManualMockStudentRepo{
		students: make(map[uuid.UUID]*models.Students),
		byUserID: make(map[uuid.UUID]*models.Students),
	}
}


func (m *ManualMockStudentRepo) Create(ctx context.Context, student *models.Students) error {
	if student.StudentID == "" {
		return errors.New("student_id (NIM) empty")
	}

	// Cek duplikat user_id
	if _, exists := m.byUserID[student.UserID]; exists {
		return errors.New("student profile already exists for this user")
	}

	// Cek duplikat NIM
	for _, s := range m.students {
		if s.StudentID == student.StudentID {
			return errors.New("NIM already registered")
		}
	}

	if student.ID == uuid.Nil {
		student.ID = uuid.New()
	}

	m.students[student.ID] = student
	m.byUserID[student.UserID] = student
	return nil
}


func (m *ManualMockStudentRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Students, error) {
	if student, exists := m.byUserID[userID]; exists {
		return student, nil
	}
	return nil, nil 
}


func (m *ManualMockStudentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Students, error) {
	if student, exists := m.students[id]; exists {
		return student, nil
	}
	return nil, nil
}


func (m *ManualMockStudentRepo) GetAll(ctx context.Context) ([]models.Students, error) {
	var result []models.Students
	for _, s := range m.students {
		result = append(result, *s)
	}
	return result, nil
}


func (m *ManualMockStudentRepo) Update(ctx context.Context, student *models.Students) error {
	if _, exists := m.students[student.ID]; !exists {
		return errors.New("student not found")
	}

	m.students[student.ID] = student
	m.byUserID[student.UserID] = student
	return nil
}


func (m *ManualMockStudentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	student, exists := m.students[id]
	if !exists {
		return errors.New("student not found")
	}

	delete(m.byUserID, student.UserID)
	delete(m.students, id)
	return nil
}


func (m *ManualMockStudentRepo) AssignAdvisor(ctx context.Context, studentID uuid.UUID, advisorID uuid.UUID) error {
	student, exists := m.students[studentID]
	if !exists {
		return errors.New("student not found")
	}

	student.AdvisorID = &advisorID
	return nil
}


func (m *ManualMockStudentRepo) GetByAdvisorID(ctx context.Context, advisorID uuid.UUID) ([]models.Students, error) {
	var result []models.Students
	for _, s := range m.students {
		if s.AdvisorID != nil && *s.AdvisorID == advisorID {
			result = append(result, *s)
		}
	}
	return result, nil
}
