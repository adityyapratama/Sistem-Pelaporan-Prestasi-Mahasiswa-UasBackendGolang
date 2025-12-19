package repository

import (
	"context"
	"database/sql"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

type ReportRepo struct {
	pgDB *sql.DB
}

func NewReportRepository(pgDB *sql.DB) *ReportRepo {
	return &ReportRepo{pgDB: pgDB}
}

func (r *ReportRepo) GetStatistics(ctx context.Context) (*models.StatisticsResponse, error) {
	var stats models.StatisticsResponse
	err := r.pgDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM students`).Scan(&stats.TotalStudents)
	if err != nil {
		return nil, err
	}
	err = r.pgDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM lectures`).Scan(&stats.TotalLectures)
	if err != nil {
		return nil, err
	}



	err = r.pgDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM achievement_references WHERE deleted_at IS NULL`).Scan(&stats.TotalAchievements)
	if err != nil {
		return nil, err
	}

	
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN LOWER(status) = 'draft' THEN 1 ELSE 0 END), 0) as draft,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'submitted' THEN 1 ELSE 0 END), 0) as submitted,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'verified' THEN 1 ELSE 0 END), 0) as verified,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'rejected' THEN 1 ELSE 0 END), 0) as rejected
		FROM achievement_references 
		WHERE deleted_at IS NULL
	`
	err = r.pgDB.QueryRowContext(ctx, query).Scan(
		&stats.AchievementStats.Draft,
		&stats.AchievementStats.Submitted,
		&stats.AchievementStats.Verified,
		&stats.AchievementStats.Rejected,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *ReportRepo) GetAchievementStatsByStudentID(ctx context.Context, studentID uuid.UUID) (*models.AchievementStatistics, int, error) {
	var stats models.AchievementStatistics
	var total int

	query := `
		SELECT 
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'draft' THEN 1 ELSE 0 END), 0) as draft,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'submitted' THEN 1 ELSE 0 END), 0) as submitted,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'verified' THEN 1 ELSE 0 END), 0) as verified,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'rejected' THEN 1 ELSE 0 END), 0) as rejected
		FROM achievement_references 
		WHERE student_id = $1 AND deleted_at IS NULL
	`
	err := r.pgDB.QueryRowContext(ctx, query, studentID).Scan(
		&total,
		&stats.Draft,
		&stats.Submitted,
		&stats.Verified,
		&stats.Rejected,
	)
	if err != nil {
		return nil, 0, err
	}

	return &stats, total, nil
}
