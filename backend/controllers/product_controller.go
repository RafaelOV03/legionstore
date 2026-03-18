package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetProducts devuelve todos los productos
func GetProducts(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description, 
		       p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
		       COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.producto_id
		WHERE p.activo = 1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
	}
	defer rows.Close()

	type ProductWithStock struct {
		models.Product
		StockTotal int `json:"stock_total"`
	}

	var products []ProductWithStock
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo,
			&product.StockTotal)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct devuelve un producto por su id
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	var product models.Product
	var activo int
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
		&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}
	product.Activo = activo == 1

	c.JSON(http.StatusOK, product)
}

// CreateProduct crea un nuevo producto
func CreateProduct(c *gin.Context) {
	var reqData struct {
		Codigo       string   `json:"codigo" binding:"required"`
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description"`
		PrecioCompra float64  `json:"precio_compra"`
		PrecioVenta  float64  `json:"precio_venta" binding:"required"`
		Category     string   `json:"category" binding:"required"`
		Brand        string   `json:"brand"`
		ImageURL     string   `json:"image_url"`
		Images       []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir array de images a JSON string
	imagesJSON := "[]"
	if len(reqData.Images) > 0 {
		imagesBytes, _ := json.Marshal(reqData.Images)
		imagesJSON = string(imagesBytes)
	}

	result, err := database.DB.Exec(`
		INSERT INTO products (codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, reqData.Codigo, reqData.Name, reqData.Description, reqData.PrecioCompra, reqData.PrecioVenta, reqData.Category,
		reqData.Brand, reqData.ImageURL, imagesJSON)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	productID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product id"})
		return
	}

	// Devolver el producto creado
	product := models.Product{
		ID:           productID,
		Codigo:       reqData.Codigo,
		Name:         reqData.Name,
		Description:  reqData.Description,
		PrecioCompra: reqData.PrecioCompra,
		PrecioVenta:  reqData.PrecioVenta,
		Category:     reqData.Category,
		Brand:        reqData.Brand,
		ImageURL:     reqData.ImageURL,
		Images:       imagesJSON,
		Activo:       true,
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct actualiza un producto existente
func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	// Verificar que el producto existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updateData struct {
		Codigo       *string          `json:"codigo"`
		Name         *string          `json:"name"`
		Description  *string          `json:"description"`
		PrecioCompra *float64         `json:"precio_compra"`
		PrecioVenta  *float64         `json:"precio_venta"`
		Category     *string          `json:"category"`
		Brand        *string          `json:"brand"`
		ImageURL     *string          `json:"image_url"`
		Images       *json.RawMessage `json:"images"`
		Activo       *bool            `json:"activo"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos proporcionados
	updates := []string{}
	args := []interface{}{}

	if updateData.Codigo != nil {
		updates = append(updates, "codigo = ?")
		args = append(args, *updateData.Codigo)
	}
	if updateData.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *updateData.Name)
	}
	if updateData.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *updateData.Description)
	}
	if updateData.PrecioCompra != nil {
		updates = append(updates, "precio_compra = ?")
		args = append(args, *updateData.PrecioCompra)
	}
	if updateData.PrecioVenta != nil {
		updates = append(updates, "precio_venta = ?")
		args = append(args, *updateData.PrecioVenta)
	}
	if updateData.Category != nil {
		updates = append(updates, "category = ?")
		args = append(args, *updateData.Category)
	}
	if updateData.Brand != nil {
		updates = append(updates, "brand = ?")
		args = append(args, *updateData.Brand)
	}
	if updateData.ImageURL != nil {
		updates = append(updates, "image_url = ?")
		args = append(args, *updateData.ImageURL)
	}
	if updateData.Images != nil {
		imagesStr := string(*updateData.Images)
		updates = append(updates, "images = ?")
		args = append(args, imagesStr)
	}
	if updateData.Activo != nil {
		activo := 0
		if *updateData.Activo {
			activo = 1
		}
		updates = append(updates, "activo = ?")
		args = append(args, activo)
	}

	if len(updates) > 0 {
		updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, id)

		query := "UPDATE products SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		_, err := database.DB.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}
	}

	// Recargar el producto para obtener los datos actualizados
	var product models.Product
	var activo int
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
		&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)
	product.Activo = activo == 1

	c.JSON(http.StatusOK, product)
}

// DeleteProduct elimina un producto (soft delete - desactiva)
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	// Verificar que el producto existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Soft delete - solo desactivar
	_, err = database.DB.Exec("UPDATE products SET activo = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetRandomProducts devuelve productos aleatorios
func GetRandomProducts(c *gin.Context) {
	limit := 8
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE activo = 1
		ORDER BY RANDOM()
		LIMIT ?
	`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random products"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// GetProductsByCategory devuelve productos de una categoría específica
func GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")

	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE category = ? AND activo = 1
		ORDER BY created_at DESC
	`, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products by category"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}
