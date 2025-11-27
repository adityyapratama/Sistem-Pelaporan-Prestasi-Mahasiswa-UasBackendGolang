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

type RoleRepository interface{
	GetAll(ctx context.Context) ([]models.Role, error) 
    GetByID (ctx context.Context, id uuid.UUID) (*models.Role, error)
}



type PermissionRepository interface {
	Create(ctx context.Context, Permission *models.Permission) error
	GetByName(ctx context.Context , name string) (*models.Permission, error)
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

}

