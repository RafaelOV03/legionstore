package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/middleware"
	"smartech/backend/repositories"
	"smartech/backend/services"

	"github.com/gin-gonic/gin"
)

func getAuthService() *services.AuthService {
	userRepo := repositories.NewUserRepository(database.DB)
	roleRepo := repositories.NewRoleRepository(database.DB)
	return services.NewAuthService(userRepo, roleRepo)
}

// Register registra un nuevo usuario
func Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, permissions, err := getAuthService().Register(services.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err == services.ErrEmailAlreadyUsed {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	if err == services.ErrDefaultRoleMissing {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Default role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generar token JWT
	token, err := middleware.GenerateToken(uint(user.ID), user.Email, uint(user.Role.ID), user.Role.Name, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login autentica un usuario
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, permissions, err := getAuthService().Login(services.LoginInput{Email: req.Email, Password: req.Password})
	if err == services.ErrInvalidCredentials {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	// Generar token JWT
	token, err := middleware.GenerateToken(uint(user.ID), user.Email, uint(user.Role.ID), user.Role.Name, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// GetCurrentUser obtiene la información del usuario autenticado
func GetCurrentUser(c *gin.Context) {
	userid, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := getAuthService().GetCurrentUser(userID)
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

// Logout cierra sesión (en implementación con JWT, principalmente es del lado del cliente)
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
