package controllers

import (
	"smartech/backend/database"

	"github.com/gin-gonic/gin"
)

// logAuditoria registra una acción en el log de auditoría
// Esta función está disponible para todos los controladores
func logAuditoria(c *gin.Context, accion, entidad string, entidadID int64, valorAnterior, valorNuevo string) {
	userID, exists := c.Get("userid")
	if !exists {
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		// Intentar conversión desde uint
		if uidUint, ok := userID.(uint); ok {
			uid = int64(uidUint)
		} else {
			return
		}
	}

	ipAddress := c.ClientIP()

	database.DB.Exec(`
		INSERT INTO logs_auditoria (usuario_id, accion, entidad, entidad_id, valor_anterior, valor_nuevo, ip_address)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		uid, accion, entidad, entidadID, valorAnterior, valorNuevo, ipAddress)
}
