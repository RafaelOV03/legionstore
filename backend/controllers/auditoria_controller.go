package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLogs obtiene los logs de auditoría
func GetLogs(c *gin.Context) {
	accion := c.Query("accion")
	entidad := c.Query("entidad")
	usuarioID := c.Query("usuario_id")
	fechaDesde := c.Query("fecha_desde")
	fechaHasta := c.Query("fecha_hasta")
	limit := c.DefaultQuery("limit", "100")

	query := `
		SELECT l.id, l.created_at, l.usuario_id, l.accion, l.entidad, l.entidad_id, 
		       l.valor_anterior, l.valor_nuevo, l.ip_address, u.name as usuario_nombre
		FROM logs_auditoria l
		INNER JOIN users u ON l.usuario_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}

	if accion != "" {
		query += " AND l.accion = ?"
		args = append(args, accion)
	}
	if entidad != "" {
		query += " AND l.entidad = ?"
		args = append(args, entidad)
	}
	if usuarioID != "" {
		query += " AND l.usuario_id = ?"
		args = append(args, usuarioID)
	}
	if fechaDesde != "" {
		query += " AND l.created_at >= ?"
		args = append(args, fechaDesde)
	}
	if fechaHasta != "" {
		query += " AND l.created_at <= ?"
		args = append(args, fechaHasta+" 23:59:59")
	}

	query += " ORDER BY l.created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	defer rows.Close()

	type LogView struct {
		models.LogAuditoria
		UsuarioNombre string `json:"usuario_nombre"`
	}

	var logs []LogView
	for rows.Next() {
		var l LogView
		err := rows.Scan(&l.ID, &l.CreatedAt, &l.UsuarioID, &l.Accion, &l.Entidad, &l.EntidadID,
			&l.ValorAnterior, &l.ValorNuevo, &l.IPAddress, &l.UsuarioNombre)
		if err != nil {
			continue
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, logs)
}

// GetLogStats obtiene estadísticas de acciones
func GetLogStats(c *gin.Context) {
	// Acciones por tipo
	rows, err := database.DB.Query(`
		SELECT accion, COUNT(*) as total
		FROM logs_auditoria
		WHERE created_at >= date('now', '-30 days')
		GROUP BY accion
		ORDER BY total DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	defer rows.Close()

	type AccionStats struct {
		Accion string `json:"accion"`
		Total  int    `json:"total"`
	}

	var accionesPorTipo []AccionStats
	for rows.Next() {
		var a AccionStats
		rows.Scan(&a.Accion, &a.Total)
		accionesPorTipo = append(accionesPorTipo, a)
	}

	// Actividad por usuario
	userRows, _ := database.DB.Query(`
		SELECT u.name, COUNT(*) as total
		FROM logs_auditoria l
		INNER JOIN users u ON l.usuario_id = u.id
		WHERE l.created_at >= date('now', '-30 days')
		GROUP BY l.usuario_id
		ORDER BY total DESC
		LIMIT 10
	`)
	defer userRows.Close()

	type UserStats struct {
		Usuario string `json:"usuario"`
		Total   int    `json:"total"`
	}

	var actividadPorUsuario []UserStats
	for userRows.Next() {
		var u UserStats
		userRows.Scan(&u.Usuario, &u.Total)
		actividadPorUsuario = append(actividadPorUsuario, u)
	}

	// Actividad por día (últimos 7 días)
	dayRows, _ := database.DB.Query(`
		SELECT date(created_at) as fecha, COUNT(*) as total
		FROM logs_auditoria
		WHERE created_at >= date('now', '-7 days')
		GROUP BY date(created_at)
		ORDER BY fecha
	`)
	defer dayRows.Close()

	type DayStats struct {
		Fecha string `json:"fecha"`
		Total int    `json:"total"`
	}

	var actividadPorDia []DayStats
	for dayRows.Next() {
		var d DayStats
		dayRows.Scan(&d.Fecha, &d.Total)
		actividadPorDia = append(actividadPorDia, d)
	}

	c.JSON(http.StatusOK, gin.H{
		"acciones_por_tipo":     accionesPorTipo,
		"actividad_por_usuario": actividadPorUsuario,
		"actividad_por_dia":     actividadPorDia,
	})
}

