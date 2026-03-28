package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getRoleService() *services.RoleService {
	repo := repositories.NewRoleRepository(database.DB)
	return services.NewRoleService(repo)
}

// GetRoles obtiene todos los roles con sus permisos
func GetRoles(c *gin.Context) {
	roles, err := getRoleService().ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRole obtiene un rol por id con sus permisos
func GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	role, err := getRoleService().GetRole(id)
	if err == services.ErrRoleNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

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

	role, err := getRoleService().CreateRole(services.CreateRoleInput{
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: req.PermissionIDs,
	})
	if err == services.ErrRoleNameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// UpdateRole actualiza un rol existente
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
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

	role, err := getRoleService().UpdateRole(id, services.UpdateRoleInput{
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: req.PermissionIDs,
	})
	if err == services.ErrRoleNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err == services.ErrSystemRoleLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify system roles"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole elimina un rol
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role id"})
		return
	}

	err = getRoleService().DeleteRole(id)
	if err == services.ErrRoleNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	if err == services.ErrSystemRoleLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete system roles"})
		return
	}
	if err == services.ErrRoleHasUsers {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete role with assigned users"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// GetPermissions obtiene todos los permisos disponibles
func GetPermissions(c *gin.Context) {
	permissions, err := getRoleService().ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
