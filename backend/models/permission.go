package models

import (
	"time"
)

// Permission representa un permiso del sistema
type Permission struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"` // Ej: "products", "orders", "users"
	Action      string    `json:"action"`   // Ej: "create", "read", "update", "delete"
}

// Role representa un rol del sistema
type Role struct {
	ID          int64        `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	IsSystem    bool         `json:"is_system"` // Roles del sistema no se pueden eliminar
}
