package repositories

import (
	"database/sql"
	"smartech/backend/models"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(userID int64) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name, email, password, role_id
		FROM users
		WHERE id = ?
	`, userID).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, name, email, password, role_id
		FROM users
		WHERE email = ?
	`, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) ListUsersWithRoleBasic() ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name, r.description, r.is_system
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
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
		user.Role = &role
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) CountByEmail(email string, excludeUserID *int64) (int, error) {
	var count int
	if excludeUserID != nil {
		err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", email, *excludeUserID).Scan(&count)
		return count, err
	}

	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count, err
}

func (r *UserRepository) CountByID(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	return count, err
}

func (r *UserRepository) CountRolesByID(roleID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM roles WHERE id = ?", roleID).Scan(&count)
	return count, err
}

func (r *UserRepository) InsertUser(user models.User) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO users (name, email, password, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, user.Name, user.Email, user.Password, user.RoleID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) UpdateUserFields(userID int64, name, email, password *string, roleID *int64) error {
	updates := make([]string, 0)
	args := make([]interface{}, 0)

	if name != nil && *name != "" {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if email != nil && *email != "" {
		updates = append(updates, "email = ?")
		args = append(args, *email)
	}
	if password != nil && *password != "" {
		updates = append(updates, "password = ?")
		args = append(args, *password)
	}
	if roleID != nil && *roleID != 0 {
		updates = append(updates, "role_id = ?")
		args = append(args, *roleID)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, userID)

	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *UserRepository) DeleteUser(userID int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}
