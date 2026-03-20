package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/middleware"
	"smartech/backend/models"
	"smartech/backend/validation"
	"time"

	"github.com/gin-gonic/gin"
)

// Register registra un nuevo usuario
func Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	// Verificar si el email ya existe
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		apiErr := errors.NewDatabaseError("Check email existence", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	if count > 0 {
		apiErr := errors.NewConflict("Email already registered")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Obtener el rol de usuario por defecto
	var userRole models.Role
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE name = ?
	`, "usuario").Scan(&userRole.ID, &userRole.CreatedAt, &userRole.UpdatedAt, &userRole.Name, &userRole.Description, new(int))

	if err == sql.ErrNoRows {
		apiErr := errors.NewInternal("Default role not found")
		c.JSON(apiErr.Code, apiErr)
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch default role", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Crear nuevo usuario
	user := models.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: userRole.ID,
	}

	// Encriptar contraseña
	if err := user.HashPassword(req.Password); err != nil {
		apiErr := errors.NewInternal("Failed to hash password")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Guardar en la base de datos
	result, err := database.DB.Exec(`
		INSERT INTO users (name, email, password, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, user.Name, user.Email, user.Password, user.RoleID)
	if err != nil {
		apiErr := errors.NewDatabaseError("Create user", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		apiErr := errors.NewDatabaseError("Get user ID", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	user.ID = userID

	// Obtener permisos del rol
	userRole.Permissions = getRolePermissions(userRole.ID)

	// Extraer nombres de permisos
	permissions := make([]string, len(userRole.Permissions))
	for i, perm := range userRole.Permissions {
		permissions[i] = perm.Name
	}

	// Generar token JWT
	token, err := middleware.GenerateToken(uint(user.ID), user.Email, uint(userRole.ID), userRole.Name, permissions)
	if err != nil {
		apiErr := errors.NewInternal("Failed to generate JWT token")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Cargar user completo con role y permisos
	user.Role = &userRole

	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login autentica un usuario
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	// Buscar usuario por email
	var user models.User
	var isSystem int
	err := database.DB.QueryRow(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name, r.description, r.is_system
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.email = ?
	`, req.Email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.RoleID,
		new(int64), new(time.Time), new(time.Time), new(string), new(string), &isSystem,
	)

	if err == sql.ErrNoRows {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch user by email", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar contraseña
	if err := user.CheckPassword(req.Password); err != nil {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Obtener rol completo con permisos
	var role models.Role
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, user.RoleID).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	role.IsSystem = isSystem == 1
	role.Permissions = getRolePermissions(role.ID)
	user.Role = &role

	// Extraer nombres de permisos
	permissions := make([]string, len(role.Permissions))
	for i, perm := range role.Permissions {
		permissions[i] = perm.Name
	}

	// Generar token JWT
	token, err := middleware.GenerateToken(uint(user.ID), user.Email, uint(role.ID), role.Name, permissions)
	if err != nil {
		apiErr := errors.NewInternal("Failed to generate JWT token")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// GetCurrentUser obtiene la información del usuario autenticado
func GetCurrentUser(c *gin.Context) {
	userid, exists := c.Get("userid")
	if !exists {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var user models.User
	var isSystem int
	err := database.DB.QueryRow(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.role_id
		FROM users u
		WHERE u.id = ?
	`, userid).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.RoleID)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("User", userid)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch current user", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Obtener rol con permisos
	var role models.Role
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, name, description, is_system
		FROM roles
		WHERE id = ?
	`, user.RoleID).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name, &role.Description, &isSystem)

	role.IsSystem = isSystem == 1
	role.Permissions = getRolePermissions(role.ID)
	user.Role = &role

	c.JSON(http.StatusOK, user)
}

// Logout cierra sesión (en implementación con JWT, principalmente es del lado del cliente)
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
