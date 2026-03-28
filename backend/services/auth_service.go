package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDefaultRoleMissing = errors.New("default role not found")
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewAuthService(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *AuthService) hydrateUserRole(user *models.User) error {
	role, err := s.roleRepo.GetRoleByID(user.RoleID)
	if err != nil {
		return err
	}

	permissions, err := s.roleRepo.ListPermissionsByRole(role.ID)
	if err == nil {
		role.Permissions = permissions
	} else {
		role.Permissions = []models.Permission{}
	}
	user.Role = &role
	return nil
}

func (s *AuthService) permissionNames(user *models.User) []string {
	if user.Role == nil {
		return []string{}
	}
	permissions := make([]string, len(user.Role.Permissions))
	for i, perm := range user.Role.Permissions {
		permissions[i] = perm.Name
	}
	return permissions
}

func (s *AuthService) Register(input RegisterInput) (*models.User, []string, error) {
	count, err := s.userRepo.CountByEmail(input.Email, nil)
	if err == nil && count > 0 {
		return nil, nil, ErrEmailAlreadyUsed
	}

	role, err := s.roleRepo.GetRoleByName("usuario")
	if err == sql.ErrNoRows {
		return nil, nil, ErrDefaultRoleMissing
	}
	if err != nil {
		return nil, nil, err
	}

	user := models.User{
		Name:   input.Name,
		Email:  input.Email,
		RoleID: role.ID,
	}
	if err := user.HashPassword(input.Password); err != nil {
		return nil, nil, err
	}

	userID, err := s.userRepo.InsertUser(user)
	if err != nil {
		return nil, nil, err
	}

	createdUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if err := s.hydrateUserRole(&createdUser); err != nil {
		return nil, nil, err
	}

	return &createdUser, s.permissionNames(&createdUser), nil
}

func (s *AuthService) Login(input LoginInput) (*models.User, []string, error) {
	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err == sql.ErrNoRows {
		return nil, nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, nil, err
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := s.hydrateUserRole(&user); err != nil {
		return nil, nil, err
	}

	return &user, s.permissionNames(&user), nil
}

func (s *AuthService) GetCurrentUser(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(int64(userID))
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := s.hydrateUserRole(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
