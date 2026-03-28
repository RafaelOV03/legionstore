package repositories

import (
	"database/sql"
	"smartech/backend/models"
	"strings"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) ListRoles() ([]models.Role, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		var role models.Role
		var isSystem int
		if err := rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem); err != nil {
			continue
		}
		role.IsSystem = isSystem == 1
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) GetRoleByID(roleID int64) (models.Role, error) {
	var role models.Role
	var isSystem int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, roleID).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)
	if err != nil {
		return models.Role{}, err
	}

	role.IsSystem = isSystem == 1
	return role, nil
}

func (r *RoleRepository) GetRoleByName(name string) (models.Role, error) {
	var role models.Role
	var isSystem int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE name = ?
	`, name).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)
	if err != nil {
		return models.Role{}, err
	}

	role.IsSystem = isSystem == 1
	return role, nil
}

func (r *RoleRepository) ListPermissionsByRole(roleID int64) ([]models.Permission, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.name, p.description, p.resource, p.action
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.name
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]models.Permission, 0)
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action); err != nil {
			continue
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *RoleRepository) ListPermissionsByRoles() (map[int64][]models.Permission, error) {
	rows, err := r.db.Query(`
		SELECT rp.role_id, p.id, p.created_at, p.updated_at, p.name, p.description, p.resource, p.action
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		ORDER BY rp.role_id, p.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permMap := make(map[int64][]models.Permission)
	for rows.Next() {
		var roleID int64
		var perm models.Permission
		if err := rows.Scan(&roleID, &perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action); err != nil {
			continue
		}
		permMap[roleID] = append(permMap[roleID], perm)
	}

	return permMap, nil
}

func (r *RoleRepository) CountByName(name string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM roles WHERE name = ?", name).Scan(&count)
	return count, err
}

func (r *RoleRepository) InsertRole(name, description string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO roles (name, description, is_system, created_at, updated_at)
		VALUES (?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, name, description)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *RoleRepository) UpdateRoleFields(roleID int64, name, description *string) error {
	updates := make([]string, 0)
	args := make([]interface{}, 0)

	if name != nil && *name != "" {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *description)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, roleID)
	query := "UPDATE roles SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *RoleRepository) ReplaceRolePermissions(roleID int64, permissionIDs []int64) error {
	if _, err := r.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID); err != nil {
		return err
	}

	for _, permID := range permissionIDs {
		if _, err := r.db.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES (?, ?)
		`, roleID, permID); err != nil {
			continue
		}
	}

	return nil
}

func (r *RoleRepository) CountUsersByRole(roleID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE role_id = ?", roleID).Scan(&count)
	return count, err
}

func (r *RoleRepository) DeleteRole(roleID int64) error {
	_, err := r.db.Exec("DELETE FROM roles WHERE id = ?", roleID)
	return err
}

func (r *RoleRepository) ListPermissions() ([]models.Permission, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, name, description, resource, action
		FROM permissions
		ORDER BY resource, action
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]models.Permission, 0)
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action); err != nil {
			continue
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}
