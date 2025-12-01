package models

import (
	"time"

	"github.com/google/uuid"
)


type Students struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	StudentID    string     `json:"student_id" db:"student_id"`
	ProgramStudy string     `json:"program_study" db:"program_study"`
	AcademicYear string     `json:"academic_year" db:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id" db:"advisor_id"` 
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty" db:"updated_at"`
}

type CreateStudensRequest struct{
	StudentID        string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`	
}

type SetAdvisorRequest struct {
	AdvisorID string `json:"advisor_id"`
}