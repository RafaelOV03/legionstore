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

// GetOrdenesTrabajo obtiene todas las órdenes de trabajo
func GetOrdenesTrabajo(c *gin.Context) {
	estado := c.Query("estado")
	prioridad := c.Query("prioridad")
	tecnicoID := c.Query("tecnico_id")
	sedeID := c.Query("sede_id")

	query := `
		SELECT ot.id, ot.created_at, ot.updated_at, ot.numero_orden, ot.cliente_nombre, 
		       ot.cliente_telefono, ot.equipo, ot.num_serie, ot.marca, ot.modelo,
		       ot.problema_reportado, ot.diagnostico_tecnico, ot.solucion_aplicada, ot.estado,
		       ot.prioridad, ot.fecha_ingreso, ot.fecha_promesa, ot.fecha_entrega, ot.costo_mano_obra,
		       ot.costo_repuestos, ot.tecnico_id, ot.sede_id, ot.notas,
		       COALESCE(u.name, 'Sin asignar') as tecnico_nombre, s.nombre as sede_nombre
		FROM ordenes_trabajo ot
		LEFT JOIN users u ON ot.tecnico_id = u.id
		INNER JOIN sedes s ON ot.sede_id = s.id
		WHERE 1=1
	`
	args := []interface{}{}

	if estado != "" {
		query += " AND ot.estado = ?"
		args = append(args, estado)
	}
	if prioridad != "" {
		query += " AND ot.prioridad = ?"
		args = append(args, prioridad)
	}
	if tecnicoID != "" {
		query += " AND ot.tecnico_id = ?"
		args = append(args, tecnicoID)
	}
	if sedeID != "" {
		query += " AND ot.sede_id = ?"
		args = append(args, sedeID)
	}

	query += " ORDER BY CASE ot.prioridad WHEN 'urgente' THEN 1 WHEN 'alta' THEN 2 WHEN 'media' THEN 3 ELSE 4 END, ot.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch órdenes de trabajo"})
		return
	}
	defer rows.Close()

	type OrdenView struct {
		models.OrdenTrabajo
		TecnicoNombre string `json:"tecnico_nombre"`
		SedeNombre    string `json:"sede_nombre"`
	}

	var ordenes []OrdenView
	for rows.Next() {
		var o OrdenView
		var fechaPromesa, fechaEntrega sql.NullTime
		var tecnicoID sql.NullInt64
		var diagnostico, solucion, notas sql.NullString
		err := rows.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.NumeroOrden, &o.ClienteNombre,
			&o.ClienteTelefono, &o.Equipo, &o.NumSerie, &o.Marca, &o.Modelo,
			&o.ProblemaReportado, &diagnostico, &solucion, &o.Estado,
			&o.Prioridad, &o.FechaIngreso, &fechaPromesa, &fechaEntrega, &o.CostoManoObra,
			&o.CostoRepuestos, &tecnicoID, &o.SedeID, &notas,
			&o.TecnicoNombre, &o.SedeNombre)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		if diagnostico.Valid {
			o.DiagnosticoTecnico = diagnostico.String
		}
		if solucion.Valid {
			o.SolucionAplicada = solucion.String
		}
		if notas.Valid {
			o.Notas = notas.String
		}
		if fechaPromesa.Valid {
			o.FechaPromesa = &fechaPromesa.Time
		}
		if fechaEntrega.Valid {
			o.FechaEntrega = &fechaEntrega.Time
		}
		if tecnicoID.Valid {
			o.TecnicoID = &tecnicoID.Int64
		}
		ordenes = append(ordenes, o)
	}

	c.JSON(http.StatusOK, ordenes)
}

