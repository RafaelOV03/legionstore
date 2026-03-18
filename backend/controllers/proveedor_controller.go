package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetProveedores obtiene todos los proveedores
func GetProveedores(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, nombre, ruc, direccion, telefono, email, contacto, activo
		FROM proveedores ORDER BY nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch proveedores", "details": err.Error()})
		return
	}
	defer rows.Close()

	var proveedores []models.Proveedor
	for rows.Next() {
		var p models.Proveedor
		var activo int
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.RucNit, &p.Direccion,
			&p.Telefono, &p.Email, &p.Contacto, &activo)
		if err != nil {
			continue
		}
		p.Activo = activo == 1
		proveedores = append(proveedores, p)
	}

	c.JSON(http.StatusOK, proveedores)
}

// GetProveedor obtiene un proveedor por ID con sus deudas
func GetProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proveedor ID"})
		return
	}

	var p models.Proveedor
	var activo int
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, nombre, ruc, direccion, telefono, email, contacto, activo
		FROM proveedores WHERE id = ?`, id).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.RucNit, &p.Direccion,
			&p.Telefono, &p.Email, &p.Contacto, &activo)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch proveedor"})
		return
	}
	p.Activo = activo == 1

	// Obtener deudas pendientes
	deudaRows, _ := database.DB.Query(`
		SELECT id, created_at, updated_at, num_factura, monto, monto_pagado, fecha_vence, estado, descripcion
		FROM deudas_proveedores WHERE proveedor_id = ? AND estado != 'pagada'
		ORDER BY fecha_vence
	`, id)
	defer deudaRows.Close()

	type DeudaSimple struct {
		ID               int64      `json:"id"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        time.Time  `json:"updated_at"`
		NumeroFactura    string     `json:"numero_factura"`
		MontoTotal       float64    `json:"monto_total"`
		MontoPagado      float64    `json:"monto_pagado"`
		FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
		Estado           string     `json:"estado"`
		Notas            string     `json:"notas"`
		ProveedorID      int64      `json:"proveedor_id"`
	}

	var deudas []DeudaSimple
	for deudaRows.Next() {
		var d DeudaSimple
		var fechaVenc sql.NullTime
		deudaRows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.NumeroFactura, &d.MontoTotal,
			&d.MontoPagado, &fechaVenc, &d.Estado, &d.Notas)
		if fechaVenc.Valid {
			d.FechaVencimiento = &fechaVenc.Time
		}
		d.ProveedorID = id
		deudas = append(deudas, d)
	}

	// Calcular total deuda
	var totalDeuda float64
	for _, d := range deudas {
		totalDeuda += (d.MontoTotal - d.MontoPagado)
	}

	c.JSON(http.StatusOK, gin.H{"proveedor": p, "deudas": deudas, "total_deuda": totalDeuda})
}

// CreateProveedor crea un nuevo proveedor
func CreateProveedor(c *gin.Context) {
	var req struct {
		Nombre    string `json:"nombre" binding:"required"`
		RucNit    string `json:"ruc_nit"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
		Email     string `json:"email"`
		Contacto  string `json:"contacto"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO proveedores (nombre, ruc, direccion, telefono, email, contacto, activo)
		VALUES (?, ?, ?, ?, ?, ?, 1)`,
		req.Nombre, req.RucNit, req.Direccion, req.Telefono, req.Email, req.Contacto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proveedor"})
		return
	}

	proveedorID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "proveedor", proveedorID, "", req.Nombre)

	c.JSON(http.StatusCreated, gin.H{"id": proveedorID})
}

// UpdateProveedor actualiza un proveedor
func UpdateProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proveedor ID"})
		return
	}

	var req struct {
		Nombre    string `json:"nombre"`
		RucNit    string `json:"ruc_nit"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
		Email     string `json:"email"`
		Contacto  string `json:"contacto"`
		Activo    bool   `json:"activo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activo := 0
	if req.Activo {
		activo = 1
	}

	_, err = database.DB.Exec(`
		UPDATE proveedores SET nombre = ?, ruc = ?, direccion = ?, telefono = ?, 
		                       email = ?, contacto = ?, activo = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		req.Nombre, req.RucNit, req.Direccion, req.Telefono, req.Email, req.Contacto, activo, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update proveedor"})
		return
	}

	logAuditoria(c, "editar", "proveedor", id, "", req.Nombre)
	c.JSON(http.StatusOK, gin.H{"message": "Proveedor updated successfully"})
}

