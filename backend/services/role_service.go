package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrRoleNameExists   = errors.New("role name already exists")
	ErrSystemRoleLocked = errors.New("cannot modify system role")
	ErrRoleHasUsers     = errors.New("cannot delete role with assigned users")
)

type CreateRoleInput struct {
	Name          string
	Description   string
	PermissionIDs []int64
}

type UpdateRoleInput struct {
	Name          *string
	Description   *string
	PermissionIDs *[]int64
}

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) ListRoles() ([]models.Role, error) {
	roles, err := s.repo.ListRoles()
	if err != nil {
		return nil, err
	}

	permMap, err := s.repo.ListPermissionsByRoles()
	if err != nil {
		return nil, err
	}

	for i := range roles {
		if perms, ok := permMap[roles[i].ID]; ok {
			roles[i].Permissions = perms
		} else {
			roles[i].Permissions = []models.Permission{}
		}
	}

	return roles, nil
}

func (s *RoleService) GetRole(roleID int64) (models.Role, error) {
	role, err := s.repo.GetRoleByID(roleID)
	if err == sql.ErrNoRows {
		return models.Role{}, ErrRoleNotFound
	}
	if err != nil {
		return models.Role{}, err
	}

	permissions, err := s.repo.ListPermissionsByRole(roleID)
	if err != nil {
		return models.Role{}, err
	}
	role.Permissions = permissions
	return role, nil
}

func (s *RoleService) CreateRole(input CreateRoleInput) (models.Role, error) {
	count, err := s.repo.CountByName(input.Name)
	if err == nil && count > 0 {
		return models.Role{}, ErrRoleNameExists
	}

	roleID, err := s.repo.InsertRole(input.Name, input.Description)
	if err != nil {
		return models.Role{}, err
	}

	if len(input.PermissionIDs) > 0 {
		if err := s.repo.ReplaceRolePermissions(roleID, input.PermissionIDs); err != nil {
			return models.Role{}, err
		}
	}

	return s.GetRole(roleID)
}

func (s *RoleService) UpdateRole(roleID int64, input UpdateRoleInput) (models.Role, error) {
	role, err := s.repo.GetRoleByID(roleID)
	if err == sql.ErrNoRows {
		return models.Role{}, ErrRoleNotFound
	}
	if err != nil {
		return models.Role{}, err
	}

	if role.IsSystem {
		return models.Role{}, ErrSystemRoleLocked
	}

	if err := s.repo.UpdateRoleFields(roleID, input.Name, input.Description); err != nil {
		return models.Role{}, err
	}

	if input.PermissionIDs != nil {
		if err := s.repo.ReplaceRolePermissions(roleID, *input.PermissionIDs); err != nil {
			return models.Role{}, err
		}
	}

	return s.GetRole(roleID)
}

func (s *RoleService) DeleteRole(roleID int64) error {
	role, err := s.repo.GetRoleByID(roleID)
	if err == sql.ErrNoRows {
		return ErrRoleNotFound
	}
	if err != nil {
		return err
	}

	if role.IsSystem {
		return ErrSystemRoleLocked
	}

	userCount, err := s.repo.CountUsersByRole(roleID)
	if err == nil && userCount > 0 {
		return ErrRoleHasUsers
	}

	return s.repo.DeleteRole(roleID)
}

func (s *RoleService) ListPermissions() ([]models.Permission, error) {
	return s.repo.ListPermissions()
}
