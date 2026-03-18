package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetInsumos obtiene todos los insumos
func GetInsumos(c *gin.Context) {
	categoria := c.Query("categoria")
	bajoStock := c.Query("bajo_stock")

	query := `
		SELECT id, created_at, updated_at, codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo
		FROM insumos WHERE 1=1
	`
	args := []interface{}{}

	if categoria != "" {
		query += " AND categoria = ?"
		args = append(args, categoria)
	}
	if bajoStock == "true" {
		query += " AND stock <= stock_minimo"
	}

	query += " ORDER BY nombre"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insumos"})
		return
	}
	defer rows.Close()

	var insumos []models.Insumo
	for rows.Next() {
		var i models.Insumo
		var activo int
		err := rows.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Codigo, &i.Nombre, &i.Descripcion,
			&i.Categoria, &i.UnidadMedida, &i.Stock, &i.StockMinimo, &i.Costo, &i.SedeID, &activo)
		if err != nil {
			continue
		}
		i.Activo = activo == 1
		insumos = append(insumos, i)
	}

	c.JSON(http.StatusOK, insumos)
}

// GetInsumo obtiene un insumo por ID
func GetInsumo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insumo ID"})
		return
	}

	var i models.Insumo
	var activo int
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo
		FROM insumos WHERE id = ?`, id).
		Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Codigo, &i.Nombre, &i.Descripcion,
			&i.Categoria, &i.UnidadMedida, &i.Stock, &i.StockMinimo, &i.Costo, &i.SedeID, &activo)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insumo"})
		return
	}
	i.Activo = activo == 1

	c.JSON(http.StatusOK, i)
}

// CreateInsumo crea un nuevo insumo
func CreateInsumo(c *gin.Context) {
	var req struct {
		Codigo       string  `json:"codigo" binding:"required"`
		Nombre       string  `json:"nombre" binding:"required"`
		Descripcion  string  `json:"descripcion"`
		Categoria    string  `json:"categoria"`
		UnidadMedida string  `json:"unidad_medida"`
		Stock        int     `json:"stock"`
		StockMinimo  int     `json:"stock_minimo"`
		Costo        float64 `json:"costo"`
		SedeID       int64   `json:"sede_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO insumos (codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		req.Codigo, req.Nombre, req.Descripcion, req.Categoria, req.UnidadMedida, req.Stock, req.StockMinimo, req.Costo, req.SedeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create insumo"})
		return
	}

	insumoID, _ := result.LastInsertId()
	logAuditoria(c, "crear", "insumo", insumoID, "", req.Nombre)

	c.JSON(http.StatusCreated, gin.H{"id": insumoID})
}

// UpdateInsumo actualiza un insumo
func UpdateInsumo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insumo ID"})
		return
	}

	var req struct {
		Codigo       string  `json:"codigo"`
		Nombre       string  `json:"nombre"`
		Descripcion  string  `json:"descripcion"`
		Categoria    string  `json:"categoria"`
		UnidadMedida string  `json:"unidad_medida"`
		Stock        int     `json:"stock"`
		StockMinimo  int     `json:"stock_minimo"`
		Costo        float64 `json:"costo"`
		Activo       bool    `json:"activo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activo := 0
	if req.Activo {
		activo = 1
	}

	_, err = database.DB.Exec(`
		UPDATE insumos SET codigo = ?, nombre = ?, descripcion = ?, categoria = ?, unidad_medida = ?, stock = ?,
		                   stock_minimo = ?, costo = ?, activo = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		req.Codigo, req.Nombre, req.Descripcion, req.Categoria, req.UnidadMedida, req.Stock,
		req.StockMinimo, req.Costo, activo, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update insumo"})
		return
	}

	logAuditoria(c, "editar", "insumo", id, "", req.Nombre)

	c.JSON(http.StatusOK, gin.H{"message": "Insumo updated successfully"})
}