// DeleteProveedor elimina un proveedor
func DeleteProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proveedor ID"})
		return
	}

	// Verificar si tiene deudas pendientes
	var deudasPendientes int
	database.DB.QueryRow("SELECT COUNT(*) FROM deudas_proveedores WHERE proveedor_id = ? AND estado != 'pagada'", id).Scan(&deudasPendientes)
	if deudasPendientes > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar un proveedor con deudas pendientes"})
		return
	}

	database.DB.Exec("DELETE FROM pagos_proveedores WHERE deuda_id IN (SELECT id FROM deudas_proveedores WHERE proveedor_id = ?)", id)
	database.DB.Exec("DELETE FROM deudas_proveedores WHERE proveedor_id = ?", id)
	database.DB.Exec("DELETE FROM proveedores WHERE id = ?", id)

	logAuditoria(c, "eliminar", "proveedor", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Proveedor deleted successfully"})
}

// GetDeudas obtiene todas las deudas pendientes de proveedores
func GetDeudas(c *gin.Context) {
	estado := c.Query("estado")
	proveedorID := c.Query("proveedor_id")

	query := `
		SELECT d.id, d.created_at, d.updated_at, d.proveedor_id, d.num_factura, d.monto,
		       d.monto_pagado, d.fecha_vence, d.estado, d.descripcion,
		       p.nombre as proveedor_nombre
		FROM deudas_proveedores d
		INNER JOIN proveedores p ON d.proveedor_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}

	if estado != "" {
		query += " AND d.estado = ?"
		args = append(args, estado)
	} else {
		query += " AND d.estado != 'pagada'"
	}
	if proveedorID != "" {
		query += " AND d.proveedor_id = ?"
		args = append(args, proveedorID)
	}

	query += " ORDER BY d.fecha_vence"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deudas"})
		return
	}
	defer rows.Close()

	type DeudaView struct {
		ID               int64      `json:"id"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        time.Time  `json:"updated_at"`
		ProveedorID      int64      `json:"proveedor_id"`
		NumeroFactura    string     `json:"numero_factura"`
		MontoTotal       float64    `json:"monto_total"`
		MontoPagado      float64    `json:"monto_pagado"`
		MontoPendiente   float64    `json:"monto_pendiente"`
		FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
		Estado           string     `json:"estado"`
		Notas            string     `json:"notas"`
		Proveedor        struct {
			Nombre string `json:"nombre"`
		} `json:"proveedor"`
	}

	var deudas []DeudaView
	for rows.Next() {
		var d DeudaView
		var fechaVenc sql.NullTime
		var provNombre string
		err := rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ProveedorID, &d.NumeroFactura,
			&d.MontoTotal, &d.MontoPagado, &fechaVenc, &d.Estado, &d.Notas,
			&provNombre)
		if err != nil {
			continue
		}
		if fechaVenc.Valid {
			d.FechaVencimiento = &fechaVenc.Time
		}
		d.MontoPendiente = d.MontoTotal - d.MontoPagado
		d.Proveedor.Nombre = provNombre
		deudas = append(deudas, d)
	}

	c.JSON(http.StatusOK, deudas)
}