// GetOrdenTrabajo obtiene una orden de trabajo por ID
func GetOrdenTrabajo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var o models.OrdenTrabajo
	var fechaPromesa, fechaEntrega sql.NullTime
	var tecnicoID sql.NullInt64
	var diagnostico, solucion, notas sql.NullString

	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, numero_orden, cliente_nombre, 
		       cliente_telefono, equipo, num_serie, marca, modelo,
		       problema_reportado, diagnostico_tecnico, solucion_aplicada, estado,
		       prioridad, fecha_ingreso, fecha_promesa, fecha_entrega, costo_mano_obra,
		       costo_repuestos, tecnico_id, sede_id, notas
		FROM ordenes_trabajo WHERE id = ?`, id).
		Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.NumeroOrden, &o.ClienteNombre,
			&o.ClienteTelefono, &o.Equipo, &o.NumSerie, &o.Marca, &o.Modelo,
			&o.ProblemaReportado, &diagnostico, &solucion, &o.Estado,
			&o.Prioridad, &o.FechaIngreso, &fechaPromesa, &fechaEntrega, &o.CostoManoObra,
			&o.CostoRepuestos, &tecnicoID, &o.SedeID, &notas)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orden"})
		return
	}

	if diagnostico.Valid {
		o.DiagnosticoTecnico = diagnostico.String
	}
	if solucion.Valid {
		o.SolucionAplicada = solucion.String
	}
	if notas.Valid {
		o.Notas = notas.String
	}
	if fechaPromesa.Valid {
		o.FechaPromesa = &fechaPromesa.Time
	}
	if fechaEntrega.Valid {
		o.FechaEntrega = &fechaEntrega.Time
	}
	if tecnicoID.Valid {
		o.TecnicoID = &tecnicoID.Int64
	}

	// Obtener insumos utilizados
	insumoRows, _ := database.DB.Query(`
		SELECT io.id, io.insumo_id, io.cantidad, i.nombre, i.codigo
		FROM insumos_orden io
		INNER JOIN insumos i ON io.insumo_id = i.id
		WHERE io.orden_id = ?
	`, id)
	defer insumoRows.Close()

	type InsumoOrden struct {
		ID           int64  `json:"id"`
		InsumoID     int64  `json:"insumo_id"`
		Cantidad     int    `json:"cantidad"`
		InsumoNombre string `json:"insumo_nombre"`
		InsumoCodigo string `json:"insumo_codigo"`
	}

	var insumos []InsumoOrden
	for insumoRows.Next() {
		var ins InsumoOrden
		insumoRows.Scan(&ins.ID, &ins.InsumoID, &ins.Cantidad, &ins.InsumoNombre, &ins.InsumoCodigo)
		insumos = append(insumos, ins)
	}

	// Obtener trazabilidad
	trazRows, _ := database.DB.Query(`
		SELECT t.id, t.created_at, t.accion, t.detalle, u.name
		FROM trazabilidad t
		INNER JOIN users u ON t.usuario_id = u.id
		WHERE t.orden_trabajo_id = ?
		ORDER BY t.created_at DESC
	`, id)
	defer trazRows.Close()

	type TrazItem struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Accion    string    `json:"accion"`
		Detalle   string    `json:"detalle"`
		Usuario   string    `json:"usuario"`
	}

	var trazabilidad []TrazItem
	for trazRows.Next() {
		var t TrazItem
		trazRows.Scan(&t.ID, &t.CreatedAt, &t.Accion, &t.Detalle, &t.Usuario)
		trazabilidad = append(trazabilidad, t)
	}

	c.JSON(http.StatusOK, gin.H{"orden": o, "insumos": insumos, "trazabilidad": trazabilidad})
}

// CreateOrdenTrabajo crea una nueva orden de trabajo
func CreateOrdenTrabajo(c *gin.Context) {
	var req struct {
		ClienteNombre     string `json:"cliente_nombre" binding:"required"`
		ClienteTelefono   string `json:"cliente_telefono"`
		Equipo            string `json:"equipo" binding:"required"`
		Marca             string `json:"marca"`
		Modelo            string `json:"modelo"`
		NumSerie          string `json:"num_serie"`
		ProblemaReportado string `json:"problema_reportado" binding:"required"`
		Prioridad         string `json:"prioridad"`
		SedeID            int64  `json:"sede_id" binding:"required"`
		TecnicoID         *int64 `json:"tecnico_id"`
		FechaPromesa      string `json:"fecha_promesa"`
		Notas             string `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	if req.Prioridad == "" {
		req.Prioridad = "media"
	}

	// Generar número de orden
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo").Scan(&count)
	numeroOrden := fmt.Sprintf("OT-%d-%04d", time.Now().Year(), count+1)

	var fechaPromesa interface{}
	if req.FechaPromesa != "" {
		fechaPromesa = req.FechaPromesa
	} else {
		fechaPromesa = nil
	}

	result, err := database.DB.Exec(`
		INSERT INTO ordenes_trabajo (numero_orden, cliente_nombre, cliente_telefono, equipo, 
		                             marca, modelo, num_serie, problema_reportado, 
		                             estado, prioridad, fecha_ingreso, fecha_promesa, tecnico_id, 
		                             sede_id, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'recibido', ?, CURRENT_TIMESTAMP, ?, ?, ?, ?)`,
		numeroOrden, req.ClienteNombre, req.ClienteTelefono, req.Equipo,
		req.Marca, req.Modelo, req.NumSerie, req.ProblemaReportado,
		req.Prioridad, fechaPromesa, req.TecnicoID, req.SedeID, req.Notas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create orden: " + err.Error()})
		return
	}

	ordenID, _ := result.LastInsertId()

	// Registrar en trazabilidad
	database.DB.Exec(`
		INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
		VALUES (?, 'ingreso', 'Equipo ingresado al servicio técnico', ?)
	`, ordenID, userID)

	logAuditoria(c, "crear", "orden_trabajo", ordenID, "", numeroOrden)

	c.JSON(http.StatusCreated, gin.H{"id": ordenID, "numero_orden": numeroOrden})
}

