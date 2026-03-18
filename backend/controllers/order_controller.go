package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"smartech/backend/database"
	"smartech/backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PayPalConfig contiene la configuración de PayPal
type PayPalConfig struct {
	Clientid string
	Secret   string
	BaseURL  string // sandbox o production
}

var paypalConfig PayPalConfig

func init() {
	// Configuración de PayPal (usar variables de entorno en producción)
	paypalConfig = PayPalConfig{
		Clientid: os.Getenv("PAYPAL_CLIENT_id"),
		Secret:   os.Getenv("PAYPAL_SECRET"),
		BaseURL:  "https://api-m.sandbox.paypal.com", // Sandbox por defecto
	}

	// Si no hay variables de entorno, usar valores por defecto para desarrollo
	if paypalConfig.Clientid == "" {
		paypalConfig.Clientid = "Aec0f6y57ztt6tvhADUXj1GAZqlH_vF0hGyK89cWblFEs3Gq3kunXw1iaqGeKe32CjPMuxy3PMuzLz6A" // Reemplazar con tu Client id de Sandbox
	}
	if paypalConfig.Secret == "" {
		paypalConfig.Secret = "EAF3VXSIZfDhEr522LM9UAyVOGn0kVvyiJBaJKPah4iRweg6MpAv2IZR8TDNhqj-9GUZuS17P066shVf" // Reemplazar con tu Secret de Sandbox
	}
}

