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

// GetRMAs obtiene todas las RMAs
func GetRMAs(c *gin.Context) {
	estado := c.Query("estado")
	sedeID := c.Query("sede_id")

	query := `
		SELECT r.id, r.created_at, r.updated_at, r.numero_rma, r.producto_id, r.cliente_nombre, 
		       r.cliente_telefono, r.cliente_email, r.num_serie, r.fecha_compra, r.motivo_devolucion,
		       r.diagnostico, r.estado, r.solucion, r.fecha_resolucion, r.usuario_id, r.sede_id, r.notas,
		       p.name as producto_nombre, p.brand as producto_marca,
		       u.name as usuario_nombre,
		       s.nombre as sede_nombre
		FROM rmas r
		INNER JOIN products p ON r.producto_id = p.id
		INNER JOIN users u ON r.usuario_id = u.id
		INNER JOIN sedes s ON r.sede_id = s.id
		WHERE 1=1
	`
	args := []interface{}{}

	if estado != "" {
		query += " AND r.estado = ?"
		args = append(args, estado)
	}
	if sedeID != "" {
		query += " AND r.sede_id = ?"
		args = append(args, sedeID)
	}

	query += " ORDER BY r.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMAs"})
		return
	}
	defer rows.Close()

	type RMAView struct {
		models.RMA
		ProductoNombre string `json:"producto_nombre"`
		ProductoMarca  string `json:"producto_marca"`
		UsuarioNombre  string `json:"usuario_nombre"`
		SedeNombre     string `json:"sede_nombre"`
	}

	var rmas []RMAView
	for rows.Next() {
		var rma RMAView
		var fechaCompra, fechaResolucion sql.NullTime
		err := rows.Scan(&rma.ID, &rma.CreatedAt, &rma.UpdatedAt, &rma.NumeroRMA, &rma.ProductoID,
			&rma.ClienteNombre, &rma.ClienteTelefono, &rma.ClienteEmail, &rma.NumSerie, &fechaCompra,
			&rma.MotivoDevolucion, &rma.Diagnostico, &rma.Estado, &rma.Solucion, &fechaResolucion,
			&rma.UsuarioID, &rma.SedeID, &rma.Notas,
			&rma.ProductoNombre, &rma.ProductoMarca, &rma.UsuarioNombre, &rma.SedeNombre)
		if err != nil {
			continue
		}
		if fechaCompra.Valid {
			rma.FechaCompra = fechaCompra.Time
		}
		if fechaResolucion.Valid {
			rma.FechaResolucion = &fechaResolucion.Time
		}
		rmas = append(rmas, rma)
	}

	c.JSON(http.StatusOK, rmas)
}

// GetRMA obtiene una RMA por ID
func GetRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RMA ID"})
		return
	}

	var rma models.RMA
	var fechaCompra, fechaResolucion sql.NullTime
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, numero_rma, producto_id, cliente_nombre, 
		       cliente_telefono, cliente_email, num_serie, fecha_compra, motivo_devolucion,
		       diagnostico, estado, solucion, fecha_resolucion, usuario_id, sede_id, notas
		FROM rmas WHERE id = ?`, id).
		Scan(&rma.ID, &rma.CreatedAt, &rma.UpdatedAt, &rma.NumeroRMA, &rma.ProductoID,
			&rma.ClienteNombre, &rma.ClienteTelefono, &rma.ClienteEmail, &rma.NumSerie, &fechaCompra,
			&rma.MotivoDevolucion, &rma.Diagnostico, &rma.Estado, &rma.Solucion, &fechaResolucion,
			&rma.UsuarioID, &rma.SedeID, &rma.Notas)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "RMA not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMA"})
		return
	}

	if fechaCompra.Valid {
		rma.FechaCompra = fechaCompra.Time
	}
	if fechaResolucion.Valid {
		rma.FechaResolucion = &fechaResolucion.Time
	}

	// Obtener historial
	histRows, _ := database.DB.Query(`
		SELECT h.id, h.created_at, h.estado_anterior, h.estado_nuevo, h.comentario, u.name
		FROM historial_rmas h
		INNER JOIN users u ON h.usuario_id = u.id
		WHERE h.rma_id = ?
		ORDER BY h.created_at DESC
	`, id)
	defer histRows.Close()

	type HistorialItem struct {
		ID             int64     `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		EstadoAnterior string    `json:"estado_anterior"`
		EstadoNuevo    string    `json:"estado_nuevo"`
		Comentario     string    `json:"comentario"`
		Usuario        string    `json:"usuario"`
	}

	var historial []HistorialItem
	for histRows.Next() {
		var h HistorialItem
		histRows.Scan(&h.ID, &h.CreatedAt, &h.EstadoAnterior, &h.EstadoNuevo, &h.Comentario, &h.Usuario)
		historial = append(historial, h)
	}

	c.JSON(http.StatusOK, gin.H{"rma": rma, "historial": historial})
}

