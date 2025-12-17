package models

import "github.com/google/uuid"

type Permission struct{
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`         
	Resource    string    `json:"resource" db:"resource"` 
	Action      string    `json:"action" db:"action"`    
	Description string    `json:"description" db:"description"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type AssignPermissionRequest struct {
	RoleID string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}


