package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCotizacionService() *services.CotizacionService {
	repo := repositories.NewCotizacionRepository(database.DB)
	return services.NewCotizacionService(repo)
}

// GetCotizaciones obtiene todas las cotizaciones
func GetCotizaciones(c *gin.Context) {
	estado := c.Query("estado")
	sedeID := c.Query("sede_id")

	items, err := getCotizacionService().ListCotizaciones(estado, sedeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cotizaciones"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetCotizacion obtiene una cotización por ID con sus items
func GetCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cotización ID"})
		return
	}

	cot, items, err := getCotizacionService().GetCotizacion(id)
	if err == services.ErrCotizacionNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cotización"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cotizacion": cot, "items": items})
}

// CreateCotizacion crea una nueva cotización
func CreateCotizacion(c *gin.Context) {
	var req struct {
		ClienteNombre   string  `json:"cliente_nombre" binding:"required"`
		ClienteTelefono string  `json:"cliente_telefono"`
		ClienteEmail    string  `json:"cliente_email"`
		ValidezDias     int     `json:"validez_dias"`
		Descuento       float64 `json:"descuento"`
		Notas           string  `json:"notas"`
		SedeID          int64   `json:"sede_id" binding:"required"`
		Items           []struct {
			ProductoID     int64   `json:"producto_id"`
			Cantidad       int     `json:"cantidad"`
			PrecioUnitario float64 `json:"precio_unitario"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	inItems := make([]services.CotizacionCreateItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		inItems = append(inItems, services.CotizacionCreateItemInput{ProductoID: it.ProductoID, Cantidad: it.Cantidad, PrecioUnitario: it.PrecioUnitario})
	}

	id, numero, total, err := getCotizacionService().CreateCotizacion(services.CreateCotizacionInput{
		ClienteNombre:   req.ClienteNombre,
		ClienteTelefono: req.ClienteTelefono,
		ClienteEmail:    req.ClienteEmail,
		ValidezDias:     req.ValidezDias,
		Descuento:       req.Descuento,
		Notas:           req.Notas,
		SedeID:          req.SedeID,
		UsuarioID:       userID,
		Items:           inItems,
	})
	if err == services.ErrCotizacionSinItems {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe incluir al menos un item"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cotización: " + err.Error()})
		return
	}

	logAuditoria(c, "crear", "cotizacion", id, "", numero)
	c.JSON(http.StatusCreated, gin.H{"id": id, "numero_cotizacion": numero, "total": total})
}

// UpdateCotizacionEstado actualiza el estado de una cotización
func UpdateCotizacionEstado(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cotización ID"})
		return
	}

	var req struct {
		Estado string `json:"estado" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estadoAnterior, err := getCotizacionService().UpdateEstado(id, req.Estado)
	if err == services.ErrEstadoCotizacionBad {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cotización"})
		return
	}

	logAuditoria(c, "editar", "cotizacion", id, estadoAnterior, req.Estado)
	c.JSON(http.StatusOK, gin.H{"message": "Cotización updated successfully"})
}

// DeleteCotizacion elimina una cotización
func DeleteCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cotización ID"})
		return
	}

	if err := getCotizacionService().DeleteCotizacion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cotización"})
		return
	}

	logAuditoria(c, "eliminar", "cotizacion", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Cotización deleted successfully"})
}

// ConvertirCotizacionAVenta convierte una cotización aprobada a venta
func ConvertirCotizacionAVenta(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cotización ID"})
		return
	}

	userID, _ := c.Get("userid")
	ventaID, numeroVenta, itemsJSON, err := getCotizacionService().ConvertirAVenta(id, userID)
	if err == services.ErrCotizacionNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización not found"})
		return
	}
	if err == services.ErrCotizacionNoAprobada {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden convertir cotizaciones aprobadas"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create venta"})
		return
	}

	logAuditoria(c, "convertir_venta", "cotizacion", id, "", itemsJSON)
	c.JSON(http.StatusCreated, gin.H{"venta_id": ventaID, "numero_venta": numeroVenta})
}

// GenerarPDFCotizacion genera un PDF de la cotización (retorna datos para frontend)
func GenerarPDFCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cotización ID"})
		return
	}

	cot, items, err := getCotizacionService().GetPDFData(id)
	if err == services.ErrCotizacionNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cotización"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cotizacion": cot, "items": items})
}