// getPayPalAccessToken obtiene el token de acceso de PayPal
func getPayPalAccessToken() (string, error) {
	url := fmt.Sprintf("%s/v1/oauth2/token", paypalConfig.BaseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(paypalConfig.Clientid, paypalConfig.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token in response")
	}

	return token, nil
}

// CreateOrder crea una orden de PayPal
func CreateOrder(c *gin.Context) {
	var request struct {
		SessionID string `json:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el carrito con items
	var cart models.Cart
	err := database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, request.SessionID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	// Obtener items del carrito con productos
	cart.CartItems = getCartItemsWithProducts(cart.ID)

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calcular el total
	var total float64
	var items []map[string]interface{}

	for _, item := range cart.CartItems {
		itemTotal := float64(item.Quantity) * item.Product.PrecioVenta
		total += itemTotal

		items = append(items, map[string]interface{}{
			"name": item.Product.Name,
			"unit_amount": map[string]interface{}{
				"currency_code": "USD",
				"value":         fmt.Sprintf("%.2f", item.Product.PrecioVenta),
			},
			"quantity": fmt.Sprintf("%d", item.Quantity),
		})
	}

	// Crear orden en PayPal
	token, err := getPayPalAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal token"})
		return
	}

	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", total),
					"breakdown": map[string]interface{}{
						"item_total": map[string]interface{}{
							"currency_code": "USD",
							"value":         fmt.Sprintf("%.2f", total),
						},
					},
				},
				"items": items,
			},
		},
	}

	orderJSON, _ := json.Marshal(orderData)

	url := fmt.Sprintf("%s/v2/checkout/orders", paypalConfig.BaseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(orderJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PayPal order"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var paypalOrder map[string]interface{}
	if err := json.Unmarshal(body, &paypalOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse PayPal response"})
		return
	}

	if resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PayPal order creation failed", "details": paypalOrder})
		return
	}

	// Extraer el id de la orden de PayPal
	paypalOrderid, ok := paypalOrder["id"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal order id", "response": paypalOrder})
		return
	}

	log.Printf("Creating order with PayPal id: %s", paypalOrderid)

	// Guardar orden en la base de datos
	result, err := database.DB.Exec(`
		INSERT INTO orders (session_id, paypal_order_id, status, total_amount, currency, finalized, created_at, updated_at)
		VALUES (?, ?, 'PENDING', ?, 'USD', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, request.SessionID, paypalOrderid, total)

	if err != nil {
		log.Printf("Failed to save order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order", "details": err.Error()})
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order id"})
		return
	}

	// Guardar items de la orden
	for _, item := range cart.CartItems {
		_, err := database.DB.Exec(`
			INSERT INTO order_items (order_id, product_id, product_name, quantity, price, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, orderID, item.ProductID, item.Product.Name, item.Quantity, item.Product.PrecioVenta)
		if err != nil {
			log.Printf("Failed to save order item: %v", err)
		}
	}

	// Crear objeto order para respuesta
	order := models.Order{
		ID:            orderID,
		SessionID:     request.SessionID,
		PayPalOrderID: paypalOrderid,
		Status:        "PENDING",
		TotalAmount:   total,
		Currency:      "USD",
		Finalized:     false,
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     paypalOrder["id"],
		"order":  order,
		"paypal": paypalOrder,
	})
}

// CaptureOrder captura el pago de una orden de PayPal
func CaptureOrder(c *gin.Context) {
	orderid := c.Param("id")
	log.Printf("Capturing order with PayPal id: %s", orderid)

	// Obtener token de PayPal
	token, err := getPayPalAccessToken()
	if err != nil {
		log.Printf("Failed to get PayPal token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal token"})
		return
	}

	// Capturar el pago en PayPal
	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypalConfig.BaseURL, orderid)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to capture PayPal order"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var captureResult map[string]interface{}
	if err := json.Unmarshal(body, &captureResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse PayPal response"})
		return
	}

	if resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PayPal capture failed", "details": captureResult})
		return
	}

	// Actualizar orden en la base de datos
	var order models.Order
	log.Printf("Searching for order with paypal_order_id = %s", orderid)
	err = database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount, 
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE paypal_order_id = ?
	`, orderid).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &order.Finalized)

	if err == sql.ErrNoRows {
		log.Printf("Order not found in database: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found", "paypal_order_id": orderid})
		return
	}
	if err != nil {
		log.Printf("Failed to fetch order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}

	now := time.Now()
	payerEmail := ""
	payerName := ""

	// Extraer información del pagador
	if payer, ok := captureResult["payer"].(map[string]interface{}); ok {
		if email, ok := payer["email_address"].(string); ok {
			payerEmail = email
		}
		if name, ok := payer["name"].(map[string]interface{}); ok {
			if givenName, ok := name["given_name"].(string); ok {
				payerName = givenName
				if surname, ok := name["surname"].(string); ok {
					payerName += " " + surname
				}
			}
		}
	}

	// Actualizar orden
	_, err = database.DB.Exec(`
		UPDATE orders
		SET status = 'COMPLETED', completed_at = ?, payer_email = ?, payer_name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, now, payerEmail, payerName, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	order.Status = "COMPLETED"
	order.CompletedAt = &now
	order.PayerEmail = payerEmail
	order.PayerName = payerName

	// Actualizar stock de productos
	rows, err := database.DB.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", order.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var productID int64
			var quantity int
			if err := rows.Scan(&productID, &quantity); err == nil {
				database.DB.Exec(`
					UPDATE products
					SET stock_quantity = stock_quantity - ?
					WHERE id = ?
				`, quantity, productID)
			}
		}
	}

	// Limpiar el carrito
	var cartID int64
	err = database.DB.QueryRow("SELECT id FROM carts WHERE session_id = ?", order.SessionID).Scan(&cartID)
	if err == nil {
		database.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
	}

	c.JSON(http.StatusOK, gin.H{
		"order":   order,
		"capture": captureResult,
	})
}

// getOrderItemsWithProducts obtiene los items de una orden
func getOrderItemsWithProducts(orderID int64) []models.OrderItem {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, order_id, product_id, product_name, quantity, price
		FROM order_items
		WHERE order_id = ?
	`, orderID)
	if err != nil {
		return []models.OrderItem{}
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.OrderID,
			&item.ProductID, &item.ProductName, &item.Quantity, &item.Price)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items
}

// GetOrder obtiene los detalles de una orden
func GetOrder(c *gin.Context) {
	orderid := c.Param("id")

	var order models.Order
	var finalized int
	err := database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE paypal_order_id = ?
	`, orderid).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &finalized)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}

	order.Finalized = finalized == 1
	order.OrderItems = getOrderItemsWithProducts(order.ID)

	c.JSON(http.StatusOK, order)
}

// GetOrders obtiene todas las órdenes de una sesión
func GetOrders(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE session_id = ?
		ORDER BY created_at DESC
	`, sessionid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var finalized int
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
			&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
			&order.CompletedAt, &finalized)
		if err != nil {
			continue
		}
		order.Finalized = finalized == 1
		order.OrderItems = getOrderItemsWithProducts(order.ID)
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

// GetAllOrders obtiene todas las órdenes del sistema (solo para administradores)
func GetAllOrders(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var finalized int
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
			&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
			&order.CompletedAt, &finalized)
		if err != nil {
			continue
		}
		order.Finalized = finalized == 1
		order.OrderItems = getOrderItemsWithProducts(order.ID)
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

// FinalizeOrder marca una orden como finalizada
func FinalizeOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}

	// Verificar que la orden existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", orderID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE orders SET finalized = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al finalizar orden"})
		return
	}

	// Obtener la orden actualizada
	var order models.Order
	var finalized int
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE id = ?
	`, orderID).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &finalized)

	order.Finalized = finalized == 1

	c.JSON(http.StatusOK, order)
}

// UpdateOrder actualiza el estado de una orden (solo administradores)
func UpdateOrder(c *gin.Context) {
	orderid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}

	// Verificar que la orden existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", orderid).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var req struct {
		Status      string  `json:"status"`
		TotalAmount float64 `json:"total_amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos
	updates := []string{}
	args := []interface{}{}

	if req.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}
	if req.TotalAmount > 0 {
		updates = append(updates, "total_amount = ?")
		args = append(args, req.TotalAmount)
	}

	if len(updates) > 0 {
		updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, orderid)

		query := "UPDATE orders SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		_, err := database.DB.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}
	}

	// Recargar con items
	var order models.Order
	var finalized int
	database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE id = ?
	`, orderid).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &finalized)

	order.Finalized = finalized == 1
	order.OrderItems = getOrderItemsWithProducts(order.ID)

	c.JSON(http.StatusOK, order)
}

// DeleteOrder elimina una orden (solo administradores)
func DeleteOrder(c *gin.Context) {
	orderid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}

	// Verificar que la orden existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", orderid).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Eliminar items de la orden primero
	_, err = database.DB.Exec("DELETE FROM order_items WHERE order_id = ?", orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order items"})
		return
	}

	// Eliminar la orden
	_, err = database.DB.Exec("DELETE FROM orders WHERE id = ?", orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// GetPayPalConfig devuelve el Client id para el frontend
func GetPayPalConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"clientId": paypalConfig.Clientid,
	})
}
