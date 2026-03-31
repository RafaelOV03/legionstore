package controllers

import (
<<<<<<< HEAD
	"database/sql"
=======
	"net/http"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/models"
<<<<<<< HEAD
	"smartech/backend/validation"
=======
	"smartech/backend/repositories"
	"smartech/backend/services"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProveedorService() *services.ProveedorService {
	repo := repositories.NewProveedorRepository(database.DB)
	return services.NewProveedorService(repo)
}

// GetProveedores obtiene todos los proveedores
func GetProveedores(c *gin.Context) {
	proveedores, err := getProveedorService().ListProveedores()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch providers", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
<<<<<<< HEAD
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

	c.JSON(200, proveedores)
=======
	c.JSON(http.StatusOK, proveedores)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetProveedor obtiene un proveedor por ID con sus deudas
func GetProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid proveedor ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	var p models.Proveedor
	var activo int
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, nombre, ruc, direccion, telefono, email, contacto, activo
		FROM proveedores WHERE id = ?`, id).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.RucNit, &p.Direccion,
			&p.Telefono, &p.Email, &p.Contacto, &activo)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("Proveedor", id)
		c.JSON(apiErr.Code, apiErr)
=======
	proveedor, deudas, totalDeuda, err := getProveedorService().GetProveedor(id)
	if err == services.ErrProveedorNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch provider", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
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

	c.JSON(200, gin.H{"proveedor": p, "deudas": deudas, "total_deuda": totalDeuda})
=======
	c.JSON(http.StatusOK, gin.H{"proveedor": proveedor, "deudas": deudas, "total_deuda": totalDeuda})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// CreateProveedor crea un nuevo proveedor
func CreateProveedor(c *gin.Context) {
	var req struct {
		Nombre    string `json:"nombre" validate:"required,min=3"`
		RucNit    string `json:"ruc_nit"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
		Email     string `json:"email"`
		Contacto  string `json:"contacto"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	id, err := getProveedorService().CreateProveedor(models.Proveedor{
		Nombre:    req.Nombre,
		RucNit:    req.RucNit,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
		Email:     req.Email,
		Contacto:  req.Contacto,
		Activo:    true,
	})
	if err != nil {
		apiErr := errors.NewDatabaseError("Insert provider", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	proveedorID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "proveedor", proveedorID, "", req.Nombre)

	c.JSON(201, gin.H{"id": proveedorID})
=======
	logAuditoria(c, "crear", "proveedor", id, "", req.Nombre)
	c.JSON(http.StatusCreated, gin.H{"id": id})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// UpdateProveedor actualiza un proveedor
func UpdateProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid proveedor ID")
		c.JSON(apiErr.Code, apiErr)
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
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var proveedor models.Proveedor
	err = database.DB.QueryRow("SELECT id FROM proveedores WHERE id = ?", id).Scan(&proveedor.ID)
	if err != nil {
		apiErr := errors.NewNotFound("Proveedor", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	err = getProveedorService().UpdateProveedor(id, models.Proveedor{
		Nombre:    req.Nombre,
		RucNit:    req.RucNit,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
		Email:     req.Email,
		Contacto:  req.Contacto,
		Activo:    req.Activo,
	})
	if err != nil {
		apiErr := errors.NewDatabaseError("Update provider", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "editar", "proveedor", id, "", req.Nombre)
	c.JSON(200, gin.H{"message": "Proveedor updated successfully"})
}

// DeleteProveedor elimina un proveedor
func DeleteProveedor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid proveedor ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar si existe
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM proveedores WHERE id = ?", id).Scan(&exists)
	if err != nil {
		apiErr := errors.NewNotFound("Proveedor", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Verificar si tiene deudas pendientes
	var deudasPendientes int
	database.DB.QueryRow("SELECT COUNT(*) FROM deudas_proveedores WHERE proveedor_id = ? AND estado != 'pagada'", id).Scan(&deudasPendientes)
	if deudasPendientes > 0 {
		apiErr := errors.NewConflict("Provider has pending debts")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("DELETE FROM pagos_proveedores WHERE deuda_id IN (SELECT id FROM deudas_proveedores WHERE proveedor_id = ?)", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete provider payments", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("DELETE FROM deudas_proveedores WHERE proveedor_id = ?", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete provider debts", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("DELETE FROM proveedores WHERE id = ?", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete provider", err)
		c.JSON(apiErr.Code, apiErr)
=======
	err = getProveedorService().DeleteProveedor(id)
	if err == services.ErrProveedorHasPendingDeudas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar un proveedor con deudas pendientes"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete proveedor"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}

	logAuditoria(c, "eliminar", "proveedor", id, "", "")
	c.JSON(200, gin.H{"message": "Proveedor deleted successfully"})
}

// GetDeudas obtiene todas las deudas pendientes de proveedores
func GetDeudas(c *gin.Context) {
	estado := c.Query("estado")
	proveedorID := c.Query("proveedor_id")

	deudas, err := getProveedorService().ListDeudas(estado, proveedorID)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch provider debts", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, deudas)
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
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var fechaVenc interface{}
	if req.FechaVencimiento != "" {
		fechaVenc = req.FechaVencimiento
	} else {
		fechaVenc = nil
	}

	deudaID, err := getProveedorService().CreateDeuda(req.ProveedorID, req.NumeroFactura, req.MontoTotal, fechaVenc, req.Notas)
	if err != nil {
		apiErr := errors.NewDatabaseError("Create provider debt", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "crear", "deuda_proveedor", deudaID, "", req.NumeroFactura)
<<<<<<< HEAD

	c.JSON(201, gin.H{"id": deudaID})
=======
	c.JSON(http.StatusCreated, gin.H{"id": deudaID})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// RegistrarPago registra un pago a una deuda
func RegistrarPago(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid deuda ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var req struct {
		Monto            float64 `json:"monto" binding:"required"`
		MetodoPago       string  `json:"metodo_pago"`
		NumeroReferencia string  `json:"numero_referencia"`
		Notas            string  `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	userID, _ := c.Get("userid")

<<<<<<< HEAD
	// Verificar deuda
	var montoTotal, montoPagado float64
	err = database.DB.QueryRow("SELECT monto, monto_pagado FROM deudas_proveedores WHERE id = ?", id).Scan(&montoTotal, &montoPagado)
	if err != nil {
		apiErr := errors.NewNotFound("Deuda", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	saldoPendiente := montoTotal - montoPagado
	if req.Monto > saldoPendiente {
		apiErr := errors.NewBadRequest("El monto excede el saldo pendiente")
		c.JSON(apiErr.Code, apiErr)
=======
	nuevoSaldo, estado, err := getProveedorService().RegistrarPago(id, req.Monto, req.MetodoPago, req.NumeroReferencia, userID)
	if err == services.ErrDeudaNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deuda not found"})
		return
	}
	if err == services.ErrPagoExcedeSaldo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El monto excede el saldo pendiente"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
<<<<<<< HEAD
		tx.Rollback()
		apiErr := errors.NewDatabaseError("Register payment", err)
		c.JSON(apiErr.Code, apiErr)
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
		apiErr := errors.NewDatabaseError("Update debt", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	tx.Commit()

	logAuditoria(c, "pago", "deuda_proveedor", id, "", strconv.FormatFloat(req.Monto, 'f', 2, 64))

	c.JSON(200, gin.H{
		"message":     "Pago registered successfully",
		"nuevo_saldo": montoTotal - nuevoMontoPagado,
		"estado":      nuevoEstado,
	})
=======
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register pago"})
		return
	}

	logAuditoria(c, "pago", "deuda_proveedor", id, "", strconv.FormatFloat(req.Monto, 'f', 2, 64))
	c.JSON(http.StatusOK, gin.H{"message": "Pago registered successfully", "nuevo_saldo": nuevoSaldo, "estado": estado})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetPagosDeuda obtiene los pagos de una deuda
func GetPagosDeuda(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid deuda ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	pagos, err := getProveedorService().ListPagosDeuda(id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch provider payments", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, pagos)
}

// GetResumenDeudas obtiene un resumen de deudas por proveedor
func GetResumenDeudas(c *gin.Context) {
	resumen, totalDeuda, pendientes, vencidas, pagadas, err := getProveedorService().ResumenDeudas()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch debt summary", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, gin.H{
		"proveedores":   resumen,
		"total_general": totalDeuda,
		"total_deuda":   totalDeuda,
		"pendientes":    pendientes,
		"vencidas":      vencidas,
		"pagadas":       pagadas,
	})
}
