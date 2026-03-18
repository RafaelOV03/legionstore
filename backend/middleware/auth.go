package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "smartech-secret-key-2024" // Clave por defecto (cambiar en producción)
	}
	return secret
}

// Claims estructura para el payload del JWT
type Claims struct {
	Userid      uint     `json:"user_id"`
	Email       string   `json:"email"`
	Roleid      uint     `json:"role_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateToken genera un token JWT para un usuario
func GenerateToken(userid uint, email string, roleid uint, roleName string, permissions []string) (string, error) {
	claims := Claims{
		Userid:      userid,
		Email:       email,
		Roleid:      roleid,
		RoleName:    roleName,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida un token JWT y retorna las claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// AuthMiddleware verifica que el usuario esté autenticado
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("userid", claims.Userid)
		c.Set("email", claims.Email)
		c.Set("roleid", claims.Roleid)
		c.Set("roleName", claims.RoleName)
		c.Set("permissions", claims.Permissions)
		c.Next()
	}
}

// RequirePermission verifica que el usuario tenga un permiso específico
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		permList, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene el permiso
		hasPermission := false
		for _, p := range permList {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole verifica que el usuario tenga un rol específico
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("roleName")
		if !exists || userRole != roleName {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminMiddleware verifica que el usuario sea administrador (deprecated, usar RequireRole)
func AdminMiddleware() gin.HandlerFunc {
	return RequireRole("administrador")
}
