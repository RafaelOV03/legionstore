package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTraspasos obtiene todos los traspasos
func GetTraspasos(c *gin.Context) {
	estado := c.Query("estado")
	sedeOrigenID := c.Query("sede_origen_id")
	sedeDestinoID := c.Query("sede_destino_id")

	query := `
		SELECT t.id, t.created_at, t.updated_at, t.numero_traspaso, t.sede_origen_id, t.sede_destino_id,
		       t.estado, t.fecha_envio, t.fecha_recepcion, t.notas, t.usuario_envia_id, t.usuario_recibe_id,
		       so.nombre as sede_origen_nombre, sd.nombre as sede_destino_nombre,
		       uo.name as usuario_origen_nombre
		FROM traspasos t
		INNER JOIN sedes so ON t.sede_origen_id = so.id
		INNER JOIN sedes sd ON t.sede_destino_id = sd.id
		INNER JOIN users uo ON t.usuario_envia_id = uo.id
		WHERE 1=1
	`
	args := []interface{}{}

	if estado != "" {
		query += " AND t.estado = ?"
		args = append(args, estado)
	}
	if sedeOrigenID != "" {
		query += " AND t.sede_origen_id = ?"
		args = append(args, sedeOrigenID)
	}
	if sedeDestinoID != "" {
		query += " AND t.sede_destino_id = ?"
		args = append(args, sedeDestinoID)
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traspasos"})
		return
	}
	defer rows.Close()

	type TraspasoView struct {
		models.Traspaso
		SedeOrigenNombre    string `json:"sede_origen_nombre"`
		SedeDestinoNombre   string `json:"sede_destino_nombre"`
		UsuarioOrigenNombre string `json:"usuario_origen_nombre"`
	}

	var traspasos []TraspasoView
	for rows.Next() {
		var t TraspasoView
		var fechaEnvio, fechaRecepcion sql.NullTime
		var usuarioRecibeID sql.NullInt64
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.NumeroTraspaso, &t.SedeOrigenID, &t.SedeDestinoID,
			&t.Estado, &fechaEnvio, &fechaRecepcion, &t.Notas, &t.UsuarioEnviaID, &usuarioRecibeID,
			&t.SedeOrigenNombre, &t.SedeDestinoNombre, &t.UsuarioOrigenNombre)
		if err != nil {
			continue
		}
		if fechaEnvio.Valid {
			t.FechaEnvio = &fechaEnvio.Time
		}
		if fechaRecepcion.Valid {
			t.FechaRecepcion = &fechaRecepcion.Time
		}
		if usuarioRecibeID.Valid {
			t.UsuarioRecibeID = &usuarioRecibeID.Int64
		}
		traspasos = append(traspasos, t)
	}

	c.JSON(http.StatusOK, traspasos)
}

