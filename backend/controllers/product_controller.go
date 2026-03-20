package controllers

import (
	"encoding/json"
	"smartech/backend/database"
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/repository"
	"smartech/backend/validation"
	"strconv"

	"github.com/gin-gonic/gin"
)

var productRepo *repository.ProductRepository

// InitProductRepository inicializa el repositorio de productos
func InitProductRepository() {
	productRepo = repository.NewProductRepository(database.DB)
}

// GetProducts devuelve todos los productos
func GetProducts(c *gin.Context) {
	products, err := productRepo.GetAll()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch products", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, products)
}

// GetProduct devuelve un producto por su id
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid product id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	product, err := productRepo.GetByID(id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch product", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	if product == nil {
		apiErr := errors.NewNotFound("Product", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, product)
}

// CreateProduct crea un nuevo producto
func CreateProduct(c *gin.Context) {
	var reqData struct {
		Codigo       string   `json:"codigo" validate:"required,min=1"`
		Name         string   `json:"name" validate:"required,min=3"`
		Description  string   `json:"description"`
		PrecioCompra float64  `json:"precio_compra" validate:"gte=0"`
		PrecioVenta  float64  `json:"precio_venta" validate:"required,gt=0"`
		Category     string   `json:"category" validate:"required,min=1"`
		Brand        string   `json:"brand"`
		ImageURL     string   `json:"image_url"`
		Images       []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	validationErrors := validation.ValidateStruct(reqData)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

	// Convertir array de images a JSON string
	imagesJSON := "[]"
	if len(reqData.Images) > 0 {
		imagesBytes, _ := json.Marshal(reqData.Images)
		imagesJSON = string(imagesBytes)
	}

	product := &models.Product{
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

	err := productRepo.Create(product)
	if err != nil {
		apiErr := errors.NewDatabaseError("Create product", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(201, product)
}

// UpdateProduct actualiza un producto existente
func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid product id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar que el producto existe
	existing, err := productRepo.GetByID(id)
	if err != nil || existing == nil {
		apiErr := errors.NewNotFound("Product", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	var updateData struct {
		Codigo       string  `json:"codigo"`
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		PrecioCompra float64 `json:"precio_compra"`
		PrecioVenta  float64 `json:"precio_venta"`
		Category     string  `json:"category"`
		Brand        string  `json:"brand"`
		ImageURL     string  `json:"image_url"`
		Images       string  `json:"images"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Usar valores existentes si no se proporcionan
	if updateData.Codigo == "" {
		updateData.Codigo = existing.Codigo
	}
	if updateData.Name == "" {
		updateData.Name = existing.Name
	}
	if updateData.PrecioVenta == 0 {
		updateData.PrecioVenta = existing.PrecioVenta
	}
	if updateData.Category == "" {
		updateData.Category = existing.Category
	}

	product := &models.Product{
		Codigo:       updateData.Codigo,
		Name:         updateData.Name,
		Description:  updateData.Description,
		PrecioCompra: updateData.PrecioCompra,
		PrecioVenta:  updateData.PrecioVenta,
		Category:     updateData.Category,
		Brand:        updateData.Brand,
		ImageURL:     updateData.ImageURL,
		Images:       updateData.Images,
	}

	err = productRepo.Update(id, product)
	if err != nil {
		apiErr := errors.NewDatabaseError("Update product", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Recargar para devolver datos actualizado
	updated, _ := productRepo.GetByID(id)
	c.JSON(200, updated)
}

// DeleteProduct elimina un producto (soft delete)
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid product id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Verificar que el producto existe
	existing, err := productRepo.GetByID(id)
	if err != nil || existing == nil {
		apiErr := errors.NewNotFound("Product", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	err = productRepo.Delete(id)
	if err != nil {
		apiErr := errors.NewDatabaseError("Delete product", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}

// GetRandomProducts devuelve productos aleatorios
func GetRandomProducts(c *gin.Context) {
	limit := 8
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	products, err := productRepo.GetRandomProducts(limit)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch random products", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, products)
}

// GetProductsByCategory devuelve productos de una categoría específica
func GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		apiErr := errors.NewBadRequest("Category parameter is required")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	products, err := productRepo.GetByCategory(category)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch products by category", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, products)
}
