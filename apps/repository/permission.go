package repository

import (
	"context"
	"database/sql"
	"uas-pelaporan-prestasi-mahasiswa/apps/models"
	"github.com/google/uuid"
)


type PostgresPermissionRepository struct {
	db *sql.DB
}

func NewPostgresPermissionRepository (db *sql.DB) *PostgresPermissionRepository {
	return &PostgresPermissionRepository{db:db}
}



func (r *PostgresPermissionRepository)Create(ctx context.Context, p *models.Permission) error{
	query := `INSERT INTO permissions(name, resource, action, description)
				VALUES ($1,$2,$3,$4)
				RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.Name, p.Resource, p.Action, p.Description).Scan(&p.ID)
	return err
}

func (r *PostgresPermissionRepository) GetAll(ctx context.Context) ([]models.Permission, error) {
	query := `SELECT id, name, resource, action, description FROM permissions`
	
    rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return nil dan error
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
    
    // Return data dan nil (error kosong)
	return permissions, nil 
}

func ( r *PostgresPermissionRepository) GetByName(ctx context.Context , name string) (*models.Permission, error){
	query := `SELECT name, resource, action, description 
			FROM permissions
			where name = $1`
	var p models.Permission
	err := r.db.QueryRowContext(ctx, query, name).Scan(p.ID,p.Name, p.Resource, p.Action, p.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	return  &p, nil
}

func ( r *PostgresPermissionRepository)AssignToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `INSERT INTO role_permissions(role_id, permission_id)	
			  VALUES ($1,$2)
			  ON CONFLICT (role_id, permission_id) DO NOTHING
			  `
	_, err :=r.db.ExecContext(ctx,query,roleID,permissionID)
	return err
}

func (r *PostgresPermissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	query :=`SELECT p.id, p.name, p.resource, p.action, p.description
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1`

		rows, err := r.db.QueryContext(ctx, query, roleID)
		if err!= nil{
			return nil, err
		}

		defer rows.Close()

		var permissions []models.Permission
		for rows.Next() {
			var p models.Permission
			if err := rows.Scan(&p.ID,&p.Name,&p.Resource,&p.Description) ; err != nil{
				return nil, err
			}
			permissions = append(permissions, p)

		}
		return  permissions, nil

}


