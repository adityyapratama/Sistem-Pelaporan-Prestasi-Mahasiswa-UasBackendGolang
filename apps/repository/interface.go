package repository

import (
	"context"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, User *models.User) error
	Update(ctx context.Context, User *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
}

type RoleRepository interface {
	GetAll(ctx context.Context) ([]models.Role, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
}

type PermissionRepository interface {
	Create(ctx context.Context, Permission *models.Permission) error
	GetByName(ctx context.Context, name string) (*models.Permission, error)
	GetAll(ctx context.Context) ([]models.Permission, error)
	AssignToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
}

type StudentsRepository interface {
	Create(ctx context.Context, student *models.Students) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Students, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Students, error)
	GetAll(ctx context.Context) ([]models.Students, error)
	Update(ctx context.Context, student *models.Students) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignAdvisor(ctx context.Context, studentID uuid.UUID, advisorID uuid.UUID) error
	GetByAdvisorID(ctx context.Context, advisorID uuid.UUID) ([]models.Students, error)
}

type LectureRepository interface {
	Create(ctx context.Context, Lecture *models.Lecture) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecture, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Lecture, error)
	GetAll(ctx context.Context) ([]models.Lecture, error)
	Update(ctx context.Context, lecture *models.Lecture) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AchievementRepository interface {
	CreateDetail(ctx context.Context, detail *models.AchievementDetail) error
	GetDetailByID(ctx context.Context, mongoID string) (*models.AchievementDetail, error)

	CreateReference(ctx context.Context, ref *models.AchievementReference) error
	GetReferenceByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error)

	UpdateStatus(ctx context.Context, ref *models.AchievementReference) error

	GetAll(ctx context.Context, status string) ([]models.AchievementReference, error)
	GetAllByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error)

	UpdateDetail(ctx context.Context, mongoID string, updateData *models.AchievementDetail) error

	GetAllDetailsFromMongo(ctx context.Context) ([]models.AchievementDetail, error)

	SoftDelete(ctx context.Context, id uuid.UUID) error
	Submit(ctx context.Context, id uuid.UUID) error
}
