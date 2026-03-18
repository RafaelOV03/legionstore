package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User representa un usuario del sistema
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No exponer en JSON
	RoleID    int64     `json:"role_id"`
	Role      *Role     `json:"role,omitempty"` // Relación con el rol
}

// HashPassword genera un hash bcrypt de la contraseña
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica si la contraseña proporcionada coincide con el hash almacenado
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
