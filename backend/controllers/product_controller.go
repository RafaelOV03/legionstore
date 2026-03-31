package controllers

import (
	"encoding/json"
	"smartech/backend/database"
<<<<<<< HEAD
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/repository"
	"smartech/backend/validation"
=======
	"smartech/backend/repositories"
	"smartech/backend/services"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	"strconv"

	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
var productRepo *repository.ProductRepository

// InitProductRepository inicializa el repositorio de productos
func InitProductRepository() {
	productRepo = repository.NewProductRepository(database.DB)
=======
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
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetProducts devuelve todos los productos
func GetProducts(c *gin.Context) {
<<<<<<< HEAD
	products, err := productRepo.GetAll()
=======
	products, err := getProductService().ListProducts()
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch products", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	c.JSON(200, products)
=======
	c.JSON(http.StatusOK, products)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetProduct devuelve un producto por su id
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewBadRequest("Invalid product id format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	product, err := productRepo.GetByID(id)
=======
	product, err := getProductService().GetProduct(id)
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
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
<<<<<<< HEAD
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
=======
	var reqData createProductRequest
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7

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

<<<<<<< HEAD
	// Convertir array de images a JSON string
	imagesJSON := "[]"
	if len(reqData.Images) > 0 {
		imagesBytes, _ := json.Marshal(reqData.Images)
		imagesJSON = string(imagesBytes)
	}

	product := &models.Product{
=======
	product, err := getProductService().CreateProduct(repositories.ProductCreateInput{
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
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

<<<<<<< HEAD
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
=======
	var updateData updateProductRequest
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7

	if err := c.ShouldBindJSON(&updateData); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}
<<<<<<< HEAD

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
=======
	product, err := getProductService().UpdateProduct(id, updateData.toRepositoryInput())
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7

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

<<<<<<< HEAD
	// Verificar que el producto existe
	existing, err := productRepo.GetByID(id)
	if err != nil || existing == nil {
		apiErr := errors.NewNotFound("Product", id)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	err = productRepo.Delete(id)
=======
	err = getProductService().DeleteProduct(id)
	if err == repositories.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
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

<<<<<<< HEAD
	products, err := productRepo.GetRandomProducts(limit)
=======
	products, err := getProductService().GetRandomProducts(limit)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch random products", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	c.JSON(200, products)
=======
	c.JSON(http.StatusOK, products)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetProductsByCategory devuelve productos de una categoría específica
func GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
<<<<<<< HEAD
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
=======

	products, err := getProductService().GetProductsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products by category"})
		return
	}
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7

	c.JSON(200, products)
}
