package controllers

import (
<<<<<<< HEAD
	"database/sql"
	"fmt"
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/validation"
=======
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
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
		apiErr := errors.NewDatabaseError("Fetch RMAs", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, rmas)
}

// GetRMA obtiene una RMA por ID
func GetRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid RMA ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	var rma models.RMA
	var fechaCompra, fechaResolucion sql.NullTime
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, numero_rma, producto_id, cliente_nombre, 
		       cliente_telefono, cliente_email, num_serie, fecha_compra, motivo_devolucion,
		       diagnostico, estado, solucion, fecha_resolucion, usuario_id, sede_id, notas
		FROM rmas WHERE id = ?`, id).
		Scan(&rma.ID, &rma.CreatedAt, &rma.UpdatedAt, &rma.NumeroRMA, &rma.ProductoID,
			&rma.ClienteNombre, &rma.ClienteTelefono, &rma.ClienteEmail, &rma.NumSerie, &fechaCompra,
			&rma.MotivoDevolucion, &rma.Diagnostico, &rma.Estado, &rma.Solucion, &fechaResolucion,
			&rma.UsuarioID, &rma.SedeID, &rma.Notas)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("RMA", id)
		c.JSON(apiErr.Code, apiErr)
=======
	rma, historial, err := getRMAService().GetRMA(id)
	if err == services.ErrRMANotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "RMA not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch RMA", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	if fechaCompra.Valid {
		rma.FechaCompra = fechaCompra.Time
	}
	if fechaResolucion.Valid {
		rma.FechaResolucion = &fechaResolucion.Time
	}

	// Obtener historial
	histRows, _ := database.DB.Query(`
		SELECT h.id, h.created_at, h.estado_anterior, h.estado_nuevo, h.comentario, u.name
		FROM historial_rmas h
		INNER JOIN users u ON h.usuario_id = u.id
		WHERE h.rma_id = ?
		ORDER BY h.created_at DESC
	`, id)
	defer histRows.Close()

	type HistorialItem struct {
		ID             int64     `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		EstadoAnterior string    `json:"estado_anterior"`
		EstadoNuevo    string    `json:"estado_nuevo"`
		Comentario     string    `json:"comentario"`
		Usuario        string    `json:"usuario"`
	}

	var historial []HistorialItem
	for histRows.Next() {
		var h HistorialItem
		histRows.Scan(&h.ID, &h.CreatedAt, &h.EstadoAnterior, &h.EstadoNuevo, &h.Comentario, &h.Usuario)
		historial = append(historial, h)
	}

	c.JSON(200, gin.H{"rma": rma, "historial": historial})
=======
	c.JSON(http.StatusOK, gin.H{"rma": rma, "historial": historial})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// CreateRMA crea una nueva RMA
func CreateRMA(c *gin.Context) {
	var req struct {
		ProductoID       int64  `json:"producto_id" validate:"required,gt=0"`
		ClienteNombre    string `json:"cliente_nombre" validate:"required,min=3"`
		ClienteTelefono  string `json:"cliente_telefono"`
		ClienteEmail     string `json:"cliente_email"`
		NumSerie         string `json:"num_serie"`
		FechaCompra      string `json:"fecha_compra"`
		MotivoDevolucion string `json:"motivo_devolucion" validate:"required"`
		SedeID           int64  `json:"sede_id" validate:"required,gt=0"`
		Notas            string `json:"notas"`
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

<<<<<<< HEAD
	// Obtener usuario del contexto
	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Generar número de RMA
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&count)
	numeroRMA := fmt.Sprintf("RMA-%d-%04d", time.Now().Year(), count+1)

	var fechaCompra interface{}
	if req.FechaCompra != "" {
		fechaCompra = req.FechaCompra
	} else {
		fechaCompra = nil
	}

	result, err := database.DB.Exec(`
		INSERT INTO rmas (numero_rma, producto_id, cliente_nombre, cliente_telefono, cliente_email, 
		                  num_serie, fecha_compra, motivo_devolucion, estado, usuario_id, sede_id, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'recibido', ?, ?, ?)`,
		numeroRMA, req.ProductoID, req.ClienteNombre, req.ClienteTelefono, req.ClienteEmail,
		req.NumSerie, fechaCompra, req.MotivoDevolucion, userID, req.SedeID, req.Notas)

=======
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
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Insert RMA", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "crear", "rma", rmaID, "", numeroRMA)

	c.JSON(201, gin.H{"id": rmaID, "numero_rma": numeroRMA})
}

// UpdateRMA actualiza una RMA
func UpdateRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid RMA ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var req struct {
		Diagnostico string `json:"diagnostico"`
		Estado      string `json:"estado" validate:"oneof=recibido en_revision resuelto rechazado pendiente"`
		Solucion    string `json:"solucion"`
		Notas       string `json:"notas"`
		Comentario  string `json:"comentario"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	if req.Estado != "" {
		validationErrors := validation.ValidateStruct(req)
		if len(validationErrors) > 0 {
			c.JSON(422, validationErrors.ToAPIError())
			return
		}
	}

	userID, _ := c.Get("userid")

<<<<<<< HEAD
	// Obtener estado anterior
	var estadoAnterior string
	err = database.DB.QueryRow("SELECT estado FROM rmas WHERE id = ?", id).Scan(&estadoAnterior)
	if err != nil {
		apiErr := errors.NewNotFound("RMA", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Si el estado cambia a 'resuelto', registrar fecha de resolución
	var fechaResolucion interface{}
	if req.Estado == "resuelto" && estadoAnterior != "resuelto" {
		fechaResolucion = time.Now()
	} else {
		fechaResolucion = nil
=======
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
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Update RMA", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Registrar cambio de estado en historial
	if req.Estado != estadoAnterior && req.Estado != "" {
		database.DB.Exec(`INSERT INTO historial_rmas (rma_id, estado_anterior, estado_nuevo, comentario, usuario_id)
			VALUES (?, ?, ?, ?, ?)`, id, estadoAnterior, req.Estado, req.Comentario, userID)
	}

=======
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	logAuditoria(c, "editar", "rma", id, estadoAnterior, req.Estado)

	c.JSON(200, gin.H{"message": "RMA updated successfully"})
}

// DeleteRMA elimina una RMA
func DeleteRMA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid RMA ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar que existe
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM rmas WHERE id = ?", id).Scan(&exists)
	if err != nil {
		apiErr := errors.NewNotFound("RMA", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Eliminar historial primero
	_, err = database.DB.Exec("DELETE FROM historial_rmas WHERE rma_id = ?", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete RMA history", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("DELETE FROM rmas WHERE id = ?", id)
=======
	err = getRMAService().DeleteRMA(id)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete RMA", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "eliminar", "rma", id, "", "")
	c.JSON(200, gin.H{"message": "RMA deleted successfully"})
}

// GetRMAStats obtiene estadísticas de RMAs
func GetRMAStats(c *gin.Context) {
	stats, err := getRMAService().Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch RMA stats"})
		return
	}

<<<<<<< HEAD
	var stats Stats
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&stats.Total)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'recibido'").Scan(&stats.Recibidos)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'en_revision'").Scan(&stats.EnRevision)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'resuelto'").Scan(&stats.Resueltos)
	database.DB.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'rechazado'").Scan(&stats.Rechazados)

	c.JSON(200, stats)
=======
	c.JSON(http.StatusOK, stats)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}
