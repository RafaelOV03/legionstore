package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getTraspasoService() *services.TraspasoService {
	repo := repositories.NewTraspasoRepository(database.DB)
	return services.NewTraspasoService(repo)
}

// GetTraspasos obtiene todos los traspasos
func GetTraspasos(c *gin.Context) {
	estado := c.Query("estado")
	sedeOrigenID := c.Query("sede_origen_id")
	sedeDestinoID := c.Query("sede_destino_id")

	traspasos, err := getTraspasoService().ListTraspasos(estado, sedeOrigenID, sedeDestinoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traspasos"})
		return
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

	t, items, err := getTraspasoService().GetTraspaso(id)
	if err == services.ErrTraspasoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traspaso"})
		return
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

	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	items := make([]services.TraspasoCreateItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, services.TraspasoCreateItemInput{ProductoID: item.ProductoID, Cantidad: item.Cantidad})
	}

	traspasoID, numeroTraspaso, err := getTraspasoService().CreateTraspaso(services.CreateTraspasoInput{
		SedeOrigenID:  req.SedeOrigenID,
		SedeDestinoID: req.SedeDestinoID,
		Notas:         req.Notas,
		Items:         items,
		UsuarioID:     userID,
	})
	if err == services.ErrTraspasoSedesIguales {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Las sedes de origen y destino deben ser diferentes"})
		return
	}
	if err == services.ErrTraspasoSinItems {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe incluir al menos un item"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err = getTraspasoService().EnviarTraspaso(id)
	if err == services.ErrTraspasoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err == services.ErrTraspasoSoloPendienteEnviar {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden enviar traspasos pendientes"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update traspaso"})
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
	_ = c.ShouldBindJSON(&req)

	userID, _ := c.Get("userid")
	items := make([]services.TraspasoRecibirItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, services.TraspasoRecibirItemInput{ItemID: item.ItemID, CantidadRecibida: item.CantidadRecibida})
	}

	err = getTraspasoService().RecibirTraspaso(services.RecibirTraspasoInput{
		ID:        id,
		UsuarioID: userID,
		Notas:     req.Notas,
		Items:     items,
	})
	if err == services.ErrTraspasoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err == services.ErrTraspasoSoloEnviadoRecibir {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden recibir traspasos enviados"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update traspaso"})
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

	err = getTraspasoService().CancelarTraspaso(id)
	if err == services.ErrTraspasoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err == services.ErrTraspasoSoloPendienteCancelar {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden cancelar traspasos pendientes"})
		return
	}
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

	err = getTraspasoService().DeleteTraspaso(id)
	if err == services.ErrTraspasoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traspaso not found"})
		return
	}
	if err == services.ErrTraspasoSoloCanceladoEliminar {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden eliminar traspasos cancelados"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete traspaso"})
		return
	}

	logAuditoria(c, "eliminar", "traspaso", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Traspaso deleted successfully"})
}
