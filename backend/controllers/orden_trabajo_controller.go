package controllers

import (
	"fmt"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getOrdenTrabajoService() *services.OrdenTrabajoService {
	repo := repositories.NewOrdenTrabajoRepository(database.DB)
	return services.NewOrdenTrabajoService(repo)
}

// GetOrdenesTrabajo obtiene todas las órdenes de trabajo
func GetOrdenesTrabajo(c *gin.Context) {
	estado := c.Query("estado")
	prioridad := c.Query("prioridad")
	tecnicoID := c.Query("tecnico_id")
	sedeID := c.Query("sede_id")

	ordenes, err := getOrdenTrabajoService().ListOrdenesTrabajo(estado, prioridad, tecnicoID, sedeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch órdenes de trabajo"})
		return
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

	o, insumos, trazabilidad, err := getOrdenTrabajoService().GetOrdenTrabajo(id)
	if err == services.ErrOrdenTrabajoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orden"})
		return
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

	ordenID, numeroOrden, err := getOrdenTrabajoService().CreateOrdenTrabajo(services.CreateOrdenTrabajoInput{
		ClienteNombre:     req.ClienteNombre,
		ClienteTelefono:   req.ClienteTelefono,
		Equipo:            req.Equipo,
		Marca:             req.Marca,
		Modelo:            req.Modelo,
		NumSerie:          req.NumSerie,
		ProblemaReportado: req.ProblemaReportado,
		Prioridad:         req.Prioridad,
		SedeID:            req.SedeID,
		TecnicoID:         req.TecnicoID,
		FechaPromesa:      req.FechaPromesa,
		Notas:             req.Notas,
		UsuarioID:         userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create orden: " + err.Error()})
		return
	}

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
	estadoAnterior, estadoNuevo, err := getOrdenTrabajoService().UpdateOrdenTrabajo(id, services.UpdateOrdenTrabajoInput{
		DiagnosticoTecnico: req.DiagnosticoTecnico,
		SolucionAplicada:   req.SolucionAplicada,
		Estado:             req.Estado,
		Prioridad:          req.Prioridad,
		TecnicoID:          req.TecnicoID,
		CostoManoObra:      req.CostoManoObra,
		CostoRepuestos:     req.CostoRepuestos,
		Notas:              req.Notas,
		FechaPromesa:       req.FechaPromesa,
		UsuarioID:          userID,
	})
	if err == services.ErrOrdenTrabajoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update orden"})
		return
	}

	logAuditoria(c, "editar", "orden_trabajo", id, estadoAnterior, estadoNuevo)
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
	tecnicoNombre, err := getOrdenTrabajoService().AsignarTecnico(id, req.TecnicoID, userID)
	if err == services.ErrOrdenTrabajoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign tecnico"})
		return
	}

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
	insumoNombre, stockActual, err := getOrdenTrabajoService().AgregarInsumo(id, req.InsumoID, req.Cantidad, userID)
	if err == services.ErrOTInsumoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}
	if err == services.ErrInsumoSinStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stock insuficiente de %s (disponible: %d)", insumoNombre, stockActual)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add insumo"})
		return
	}

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
	err = getOrdenTrabajoService().RegistrarTrazabilidad(id, req.Accion, req.Detalle, userID)
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

	err = getOrdenTrabajoService().DeleteOrdenTrabajo(id)
	if err == services.ErrOrdenTrabajoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}
	if err == services.ErrOrdenTrabajoDeleteBlocked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden eliminar órdenes recibidas o canceladas"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete orden"})
		return
	}

	logAuditoria(c, "eliminar", "orden_trabajo", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Orden deleted successfully"})
}

// GetOrdenesStats obtiene estadísticas de órdenes de trabajo
func GetOrdenesStats(c *gin.Context) {
	stats, err := getOrdenTrabajoService().Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetTecnicos obtiene la lista de técnicos (usuarios con rol técnico)
func GetTecnicos(c *gin.Context) {
	tecnicos, err := getOrdenTrabajoService().ListTecnicos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tecnicos"})
		return
	}

	c.JSON(http.StatusOK, tecnicos)
}