// GetReporteGanancias obtiene el reporte de ganancias
func GetReporteGanancias(c *gin.Context) {
	fechaDesde := c.DefaultQuery("fecha_desde", "")
	fechaHasta := c.DefaultQuery("fecha_hasta", "")
	sedeID := c.Query("sede_id")

	// Ventas totales
	ventasQuery := `
		SELECT COALESCE(SUM(total), 0) FROM ventas WHERE 1=1
	`
	ventasArgs := []interface{}{}

	if fechaDesde != "" {
		ventasQuery += " AND created_at >= ?"
		ventasArgs = append(ventasArgs, fechaDesde)
	}
	if fechaHasta != "" {
		ventasQuery += " AND created_at <= ?"
		ventasArgs = append(ventasArgs, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		ventasQuery += " AND sede_id = ?"
		ventasArgs = append(ventasArgs, sedeID)
	}

	var totalVentas float64
	database.DB.QueryRow(ventasQuery, ventasArgs...).Scan(&totalVentas)

	// Costo de productos vendidos
	costoQuery := `
		SELECT COALESCE(SUM(vi.cantidad * p.precio_compra), 0)
		FROM venta_items vi
		INNER JOIN products p ON vi.producto_id = p.id
		INNER JOIN ventas v ON vi.venta_id = v.id
		WHERE 1=1
	`
	costoArgs := []interface{}{}

	if fechaDesde != "" {
		costoQuery += " AND v.created_at >= ?"
		costoArgs = append(costoArgs, fechaDesde)
	}
	if fechaHasta != "" {
		costoQuery += " AND v.created_at <= ?"
		costoArgs = append(costoArgs, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		costoQuery += " AND v.sede_id = ?"
		costoArgs = append(costoArgs, sedeID)
	}

	var costoProductos float64
	database.DB.QueryRow(costoQuery, costoArgs...).Scan(&costoProductos)

	// Ingresos por servicios técnicos
	serviciosQuery := `
		SELECT COALESCE(SUM(costo_servicio + costo_repuestos), 0)
		FROM ordenes_trabajo
		WHERE estado = 'entregado'
	`
	serviciosArgs := []interface{}{}

	if fechaDesde != "" {
		serviciosQuery += " AND fecha_entrega >= ?"
		serviciosArgs = append(serviciosArgs, fechaDesde)
	}
	if fechaHasta != "" {
		serviciosQuery += " AND fecha_entrega <= ?"
		serviciosArgs = append(serviciosArgs, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		serviciosQuery += " AND sede_id = ?"
		serviciosArgs = append(serviciosArgs, sedeID)
	}

	var totalServicios float64
	database.DB.QueryRow(serviciosQuery, serviciosArgs...).Scan(&totalServicios)

	gananciaProductos := totalVentas - costoProductos
	gananciaTotal := gananciaProductos + totalServicios

	// Ventas por categoría
	catRows, _ := database.DB.Query(`
		SELECT p.category, COALESCE(SUM(vi.subtotal), 0) as total
		FROM venta_items vi
		INNER JOIN products p ON vi.producto_id = p.id
		INNER JOIN ventas v ON vi.venta_id = v.id
		GROUP BY p.category
		ORDER BY total DESC
	`)
	defer catRows.Close()

	type CatStats struct {
		Categoria string  `json:"categoria"`
		Total     float64 `json:"total"`
	}

	var ventasPorCategoria []CatStats
	for catRows.Next() {
		var cs CatStats
		catRows.Scan(&cs.Categoria, &cs.Total)
		ventasPorCategoria = append(ventasPorCategoria, cs)
	}

	// Ventas por sede
	sedeRows, _ := database.DB.Query(`
		SELECT s.nombre, COALESCE(SUM(v.total), 0) as total
		FROM ventas v
		INNER JOIN sedes s ON v.sede_id = s.id
		GROUP BY v.sede_id
		ORDER BY total DESC
	`)
	defer sedeRows.Close()

	type SedeStats struct {
		Sede  string  `json:"sede"`
		Total float64 `json:"total"`
	}

	var ventasPorSede []SedeStats
	for sedeRows.Next() {
		var ss SedeStats
		sedeRows.Scan(&ss.Sede, &ss.Total)
		ventasPorSede = append(ventasPorSede, ss)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_ventas":         totalVentas,
		"costo_productos":      costoProductos,
		"ganancia_productos":   gananciaProductos,
		"total_servicios":      totalServicios,
		"ganancia_total":       gananciaTotal,
		"ventas_por_categoria": ventasPorCategoria,
		"ventas_por_sede":      ventasPorSede,
	})
}

// GetSegmentaciones obtiene las segmentaciones de clientes
func GetSegmentaciones(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, nombre, descripcion, criterios, activa
		FROM segmentaciones ORDER BY nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch segmentaciones"})
		return
	}
	defer rows.Close()

	var segmentaciones []models.Segmentacion
	for rows.Next() {
		var s models.Segmentacion
		var activa int
		err := rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Nombre, &s.Descripcion, &s.Criterios, &activa)
		if err != nil {
			continue
		}
		s.Activa = activa == 1
		segmentaciones = append(segmentaciones, s)
	}

	c.JSON(http.StatusOK, segmentaciones)
}

