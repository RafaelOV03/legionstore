package repository

import (
	"database/sql"
	"smartech/backend/models"
)

// UserRepository maneja toda la lógica de acceso a datos de usuarios
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll obtiene todos los usuarios
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name
		FROM usuarios u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, DBError(err, "GetAll users")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var role models.Role
		var roleID sql.NullInt64

		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &roleID,
			&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			continue
		}

		if roleID.Valid {
			user.RoleID = roleID.Int64
			user.Role = &role
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	var role models.Role
	var roleID sql.NullInt64

	err := r.db.QueryRow(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name
		FROM usuarios u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &roleID,
		&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "GetByID user")
	}

	if roleID.Valid {
		user.RoleID = roleID.Int64
		user.Role = &role
	}

	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	var role models.Role
	var roleID sql.NullInt64

	err := r.db.QueryRow(`
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name
		FROM usuarios u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = ?
	`, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &roleID,
		&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "GetByEmail user")
	}

	if roleID.Valid {
		user.RoleID = roleID.Int64
		user.Role = &role
	}

	return &user, nil
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(user *models.User) error {
	result, err := r.db.Exec(`
		INSERT INTO usuarios (name, email, password, role_id)
		VALUES (?, ?, ?, ?)
	`, user.Name, user.Email, user.Password, user.RoleID)

	if err != nil {
		return DBError(err, "Create user")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return DBError(err, "Get last insert id")
	}

	user.ID = id
	return nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(id int64, user *models.User) error {
	result, err := r.db.Exec(`
		UPDATE usuarios
		SET name = ?, email = ?, role_id = ?
		WHERE id = ?
	`, user.Name, user.Email, user.RoleID, id)

	if err != nil {
		return DBError(err, "Update user")
	}

	return CheckRowsAffected(result, "Update user")
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *UserRepository) UpdatePassword(id int64, hashedPassword string) error {
	result, err := r.db.Exec(`
		UPDATE usuarios SET password = ? WHERE id = ?
	`, hashedPassword, id)

	if err != nil {
		return DBError(err, "Update user password")
	}

	return CheckRowsAffected(result, "Update user password")
}

// Delete elimina un usuario
func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec(`
		DELETE FROM usuarios WHERE id = ?
	`, id)

	if err != nil {
		return DBError(err, "Delete user")
	}

	return CheckRowsAffected(result, "Delete user")
}

// GetByRole obtiene todos los usuarios con un rol específico
func (r *UserRepository) GetByRole(roleID int64) ([]models.User, error) {
	query := `
		SELECT u.id, u.created_at, u.updated_at, u.name, u.email, u.password, u.role_id,
		       r.id, r.created_at, r.updated_at, r.name
		FROM usuarios u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.role_id = ?
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, DBError(err, "GetByRole users")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var role models.Role
		var roleID sql.NullInt64

		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &roleID,
			&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			continue
		}

		if roleID.Valid {
			user.RoleID = roleID.Int64
			user.Role = &role
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// Count obtiene el total de usuarios
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&count)
	if err != nil {
		return 0, DBError(err, "Count users")
	}
	return count, nil
}
