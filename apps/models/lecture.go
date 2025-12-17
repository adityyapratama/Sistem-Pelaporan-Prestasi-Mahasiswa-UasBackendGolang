package models

import (
	"github.com/google/uuid"
	"time"
)

type Lecture struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	LecturerID string    `json:"lecturer_id" db:"lecturer_id"` 
	Department string    `json:"department" db:"department"`   
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	User *User `json:"user,omitempty" db:"-"`
}

type CreateLectureRequest struct {
	UserID     string `json:"user_id"`
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}