// GetTraspaso obtiene un traspaso por ID con sus items
func GetTraspaso(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traspaso ID"})
		return
	}

	var t models.Traspaso
	var fechaEnvio, fechaRecepcion sql.NullTime
	var usuarioRecibeID sql.NullInt64

	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, numero_traspaso, sede_origen_id, sede_destino_id,
		       estado, fecha_envio, fecha_recepcion, notas, usuario_envia_id, usuario_recibe_id
		FROM traspasos WHERE id = ?`, id).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.NumeroTraspaso, &t.SedeOrigenID, &t.SedeDestinoID,
			&t.Estado, &fechaEnvio, &fechaRecepcion, &t.Notas, &t.UsuarioEnviaID, &usuarioRecibeID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traspaso"})
		return
	}

	if fechaEnvio.Valid {
		t.FechaEnvio = &fechaEnvio.Time
	}
	if fechaRecepcion.Valid {
		t.FechaRecepcion = &fechaRecepcion.Time
	}
	if usuarioRecibeID.Valid {
		t.UsuarioRecibeID = &usuarioRecibeID.Int64
	}

	// Obtener items
	itemRows, _ := database.DB.Query(`
		SELECT ti.id, ti.producto_id, ti.cantidad, ti.cantidad_recibida,
		       p.name, p.brand, p.codigo
		FROM traspaso_items ti
		INNER JOIN products p ON ti.producto_id = p.id
		WHERE ti.traspaso_id = ?
	`, id)
	defer itemRows.Close()

	type ItemView struct {
		ID               int64  `json:"id"`
		ProductoID       int64  `json:"producto_id"`
		Cantidad         int    `json:"cantidad"`
		CantidadRecibida int    `json:"cantidad_recibida"`
		ProductoNombre   string `json:"producto_nombre"`
		ProductoMarca    string `json:"producto_marca"`
		ProductoCodigo   string `json:"producto_codigo"`
	}

	var items []ItemView
	for itemRows.Next() {
		var item ItemView
		var cantRecibida sql.NullInt64
		itemRows.Scan(&item.ID, &item.ProductoID, &item.Cantidad, &cantRecibida,
			&item.ProductoNombre, &item.ProductoMarca, &item.ProductoCodigo)
		if cantRecibida.Valid {
			item.CantidadRecibida = int(cantRecibida.Int64)
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{"traspaso": t, "items": items})
}

// CreateTraspaso crea un nuevo traspaso
func CreateTraspaso(c *gin.Context) {
	var req struct {
		SedeOrigenID  int64  `json:"sede_origen_id" binding:"required"`
		SedeDestinoID int64  `json:"sede_destino_id" binding:"required"`
		Notas         string `json:"notas"`
		Items         []struct {
			ProductoID int64 `json:"producto_id"`
			Cantidad   int   `json:"cantidad"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SedeOrigenID == req.SedeDestinoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Las sedes de origen y destino deben ser diferentes"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe incluir al menos un item"})
		return
	}

	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Verificar que las sedes existan
	var sedeOrigenExists, sedeDestinoExists int
	database.DB.QueryRow("SELECT COUNT(*) FROM sedes WHERE id = ?", req.SedeOrigenID).Scan(&sedeOrigenExists)
	database.DB.QueryRow("SELECT COUNT(*) FROM sedes WHERE id = ?", req.SedeDestinoID).Scan(&sedeDestinoExists)
	
	if sedeOrigenExists == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Sede origen con ID %d no existe", req.SedeOrigenID)})
		return
	}
	if sedeDestinoExists == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Sede destino con ID %d no existe", req.SedeDestinoID)})
		return
	}

	// Verificar stock disponible en sede origen
	for _, item := range req.Items {
		var stockDisponible int
		err := database.DB.QueryRow(`
			SELECT COALESCE(cantidad, 0) FROM stock_sedes 
			WHERE producto_id = ? AND sede_id = ?`, item.ProductoID, req.SedeOrigenID).Scan(&stockDisponible)
		
		var nombreProducto string
		database.DB.QueryRow("SELECT name FROM products WHERE id = ?", item.ProductoID).Scan(&nombreProducto)
		
		if err == sql.ErrNoRows || stockDisponible == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("El producto '%s' no tiene stock en la sede origen", nombreProducto),
			})
			return
		}
		if stockDisponible < item.Cantidad {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Stock insuficiente de %s en sede origen (disponible: %d, solicitado: %d)",
					nombreProducto, stockDisponible, item.Cantidad),
			})
			return
		}
	}

	// Generar número de traspaso
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM traspasos").Scan(&count)
	numeroTraspaso := fmt.Sprintf("TRP-%d-%04d", time.Now().Year(), count+1)

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	result, err := tx.Exec(`
		INSERT INTO traspasos (numero_traspaso, sede_origen_id, sede_destino_id, estado, notas, usuario_envia_id)
		VALUES (?, ?, ?, 'pendiente', ?, ?)`,
		numeroTraspaso, req.SedeOrigenID, req.SedeDestinoID, req.Notas, userID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create traspaso: " + err.Error()})
		return
	}

	traspasoID, _ := result.LastInsertId()

	// Insertar items
	for _, item := range req.Items {
		_, err = tx.Exec(`
			INSERT INTO traspaso_items (traspaso_id, producto_id, cantidad)
			VALUES (?, ?, ?)`,
			traspasoID, item.ProductoID, item.Cantidad)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add items: " + err.Error()})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	logAuditoria(c, "crear", "traspaso", traspasoID, "", numeroTraspaso)

	c.JSON(http.StatusCreated, gin.H{"id": traspasoID, "numero_traspaso": numeroTraspaso})
}