// UpdateOrdenTrabajo actualiza una orden de trabajo
func UpdateOrdenTrabajo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var req struct {
		DiagnosticoTecnico string  `json:"diagnostico_tecnico"`
		SolucionAplicada   string  `json:"solucion_aplicada"`
		Estado             string  `json:"estado"`
		Prioridad          string  `json:"prioridad"`
		TecnicoID          *int64  `json:"tecnico_id"`
		CostoManoObra      float64 `json:"costo_mano_obra"`
		CostoRepuestos     float64 `json:"costo_repuestos"`
		Notas              string  `json:"notas"`
		FechaPromesa       string  `json:"fecha_promesa"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	// Obtener estado anterior
	var estadoAnterior string
	database.DB.QueryRow("SELECT estado FROM ordenes_trabajo WHERE id = ?", id).Scan(&estadoAnterior)

	var fechaPromesa interface{}
	if req.FechaPromesa != "" {
		fechaPromesa = req.FechaPromesa
	} else {
		fechaPromesa = nil
	}

	var fechaEntrega interface{}
	if req.Estado == "entregado" && estadoAnterior != "entregado" {
		fechaEntrega = time.Now()
	} else {
		fechaEntrega = nil
	}

	_, err = database.DB.Exec(`
		UPDATE ordenes_trabajo SET 
		    diagnostico_tecnico = ?, solucion_aplicada = ?, estado = ?, prioridad = ?, tecnico_id = ?,
		    costo_mano_obra = ?, costo_repuestos = ?, notas = ?, fecha_promesa = COALESCE(?, fecha_promesa),
		    fecha_entrega = COALESCE(?, fecha_entrega), updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		req.DiagnosticoTecnico, req.SolucionAplicada, req.Estado, req.Prioridad, req.TecnicoID,
		req.CostoManoObra, req.CostoRepuestos, req.Notas, fechaPromesa, fechaEntrega, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update orden"})
		return
	}

	// Registrar cambio de estado en trazabilidad
	if req.Estado != estadoAnterior {
		database.DB.Exec(`
			INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
			VALUES (?, ?, ?, ?)
		`, id, "cambio_estado", fmt.Sprintf("Estado cambiado de %s a %s", estadoAnterior, req.Estado), userID)
	}

	logAuditoria(c, "editar", "orden_trabajo", id, estadoAnterior, req.Estado)

	c.JSON(http.StatusOK, gin.H{"message": "Orden updated successfully"})
}

// AsignarTecnico asigna un técnico a una orden
func AsignarTecnico(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var req struct {
		TecnicoID int64 `json:"tecnico_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	// Obtener nombre del técnico
	var tecnicoNombre string
	database.DB.QueryRow("SELECT name FROM users WHERE id = ?", req.TecnicoID).Scan(&tecnicoNombre)

	_, err = database.DB.Exec("UPDATE ordenes_trabajo SET tecnico_id = ?, estado = CASE WHEN estado = 'recibido' THEN 'en_diagnostico' ELSE estado END, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.TecnicoID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign tecnico"})
		return
	}

	database.DB.Exec(`
		INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
		VALUES (?, 'asignacion', ?, ?)
	`, id, fmt.Sprintf("Técnico asignado: %s", tecnicoNombre), userID)

	logAuditoria(c, "asignar_tecnico", "orden_trabajo", id, "", tecnicoNombre)

	c.JSON(http.StatusOK, gin.H{"message": "Técnico assigned successfully"})
}

// AgregarInsumo agrega un insumo a una orden de trabajo
func AgregarInsumo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var req struct {
		InsumoID int64 `json:"insumo_id" binding:"required"`
		Cantidad int   `json:"cantidad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	// Verificar stock de insumo
	var stockActual int
	var insumoNombre string
	err = database.DB.QueryRow("SELECT nombre, stock FROM insumos WHERE id = ?", req.InsumoID).Scan(&insumoNombre, &stockActual)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}

	if stockActual < req.Cantidad {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stock insuficiente de %s (disponible: %d)", insumoNombre, stockActual)})
		return
	}

	tx, _ := database.DB.Begin()

	// Registrar uso de insumo
	_, err = tx.Exec(`
		INSERT INTO insumos_orden (orden_id, insumo_id, cantidad)
		VALUES (?, ?, ?)
	`, id, req.InsumoID, req.Cantidad)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add insumo"})
		return
	}

	// Descontar stock
	_, err = tx.Exec("UPDATE insumos SET stock = stock - ? WHERE id = ?", req.Cantidad, req.InsumoID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	tx.Commit()

	// Registrar en trazabilidad
	database.DB.Exec(`
		INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
		VALUES (?, 'insumo_usado', ?, ?)
	`, id, fmt.Sprintf("Insumo utilizado: %s x%d", insumoNombre, req.Cantidad), userID)

	c.JSON(http.StatusOK, gin.H{"message": "Insumo added successfully"})
}

