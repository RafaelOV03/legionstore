package controllers

import (
	"encoding/json"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
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

type updateProductRequest struct {
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

func (r updateProductRequest) toRepositoryInput() repositories.ProductUpdateInput {
	var images *string
	if r.Images != nil {
		imagesValue := string(*r.Images)
		images = &imagesValue
	}

	return repositories.ProductUpdateInput{
		Codigo:       r.Codigo,
		Name:         r.Name,
		Description:  r.Description,
		PrecioCompra: r.PrecioCompra,
		PrecioVenta:  r.PrecioVenta,
		Category:     r.Category,
		Brand:        r.Brand,
		ImageURL:     r.ImageURL,
		Images:       images,
		Activo:       r.Activo,
	}
}

func getProductService() *services.ProductService {
	repo := repositories.NewProductRepository(database.DB)
	return services.NewProductService(repo)
}

// GetProducts devuelve todos los productos
func GetProducts(c *gin.Context) {
	products, err := getProductService().ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
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

	product, err := getProductService().GetProduct(id)
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct crea un nuevo producto
func CreateProduct(c *gin.Context) {
	var reqData createProductRequest

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := getProductService().CreateProduct(repositories.ProductCreateInput{
		Codigo:       reqData.Codigo,
		Name:         reqData.Name,
		Description:  reqData.Description,
		PrecioCompra: reqData.PrecioCompra,
		PrecioVenta:  reqData.PrecioVenta,
		Category:     reqData.Category,
		Brand:        reqData.Brand,
		ImageURL:     reqData.ImageURL,
		Images:       reqData.Images,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
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

	var updateData updateProductRequest

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := getProductService().UpdateProduct(id, updateData.toRepositoryInput())
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct elimina un producto (soft delete - desactiva)
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	err = getProductService().DeleteProduct(id)
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
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

	products, err := getProductService().GetRandomProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductsByCategory devuelve productos de una categoría específica
func GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")

	products, err := getProductService().GetProductsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products by category"})
		return
	}

	c.JSON(http.StatusOK, products)
}
