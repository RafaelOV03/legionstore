package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getRMAService() *services.RMAService {
	repo := repositories.NewRMARepository(database.DB)
	return services.NewRMAService(repo)
}

// GetRMAs obtiene todas las RMAs
func GetRMAs(c *gin.Context) {
	estado := c.Query("estado")
	sedeID := c.Query("sede_id")

	rmas, err := getRMAService().ListRMAs(estado, sedeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMAs"})
		return
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

	rma, historial, err := getRMAService().GetRMA(id)
	if err == services.ErrRMANotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "RMA not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMA"})
		return
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

	userID, _ := c.Get("userid")
	rmaID, numeroRMA, err := getRMAService().CreateRMA(services.CreateRMAInput{
		ProductoID:       req.ProductoID,
		ClienteNombre:    req.ClienteNombre,
		ClienteTelefono:  req.ClienteTelefono,
		ClienteEmail:     req.ClienteEmail,
		NumSerie:         req.NumSerie,
		FechaCompra:      req.FechaCompra,
		MotivoDevolucion: req.MotivoDevolucion,
		SedeID:           req.SedeID,
		Notas:            req.Notas,
	}, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create RMA: " + err.Error()})
		return
	}

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

	estadoAnterior, err := getRMAService().UpdateRMA(id, services.UpdateRMAInput{
		Diagnostico: req.Diagnostico,
		Estado:      req.Estado,
		Solucion:    req.Solucion,
		Notas:       req.Notas,
		Comentario:  req.Comentario,
	}, userID)
	if err == services.ErrRMANotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "RMA not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update RMA"})
		return
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

	err = getRMAService().DeleteRMA(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete RMA"})
		return
	}

	logAuditoria(c, "eliminar", "rma", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "RMA deleted successfully"})
}

// GetRMAStats obtiene estadísticas de RMAs
func GetRMAStats(c *gin.Context) {
	stats, err := getRMAService().Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMA stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
