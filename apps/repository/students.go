package repository

import (
	"context"
	"database/sql"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

type PostStudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *PostStudentRepository {
	return &PostStudentRepository{db: db}
}

func (r *PostStudentRepository) Create(ctx context.Context, student *models.Students) error {
	query := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, 
		student.UserID, 
		student.StudentID, 
		student.ProgramStudy, 
		student.AcademicYear,
	).Scan(&student.ID, &student.CreatedAt)
	
	return err
}



func (r *PostStudentRepository) GetByUserID(ctx context.Context, UserID uuid.UUID) (*models.Students, error) {
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
	          FROM students WHERE user_id = $1`
	
	var s models.Students
	err := r.db.QueryRowContext(ctx, query, UserID).Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}


func (r *PostStudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Students, error) {
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id = $1`
	var s models.Students
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
	) ; err !=nil{
		if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    
	return nil,err}
	}
	
	 return &s, nil
}


func (r *PostStudentRepository) GetAll(ctx context.Context) ([]models.Students, error) {
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Students
	for rows.Next() {
		var s models.Students
		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}


func (r *PostStudentRepository) Update(ctx context.Context, s *models.Students) error {
	query := `
		UPDATE students 
		SET student_id = $1, program_study = $2, academic_year = $3, advisor_id =$4 
		WHERE id = $5
	`
	result, err := r.db.ExecContext(ctx, query, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID,s.ID)
	if err != nil {
		return err
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("student not found")
	}
	return nil
}


func (r *PostStudentRepository) AssignAdvisor(ctx context.Context, studentID uuid.UUID, advisorID uuid.UUID)error{
	query :=`UPDATE students SET advisor_id =$1 WHERE id =$2 `
	result, err := r.db.ExecContext(ctx,query,advisorID,studentID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("student not found")
	}
	return nil

}

func (r *PostStudentRepository) GetByAdvisorID(ctx context.Context, advisorID uuid.UUID) ([]models.Students, error){
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
              FROM students WHERE advisor_id = $1`
	 rows, err := r.db.QueryContext(ctx, query, advisorID)
    if err != nil {
        return nil, err
    }
	defer rows.Close()
	var student []models.Students
	for rows.Next(){
		var s models.Students
		if err:=rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt) ; err!=nil{
			return  nil, err
		}
		student= append(student, s)
	}
	return  student,nil
}