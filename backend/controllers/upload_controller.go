package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadProductImage sube una imagen de producto
func UploadProductImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Validar tipo de archivo
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Crear directorio si no existe
	uploadDir := "./uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generar nombre único para el archivo
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strings.ReplaceAll(file.Filename, " ", "_"), "")
	if len(filename) > 100 {
		filename = fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	}
	filepath := filepath.Join(uploadDir, filename)

	// Guardar archivo
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Retornar URL de la imagen
	imageURL := fmt.Sprintf("/uploads/products/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": imageURL})
}

// DeleteProductImage elimina una imagen de producto
func DeleteProductImage(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extraer el nombre del archivo de la URL
	// Manejar URLs completas (http://...) o relativas (/uploads/...)
	urlPath := request.URL
	if strings.Contains(urlPath, "http://") || strings.Contains(urlPath, "https://") {
		// Extraer solo la parte de la ruta
		parts := strings.Split(urlPath, "/")
		// Buscar "uploads" en la URL
		uploadsIndex := -1
		for i, part := range parts {
			if part == "uploads" {
				uploadsIndex = i
				break
			}
		}
		if uploadsIndex == -1 || uploadsIndex+2 >= len(parts) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image URL format"})
			return
		}
		urlPath = "/" + strings.Join(parts[uploadsIndex:], "/")
	}

	parts := strings.Split(urlPath, "/")
	if len(parts) < 3 || parts[1] != "uploads" || parts[2] != "products" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image URL path"})
		return
	}

	filename := parts[len(parts)-1]
	filepath := filepath.Join("./uploads/products", filename)

	// Eliminar archivo
	if err := os.Remove(filepath); err != nil {
		// No es un error crítico si el archivo no existe
		c.JSON(http.StatusOK, gin.H{"message": "Image removed from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// ServeUploadedFile sirve archivos estáticos del directorio uploads
func ServeUploads(c *gin.Context) {
	filepath := c.Param("filepath")
	fullPath := "./uploads/" + filepath

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(fullPath)
}

// GetProductImages obtiene todas las imágenes de un producto
func GetProductImages(c *gin.Context) {
	productid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	var images string
	err = database.DB.QueryRow("SELECT images FROM products WHERE id = ?", productid).Scan(&images)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	// Parsear el JSON de imágenes
	var imageList []string
	if images != "" {
		json.Unmarshal([]byte(images), &imageList)
	}

	c.JSON(http.StatusOK, gin.H{"images": imageList})
}

// UpdateProductImages actualiza las imágenes de un producto
func UpdateProductImages(c *gin.Context) {
	productid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	// Verificar que el producto existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productid).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var request struct {
		Images []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a JSON
	imagesJSON, err := json.Marshal(request.Images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode images"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE products SET images = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, string(imagesJSON), productid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update images"})
		return
	}

	// Obtener el producto actualizado
	var product models.Product
	var activo int
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, productid).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
		&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)
	product.Activo = activo == 1

	c.JSON(http.StatusOK, product)
}
