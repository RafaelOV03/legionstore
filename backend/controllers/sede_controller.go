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

func getSedeService() *services.SedeService {
	repo := repositories.NewSedeRepository(database.DB)
	return services.NewSedeService(repo)
}

// GetSedes obtiene todas las sedes
func GetSedes(c *gin.Context) {
	sedes, err := getSedeService().ListSedes()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch locations", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, sedes)
}

// GetSede obtiene una sede por ID
func GetSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid sede ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	var sede models.Sede
	var activa int
	err = database.DB.QueryRow(`SELECT id, created_at, updated_at, nombre, direccion, telefono, activa FROM sedes WHERE id = ?`, id).
		Scan(&sede.ID, &sede.CreatedAt, &sede.UpdatedAt, &sede.Nombre, &sede.Direccion, &sede.Telefono, &activa)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("Sede", id)
		c.JSON(apiErr.Code, apiErr)
=======
	sede, err := getSedeService().GetSede(id)
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch location", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
<<<<<<< HEAD

	sede.Activa = activa == 1
	c.JSON(200, sede)
=======
	c.JSON(http.StatusOK, sede)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// CreateSede crea una nueva sede
func CreateSede(c *gin.Context) {
	var req struct {
		Nombre    string `json:"nombre" validate:"required,min=3"`
		Direccion string `json:"direccion" validate:"required,min=5"`
		Telefono  string `json:"telefono" validate:"required"`
	}

	var sede models.Sede
	if err := c.ShouldBindJSON(&sede); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	req.Nombre = sede.Nombre
	req.Direccion = sede.Direccion
	req.Telefono = sede.Telefono
	validationErrors := validation.ValidateStruct(req)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	createdSede, err := getSedeService().CreateSede(sede)
	if err != nil {
		apiErr := errors.NewDatabaseError("Insert location", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Log de auditoría
	logAuditoria(c, "crear", "sede", createdSede.ID, "", createdSede.Nombre)

<<<<<<< HEAD
	c.JSON(201, sede)
=======
	c.JSON(http.StatusCreated, createdSede)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// UpdateSede actualiza una sede
func UpdateSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid sede ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var sede models.Sede
	if err := c.ShouldBindJSON(&sede); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM sedes WHERE id = ?", id).Scan(&exists)
	if err != nil {
		apiErr := errors.NewNotFound("Sede", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	updatedSede, err := getSedeService().UpdateSede(id, sede)
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Update location", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	logAuditoria(c, "editar", "sede", id, "", updatedSede.Nombre)

<<<<<<< HEAD
	sede.ID = id
	logAuditoria(c, "editar", "sede", id, "", sede.Nombre)

	c.JSON(200, sede)
=======
	c.JSON(http.StatusOK, updatedSede)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// DeleteSede elimina una sede
func DeleteSede(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid sede ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM sedes WHERE id = ?", id).Scan(&exists)
	if err != nil {
		apiErr := errors.NewNotFound("Sede", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	// Verificar que no haya usuarios o stock asociados
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE sede_id = ?", id).Scan(&count)
	if count > 0 {
		apiErr := errors.NewConflict("Location has associated users")
		c.JSON(apiErr.Code, apiErr)
=======
	err = getSedeService().DeleteSede(id)
	if err == repositories.ErrSedeHasAssociatedUsers {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar la sede, tiene usuarios asociados"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err == repositories.ErrSedeNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sede not found"})
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete location", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "eliminar", "sede", id, "", "")
	c.JSON(200, gin.H{"message": "Sede deleted successfully"})
}
<<<<<<< HEAD

// GetStockMultisede obtiene el stock de todos los productos en todas las sedes
func GetStockMultisede(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT ss.id, ss.sede_id, ss.producto_id, ss.cantidad, ss.stock_minimo, ss.stock_maximo,
		       s.nombre as sede_nombre,
		       p.codigo, p.name, p.category, p.brand
		FROM stock_sedes ss
		INNER JOIN sedes s ON ss.sede_id = s.id
		INNER JOIN products p ON ss.producto_id = p.id
		WHERE p.activo = 1 AND s.activa = 1
		ORDER BY p.name, s.nombre
	`)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch stock", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	defer rows.Close()

	type StockItem struct {
		ID          int64  `json:"id"`
		SedeID      int64  `json:"sede_id"`
		ProductoID  int64  `json:"producto_id"`
		Cantidad    int    `json:"cantidad"`
		StockMinimo int    `json:"stock_minimo"`
		StockMaximo int    `json:"stock_maximo"`
		SedeNombre  string `json:"sede_nombre"`
		Codigo      string `json:"codigo"`
		Producto    string `json:"producto"`
		Categoria   string `json:"categoria"`
		Marca       string `json:"marca"`
	}

	var items []StockItem
	for rows.Next() {
		var item StockItem
		err := rows.Scan(&item.ID, &item.SedeID, &item.ProductoID, &item.Cantidad, &item.StockMinimo, &item.StockMaximo,
			&item.SedeNombre, &item.Codigo, &item.Producto, &item.Categoria, &item.Marca)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}

// GetStockBySede obtiene el stock de una sede específica
func GetStockBySede(c *gin.Context) {
	sedeID, err := strconv.ParseInt(c.Param("sede_id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid sede ID")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	rows, err := database.DB.Query(`
		SELECT ss.id, ss.producto_id, ss.cantidad, ss.stock_minimo, ss.stock_maximo,
		       p.codigo, p.name, p.description, p.precio_venta, p.category, p.brand, p.image_url
		FROM stock_sedes ss
		INNER JOIN products p ON ss.producto_id = p.id
		WHERE ss.sede_id = ? AND p.activo = 1
		ORDER BY p.name
	`, sedeID)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch stock by location", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	defer rows.Close()

	type StockProducto struct {
		ID          int64   `json:"id"`
		ProductoID  int64   `json:"producto_id"`
		Cantidad    int     `json:"cantidad"`
		StockMinimo int     `json:"stock_minimo"`
		StockMaximo int     `json:"stock_maximo"`
		Codigo      string  `json:"codigo"`
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
		Categoria   string  `json:"categoria"`
		Marca       string  `json:"marca"`
		ImageURL    string  `json:"image_url"`
	}

	var items []StockProducto
	for rows.Next() {
		var item StockProducto
		err := rows.Scan(&item.ID, &item.ProductoID, &item.Cantidad, &item.StockMinimo, &item.StockMaximo,
			&item.Codigo, &item.Nombre, &item.Descripcion, &item.Precio, &item.Categoria, &item.Marca, &item.ImageURL)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}

// UpdateStock actualiza el stock de un producto en una sede
func UpdateStock(c *gin.Context) {
	var req struct {
		SedeID      int64 `json:"sede_id" binding:"required"`
		ProductoID  int64 `json:"producto_id" binding:"required"`
		Cantidad    int   `json:"cantidad"`
		StockMinimo int   `json:"stock_minimo"`
		StockMaximo int   `json:"stock_maximo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar si existe el registro
	var existingID int64
	err := database.DB.QueryRow("SELECT id FROM stock_sedes WHERE sede_id = ? AND producto_id = ?", req.SedeID, req.ProductoID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Crear nuevo registro
		_, err = database.DB.Exec(`INSERT INTO stock_sedes (sede_id, producto_id, cantidad, stock_minimo, stock_maximo) VALUES (?, ?, ?, ?, ?)`,
			req.SedeID, req.ProductoID, req.Cantidad, req.StockMinimo, req.StockMaximo)
	} else {
		// Actualizar existente
		_, err = database.DB.Exec(`UPDATE stock_sedes SET cantidad = ?, stock_minimo = ?, stock_maximo = ?, updated_at = CURRENT_TIMESTAMP WHERE sede_id = ? AND producto_id = ?`,
			req.Cantidad, req.StockMinimo, req.StockMaximo, req.SedeID, req.ProductoID)
	}

	if err != nil {
		apiErr := errors.NewDatabaseError("Update stock", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	logAuditoria(c, "actualizar_stock", "stock_sedes", req.ProductoID, "", "")
	c.JSON(200, gin.H{"message": "Stock updated successfully"})
}
=======
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
