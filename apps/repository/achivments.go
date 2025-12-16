package repository


import (
	"context"
	"database/sql"
	"errors"
	"time"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepo struct {
	pgDB    *sql.DB
	mongoDB *mongo.Database
}

func NewHAchievementRepo(pgDB *sql.DB, mongoDB *mongo.Database) *AchievementRepo {
	return &AchievementRepo{
		pgDB:    pgDB,
		mongoDB: mongoDB,
	}
}

func (r *AchievementRepo) CreateDetail(ctx context.Context, detail *models.AchievementDetail) error {
	collection := r.mongoDB.Collection("achievements")
	detail.ID = primitive.NewObjectID()
	detail.CreatedAt = time.Now()
	detail.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, detail)
	return err
}

func (r *AchievementRepo) GetDetailByID(ctx context.Context, mongoID string) (*models.AchievementDetail, error) {
	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	var detail models.AchievementDetail
	err = r.mongoDB.Collection("achievements").FindOne(ctx, bson.M{"_id": objID}).Decode(&detail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

func (r *AchievementRepo) CreateReference(ctx context.Context, ref *models.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	if ref.Status == "" {
		ref.Status = "draft"
	}
	err := r.pgDB.QueryRowContext(ctx, query, ref.StudentID, ref.MongoAchievementID, ref.Status).Scan(&ref.ID, &ref.CreatedAt, &ref.UpdatedAt)
	return err
}

func (r *AchievementRepo) GetReferenceByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at FROM achievement_references WHERE id = $1`
	var model models.AchievementReference
	err := r.pgDB.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.StudentID,
		&model.MongoAchievementID,
		&model.Status,
		&model.SubmittedAt,
		&model.VerifiedAt,
		&model.VerifiedBy,
		&model.RejectionNote,
		&model.CreatedAt,
		&model.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *AchievementRepo) GetAllByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at FROM achievement_references WHERE student_id = $1 ORDER BY created_at DESC`
	rows, err := r.pgDB.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var refs []models.AchievementReference
	for rows.Next() {
		var r models.AchievementReference
		if err := rows.Scan(&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		refs = append(refs, r)
	}
	return refs, nil
}

func (r *AchievementRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE achievement_references SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.pgDB.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no achievement status updated")
	}
	return nil
}
