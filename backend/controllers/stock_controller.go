package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getStockService() *services.StockService {
	repo := repositories.NewStockRepository(database.DB)
	return services.NewStockService(repo)
}

// GetStockMultisede obtiene el stock de todos los productos por sede
func GetStockMultisede(c *gin.Context) {
	items, err := getStockService().GetStockMultisede()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener stock multisede"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetStockBySede obtiene el stock de productos para una sede específica
func GetStockBySede(c *gin.Context) {
	sedeID, err := strconv.ParseInt(c.Param("sede_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sede inválido"})
		return
	}

	items, err := getStockService().GetStockBySede(sedeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener stock de sede"})
		return
	}

	c.JSON(http.StatusOK, items)
}

type updateStockRequest struct {
	SedeID      int64 `json:"sede_id" binding:"required"`
	ProductoID  int64 `json:"producto_id" binding:"required"`
	Cantidad    int   `json:"cantidad"`
	StockMinimo int   `json:"stock_minimo"`
	StockMaximo int   `json:"stock_maximo"`
}

// UpdateStock crea o actualiza el stock de un producto en una sede
func UpdateStock(c *gin.Context) {
	var req updateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	input := repositories.StockUpdateInput{
		SedeID:      req.SedeID,
		ProductoID:  req.ProductoID,
		Cantidad:    req.Cantidad,
		StockMinimo: req.StockMinimo,
		StockMaximo: req.StockMaximo,
	}

	if err := getStockService().UpdateStock(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock actualizado correctamente"})
}