// EnviarTraspaso marca un traspaso como enviado y descuenta stock de origen
func EnviarTraspaso(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traspaso ID"})
		return
	}

	var estado string
	var sedeOrigenID int64
	err = database.DB.QueryRow("SELECT estado, sede_origen_id FROM traspasos WHERE id = ?", id).Scan(&estado, &sedeOrigenID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}

	if estado != "pendiente" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden enviar traspasos pendientes"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Obtener items y descontar stock
	itemRows, _ := tx.Query("SELECT producto_id, cantidad FROM traspaso_items WHERE traspaso_id = ?", id)
	type Item struct {
		ProductoID int64
		Cantidad   int
	}
	var items []Item
	for itemRows.Next() {
		var item Item
		itemRows.Scan(&item.ProductoID, &item.Cantidad)
		items = append(items, item)
	}
	itemRows.Close()

	for _, item := range items {
		_, err = tx.Exec(`
			UPDATE stock_sedes SET cantidad = cantidad - ? WHERE producto_id = ? AND sede_id = ?`,
			item.Cantidad, item.ProductoID, sedeOrigenID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}
	}

	now := time.Now()
	_, err = tx.Exec("UPDATE traspasos SET estado = 'enviado', fecha_envio = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", now, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update traspaso"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	logAuditoria(c, "enviar", "traspaso", id, "pendiente", "enviado")

	c.JSON(http.StatusOK, gin.H{"message": "Traspaso enviado successfully"})
}

