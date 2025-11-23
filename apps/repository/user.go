package repository

import (
	"context"
	"database/sql"
	"errors"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"github.com/google/uuid"
)


type PostUserRepository struct{
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *PostUserRepository {
	return &PostUserRepository{db:db}
}

func (r *PostUserRepository)GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error){
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.is_active, u.role_id,
		       r.id, r.name, r.description
			   FROM users u
			   LEFT JOIN roles r ON u.role_id = r.id
			   WHERE u.username = $1 OR u.email =$1
			   `
	var user models.User
	var role models.Role

	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.IsActive, &user.RoleID,
		&role.ID, &role.Name, &role.Description,
	)

	if err != nil{
		if err ==sql.ErrNoRows{
			return  nil, errors.New("user tidak di temukan")
		}
		return  nil,err
	}

	user.Role = &role
	return &user, err
}

func (r *PostUserRepository)GetAll(ctx context.Context) ([]models.User, error) {
query := `
		SELECT u.id, u.username, u.email, u.full_name, u.is_active, u.role_id,
		       r.id, r.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var role models.Role
		
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.IsActive, &u.RoleID, &role.ID, &role.Name); err != nil {
			return nil, err
		}
		u.Role = &role
		users = append(users, u)
	}
	return users, nil
}

func( r *PostUserRepository)GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
query := `
		SELECT u.id, u.username, u.email, u.full_name, u.is_active, u.role_id,
		       r.id, r.name, r.description
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	var user models.User
	var role models.Role

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName, &user.IsActive, &user.RoleID,
		&role.ID, &role.Name, &role.Description,
	)

	if err != nil {
		return nil, err
	}
	user.Role = &role
	return &user, nil
}

func( r *PostUserRepository) Create(ctx context.Context, User *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, is_active
	`
	// Kita kembalikan ID yang baru dibuat ke struct user
	err := r.db.QueryRowContext(ctx, query,
		User.Username, User.Email, User.PasswordHash, User.FullName, User.RoleID,
	).Scan(&User.ID, &User.CreatedAt, &User.UpdatedAt, &User.IsActive)

	return err
}


func( r *PostUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET full_name = $1, username = $2, email = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.FullName, user.Username, user.Email, user.ID)
	return err
}

func( r *PostUserRepository)Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users where id = $1`
	_,err := r.db.ExecContext(ctx, query, id)
	return	err
}

func (r *PostUserRepository)UpdateRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	query := "UPDATE users SET role_id = $1, updated_at = NOW() WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, roleID, userID)
	return err
}