// CreateDeuda crea una nueva deuda de proveedor
func CreateDeuda(c *gin.Context) {
	var req struct {
		ProveedorID      int64   `json:"proveedor_id" binding:"required"`
		NumeroFactura    string  `json:"numero_factura" binding:"required"`
		MontoTotal       float64 `json:"monto_total" binding:"required"`
		FechaFactura     string  `json:"fecha_factura" binding:"required"`
		FechaVencimiento string  `json:"fecha_vencimiento"`
		Notas            string  `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fechaVenc interface{}
	if req.FechaVencimiento != "" {
		fechaVenc = req.FechaVencimiento
	} else {
		fechaVenc = nil
	}

	result, err := database.DB.Exec(`
		INSERT INTO deudas_proveedores (proveedor_id, num_factura, monto, monto_pagado, 
		                                fecha_vence, estado, descripcion)
		VALUES (?, ?, ?, 0, ?, 'pendiente', ?)`,
		req.ProveedorID, req.NumeroFactura, req.MontoTotal, fechaVenc, req.Notas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deuda"})
		return
	}

	deudaID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "deuda_proveedor", deudaID, "", req.NumeroFactura)

	c.JSON(http.StatusCreated, gin.H{"id": deudaID})
}

// RegistrarPago registra un pago a una deuda
func RegistrarPago(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deuda ID"})
		return
	}

	var req struct {
		Monto            float64 `json:"monto" binding:"required"`
		MetodoPago       string  `json:"metodo_pago"`
		NumeroReferencia string  `json:"numero_referencia"`
		Notas            string  `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userid")

	// Verificar deuda
	var montoTotal, montoPagado float64
	err = database.DB.QueryRow("SELECT monto, monto_pagado FROM deudas_proveedores WHERE id = ?", id).Scan(&montoTotal, &montoPagado)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deuda not found"})
		return
	}

	saldoPendiente := montoTotal - montoPagado
	if req.Monto > saldoPendiente {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El monto excede el saldo pendiente"})
		return
	}

	tx, _ := database.DB.Begin()

	// Registrar pago
	_, err = tx.Exec(`
		INSERT INTO pagos_proveedores (deuda_id, monto, metodo, referencia, usuario_id)
		VALUES (?, ?, ?, ?, ?)`,
		id, req.Monto, req.MetodoPago, req.NumeroReferencia, userID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register pago"})
		return
	}

	// Actualizar monto pagado
	nuevoMontoPagado := montoPagado + req.Monto
	nuevoEstado := "pendiente"
	if nuevoMontoPagado >= montoTotal {
		nuevoEstado = "pagada"
	} else if nuevoMontoPagado > 0 {
		nuevoEstado = "parcial"
	}

	_, err = tx.Exec(`
		UPDATE deudas_proveedores SET monto_pagado = ?, estado = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		nuevoMontoPagado, nuevoEstado, id)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deuda"})
		return
	}

	tx.Commit()

	logAuditoria(c, "pago", "deuda_proveedor", id, "", strconv.FormatFloat(req.Monto, 'f', 2, 64))

	c.JSON(http.StatusOK, gin.H{
		"message":     "Pago registered successfully",
		"nuevo_saldo": montoTotal - nuevoMontoPagado,
		"estado":      nuevoEstado,
	})
}

// GetPagosDeuda obtiene los pagos de una deuda
func GetPagosDeuda(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deuda ID"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT pp.id, pp.created_at, pp.monto, pp.metodo, pp.referencia, u.name
		FROM pagos_proveedores pp
		INNER JOIN users u ON pp.usuario_id = u.id
		WHERE pp.deuda_id = ?
		ORDER BY pp.created_at DESC
	`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pagos"})
		return
	}
	defer rows.Close()

	type PagoView struct {
		ID               int64   `json:"id"`
		CreatedAt        string  `json:"created_at"`
		Monto            float64 `json:"monto"`
		FechaPago        string  `json:"fecha_pago"`
		MetodoPago       string  `json:"metodo_pago"`
		NumeroReferencia string  `json:"numero_referencia"`
		UsuarioNombre    string  `json:"usuario_nombre"`
	}

	var pagos []PagoView
	for rows.Next() {
		var p PagoView
		rows.Scan(&p.ID, &p.CreatedAt, &p.Monto, &p.MetodoPago, &p.NumeroReferencia, &p.UsuarioNombre)
		p.FechaPago = p.CreatedAt // Use created_at as fecha_pago for frontend
		pagos = append(pagos, p)
	}

	c.JSON(http.StatusOK, pagos)
}

// GetResumenDeudas obtiene un resumen de deudas por proveedor
func GetResumenDeudas(c *gin.Context) {
	// Obtener estadísticas generales
	var totalDeuda float64
	var pendientes, vencidas, pagadas int

	database.DB.QueryRow(`
		SELECT COALESCE(SUM(monto - monto_pagado), 0) FROM deudas_proveedores WHERE estado != 'pagada'
	`).Scan(&totalDeuda)

	database.DB.QueryRow(`
		SELECT COUNT(*) FROM deudas_proveedores WHERE estado = 'pendiente'
	`).Scan(&pendientes)

	database.DB.QueryRow(`
		SELECT COUNT(*) FROM deudas_proveedores WHERE estado != 'pagada' AND fecha_vence < DATE('now')
	`).Scan(&vencidas)

	database.DB.QueryRow(`
		SELECT COUNT(*) FROM deudas_proveedores WHERE estado = 'pagada'
	`).Scan(&pagadas)

	// Obtener resumen por proveedor
	rows, err := database.DB.Query(`
		SELECT p.id, p.nombre, 
		       COUNT(d.id) as num_facturas,
		       COALESCE(SUM(d.monto - d.monto_pagado), 0) as saldo_total,
		       MIN(d.fecha_vence) as proxima_fecha
		FROM proveedores p
		LEFT JOIN deudas_proveedores d ON p.id = d.proveedor_id AND d.estado != 'pagada'
		WHERE p.activo = 1
		GROUP BY p.id, p.nombre
		HAVING saldo_total > 0
		ORDER BY saldo_total DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resumen"})
		return
	}
	defer rows.Close()

	type ResumenProveedor struct {
		ProveedorID  int64   `json:"proveedor_id"`
		Nombre       string  `json:"nombre"`
		NumFacturas  int     `json:"num_facturas"`
		SaldoTotal   float64 `json:"saldo_total"`
		ProximaFecha *string `json:"proxima_fecha"`
	}

	var resumen []ResumenProveedor
	for rows.Next() {
		var r ResumenProveedor
		var fecha sql.NullString
		rows.Scan(&r.ProveedorID, &r.Nombre, &r.NumFacturas, &r.SaldoTotal, &fecha)
		if fecha.Valid {
			r.ProximaFecha = &fecha.String
		}
		resumen = append(resumen, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"proveedores":   resumen,
		"total_general": totalDeuda,
		"total_deuda":   totalDeuda,
		"pendientes":    pendientes,
		"vencidas":      vencidas,
		"pagadas":       pagadas,
	})
}
