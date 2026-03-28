package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAuditoriaService() *services.AuditoriaService {
	repo := repositories.NewAuditoriaRepository(database.DB)
	return services.NewAuditoriaService(repo)
}

// GetLogs obtiene los logs de auditoría
func GetLogs(c *gin.Context) {
	logs, err := getAuditoriaService().ListLogs(services.GetLogsInput{
		Accion:     c.Query("accion"),
		Entidad:    c.Query("entidad"),
		UsuarioID:  c.Query("usuario_id"),
		FechaDesde: c.Query("fecha_desde"),
		FechaHasta: c.Query("fecha_hasta"),
		Limit:      c.DefaultQuery("limit", "100"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetLogStats obtiene estadísticas de acciones
func GetLogStats(c *gin.Context) {
	stats, err := getAuditoriaService().LogStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetReporteGanancias obtiene el reporte de ganancias
func GetReporteGanancias(c *gin.Context) {
	reporte, err := getAuditoriaService().ReporteGanancias(services.GananciasInput{
		FechaDesde: c.DefaultQuery("fecha_desde", ""),
		FechaHasta: c.DefaultQuery("fecha_hasta", ""),
		SedeID:     c.Query("sede_id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reporte de ganancias"})
		return
	}

	c.JSON(http.StatusOK, reporte)
}

// GetSegmentaciones obtiene las segmentaciones de clientes
func GetSegmentaciones(c *gin.Context) {
	segmentaciones, err := getAuditoriaService().ListSegmentaciones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch segmentaciones"})
		return
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

	segID, err := getAuditoriaService().CreateSegmentacion(services.CreateSegmentacionInput{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Criterios:   req.Criterios,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create segmentación"})
		return
	}

	logAuditoria(c, "crear", "segmentacion", segID, "", req.Nombre)

	c.JSON(http.StatusCreated, gin.H{"id": segID})
}

// GetPromociones obtiene las promociones
func GetPromociones(c *gin.Context) {
	promociones, err := getAuditoriaService().ListPromociones(c.Query("activas") == "true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promociones"})
		return
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

	promoID, err := getAuditoriaService().CreatePromocion(services.CreatePromocionInput{
		Nombre:         req.Nombre,
		Descripcion:    req.Descripcion,
		Tipo:           req.Tipo,
		Valor:          req.Valor,
		FechaInicio:    req.FechaInicio,
		FechaFin:       req.FechaFin,
		ProductosIDs:   req.ProductosIDs,
		SegmentacionID: req.SegmentacionID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promoción"})
		return
	}

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

	err = getAuditoriaService().UpdatePromocion(id, services.UpdatePromocionInput{
		Nombre:         req.Nombre,
		Descripcion:    req.Descripcion,
		Tipo:           req.Tipo,
		Valor:          req.Valor,
		FechaInicio:    req.FechaInicio,
		FechaFin:       req.FechaFin,
		ProductosIDs:   req.ProductosIDs,
		SegmentacionID: req.SegmentacionID,
		Activa:         req.Activa,
	})
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

	err = getAuditoriaService().DeletePromocion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete promoción"})
		return
	}

	logAuditoria(c, "eliminar", "promocion", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Promoción deleted successfully"})
}
