package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementDetail struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"student_id" json:"student_id"` // Disimpan sebagai string UUID
	AchievementType string                 `bson:"achievement_type" json:"achievement_type"` // competition, organization, etc
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	
	// Field Dinamis: Isinya bisa apa saja (Juara, Tingkat, Tanggal, dll)
	Details         map[string]interface{} `bson:"details" json:"details"`
	
	Attachments     []string               `bson:"attachments" json:"attachments"` // Array URL File
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}


type AchievementReference struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	StudentID          uuid.UUID  `json:"student_id" db:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id" db:"mongo_achievement_id"` // ID dari Mongo
	Status             string     `json:"status" db:"status"` // draft, submitted, verified, rejected
	SubmittedAt        *time.Time `json:"submitted_at" db:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at" db:"verified_at"`
	VerifiedBy         *uuid.UUID `json:"verified_by" db:"verified_by"`
	RejectionNote      *string    `json:"rejection_note" db:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateAchievementRequest struct {
	Type        string                 `json:"type" validate:"required"` // competition, etc
	Title       string                 `json:"title" validate:"required"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"` // JSON Object bebas
	Tags        []string               `json:"tags"`

}