// RegistrarTrazabilidad registra una entrada manual en la trazabilidad
func RegistrarTrazabilidad(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var req struct {
		Accion  string `json:"accion" binding:"required"`
		Detalle string `json:"detalle" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	_, err = database.DB.Exec(`
		INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
		VALUES (?, ?, ?, ?)
	`, id, req.Accion, req.Detalle, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register trazabilidad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trazabilidad registered successfully"})
}

// DeleteOrdenTrabajo elimina una orden de trabajo (solo si no está en proceso)
func DeleteOrdenTrabajo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orden ID"})
		return
	}

	var estado string
	err = database.DB.QueryRow("SELECT estado FROM ordenes_trabajo WHERE id = ?", id).Scan(&estado)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}

	if estado != "recibido" && estado != "cancelado" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden eliminar órdenes recibidas o canceladas"})
		return
	}

	database.DB.Exec("DELETE FROM trazabilidad WHERE orden_trabajo_id = ?", id)
	database.DB.Exec("DELETE FROM insumos_orden WHERE orden_id = ?", id)
	database.DB.Exec("DELETE FROM ordenes_trabajo WHERE id = ?", id)

	logAuditoria(c, "eliminar", "orden_trabajo", id, "", "")

	c.JSON(http.StatusOK, gin.H{"message": "Orden deleted successfully"})
}

// GetOrdenesStats obtiene estadísticas de órdenes de trabajo
func GetOrdenesStats(c *gin.Context) {
	type Stats struct {
		Total         int `json:"total"`
		Recibidos     int `json:"recibidos"`
		EnDiagnostico int `json:"en_diagnostico"`
		EnReparacion  int `json:"en_reparacion"`
		Terminados    int `json:"terminados"`
		Entregados    int `json:"entregados"`
		Urgentes      int `json:"urgentes"`
	}

	var stats Stats
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo").Scan(&stats.Total)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'recibido'").Scan(&stats.Recibidos)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'en_diagnostico'").Scan(&stats.EnDiagnostico)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'en_reparacion'").Scan(&stats.EnReparacion)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'terminado'").Scan(&stats.Terminados)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'entregado'").Scan(&stats.Entregados)
	database.DB.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE prioridad = 'urgente' AND estado NOT IN ('entregado', 'cancelado')").Scan(&stats.Urgentes)

	c.JSON(http.StatusOK, stats)
}

// GetTecnicos obtiene la lista de técnicos (usuarios con rol técnico)
// Este endpoint requiere ordenes.read en lugar de users.read
func GetTecnicos(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT u.id, u.name, u.email, u.sede_id, s.nombre as sede_nombre
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		LEFT JOIN sedes s ON u.sede_id = s.id
		WHERE r.name = 'tecnico'
		ORDER BY u.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tecnicos"})
		return
	}
	defer rows.Close()

	type Tecnico struct {
		ID         int64   `json:"id"`
		Name       string  `json:"name"`
		Email      string  `json:"email"`
		SedeID     *int64  `json:"sede_id"`
		SedeNombre *string `json:"sede_nombre"`
	}

	var tecnicos []Tecnico
	for rows.Next() {
		var t Tecnico
		var sedeID sql.NullInt64
		var sedeNombre sql.NullString
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &sedeID, &sedeNombre)
		if err != nil {
			continue
		}
		if sedeID.Valid {
			t.SedeID = &sedeID.Int64
		}
		if sedeNombre.Valid {
			t.SedeNombre = &sedeNombre.String
		}
		tecnicos = append(tecnicos, t)
	}

	c.JSON(http.StatusOK, tecnicos)
}
