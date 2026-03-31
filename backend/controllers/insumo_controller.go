package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getInsumoService() *services.InsumoService {
	repo := repositories.NewInsumoRepository(database.DB)
	return services.NewInsumoService(repo)
}

// GetInsumos obtiene todos los insumos
func GetInsumos(c *gin.Context) {
	categoria := c.Query("categoria")
	bajoStock := c.Query("bajo_stock") == "true"

	insumos, err := getInsumoService().ListInsumos(categoria, bajoStock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insumos"})
		return
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

	i, err := getInsumoService().GetInsumo(id)
	if err == services.ErrInsumoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insumo"})
		return
	}

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

	id, err := getInsumoService().CreateInsumo(models.Insumo{
		Codigo:       req.Codigo,
		Nombre:       req.Nombre,
		Descripcion:  req.Descripcion,
		Categoria:    req.Categoria,
		UnidadMedida: req.UnidadMedida,
		Stock:        req.Stock,
		StockMinimo:  req.StockMinimo,
		Costo:        req.Costo,
		SedeID:       req.SedeID,
		Activo:       true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create insumo"})
		return
	}

	logAuditoria(c, "crear", "insumo", id, "", req.Nombre)
	c.JSON(http.StatusCreated, gin.H{"id": id})
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

	err = getInsumoService().UpdateInsumo(id, models.Insumo{
		Codigo:       req.Codigo,
		Nombre:       req.Nombre,
		Descripcion:  req.Descripcion,
		Categoria:    req.Categoria,
		UnidadMedida: req.UnidadMedida,
		Stock:        req.Stock,
		StockMinimo:  req.StockMinimo,
		Costo:        req.Costo,
		Activo:       req.Activo,
	})
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
		Tipo     string `json:"tipo" binding:"required"`
		Motivo   string `json:"motivo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stockActual, nuevoStock, err := getInsumoService().AjustarStock(id, services.AjusteStockInput{Cantidad: req.Cantidad, Tipo: req.Tipo, Motivo: req.Motivo})
	if err == services.ErrInsumoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insumo not found"})
		return
	}
	if err == services.ErrStockInsuficiente {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
		return
	}
	if err == services.ErrTipoAjusteInvalido {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido (debe ser entrada o salida)"})
		return
	}
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

	err = getInsumoService().DeleteInsumo(id)
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

	comps, err := getInsumoService().ListCompatibilidades(productoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch compatibilidades"})
		return
	}

	c.JSON(http.StatusOK, comps)
}

// BuscarCompatibles busca productos compatibles con uno dado
func BuscarCompatibles(c *gin.Context) {
	productoID := c.Param("id")

	res, err := getInsumoService().BuscarCompatibles(productoID)
	if err == services.ErrProductoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch compatibles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"producto": res.Producto, "compatibles": res.Compatibles})
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

	id, err := getInsumoService().CreateCompatibilidad(services.CreateCompatibilidadInput{ProductoID: req.ProductoID, CompatibleCon: req.CompatibleCon, Notas: req.Notas})
	if err == services.ErrCompatibilidadSelf {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Un producto no puede ser compatible consigo mismo"})
		return
	}
	if err == services.ErrCompatibilidadDup {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Esta compatibilidad ya existe"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create compatibilidad"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteCompatibilidad elimina una compatibilidad
func DeleteCompatibilidad(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compatibilidad ID"})
		return
	}

	err = getInsumoService().DeleteCompatibilidad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete compatibilidad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compatibilidad deleted successfully"})
}

// GetInsumosStats obtiene estadísticas de insumos
func GetInsumosStats(c *gin.Context) {
	stats, err := getInsumoService().Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insumo stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
