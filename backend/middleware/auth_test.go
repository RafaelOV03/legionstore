package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: GenerateToken crea un token válido
func TestGenerateToken_Success(t *testing.T) {
	userid := uint(1)
	email := "test@example.com"
	roleid := uint(1)
	roleName := "admin"
	permissions := []string{"users.read", "users.write"}

	token, err := GenerateToken(userid, email, roleid, roleName, permissions)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// Test: ValidateToken valida un token válido
func TestValidateToken_Success(t *testing.T) {
	userid := uint(1)
	email := "test@example.com"
	roleid := uint(1)
	roleName := "admin"
	permissions := []string{"users.read", "users.write"}

	// Generar token
	token, err := GenerateToken(userid, email, roleid, roleName, permissions)
	require.NoError(t, err)

	// Validar token
	claims, err := ValidateToken(token)

	require.NoError(t, err)
	assert.Equal(t, userid, claims.Userid)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roleid, claims.Roleid)
	assert.Equal(t, roleName, claims.RoleName)
	assert.Equal(t, permissions, claims.Permissions)
}

// Test: ValidateToken rechaza token inválido
func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	claims, err := ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Test: ValidateToken rechaza token vacío
func TestValidateToken_EmptyToken(t *testing.T) {
	emptyToken := ""

	claims, err := ValidateToken(emptyToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Test: Token expirado es rechazado
func TestValidateToken_ExpiredToken(t *testing.T) {
	// Crear un token con expiración en el pasado
	claims := Claims{
		Userid:      1,
		Email:       "test@example.com",
		Roleid:      1,
		RoleName:    "admin",
		Permissions: []string{"users.read"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expirado hace 1 hora
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// Intentar validar token expirado
	_, err = ValidateToken(tokenString)
	assert.Error(t, err)
}

// Test: AuthMiddleware permite solicitudes con token válido
func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generar token válido
	token, err := GenerateToken(1, "test@example.com", 1, "admin", []string{"users.read"})
	require.NoError(t, err)

	// Crear request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar middleware
	middleware := AuthMiddleware()
	middleware(c)

	// Verificar que no hay error (StatusOK)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uint(1), c.GetUint("userid"))
	assert.Equal(t, "test@example.com", c.GetString("email"))
}

// Test: AuthMiddleware rechaza solicitudes sin token
func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	// Sin Authorization header

	// Ejecutar middleware
	middleware := AuthMiddleware()
	middleware(c)

	// Verificar respuesta
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: AuthMiddleware rechaza token inválido
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.string")

	// Ejecutar middleware
	middleware := AuthMiddleware()
	middleware(c)

	// Verificar respuesta
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: RequirePermission permite acceso con permiso correcto
func TestRequirePermission_HasPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Simular usuario autenticado con permisos
	permissions := []string{"users.read", "users.write"}
	c.Set("permissions", permissions)

	// Ejecutar middleware con permiso que el usuario tiene
	middleware := RequirePermission("users.read")
	middleware(c)

	// Verificar que llegó sin error (StatusOK)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test: RequirePermission rechaza acceso sin permiso
func TestRequirePermission_NoPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Simular usuario autenticado sin permisos suficientes
	permissions := []string{"products.read"}
	c.Set("permissions", permissions)

	// Ejecutar middleware pidiendo un permiso que el usuario no tiene
	middleware := RequirePermission("users.write")
	middleware(c)

	// Verificar respuesta
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Test: RequirePermission rechaza si no hay permissions set
func TestRequirePermission_NoPermissionsSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	// Sin permissions set

	// Ejecutar middleware
	middleware := RequirePermission("users.read")
	middleware(c)

	// Verificar respuesta
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Test: GenerateToken con múltiples permisos
func TestGenerateToken_MultiplePermissions(t *testing.T) {
	userid := uint(2)
	email := "manager@example.com"
	roleid := uint(2)
	roleName := "manager"
	permissions := []string{
		"users.read",
		"users.write",
		"products.read",
		"products.write",
		"orders.read",
	}

	token, err := GenerateToken(userid, email, roleid, roleName, permissions)
	require.NoError(t, err)

	claims, err := ValidateToken(token)
	require.NoError(t, err)

	assert.Equal(t, len(permissions), len(claims.Permissions))
	for i, perm := range permissions {
		assert.Equal(t, perm, claims.Permissions[i])
	}
}

// Test: Claims estructura
func TestClaimsStructure(t *testing.T) {
	claims := Claims{
		Userid:      1,
		Email:       "test@example.com",
		Roleid:      1,
		RoleName:    "admin",
		Permissions: []string{"users.read"},
	}

	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.RoleName)
	assert.NotEmpty(t, claims.Permissions)
}
