package controllers

import (
<<<<<<< HEAD
	"database/sql"
	"encoding/json"
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
		apiErr := errors.NewDatabaseError("Fetch cotizaciones", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	type CotizacionView struct {
		models.Cotizacion
		UsuarioNombre string `json:"usuario_nombre"`
		SedeNombre    string `json:"sede_nombre"`
	}

	var cotizaciones []CotizacionView
	for rows.Next() {
		var cot CotizacionView
		err := rows.Scan(&cot.ID, &cot.CreatedAt, &cot.UpdatedAt, &cot.NumeroCotizacion, &cot.ClienteNombre,
			&cot.ClienteTelefono, &cot.ClienteEmail, &cot.ValidezDias, &cot.Estado, &cot.Total,
			&cot.Descuento, &cot.Notas, &cot.UsuarioID, &cot.SedeID,
			&cot.UsuarioNombre, &cot.SedeNombre)
		if err != nil {
			continue
		}
		cotizaciones = append(cotizaciones, cot)
	}

	c.JSON(200, cotizaciones)
=======
	c.JSON(http.StatusOK, items)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetCotizacion obtiene una cotización por ID con sus items
func GetCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotization id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	var cot models.Cotizacion
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, numero_cotizacion, cliente_nombre, 
		       cliente_telefono, cliente_email, validez, estado, total, 
		       descuento, notas, usuario_id, sede_id
		FROM cotizaciones WHERE id = ?`, id).
		Scan(&cot.ID, &cot.CreatedAt, &cot.UpdatedAt, &cot.NumeroCotizacion, &cot.ClienteNombre,
			&cot.ClienteTelefono, &cot.ClienteEmail, &cot.ValidezDias, &cot.Estado, &cot.Total,
			&cot.Descuento, &cot.Notas, &cot.UsuarioID, &cot.SedeID)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("Cotization", id)
		c.JSON(apiErr.Code, apiErr)
=======
	cot, items, err := getCotizacionService().GetCotizacion(id)
	if err == services.ErrCotizacionNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch cotization", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Obtener items
	itemRows, _ := database.DB.Query(`
		SELECT ci.id, ci.producto_id, ci.cantidad, ci.precio_unit, ci.subtotal,
		       p.name, p.brand, p.codigo
		FROM cotizacion_items ci
		INNER JOIN products p ON ci.producto_id = p.id
		WHERE ci.cotizacion_id = ?
	`, id)
	defer itemRows.Close()

	type ItemView struct {
		ID             int64   `json:"id"`
		ProductoID     int64   `json:"producto_id"`
		Cantidad       int     `json:"cantidad"`
		PrecioUnitario float64 `json:"precio_unitario"`
		Subtotal       float64 `json:"subtotal"`
		ProductoNombre string  `json:"producto_nombre"`
		ProductoMarca  string  `json:"producto_marca"`
		ProductoCodigo string  `json:"producto_codigo"`
	}

	var items []ItemView
	for itemRows.Next() {
		var item ItemView
		itemRows.Scan(&item.ID, &item.ProductoID, &item.Cantidad, &item.PrecioUnitario, &item.Subtotal,
			&item.ProductoNombre, &item.ProductoMarca, &item.ProductoCodigo)
		items = append(items, item)
	}

	c.JSON(200, gin.H{"cotizacion": cot, "items": items})
=======
	c.JSON(http.StatusOK, gin.H{"cotizacion": cot, "items": items})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// CreateCotizacion crea una nueva cotización
func CreateCotizacion(c *gin.Context) {
	var req struct {
		ClienteNombre   string  `json:"cliente_nombre" validate:"required,min=3"`
		ClienteTelefono string  `json:"cliente_telefono"`
		ClienteEmail    string  `json:"cliente_email" validate:"email"`
		ValidezDias     int     `json:"validez_dias"`
		Descuento       float64 `json:"descuento" validate:"gte=0"`
		Notas           string  `json:"notas"`
		SedeID          int64   `json:"sede_id" validate:"required,gt=0"`
		Items           []struct {
			ProductoID     int64   `json:"producto_id"`
			Cantidad       int     `json:"cantidad"`
			PrecioUnitario float64 `json:"precio_unitario"`
		} `json:"items" validate:"required"`
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
	if len(req.Items) == 0 {
		apiErr := errors.NewBadRequest("Must include at least one item")
		c.JSON(apiErr.Code, apiErr)
		return
	}

=======
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	inItems := make([]services.CotizacionCreateItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		inItems = append(inItems, services.CotizacionCreateItemInput{ProductoID: it.ProductoID, Cantidad: it.Cantidad, PrecioUnitario: it.PrecioUnitario})
	}

<<<<<<< HEAD
	tx, err := database.DB.Begin()
	if err != nil {
		apiErr := errors.NewDatabaseError("Start transaction", err)
		c.JSON(apiErr.Code, apiErr)
=======
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
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
<<<<<<< HEAD
		tx.Rollback()
		apiErr := errors.NewDatabaseError("Insert cotization", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	cotID, _ := result.LastInsertId()

	// Insertar items
	for _, item := range req.Items {
		subtotal := float64(item.Cantidad) * item.PrecioUnitario
		_, err = tx.Exec(`
			INSERT INTO cotizacion_items (cotizacion_id, producto_id, cantidad, precio_unit, subtotal)
			VALUES (?, ?, ?, ?, ?)`,
			cotID, item.ProductoID, item.Cantidad, item.PrecioUnitario, subtotal)
		if err != nil {
			tx.Rollback()
			apiErr := errors.NewDatabaseError("Insert cotization items", err)
			c.JSON(apiErr.Code, apiErr)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		apiErr := errors.NewDatabaseError("Commit transaction", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "crear", "cotizacion", cotID, "", numeroCotizacion)

	c.JSON(201, gin.H{"id": cotID, "numero_cotizacion": numeroCotizacion, "total": total})
=======
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cotización: " + err.Error()})
		return
	}

	logAuditoria(c, "crear", "cotizacion", id, "", numero)
	c.JSON(http.StatusCreated, gin.H{"id": id, "numero_cotizacion": numero, "total": total})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// UpdateCotizacionEstado actualiza el estado de una cotización
func UpdateCotizacionEstado(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotización ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var req struct {
		Estado string `json:"estado" validate:"required,oneof=pendiente aprobada rechazada vencida convertida"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Validar estructura
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	var estadoAnterior string
	err = database.DB.QueryRow("SELECT estado FROM cotizaciones WHERE id = ?", id).Scan(&estadoAnterior)
	if err != nil {
		apiErr := errors.NewNotFound("Cotización", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("UPDATE cotizaciones SET estado = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Estado, id)
=======
	estadoAnterior, err := getCotizacionService().UpdateEstado(id, req.Estado)
	if err == services.ErrEstadoCotizacionBad {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Update cotization status", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "editar", "cotizacion", id, estadoAnterior, req.Estado)
<<<<<<< HEAD

	c.JSON(200, gin.H{"message": "Cotización updated successfully"})
=======
	c.JSON(http.StatusOK, gin.H{"message": "Cotización updated successfully"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// DeleteCotizacion elimina una cotización
func DeleteCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotización ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar que existe
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM cotizaciones WHERE id = ?", id).Scan(&exists)
	if err != nil {
		apiErr := errors.NewNotFound("Cotización", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Eliminar items primero
	_, err = database.DB.Exec("DELETE FROM cotizacion_items WHERE cotizacion_id = ?", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete cotization items", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = database.DB.Exec("DELETE FROM cotizaciones WHERE id = ?", id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete cotization", err)
		c.JSON(apiErr.Code, apiErr)
=======
	if err := getCotizacionService().DeleteCotizacion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cotización"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}

	logAuditoria(c, "eliminar", "cotizacion", id, "", "")
	c.JSON(200, gin.H{"message": "Cotización deleted successfully"})
}

// ConvertirCotizacionAVenta convierte una cotización aprobada a venta
func ConvertirCotizacionAVenta(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotización ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Verificar que la cotización esté aprobada
	var estado string
	var sedeID int64
	var total float64
	var clienteNombre string
	err = database.DB.QueryRow("SELECT estado, sede_id, total, cliente_nombre FROM cotizaciones WHERE id = ?", id).
		Scan(&estado, &sedeID, &total, &clienteNombre)
	if err != nil {
		apiErr := errors.NewNotFound("Cotización", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	if estado != "aprobada" {
		apiErr := errors.NewConflict("Can only convert approved quotations")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	userID, exists := c.Get("userid")
	if !exists {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		apiErr := errors.NewDatabaseError("Start transaction", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Crear venta
	var count int
	tx.QueryRow("SELECT COUNT(*) FROM ventas").Scan(&count)
	numeroVenta := fmt.Sprintf("VNT-%d-%04d", time.Now().Year(), count+1)

	result, err := tx.Exec(`
		INSERT INTO ventas (numero_venta, cliente_nombre, total, usuario_id, sede_id, cotizacion_id)
		VALUES (?, ?, ?, ?, ?, ?)`,
		numeroVenta, clienteNombre, total, userID, sedeID, id)

	if err != nil {
		tx.Rollback()
		apiErr := errors.NewDatabaseError("Create sale", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	ventaID, _ := result.LastInsertId()

	// Obtener items de la cotización y registrar en detalle de venta
	itemRows, _ := tx.Query(`
		SELECT producto_id, cantidad, precio_unit, subtotal
		FROM cotizacion_items WHERE cotizacion_id = ?
	`, id)

	type Item struct {
		ProductoID     int64
		Cantidad       int
		PrecioUnitario float64
		Subtotal       float64
	}
	var items []Item

	for itemRows.Next() {
		var item Item
		itemRows.Scan(&item.ProductoID, &item.Cantidad, &item.PrecioUnitario, &item.Subtotal)
		items = append(items, item)
	}
	itemRows.Close()

	// Insertar items en venta y descontar stock
	for _, item := range items {
		_, err = tx.Exec(`
			INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unit, subtotal)
			VALUES (?, ?, ?, ?, ?)`,
			ventaID, item.ProductoID, item.Cantidad, item.PrecioUnitario, item.Subtotal)
		if err != nil {
			tx.Rollback()
			apiErr := errors.NewDatabaseError("Insert sale items", err)
			c.JSON(apiErr.Code, apiErr)
			return
		}

		// Descontar stock de la sede
		_, err = tx.Exec(`
			UPDATE stock_sedes SET cantidad = cantidad - ? WHERE producto_id = ? AND sede_id = ?`,
			item.Cantidad, item.ProductoID, sedeID)
		if err != nil {
			tx.Rollback()
			apiErr := errors.NewDatabaseError("Update stock", err)
			c.JSON(apiErr.Code, apiErr)
			return
		}
	}

	// Actualizar estado de cotización
	tx.Exec("UPDATE cotizaciones SET estado = 'convertida', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)

	if err = tx.Commit(); err != nil {
		apiErr := errors.NewDatabaseError("Commit transaction", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	itemsJSON, _ := json.Marshal(items)
	logAuditoria(c, "convertir_venta", "cotizacion", id, "", string(itemsJSON))

	c.JSON(201, gin.H{"venta_id": ventaID, "numero_venta": numeroVenta})
=======
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
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GenerarPDFCotizacion genera un PDF de la cotización (retorna datos para frontend)
func GenerarPDFCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotización ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Obtener datos completos de la cotización para que el frontend genere el PDF
	var cot struct {
		NumeroCotizacion string    `json:"numero_cotizacion"`
		ClienteNombre    string    `json:"cliente_nombre"`
		ClienteTelefono  string    `json:"cliente_telefono"`
		ClienteEmail     string    `json:"cliente_email"`
		ValidezDias      int       `json:"validez_dias"`
		Total            float64   `json:"total"`
		Descuento        float64   `json:"descuento"`
		Notas            string    `json:"notas"`
		FechaCreacion    time.Time `json:"fecha_creacion"`
		UsuarioNombre    string    `json:"usuario_nombre"`
		SedeNombre       string    `json:"sede_nombre"`
		SedeDireccion    string    `json:"sede_direccion"`
	}

	err = database.DB.QueryRow(`
		SELECT c.numero_cotizacion, c.cliente_nombre, c.cliente_telefono, c.cliente_email,
		       c.validez, c.total, c.descuento, c.notas, c.created_at,
		       u.name, s.nombre, s.direccion
		FROM cotizaciones c
		INNER JOIN users u ON c.usuario_id = u.id
		INNER JOIN sedes s ON c.sede_id = s.id
		WHERE c.id = ?
	`, id).Scan(&cot.NumeroCotizacion, &cot.ClienteNombre, &cot.ClienteTelefono, &cot.ClienteEmail,
		&cot.ValidezDias, &cot.Total, &cot.Descuento, &cot.Notas, &cot.FechaCreacion,
		&cot.UsuarioNombre, &cot.SedeNombre, &cot.SedeDireccion)

	if err != nil {
		apiErr := errors.NewNotFound("Cotización", id)
		c.JSON(apiErr.Code, apiErr)
=======
	cot, items, err := getCotizacionService().GetPDFData(id)
	if err == services.ErrCotizacionNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cotización"})
		return
	}

	c.JSON(200, gin.H{"cotizacion": cot, "items": items})
}
