package services

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
	"smartech/backend/repositories"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already registered")
	ErrInvalidRoleID    = errors.New("invalid role id")
	ErrCannotDeleteSelf = errors.New("cannot delete your own user")
)

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	RoleID   int64
}

type UpdateUserInput struct {
	Name     *string
	Email    *string
	Password *string
	RoleID   *int64
}

type UserService struct {
	repo     *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserService(repo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *UserService {
	return &UserService{repo: repo, roleRepo: roleRepo}
}

func (s *UserService) hydrateUserRole(user *models.User, roleCache map[int64]*models.Role) {
	if user.Role == nil {
		return
	}

	if cachedRole, ok := roleCache[user.Role.ID]; ok {
		user.Role = cachedRole
		return
	}

	permissions, err := s.roleRepo.ListPermissionsByRole(user.Role.ID)
	if err != nil {
		user.Role.Permissions = []models.Permission{}
		roleCache[user.Role.ID] = user.Role
		return
	}

	user.Role.Permissions = permissions
	roleCache[user.Role.ID] = user.Role
}

func (s *UserService) getUserWithRole(userID int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetRoleByID(user.RoleID)
	if err == nil {
		permissions, permErr := s.roleRepo.ListPermissionsByRole(role.ID)
		if permErr == nil {
			role.Permissions = permissions
		}
		user.Role = &role
	}

	return &user, nil
}

func (s *UserService) ListUsers() ([]models.User, error) {
	users, err := s.repo.ListUsersWithRoleBasic()
	if err != nil {
		return nil, err
	}

	roleCache := make(map[int64]*models.Role)
	for i := range users {
		s.hydrateUserRole(&users[i], roleCache)
	}

	return users, nil
}

func (s *UserService) GetUser(userID int64) (*models.User, error) {
	return s.getUserWithRole(userID)
}

func (s *UserService) CreateUser(input CreateUserInput) (*models.User, error) {
	emailCount, err := s.repo.CountByEmail(input.Email, nil)
	if err == nil && emailCount > 0 {
		return nil, ErrEmailAlreadyUsed
	}

	roleCount, err := s.repo.CountRolesByID(input.RoleID)
	if err != nil || roleCount == 0 {
		return nil, ErrInvalidRoleID
	}

	user := models.User{
		Name:   input.Name,
		Email:  input.Email,
		RoleID: input.RoleID,
	}

	if err := user.HashPassword(input.Password); err != nil {
		return nil, err
	}

	userID, err := s.repo.InsertUser(user)
	if err != nil {
		return nil, err
	}

	return s.getUserWithRole(userID)
}

func (s *UserService) UpdateUser(userID int64, input UpdateUserInput) (*models.User, error) {
	count, err := s.repo.CountByID(userID)
	if err != nil || count == 0 {
		return nil, ErrUserNotFound
	}

	if input.Email != nil && *input.Email != "" {
		emailCount, _ := s.repo.CountByEmail(*input.Email, &userID)
		if emailCount > 0 {
			return nil, ErrEmailAlreadyUsed
		}
	}

	if input.RoleID != nil && *input.RoleID != 0 {
		roleCount, _ := s.repo.CountRolesByID(*input.RoleID)
		if roleCount == 0 {
			return nil, ErrInvalidRoleID
		}
	}

	if input.Password != nil && *input.Password != "" {
		var user models.User
		if err := user.HashPassword(*input.Password); err != nil {
			return nil, err
		}
		hashed := user.Password
		input.Password = &hashed
	}

	if err := s.repo.UpdateUserFields(userID, input.Name, input.Email, input.Password, input.RoleID); err != nil {
		return nil, err
	}

	return s.getUserWithRole(userID)
}

func (s *UserService) DeleteUser(targetUserID int64, actorUserID uint) error {
	if int64(actorUserID) == targetUserID {
		return ErrCannotDeleteSelf
	}

	count, err := s.repo.CountByID(targetUserID)
	if err != nil || count == 0 {
		return ErrUserNotFound
	}

	return s.repo.DeleteUser(targetUserID)
}
