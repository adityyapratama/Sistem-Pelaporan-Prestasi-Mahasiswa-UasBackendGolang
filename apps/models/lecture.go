package models

import (
	"time"

	"github.com/google/uuid"
)


type Lecture struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	LecturerID string    `json:"lecturer_id" db:"lecturer_id"` // NIP / NID
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type CreateLecturerRequest struct {
	UserID     string `json:"user_id"`
	LecturerID string `json:"lecturer_id"` 
	Department string `json:"department"`
}

