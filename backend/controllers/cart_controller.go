package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCartService() *services.CartService {
	repo := repositories.NewCartRepository(database.DB)
	return services.NewCartService(repo)
}

// GetCart devuelve el carrito de compras para una sesión
func GetCart(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	cart, err := getCartService().GetCart(sessionid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart agrega un producto al carrito
func AddToCart(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	var request struct {
		ProductID int64 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := getCartService().AddToCart(services.AddToCartInput{
		SessionID: sessionid,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	})
	if err == services.ErrProductNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err == services.ErrInsufficientStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// UpdateCartItem actualiza la cantidad de un producto en el carrito
func UpdateCartItem(c *gin.Context) {
	itemid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item id"})
		return
	}

	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	var request struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItem, err := getCartService().UpdateCartItem(services.UpdateCartItemInput{
		SessionID: sessionid,
		ItemID:    itemid,
		Quantity:  request.Quantity,
	})
	if err == services.ErrCartItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	if err == services.ErrCartUnauthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}
	c.JSON(http.StatusOK, cartItem)
}

// RemoveFromCart elimina un producto del carrito
func RemoveFromCart(c *gin.Context) {
	itemid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item id"})
		return
	}

	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	err = getCartService().RemoveFromCart(sessionid, itemid)
	if err == services.ErrCartItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	if err == services.ErrCartUnauthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ClearCart vacía todo el carrito
func ClearCart(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	err := getCartService().ClearCart(sessionid)
	if err == services.ErrCartNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
