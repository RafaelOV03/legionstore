package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
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

	user, err := getUserService().GetUser(id)
	if err == services.ErrUserNotFound {
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser actualiza un usuario
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
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
}

// DeleteUser elimina un usuario
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	userid, _ := c.Get("userid")
	actorUserID, _ := userid.(uint)

	err = getUserService().DeleteUser(id, actorUserID)
	if err == services.ErrCannotDeleteSelf {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own user"})
		return
	}
	if err == services.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
