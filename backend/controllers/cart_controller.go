package controllers

import (
	"database/sql"
	"net/http"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getCartItemsWithProducts obtiene los items del carrito con información del producto
func getCartItemsWithProducts(cartID int64) []models.CartItem {
	rows, err := database.DB.Query(`
		SELECT ci.id, ci.created_at, ci.updated_at, ci.cart_id, ci.product_id, ci.quantity,
		       p.id, p.created_at, p.updated_at, p.name, p.description, p.precio_venta, p.precio_compra,
		       p.category, p.brand, p.image_url, p.images, p.activo
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = ?
	`, cartID)
	if err != nil {
		return []models.CartItem{}
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var product models.Product
		var activo int
		err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.CartID, &item.ProductID, &item.Quantity,
			&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Name, &product.Description,
			&product.PrecioVenta, &product.PrecioCompra, &product.Category, &product.Brand,
			&product.ImageURL, &product.Images, &activo,
		)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		item.Product = &product
		items = append(items, item)
	}
	return items
}

// GetCart devuelve el carrito de compras para una sesión
func GetCart(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	var cart models.Cart
	err := database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, sessionid).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)

	if err == sql.ErrNoRows {
		// Si no existe, crear un nuevo carrito
		result, err := database.DB.Exec(`
			INSERT INTO carts (session_id, created_at, updated_at)
			VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, sessionid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		cartID, _ := result.LastInsertId()
		cart.ID = cartID
		cart.SessionID = sessionid
		cart.CartItems = []models.CartItem{}
		c.JSON(http.StatusOK, cart)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	// Obtener items del carrito con productos
	cart.CartItems = getCartItemsWithProducts(cart.ID)
	if cart.CartItems == nil {
		cart.CartItems = []models.CartItem{}
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

	// Verificar que el producto existe y tiene stock
	var stockQuantity int
	err := database.DB.QueryRow(`
		SELECT stock_quantity FROM products WHERE id = ?
	`, request.ProductID).Scan(&stockQuantity)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if stockQuantity < request.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Buscar o crear el carrito
	var cart models.Cart
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, sessionid).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)

	if err == sql.ErrNoRows {
		result, err := database.DB.Exec(`
			INSERT INTO carts (session_id, created_at, updated_at)
			VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, sessionid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		cartID, _ := result.LastInsertId()
		cart.ID = cartID
		cart.SessionID = sessionid
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	// Verificar si el producto ya está en el carrito
	var existingItemID int64
	var existingQuantity int
	err = database.DB.QueryRow(`
		SELECT id, quantity FROM cart_items 
		WHERE cart_id = ? AND product_id = ?
	`, cart.ID, request.ProductID).Scan(&existingItemID, &existingQuantity)

	if err == nil {
		// Actualizar la cantidad
		newQuantity := existingQuantity + request.Quantity
		_, err := database.DB.Exec(`
			UPDATE cart_items SET quantity = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, newQuantity, existingItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
	} else if err == sql.ErrNoRows {
		// Crear nuevo item
		_, err := database.DB.Exec(`
			INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, cart.ID, request.ProductID, request.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cart item"})
		return
	}

	// Devolver el carrito actualizado
	cart.CartItems = getCartItemsWithProducts(cart.ID)
	if cart.CartItems == nil {
		cart.CartItems = []models.CartItem{}
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

	// Obtener el cart item con producto
	var cartItem models.CartItem
	var product models.Product
	var activo int
	err = database.DB.QueryRow(`
		SELECT ci.id, ci.created_at, ci.updated_at, ci.cart_id, ci.product_id, ci.quantity,
		       p.id, p.created_at, p.updated_at, p.name, p.description, p.precio_venta, p.precio_compra,
		       p.category, p.brand, p.image_url, p.images, p.activo
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		WHERE ci.id = ?
	`, itemid).Scan(
		&cartItem.ID, &cartItem.CreatedAt, &cartItem.UpdatedAt, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity,
		&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Name, &product.Description,
		&product.PrecioVenta, &product.PrecioCompra, &product.Category, &product.Brand,
		&product.ImageURL, &product.Images, &activo,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	product.Activo = activo == 1
	cartItem.Product = &product

	// Verificar que el cart_item pertenece al carrito del usuario
	var cartSessionID string
	err = database.DB.QueryRow("SELECT session_id FROM carts WHERE id = ?", cartItem.CartID).Scan(&cartSessionID)
	if err != nil || cartSessionID != sessionid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Usar la cantidad solicitada
	adjustedQuantity := request.Quantity

	_, err = database.DB.Exec(`
		UPDATE cart_items SET quantity = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, adjustedQuantity, itemid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	cartItem.Quantity = adjustedQuantity

	// Si se ajustó la cantidad, informar al cliente
	if adjustedQuantity < request.Quantity {
		c.JSON(http.StatusOK, gin.H{
			"item":     cartItem,
			"adjusted": true,
			"message":  "Quantity adjusted to available stock",
		})
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

	// Obtener el cart item
	var cartID int64
	err = database.DB.QueryRow("SELECT cart_id FROM cart_items WHERE id = ?", itemid).Scan(&cartID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	// Verificar que el carrito pertenece a esta sesión
	var cartSessionID string
	err = database.DB.QueryRow("SELECT session_id FROM carts WHERE id = ?", cartID).Scan(&cartSessionID)
	if err != nil || cartSessionID != sessionid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM cart_items WHERE id = ?", itemid)
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

	var cartID int64
	err := database.DB.QueryRow("SELECT id FROM carts WHERE session_id = ?", sessionid).Scan(&cartID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
