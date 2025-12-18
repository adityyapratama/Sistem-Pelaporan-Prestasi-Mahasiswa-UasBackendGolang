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

func NewAchievementRepo(pgDB *sql.DB, mongoDB *mongo.Database) *AchievementRepo {
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
        INSERT INTO achievement_references (
            student_id, mongo_achievement_id, status, submitted_at, created_at
        )
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id, created_at, updated_at
    `

	if ref.Status == "" {
		ref.Status = "draft" 
	}

	err := r.pgDB.QueryRowContext(ctx, query, 
		ref.StudentID, 
		ref.MongoAchievementID, 
		ref.Status, 
		ref.SubmittedAt, 
	).Scan(&ref.ID, &ref.CreatedAt, &ref.UpdatedAt)
	
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

func (r *AchievementRepo) UpdateStatus(ctx context.Context, ref *models.AchievementReference) error {
	query := `
		UPDATE achievement_references 
		SET status = $1, verified_by = $2, verified_at = $3, rejection_note = $4, updated_at = NOW() 
		WHERE id = $5
	`
	
	result, err := r.pgDB.ExecContext(ctx, query, 
		ref.Status, 
		ref.VerifiedBy,    
		ref.VerifiedAt,    
		ref.RejectionNote, 
		ref.ID,
	)
	
	if err != nil {
		return err
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no achievement status updated (id not found)")
	}
	return nil
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






func (r *AchievementRepo) GetAll(ctx context.Context, status string) ([]models.AchievementReference, error) {
    var query string
    var rows *sql.Rows
    var err error

    
    if status != "" {
        query = `
            SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at 
            FROM achievement_references 
            WHERE status = $1 
            ORDER BY created_at DESC
        `
        rows, err = r.pgDB.QueryContext(ctx, query, status)
    } else {
        query = `
            SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at 
            FROM achievement_references 
            ORDER BY created_at DESC
        `
        rows, err = r.pgDB.QueryContext(ctx, query)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var refs []models.AchievementReference
    for rows.Next() {
        var r models.AchievementReference

        err := rows.Scan(
            &r.ID, 
            &r.StudentID, 
            &r.MongoAchievementID, 
            &r.Status, 
            &r.SubmittedAt,
            &r.VerifiedAt,
            &r.VerifiedBy,
            &r.RejectionNote,
            &r.CreatedAt, 
            &r.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        refs = append(refs, r)
    }
    return refs, nil
}




func (r *AchievementRepo) UpdateDetail(ctx context.Context, mongoID string, updateData *models.AchievementDetail) error {
    objID, _ := primitive.ObjectIDFromHex(mongoID)
    

    filter := bson.M{"_id": objID}
    update := bson.M{
        "$set": bson.M{
            "title":            updateData.Title,
            "description":      updateData.Description,
            "achievement_type": updateData.AchievementType,
            "details":          updateData.Details,
            "attachments":      updateData.Attachments,
            "tags":             updateData.Tags,
            "updated_at":       time.Now(),
        },
    }

    
    _, err := r.mongoDB.Collection("achievements").UpdateOne(ctx, filter, update)
    
    return err
}

func (r *AchievementRepo) GetAllDetailsFromMongo(ctx context.Context) ([]models.AchievementDetail, error) {
    cursor, err := r.mongoDB.Collection("achievements").Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var details []models.AchievementDetail
    if err = cursor.All(ctx, &details); err != nil {
        return nil, err
    }

    return details, nil
}