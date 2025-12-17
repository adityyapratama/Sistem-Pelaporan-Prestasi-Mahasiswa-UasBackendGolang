package repository

import (
	"context"
	"database/sql"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)


type PostgresLectureRepository struct {
	db *sql.DB
}

func NewPostgresLectureRepository (db *sql.DB) *PostgresLectureRepository {
	return &PostgresLectureRepository{db:db}
}


func(r *PostgresLectureRepository)Create(ctx context.Context, l *models.Lecture)error{
	query := `
		INSERT INTO lecturers (user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, 
		l.UserID, 
		l.LecturerID, 
		l.Department,
	).Scan(&l.ID, &l.CreatedAt)
	
	return err
			
}


func (r *PostgresLectureRepository) GetAll(ctx context.Context) ([]models.Lecture, error) {
    query := `
        SELECT id, user_id, lecturer_id, department, created_at 
        FROM lecturers
    `
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var lecturers []models.Lecture 

    for rows.Next() {
        var l models.Lecture
        if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
            return nil, err
        }
        lecturers = append(lecturers, l)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return lecturers, nil
}



func (r *PostgresLectureRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.Lecture, error) {
	query := `SELECT id, user_id, lecturer_id, department, created_at 
	          FROM lecturers WHERE user_id = $1`
	var l models.Lecture
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt,
    )
	if err != nil {
		return nil, err
	}
	return &l, nil
}



func (r *PostgresLectureRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecture, error) {
	query := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE user_id = $1`
	var l models.Lecture
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}