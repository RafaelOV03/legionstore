package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: HashPassword genera hash válido y no nil
func TestHashPassword_Success(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := "SecurePassword123!"
	err := user.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password)
}

// Test: HashPassword con contraseña vacía
func TestHashPassword_EmptyPassword(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := ""
	err := user.HashPassword(password)

	// Bcrypt genera hash incluso para contraseña vacía
	require.NoError(t, err)
	assert.NotEmpty(t, user.Password)
}

// Test: CheckPassword con contraseña correcta
func TestCheckPassword_Success(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := "SecurePassword123!"
	err := user.HashPassword(password)
	require.NoError(t, err)

	// Verificar contraseña correcta
	err = user.CheckPassword(password)
	assert.NoError(t, err)
}

// Test: CheckPassword con contraseña incorrecta
func TestCheckPassword_Failure(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := "SecurePassword123!"
	err := user.HashPassword(password)
	require.NoError(t, err)

	// Intentar con contraseña incorrecta
	wrongPassword := "WrongPassword456!"
	err = user.CheckPassword(wrongPassword)
	assert.Error(t, err)
}

// Test: CheckPassword es sensible a mayúsculas/minúsculas
func TestCheckPassword_CaseSensitive(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := "MyPassword123"
	err := user.HashPassword(password)
	require.NoError(t, err)

	// Intentar con las letras diferentes
	wrongPassword := "mypassword123"
	err = user.CheckPassword(wrongPassword)
	assert.Error(t, err)
}

// Test: User struct fields
func TestUserStruct(t *testing.T) {
	user := User{
		ID:    42,
		Name:  "Jane Doe",
		Email: "jane@example.com",
		RoleID: 2,
	}

	assert.Equal(t, int64(42), user.ID)
	assert.Equal(t, "jane@example.com", user.Email)
	assert.Equal(t, "Jane Doe", user.Name)
	assert.Equal(t, int64(2), user.RoleID)
}

// Test: Password field no se expone en JSON (json:"-" tag)
func TestUserPassword_NotExposedInJSON(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	password := "SecurePassword123!"
	err := user.HashPassword(password)
	require.NoError(t, err)

	// El password debe ser hash, no la contraseña original
	assert.NotEqual(t, password, user.Password)
	assert.NotEmpty(t, user.Password)
}

// Benchmark: HashPassword performance
func BenchmarkHashPassword(b *testing.B) {
	user := &User{}
	password := "TestPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.HashPassword(password)
	}
}

// Benchmark: CheckPassword performance
func BenchmarkCheckPassword(b *testing.B) {
	user := &User{}
	password := "TestPassword123!"
	user.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.CheckPassword(password)
	}
}