// CreateRMA crea una nueva RMA
func CreateRMA(c *gin.Context) {
	var req struct {
		ProductoID       int64  `json:"producto_id" binding:"required"`
		ClienteNombre    string `json:"cliente_nombre" binding:"required"`
		ClienteTelefono  string `json:"cliente_telefono"`
		ClienteEmail     string `json:"cliente_email"`
		NumSerie         string `json:"num_serie"`
		FechaCompra      string `json:"fecha_compra"`
		MotivoDevolucion string `json:"motivo_devolucion" binding:"required"`
		SedeID           int64  `json:"sede_id" binding:"required"`
		Notas            string `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener usuario del contexto
	userID, _ := c.Get("userid")

	// Generar número de RMA
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&count)
	numeroRMA := fmt.Sprintf("RMA-%d-%04d", time.Now().Year(), count+1)

	var fechaCompra interface{}
	if req.FechaCompra != "" {
		fechaCompra = req.FechaCompra
	} else {
		fechaCompra = nil
	}

	result, err := database.DB.Exec(`
		INSERT INTO rmas (numero_rma, producto_id, cliente_nombre, cliente_telefono, cliente_email, 
		                  num_serie, fecha_compra, motivo_devolucion, estado, usuario_id, sede_id, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'recibido', ?, ?, ?)`,
		numeroRMA, req.ProductoID, req.ClienteNombre, req.ClienteTelefono, req.ClienteEmail,
		req.NumSerie, fechaCompra, req.MotivoDevolucion, userID, req.SedeID, req.Notas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create RMA: " + err.Error()})
		return
	}

	rmaID, _ := result.LastInsertId()

	// Registrar en historial
	database.DB.Exec(`INSERT INTO historial_rmas (rma_id, estado_anterior, estado_nuevo, comentario, usuario_id)
		VALUES (?, '', 'recibido', 'RMA creada', ?)`, rmaID, userID)

	logAuditoria(c, "crear", "rma", rmaID, "", numeroRMA)

	c.JSON(http.StatusCreated, gin.H{"id": rmaID, "numero_rma": numeroRMA})
}

// UpdateRMA actualiza una RMA
func UpdateRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RMA ID"})
		return
	}

	var req struct {
		Diagnostico string `json:"diagnostico"`
		Estado      string `json:"estado"`
		Solucion    string `json:"solucion"`
		Notas       string `json:"notas"`
		Comentario  string `json:"comentario"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	// Obtener estado anterior
	var estadoAnterior string
	database.DB.QueryRow("SELECT estado FROM rmas WHERE id = ?", id).Scan(&estadoAnterior)

	// Si el estado cambia a 'resuelto', registrar fecha de resolución
	var fechaResolucion interface{}
	if req.Estado == "resuelto" && estadoAnterior != "resuelto" {
		fechaResolucion = time.Now()
	} else {
		fechaResolucion = nil
	}

	_, err = database.DB.Exec(`
		UPDATE rmas SET diagnostico = ?, estado = ?, solucion = ?, notas = ?, 
		                fecha_resolucion = COALESCE(?, fecha_resolucion), updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		req.Diagnostico, req.Estado, req.Solucion, req.Notas, fechaResolucion, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update RMA"})
		return
	}

	// Registrar cambio de estado en historial
	if req.Estado != estadoAnterior {
		database.DB.Exec(`INSERT INTO historial_rmas (rma_id, estado_anterior, estado_nuevo, comentario, usuario_id)
			VALUES (?, ?, ?, ?, ?)`, id, estadoAnterior, req.Estado, req.Comentario, userID)
	}

	logAuditoria(c, "editar", "rma", id, estadoAnterior, req.Estado)

	c.JSON(http.StatusOK, gin.H{"message": "RMA updated successfully"})
}

// DeleteRMA elimina una RMA
func DeleteRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RMA ID"})
		return
	}

	// Eliminar historial primero
	database.DB.Exec("DELETE FROM historial_rmas WHERE rma_id = ?", id)

	_, err = database.DB.Exec("DELETE FROM rmas WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete RMA"})
		return
	}

	logAuditoria(c, "eliminar", "rma", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "RMA deleted successfully"})
}

// GetRMAStats obtiene estadísticas de RMAs
func GetRMAStats(c *gin.Context) {
	type Stats struct {
		Total      int `json:"total"`
		Recibidos  int `json:"recibidos"`
		EnRevision int `json:"en_revision"`
		Resueltos  int `json:"resueltos"`
		Rechazados int `json:"rechazados"`
	}

	var stats Stats
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&stats.Total)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'recibido'").Scan(&stats.Recibidos)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'en_revision'").Scan(&stats.EnRevision)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'resuelto'").Scan(&stats.Resueltos)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'rechazado'").Scan(&stats.Rechazados)

	c.JSON(http.StatusOK, stats)
}