// CreateSegmentacion crea una nueva segmentación
func CreateSegmentacion(c *gin.Context) {
	var req struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
		Criterios   string `json:"criterios"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO segmentaciones (nombre, descripcion, criterios, activa)
		VALUES (?, ?, ?, 1)`,
		req.Nombre, req.Descripcion, req.Criterios)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create segmentación"})
		return
	}

	segID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "segmentacion", segID, "", req.Nombre)

	c.JSON(http.StatusCreated, gin.H{"id": segID})
}

// GetPromociones obtiene las promociones
func GetPromociones(c *gin.Context) {
	activas := c.Query("activas")

	query := `
		SELECT id, created_at, updated_at, nombre, descripcion, tipo_descuento, valor, fecha_inicio, 
		       fecha_fin, producto_ids, categorias, activa
		FROM promociones
	`

	if activas == "true" {
		query += ` WHERE activa = 1 AND fecha_inicio <= date('now') AND fecha_fin >= date('now')`
	}

	query += " ORDER BY fecha_inicio DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promociones"})
		return
	}
	defer rows.Close()

	var promociones []models.Promocion
	for rows.Next() {
		var p models.Promocion
		var activa int
		var prodIDs, categorias sql.NullString
		var fechaInicio, fechaFin sql.NullTime
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.Descripcion, &p.Tipo,
			&p.Valor, &fechaInicio, &fechaFin, &prodIDs, &categorias, &activa)
		if err != nil {
			continue
		}
		p.Activa = activa == 1
		if prodIDs.Valid {
			p.ProductosIDs = prodIDs.String
		}
		if fechaInicio.Valid {
			p.FechaInicio = fechaInicio.Time
		}
		if fechaFin.Valid {
			p.FechaFin = fechaFin.Time
		}
		promociones = append(promociones, p)
	}

	c.JSON(http.StatusOK, promociones)
}

// CreatePromocion crea una nueva promoción
func CreatePromocion(c *gin.Context) {
	var req struct {
		Nombre         string  `json:"nombre" binding:"required"`
		Descripcion    string  `json:"descripcion"`
		Tipo           string  `json:"tipo" binding:"required"` // porcentaje, monto_fijo, 2x1
		Valor          float64 `json:"valor"`
		FechaInicio    string  `json:"fecha_inicio" binding:"required"`
		FechaFin       string  `json:"fecha_fin" binding:"required"`
		ProductosIDs   string  `json:"productos_ids"` // JSON array de IDs
		SegmentacionID *int64  `json:"segmentacion_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO promociones (nombre, descripcion, tipo, valor, fecha_inicio, fecha_fin, 
		                         productos_ids, segmentacion_id, activa)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		req.Nombre, req.Descripcion, req.Tipo, req.Valor, req.FechaInicio, req.FechaFin,
		req.ProductosIDs, req.SegmentacionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promoción"})
		return
	}

	promoID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "promocion", promoID, "", req.Nombre)

	c.JSON(http.StatusCreated, gin.H{"id": promoID})
}

// UpdatePromocion actualiza una promoción
func UpdatePromocion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid promoción ID"})
		return
	}

	var req struct {
		Nombre         string  `json:"nombre"`
		Descripcion    string  `json:"descripcion"`
		Tipo           string  `json:"tipo"`
		Valor          float64 `json:"valor"`
		FechaInicio    string  `json:"fecha_inicio"`
		FechaFin       string  `json:"fecha_fin"`
		ProductosIDs   string  `json:"productos_ids"`
		SegmentacionID *int64  `json:"segmentacion_id"`
		Activa         bool    `json:"activa"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activa := 0
	if req.Activa {
		activa = 1
	}

	_, err = database.DB.Exec(`
		UPDATE promociones SET nombre = ?, descripcion = ?, tipo = ?, valor = ?, fecha_inicio = ?,
		                       fecha_fin = ?, productos_ids = ?, segmentacion_id = ?, activa = ?,
		                       updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		req.Nombre, req.Descripcion, req.Tipo, req.Valor, req.FechaInicio, req.FechaFin,
		req.ProductosIDs, req.SegmentacionID, activa, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update promoción"})
		return
	}

	logAuditoria(c, "editar", "promocion", id, "", req.Nombre)

	c.JSON(http.StatusOK, gin.H{"message": "Promoción updated successfully"})
}

// DeletePromocion elimina una promoción
func DeletePromocion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid promoción ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM promociones WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete promoción"})
		return
	}

	logAuditoria(c, "eliminar", "promocion", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Promoción deleted successfully"})
}
