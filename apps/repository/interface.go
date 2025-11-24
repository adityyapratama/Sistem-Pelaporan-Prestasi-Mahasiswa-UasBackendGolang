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
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
}

type RoleRepository interface{
	GetAll(ctx context.Context )([]models.User, error)
	GetByID (ctx context.Context, id uuid.UUID) (*models.Role, error)
}




