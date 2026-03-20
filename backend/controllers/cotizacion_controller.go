package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/validation"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCotizaciones obtiene todas las cotizaciones
func GetCotizaciones(c *gin.Context) {
	estado := c.Query("estado")
	sedeID := c.Query("sede_id")

	query := `
		SELECT c.id, c.created_at, c.updated_at, c.numero_cotizacion, c.cliente_nombre, 
		       c.cliente_telefono, c.cliente_email, c.validez, c.estado, c.total, 
		       c.descuento, c.notas, c.usuario_id, c.sede_id,
		       u.name as usuario_nombre, s.nombre as sede_nombre
		FROM cotizaciones c
		INNER JOIN users u ON c.usuario_id = u.id
		INNER JOIN sedes s ON c.sede_id = s.id
		WHERE 1=1
	`
	args := []interface{}{}

	if estado != "" {
		query += " AND c.estado = ?"
		args = append(args, estado)
	}
	if sedeID != "" {
		query += " AND c.sede_id = ?"
		args = append(args, sedeID)
	}

	query += " ORDER BY c.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch cotizaciones", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	defer rows.Close()

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
}

// GetCotizacion obtiene una cotización por ID con sus items
func GetCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotization id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

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
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch cotization", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

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

	if len(req.Items) == 0 {
		apiErr := errors.NewBadRequest("Must include at least one item")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	userID, exists := c.Get("userid")
	if !exists || userID == nil {
		apiErr := errors.ErrUnauthorized
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Generar número de cotización
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM cotizaciones").Scan(&count)
	numeroCotizacion := fmt.Sprintf("COT-%d-%04d", time.Now().Year(), count+1)

	// Calcular total
	var total float64
	for _, item := range req.Items {
		total += float64(item.Cantidad) * item.PrecioUnitario
	}
	total = total - req.Descuento

	if req.ValidezDias == 0 {
		req.ValidezDias = 30
	}

	tx, err := database.DB.Begin()
	if err != nil {
		apiErr := errors.NewDatabaseError("Start transaction", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	result, err := tx.Exec(`
		INSERT INTO cotizaciones (numero_cotizacion, cliente_nombre, cliente_telefono, cliente_email, 
		                          validez, estado, total, descuento, notas, usuario_id, sede_id)
		VALUES (?, ?, ?, ?, ?, 'pendiente', ?, ?, ?, ?, ?)`,
		numeroCotizacion, req.ClienteNombre, req.ClienteTelefono, req.ClienteEmail,
		req.ValidezDias, total, req.Descuento, req.Notas, userID, req.SedeID)

	if err != nil {
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
	if err != nil {
		apiErr := errors.NewDatabaseError("Update cotization status", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "editar", "cotizacion", id, estadoAnterior, req.Estado)

	c.JSON(200, gin.H{"message": "Cotización updated successfully"})
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
}

// GenerarPDFCotizacion genera un PDF de la cotización (retorna datos para frontend)
func GenerarPDFCotizacion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid cotización ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

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
		return
	}

	// Obtener items
	itemRows, _ := database.DB.Query(`
		SELECT p.codigo, p.name, p.brand, ci.cantidad, ci.precio_unit, ci.subtotal
		FROM cotizacion_items ci
		INNER JOIN products p ON ci.producto_id = p.id
		WHERE ci.cotizacion_id = ?
	`, id)
	defer itemRows.Close()

	type ItemPDF struct {
		Codigo         string  `json:"codigo"`
		Nombre         string  `json:"nombre"`
		Marca          string  `json:"marca"`
		Cantidad       int     `json:"cantidad"`
		PrecioUnitario float64 `json:"precio_unitario"`
		Subtotal       float64 `json:"subtotal"`
	}

	var items []ItemPDF
	for itemRows.Next() {
		var item ItemPDF
		itemRows.Scan(&item.Codigo, &item.Nombre, &item.Marca, &item.Cantidad, &item.PrecioUnitario, &item.Subtotal)
		items = append(items, item)
	}

	c.JSON(200, gin.H{"cotizacion": cot, "items": items})
}
