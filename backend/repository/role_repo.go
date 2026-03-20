package repository

import (
	"database/sql"
	"smartech/backend/models"
)

// RoleRepository maneja toda la lógica de acceso a datos de roles
type RoleRepository struct {
	db *sql.DB
}

// NewRoleRepository crea una nueva instancia de RoleRepository
func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetAll obtiene todos los roles
func (r *RoleRepository) GetAll() ([]models.Role, error) {
	query := `
		SELECT id, created_at, updated_at, name
		FROM roles
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, DBError(err, "GetAll roles")
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// GetByID obtiene un rol por su ID
func (r *RoleRepository) GetByID(id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name
		FROM roles
		WHERE id = ?
	`, id).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "GetByID role")
	}

	return &role, nil
}

// GetByName obtiene un rol por su nombre
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name
		FROM roles
		WHERE name = ?
	`, name).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "GetByName role")
	}

	return &role, nil
}

// Create crea un nuevo rol
func (r *RoleRepository) Create(role *models.Role) error {
	result, err := r.db.Exec(`
		INSERT INTO roles (name)
		VALUES (?)
	`, role.Name)

	if err != nil {
		return DBError(err, "Create role")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return DBError(err, "Get last insert id")
	}

	role.ID = id
	return nil
}

// Update actualiza un rol existente
func (r *RoleRepository) Update(id int64, role *models.Role) error {
	result, err := r.db.Exec(`
		UPDATE roles
		SET name = ?
		WHERE id = ?
	`, role.Name, id)

	if err != nil {
		return DBError(err, "Update role")
	}

	return CheckRowsAffected(result, "Update role")
}

// Delete elimina un rol
func (r *RoleRepository) Delete(id int64) error {
	result, err := r.db.Exec(`
		DELETE FROM roles WHERE id = ?
	`, id)

	if err != nil {
		return DBError(err, "Delete role")
	}

	return CheckRowsAffected(result, "Delete role")
}

// GetPermissions obtiene los permisos de un rol
func (r *RoleRepository) GetPermissions(roleID int64) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.created_at, p.updated_at, p.nombre, p.descripcion
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.nombre
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, DBError(err, "GetPermissions")
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.ID, &permission.CreatedAt, &permission.UpdatedAt, &permission.Name, &permission.Description)
		if err != nil {
			continue
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// AssignPermission asigna un permiso a un rol
func (r *RoleRepository) AssignPermission(roleID, permissionID int64) error {
	_, err := r.db.Exec(`
		INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
		VALUES (?, ?)
	`, roleID, permissionID)

	if err != nil {
		return DBError(err, "Assign permission to role")
	}

	return nil
}

// RemovePermission remueve un permiso de un rol
func (r *RoleRepository) RemovePermission(roleID, permissionID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM role_permissions
		WHERE role_id = ? AND permission_id = ?
	`, roleID, permissionID)

	if err != nil {
		return DBError(err, "Remove permission from role")
	}

	return nil
}

// Count obtiene el total de roles
func (r *RoleRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM roles`).Scan(&count)
	if err != nil {
		return 0, DBError(err, "Count roles")
	}
	return count, nil
}
