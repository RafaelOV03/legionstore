package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"smartech/backend/services"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch proveedores", "details": err.Error()})
		return
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

	proveedor, deudas, totalDeuda, err := getProveedorService().GetProveedor(id)
	if err == services.ErrProveedorNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch proveedor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"proveedor": proveedor, "deudas": deudas, "total_deuda": totalDeuda})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proveedor"})
		return
	}

	logAuditoria(c, "crear", "proveedor", id, "", req.Nombre)
	c.JSON(http.StatusCreated, gin.H{"id": id})
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

	err = getProveedorService().DeleteProveedor(id)
	if err == services.ErrProveedorHasPendingDeudas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar un proveedor con deudas pendientes"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete proveedor"})
		return
	}

	logAuditoria(c, "eliminar", "proveedor", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Proveedor deleted successfully"})
}

// GetDeudas obtiene todas las deudas pendientes de proveedores
func GetDeudas(c *gin.Context) {
	estado := c.Query("estado")
	proveedorID := c.Query("proveedor_id")

	deudas, err := getProveedorService().ListDeudas(estado, proveedorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deudas"})
		return
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

	deudaID, err := getProveedorService().CreateDeuda(req.ProveedorID, req.NumeroFactura, req.MontoTotal, fechaVenc, req.Notas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deuda"})
		return
	}

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

	nuevoSaldo, estado, err := getProveedorService().RegistrarPago(id, req.Monto, req.MetodoPago, req.NumeroReferencia, userID)
	if err == services.ErrDeudaNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deuda not found"})
		return
	}
	if err == services.ErrPagoExcedeSaldo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El monto excede el saldo pendiente"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register pago"})
		return
	}

	logAuditoria(c, "pago", "deuda_proveedor", id, "", strconv.FormatFloat(req.Monto, 'f', 2, 64))
	c.JSON(http.StatusOK, gin.H{"message": "Pago registered successfully", "nuevo_saldo": nuevoSaldo, "estado": estado})
}

// GetPagosDeuda obtiene los pagos de una deuda
func GetPagosDeuda(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deuda ID"})
		return
	}

	pagos, err := getProveedorService().ListPagosDeuda(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pagos"})
		return
	}

	c.JSON(http.StatusOK, pagos)
}

// GetResumenDeudas obtiene un resumen de deudas por proveedor
func GetResumenDeudas(c *gin.Context) {
	resumen, totalDeuda, pendientes, vencidas, pagadas, err := getProveedorService().ResumenDeudas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resumen"})
		return
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
