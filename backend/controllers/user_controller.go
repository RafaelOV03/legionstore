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

// getUserWithRole obtiene un usuario con su rol y permisos
func getUserWithRole(userID int64) (*models.User, error) {
	var user models.User
	var isSystem int

	err := database.DB.QueryRow(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id
		FROM users u
		WHERE u.id = ?
	`, userID).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.RoleID)

	if err != nil {
		return nil, err
	}

	// Obtener rol con permisos
	var role models.Role
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, user.RoleID).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	if err == nil {
		role.IsSystem = isSystem == 1
		role.Permissions = getRolePermissions(role.ID)
		user.Role = &role
	}

	return &user, nil
}

// GetUsers obtiene todos los usuarios
func GetUsers(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name, r.description, r.is_system
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	roleMap := make(map[int64]*models.Role)

	for rows.Next() {
		var user models.User
		var role models.Role
		var isSystem int

		err := rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.RoleID,
			&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem,
		)
		if err != nil {
			continue
		}

		role.IsSystem = isSystem == 1

		// Cache role to avoid loading permissions multiple times
		if cachedRole, exists := roleMap[role.ID]; exists {
			user.Role = cachedRole
		} else {
			role.Permissions = getRolePermissions(role.ID)
			rolePtr := &role
			roleMap[role.ID] = rolePtr
			user.Role = rolePtr
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// GetUser obtiene un usuario por id
func GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	user, err := getUserWithRole(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser crea un nuevo usuario (solo administradores)
func CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RoleID   int64  `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el email ya existe
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err == nil && count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Verificar que el rol existe
	err = database.DB.QueryRow("SELECT COUNT(*) FROM roles WHERE id = ?", req.RoleID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	// Crear usuario
	user := models.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleID,
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO users (name, email, password, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, user.Name, user.Email, user.Password, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user id"})
		return
	}

	// Obtener el usuario creado con role
	createdUser, _ := getUserWithRole(userID)
	c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser actualiza un usuario
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// Verificar que el usuario existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		RoleID   *int64  `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	// Actualizar campos
	if req.Name != nil && *req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Email != nil && *req.Email != "" {
		// Verificar que el email no esté en uso por otro usuario
		var emailCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", *req.Email, id).Scan(&emailCount)
		if emailCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		updates = append(updates, "email = ?")
		args = append(args, *req.Email)
	}
	if req.Password != nil && *req.Password != "" {
		var user models.User
		if err := user.HashPassword(*req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates = append(updates, "password = ?")
		args = append(args, user.Password)
	}
	if req.RoleID != nil && *req.RoleID != 0 {
		// Verificar que el rol existe
		var roleCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM roles WHERE id = ?", *req.RoleID).Scan(&roleCount)
		if roleCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
			return
		}
		updates = append(updates, "role_id = ?")
		args = append(args, *req.RoleID)
	}

	if len(updates) > 0 {
		updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, id)

		query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		_, err := database.DB.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Obtener el usuario actualizado
	updatedUser, _ := getUserWithRole(id)
	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser elimina un usuario
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// No permitir eliminar el propio usuario
	userid, _ := c.Get("userid")
	if userIDUint, ok := userid.(uint); ok && int64(userIDUint) == id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own user"})
		return
	}

	// Verificar que el usuario existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
