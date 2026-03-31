package controllers

import (
	"smartech/backend/database"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

func getRolePermissions(roleID int64) []models.Permission {
	repo := repositories.NewRoleRepository(database.DB)
	permissions, err := repo.ListPermissionsByRole(roleID)
	if err != nil {
		return []models.Permission{}
	}
	return permissions
}