// RecibirTraspaso marca un traspaso como recibido y añade stock a destino
func RecibirTraspaso(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traspaso ID"})
		return
	}

	// Acepta tanto formato con items como sin items (auto-recibe todo)
	var req struct {
		Items []struct {
			ItemID           int64 `json:"item_id"`
			CantidadRecibida int   `json:"cantidad_recibida"`
		} `json:"items"`
		Notas string `json:"notas"`
	}

	c.ShouldBindJSON(&req) // Ignorar error - items es opcional ahora

	var estado string
	var sedeDestinoID int64
	err = database.DB.QueryRow("SELECT estado, sede_destino_id FROM traspasos WHERE id = ?", id).Scan(&estado, &sedeDestinoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}

	if estado != "enviado" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden recibir traspasos enviados"})
		return
	}

	userID, _ := c.Get("userid")

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Si no se enviaron items, obtener todos los items del traspaso y recibirlos completos
	if len(req.Items) == 0 {
		itemRows, _ := tx.Query("SELECT id, producto_id, cantidad FROM traspaso_items WHERE traspaso_id = ?", id)
		type AutoItem struct {
			ID         int64
			ProductoID int64
			Cantidad   int
		}
		var autoItems []AutoItem
		for itemRows.Next() {
			var ai AutoItem
			itemRows.Scan(&ai.ID, &ai.ProductoID, &ai.Cantidad)
			autoItems = append(autoItems, ai)
		}
		itemRows.Close()

		for _, ai := range autoItems {
			// Actualizar cantidad recibida
			_, err = tx.Exec("UPDATE traspaso_items SET cantidad_recibida = ? WHERE id = ?", ai.Cantidad, ai.ID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
				return
			}

			// Añadir stock a sede destino
			var exists int
			tx.QueryRow("SELECT COUNT(*) FROM stock_sedes WHERE producto_id = ? AND sede_id = ?", ai.ProductoID, sedeDestinoID).Scan(&exists)
			if exists > 0 {
				_, err = tx.Exec(`UPDATE stock_sedes SET cantidad = cantidad + ? WHERE producto_id = ? AND sede_id = ?`,
					ai.Cantidad, ai.ProductoID, sedeDestinoID)
			} else {
				_, err = tx.Exec(`INSERT INTO stock_sedes (producto_id, sede_id, cantidad, stock_minimo) VALUES (?, ?, ?, 5)`,
					ai.ProductoID, sedeDestinoID, ai.Cantidad)
			}
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
				return
			}
		}
	} else {
		// Proceso original para items específicos
		for _, item := range req.Items {
			// Obtener producto_id del item
			var productoID int64
			tx.QueryRow("SELECT producto_id FROM traspaso_items WHERE id = ?", item.ItemID).Scan(&productoID)

			// Actualizar cantidad recibida
			_, err = tx.Exec("UPDATE traspaso_items SET cantidad_recibida = ? WHERE id = ?", item.CantidadRecibida, item.ItemID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
				return
			}

			// Añadir stock a sede destino (crear registro si no existe)
			var exists int
			tx.QueryRow("SELECT COUNT(*) FROM stock_sedes WHERE producto_id = ? AND sede_id = ?", productoID, sedeDestinoID).Scan(&exists)
			if exists > 0 {
				_, err = tx.Exec(`
					UPDATE stock_sedes SET cantidad = cantidad + ? WHERE producto_id = ? AND sede_id = ?`,
					item.CantidadRecibida, productoID, sedeDestinoID)
			} else {
				_, err = tx.Exec(`
					INSERT INTO stock_sedes (producto_id, sede_id, cantidad, stock_minimo) VALUES (?, ?, ?, 5)`,
					productoID, sedeDestinoID, item.CantidadRecibida)
			}
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
				return
			}
		}
	}

	now := time.Now()
	notasUpdate := ""
	if req.Notas != "" {
		notasUpdate = req.Notas
	}
	_, err = tx.Exec(`
		UPDATE traspasos SET estado = 'recibido', fecha_recepcion = ?, usuario_recibe_id = ?, 
		                     notas = CASE WHEN ? != '' THEN notas || ' | Recepción: ' || ? ELSE notas END,
		                     updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`, now, userID, notasUpdate, notasUpdate, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update traspaso"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	logAuditoria(c, "recibir", "traspaso", id, "enviado", "recibido")

	c.JSON(http.StatusOK, gin.H{"message": "Traspaso recibido successfully"})
}

// CancelarTraspaso cancela un traspaso pendiente
func CancelarTraspaso(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traspaso ID"})
		return
	}

	var estado string
	err = database.DB.QueryRow("SELECT estado FROM traspasos WHERE id = ?", id).Scan(&estado)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}

	if estado != "pendiente" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden cancelar traspasos pendientes"})
		return
	}

	_, err = database.DB.Exec("UPDATE traspasos SET estado = 'cancelado', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel traspaso"})
		return
	}

	logAuditoria(c, "cancelar", "traspaso", id, "pendiente", "cancelado")

	c.JSON(http.StatusOK, gin.H{"message": "Traspaso cancelled successfully"})
}

// DeleteTraspaso elimina un traspaso (solo si está cancelado)
func DeleteTraspaso(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traspaso ID"})
		return
	}

	var estado string
	err = database.DB.QueryRow("SELECT estado FROM traspasos WHERE id = ?", id).Scan(&estado)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}

	if estado != "cancelado" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden eliminar traspasos cancelados"})
		return
	}

	database.DB.Exec("DELETE FROM traspaso_items WHERE traspaso_id = ?", id)
	database.DB.Exec("DELETE FROM traspasos WHERE id = ?", id)

	logAuditoria(c, "eliminar", "traspaso", id, "", "")

	c.JSON(http.StatusOK, gin.H{"message": "Traspaso deleted successfully"})
}
