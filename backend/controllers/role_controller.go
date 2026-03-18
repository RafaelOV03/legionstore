package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRoles obtiene todos los roles con sus permisos
func GetRoles(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		ORDER BY name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	defer rows.Close()

	// First collect all roles
	var roles []models.Role
	var roleIDs []int64

	for rows.Next() {
		var role models.Role
		var isSystem int
		err := rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)
		if err != nil {
			continue
		}
		role.IsSystem = isSystem == 1
		roles = append(roles, role)
		roleIDs = append(roleIDs, role.ID)
	}

	// Load all permissions for all roles in batch
	permMap := loadAllRolePermissions(roleIDs)

	// Assign permissions to each role
	for i := range roles {
		if perms, exists := permMap[roles[i].ID]; exists {
			roles[i].Permissions = perms
		} else {
			roles[i].Permissions = []models.Permission{}
		}
	}

	c.JSON(http.StatusOK, roles)
}

// getRolePermissions obtiene todos los permisos de un rol
func getRolePermissions(roleID int64) []models.Permission {
	rows, err := database.DB.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.name, p.description, p.resource, p.action
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.name
	`, roleID)
	if err != nil {
		return []models.Permission{}
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(&perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action)
		if err != nil {
			continue
		}
		permissions = append(permissions, perm)
	}

	return permissions
}

// loadAllRolePermissions carga todos los permisos de múltiples roles en una sola query
func loadAllRolePermissions(roleIDs []int64) map[int64][]models.Permission {
	if len(roleIDs) == 0 {
		return make(map[int64][]models.Permission)
	}

	rows, err := database.DB.Query(`
		SELECT rp.role_id, p.id, p.created_at, p.updated_at, p.name, p.description, p.resource, p.action
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		ORDER BY rp.role_id, p.name
	`)
	if err != nil {
		return make(map[int64][]models.Permission)
	}
	defer rows.Close()

	permMap := make(map[int64][]models.Permission)
	for rows.Next() {
		var roleID int64
		var perm models.Permission
		err := rows.Scan(&roleID, &perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action)
		if err != nil {
			continue
		}
		permMap[roleID] = append(permMap[roleID], perm)
	}

	return permMap
}

// GetRole obtiene un rol por id con sus permisos
func GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	var role models.Role
	var isSystem int
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, id).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

	role.IsSystem = isSystem == 1
	role.Permissions = getRolePermissions(role.ID)

	c.JSON(http.StatusOK, role)
}

// CreateRole crea un nuevo rol
func CreateRole(c *gin.Context) {
	var req struct {
		Name          string  `json:"name" binding:"required"`
		Description   string  `json:"description"`
		PermissionIDs []int64 `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el rol ya existe
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM roles WHERE name = ?", req.Name).Scan(&count)
	if err == nil && count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists"})
		return
	}

	// Crear el rol
	result, err := database.DB.Exec(`
		INSERT INTO roles (name, description, is_system, created_at, updated_at)
		VALUES (?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	roleID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role id"})
		return
	}

	// Asignar permisos al rol
	if len(req.PermissionIDs) > 0 {
		for _, permID := range req.PermissionIDs {
			_, err := database.DB.Exec(`
				INSERT INTO role_permissions (role_id, permission_id)
				VALUES (?, ?)
			`, roleID, permID)
			if err != nil {
				// Log pero continuar
				continue
			}
		}
	}

	// Obtener el rol creado con sus permisos
	var role models.Role
	var isSystem int
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, roleID).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	role.IsSystem = isSystem == 1
	role.Permissions = getRolePermissions(roleID)

	c.JSON(http.StatusCreated, role)
}

// UpdateRole actualiza un rol existente
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	// Verificar si el rol existe y si es del sistema
	var isSystem int
	err = database.DB.QueryRow("SELECT is_system FROM roles WHERE id = ?", id).Scan(&isSystem)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

	if isSystem == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify system roles"})
		return
	}

	var req struct {
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		PermissionIDs *[]int64 `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construir query de actualización
	updates := []string{}
	args := []interface{}{}

	if req.Name != nil && *req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
	}

	if len(updates) > 0 {
		updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, id)

		query := "UPDATE roles SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		_, err := database.DB.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
			return
		}
	}

	// Actualizar permisos si se proporcionaron
	if req.PermissionIDs != nil {
		// Eliminar permisos existentes
		_, err := database.DB.Exec("DELETE FROM role_permissions WHERE role_id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions"})
			return
		}

		// Insertar nuevos permisos
		for _, permID := range *req.PermissionIDs {
			_, err := database.DB.Exec(`
				INSERT INTO role_permissions (role_id, permission_id)
				VALUES (?, ?)
			`, id, permID)
			if err != nil {
				continue
			}
		}
	}

	// Obtener el rol actualizado
	var role models.Role
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, id).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	role.IsSystem = isSystem == 1
	role.Permissions = getRolePermissions(id)

	c.JSON(http.StatusOK, role)
}

// DeleteRole elimina un rol
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	// Verificar si el rol existe y si es del sistema
	var isSystem int
	err = database.DB.QueryRow("SELECT is_system FROM roles WHERE id = ?", id).Scan(&isSystem)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

	if isSystem == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete system roles"})
		return
	}

	// Verificar si hay usuarios con este rol
	var userCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role_id = ?", id).Scan(&userCount)
	if err == nil && userCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete role with assigned users"})
		return
	}

	// Eliminar el rol (las relaciones en role_permissions se eliminan por CASCADE)
	_, err = database.DB.Exec("DELETE FROM roles WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// GetPermissions obtiene todos los permisos disponibles
func GetPermissions(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, name, description, resource, action
		FROM permissions
		ORDER BY resource, action
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(&perm.ID, &perm.CreatedAt, &perm.UpdatedAt, &perm.Name, &perm.Description, &perm.Resource, &perm.Action)
		if err != nil {
			continue
		}
		permissions = append(permissions, perm)
	}

	c.JSON(http.StatusOK, permissions)
}