// AjustarStockInsumo ajusta el stock de un insumo (entrada/salida)
func AjustarStockInsumo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insumo ID"})
		return
	}

	var req struct {
		Cantidad int    `json:"cantidad" binding:"required"`
		Tipo     string `json:"tipo" binding:"required"` // entrada, salida
		Motivo   string `json:"motivo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener stock actual
	var stockActual int
	err = database.DB.QueryRow("SELECT stock FROM insumos WHERE id = ?", id).Scan(&stockActual)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}

	nuevoStock := stockActual
	if req.Tipo == "entrada" {
		nuevoStock = stockActual + req.Cantidad
	} else if req.Tipo == "salida" {
		if stockActual < req.Cantidad {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
			return
		}
		nuevoStock = stockActual - req.Cantidad
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido (debe ser entrada o salida)"})
		return
	}

	_, err = database.DB.Exec("UPDATE insumos SET stock = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", nuevoStock, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	logAuditoria(c, "ajustar_stock", "insumo", id, strconv.Itoa(stockActual), strconv.Itoa(nuevoStock)+" ("+req.Motivo+")")

	c.JSON(http.StatusOK, gin.H{"message": "Stock adjusted successfully", "nuevo_stock": nuevoStock})
}

// DeleteInsumo elimina un insumo
func DeleteInsumo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insumo ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM insumos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete insumo"})
		return
	}

	logAuditoria(c, "eliminar", "insumo", id, "", "")
	c.JSON(http.StatusOK, gin.H{"message": "Insumo deleted successfully"})
}

// GetCompatibilidades obtiene las compatibilidades de productos
func GetCompatibilidades(c *gin.Context) {
	productoID := c.Query("producto_id")

	query := `
		SELECT c.id, c.producto_id, c.compatible_con_id, c.notas,
		       p1.name as producto_nombre, p1.brand as producto_marca,
		       p2.name as compatible_nombre, p2.brand as compatible_marca
		FROM compatibilidades c
		INNER JOIN products p1 ON c.producto_id = p1.id
		INNER JOIN products p2 ON c.compatible_con_id = p2.id
		WHERE 1=1
	`
	args := []interface{}{}

	if productoID != "" {
		query += " AND (c.producto_id = ? OR c.compatible_con_id = ?)"
		args = append(args, productoID, productoID)
	}

	query += " ORDER BY p1.name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch compatibilidades"})
		return
	}
	defer rows.Close()

	type CompView struct {
		ID               int64  `json:"id"`
		ProductoID       int64  `json:"producto_id"`
		CompatibleCon    int64  `json:"compatible_con"`
		Notas            string `json:"notas"`
		ProductoNombre   string `json:"producto_nombre"`
		ProductoMarca    string `json:"producto_marca"`
		CompatibleNombre string `json:"compatible_nombre"`
		CompatibleMarca  string `json:"compatible_marca"`
	}

	var comps []CompView
	for rows.Next() {
		var c CompView
		rows.Scan(&c.ID, &c.ProductoID, &c.CompatibleCon, &c.Notas,
			&c.ProductoNombre, &c.ProductoMarca, &c.CompatibleNombre, &c.CompatibleMarca)
		comps = append(comps, c)
	}

	c.JSON(http.StatusOK, comps)
}

// BuscarCompatibles busca productos compatibles con uno dado
func BuscarCompatibles(c *gin.Context) {
	productoID := c.Param("id")

	// Obtener info del producto
	var producto struct {
		ID       int64
		Name     string
		Brand    string
		Category string
	}
	err := database.DB.QueryRow("SELECT id, name, brand, category FROM products WHERE id = ?", productoID).
		Scan(&producto.ID, &producto.Name, &producto.Brand, &producto.Category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto not found"})
		return
	}

	// Buscar compatibles directos (compatible_con_id es el nombre de columna en la tabla)
	rows, _ := database.DB.Query(`
		SELECT p.id, p.name, p.brand, p.category, p.precio_venta, COALESCE(c.notas, '')
		FROM compatibilidades c
		INNER JOIN products p ON p.id = CASE WHEN c.producto_id = ? THEN c.compatible_con_id ELSE c.producto_id END
		WHERE c.producto_id = ? OR c.compatible_con_id = ?
	`, productoID, productoID, productoID)
	defer rows.Close()

	type CompatibleView struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Brand     string  `json:"brand"`
		Category  string  `json:"category"`
		Precio    float64 `json:"precio"`
		Notas     string  `json:"notas"`
		TipoMatch string  `json:"tipo_match"` // directo, misma_categoria, mismo_fabricante
	}

	var compatibles []CompatibleView
	seenIDs := make(map[int64]bool)

	for rows.Next() {
		var comp CompatibleView
		rows.Scan(&comp.ID, &comp.Name, &comp.Brand, &comp.Category, &comp.Precio, &comp.Notas)
		comp.TipoMatch = "directo"
		compatibles = append(compatibles, comp)
		seenIDs[comp.ID] = true
	}

	// Buscar productos de la misma categoría y marca
	rows2, _ := database.DB.Query(`
		SELECT id, name, brand, category, precio_venta
		FROM products
		WHERE id != ? AND activo = 1 AND (category = ? OR brand = ?)
		LIMIT 20
	`, productoID, producto.Category, producto.Brand)
	defer rows2.Close()

	for rows2.Next() {
		var comp CompatibleView
		rows2.Scan(&comp.ID, &comp.Name, &comp.Brand, &comp.Category, &comp.Precio)
		if seenIDs[comp.ID] {
			continue
		}
		if comp.Category == producto.Category {
			comp.TipoMatch = "misma_categoria"
		} else {
			comp.TipoMatch = "mismo_fabricante"
		}
		compatibles = append(compatibles, comp)
	}

	c.JSON(http.StatusOK, gin.H{"producto": producto, "compatibles": compatibles})
}

// CreateCompatibilidad crea una nueva compatibilidad
func CreateCompatibilidad(c *gin.Context) {
	var req struct {
		ProductoID    int64  `json:"producto_id" binding:"required"`
		CompatibleCon int64  `json:"compatible_con" binding:"required"`
		Notas         string `json:"notas"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ProductoID == req.CompatibleCon {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Un producto no puede ser compatible consigo mismo"})
		return
	}

	// Verificar que no exista ya
	var exists int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM compatibilidades 
		WHERE (producto_id = ? AND compatible_con_id = ?) OR (producto_id = ? AND compatible_con_id = ?)
	`, req.ProductoID, req.CompatibleCon, req.CompatibleCon, req.ProductoID).Scan(&exists)

	if exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Esta compatibilidad ya existe"})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO compatibilidades (producto_id, compatible_con_id, notas)
		VALUES (?, ?, ?)`,
		req.ProductoID, req.CompatibleCon, req.Notas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create compatibilidad"})
		return
	}

	compID, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": compID})
}

// DeleteCompatibilidad elimina una compatibilidad
func DeleteCompatibilidad(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compatibilidad ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM compatibilidades WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete compatibilidad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compatibilidad deleted successfully"})
}

// GetInsumosStats obtiene estadísticas de insumos
func GetInsumosStats(c *gin.Context) {
	type Stats struct {
		TotalInsumos    int     `json:"total_insumos"`
		BajoStock       int     `json:"bajo_stock"`
		SinStock        int     `json:"sin_stock"`
		ValorInventario float64 `json:"valor_inventario"`
	}

	var stats Stats
	database.DB.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1").Scan(&stats.TotalInsumos)
	database.DB.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1 AND stock <= stock_minimo AND stock > 0").Scan(&stats.BajoStock)
	database.DB.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1 AND stock = 0").Scan(&stats.SinStock)
	database.DB.QueryRow("SELECT COALESCE(SUM(stock * costo), 0) FROM insumos WHERE activo = 1").Scan(&stats.ValorInventario)

	c.JSON(http.StatusOK, stats)
}
