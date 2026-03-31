package controllers

import (
<<<<<<< HEAD
	"database/sql"
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/validation"
=======
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserService() *services.UserService {
	userRepo := repositories.NewUserRepository(database.DB)
	roleRepo := repositories.NewRoleRepository(database.DB)
	return services.NewUserService(userRepo, roleRepo)
}

// GetUsers obtiene todos los usuarios
func GetUsers(c *gin.Context) {
	users, err := getUserService().ListUsers()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch users", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, users)
}

// GetUser obtiene un usuario por id
func GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid user id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	user, err := getUserWithRole(id)
	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("User", id)
		c.JSON(apiErr.Code, apiErr)
=======
	user, err := getUserService().GetUser(id)
	if err == services.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch user", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, user)
}

// CreateUser crea un nuevo usuario (solo administradores)
func CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" validate:"required,min=3"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		RoleID   int64  `json:"role_id" validate:"required,gt=0"`
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

<<<<<<< HEAD
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

	// Verificar que el rol existe
	err = database.DB.QueryRow("SELECT COUNT(*) FROM roles WHERE id = ?", req.RoleID).Scan(&count)
	if err != nil || count == 0 {
		apiErr := errors.NewBadRequest("Invalid role id")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Crear usuario
	user := models.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleID,
	}

	if err := user.HashPassword(req.Password); err != nil {
		apiErr := errors.NewInternal("Failed to hash password")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO users (name, email, password, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, user.Name, user.Email, user.Password, user.RoleID)
=======
	createdUser, err := getUserService().CreateUser(services.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	})
	if err == services.ErrEmailAlreadyUsed {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	if err == services.ErrInvalidRoleID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Create user", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
<<<<<<< HEAD

	userID, err := result.LastInsertId()
	if err != nil {
		apiErr := errors.NewDatabaseError("Get user ID", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Obtener el usuario creado con role
	createdUser, _ := getUserWithRole(userID)
	c.JSON(201, createdUser)
=======
	c.JSON(http.StatusCreated, createdUser)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// UpdateUser actualiza un usuario
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid user id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Verificar que el usuario existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		apiErr := errors.NewNotFound("User", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

=======
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		RoleID   *int64  `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	updatedUser, err := getUserService().UpdateUser(id, services.UpdateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	})
	if err == services.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
<<<<<<< HEAD
	if req.Email != nil && *req.Email != "" {
		// Verificar que el email no esté en uso por otro usuario
		var emailCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", *req.Email, id).Scan(&emailCount)
		if emailCount > 0 {
			apiErr := errors.NewConflict("Email already in use")
			c.JSON(apiErr.Code, apiErr)
			return
		}
		updates = append(updates, "email = ?")
		args = append(args, *req.Email)
	}
	if req.Password != nil && *req.Password != "" {
		var user models.User
		if err := user.HashPassword(*req.Password); err != nil {
			apiErr := errors.NewInternal("Failed to hash password")
			c.JSON(apiErr.Code, apiErr)
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
			apiErr := errors.NewBadRequest("Invalid role id")
			c.JSON(apiErr.Code, apiErr)
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
			apiErr := errors.NewDatabaseError("Update user", err)
			c.JSON(apiErr.Code, apiErr)
			return
		}
	}

	// Obtener el usuario actualizado
	updatedUser, _ := getUserWithRole(id)
	c.JSON(200, updatedUser)
=======
	if err == services.ErrEmailAlreadyUsed {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}
	if err == services.ErrInvalidRoleID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// DeleteUser elimina un usuario
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid user id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	userid, _ := c.Get("userid")
<<<<<<< HEAD
	if userIDUint, ok := userid.(uint); ok && int64(userIDUint) == id {
		apiErr := errors.ErrForbidden
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar que el usuario existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		apiErr := errors.NewNotFound("User", id)
		c.JSON(apiErr.Code, apiErr)
=======
	actorUserID, _ := userid.(uint)

	err = getUserService().DeleteUser(id, actorUserID)
	if err == services.ErrCannotDeleteSelf {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own user"})
		return
	}
	if err == services.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete user", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
