package controllers

import (
	"net/http"
	"smartech/backend/database"
<<<<<<< HEAD
	"smartech/backend/errors"
	"smartech/backend/models"
	"smartech/backend/validation"
=======
	"smartech/backend/repositories"
	"smartech/backend/services"
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	"strconv"

	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
// PayPalConfig contiene la configuración de PayPal
type PayPalConfig struct {
	Clientid string
	Secret   string
	BaseURL  string // sandbox o production
}

var paypalConfig PayPalConfig

func init() {
	// Configuración de PayPal (desde variables de entorno o valores por defecto para desarrollo)
	paypalConfig = PayPalConfig{
		Clientid: os.Getenv("PAYPAL_CLIENT_ID"),
		Secret:   os.Getenv("PAYPAL_SECRET"),
		BaseURL:  os.Getenv("PAYPAL_BASE_URL"),
	}

	// Si no hay variables de entorno, usar valores por defecto para desarrollo
	if paypalConfig.Clientid == "" {
		log.Println("Warning: PAYPAL_CLIENT_ID not configured, using development default")
		paypalConfig.Clientid = "Aec0f6y57ztt6tvhADUXj1GAZqlH_vF0hGyK89cWblFEs3Gq3kunXw1iaqGeKe32CjPMuxy3PMuzLz6A"
	}
	if paypalConfig.Secret == "" {
		log.Println("Warning: PAYPAL_SECRET not configured, using development default")
		paypalConfig.Secret = "EAF3VXSIZfDhEr522LM9UAyVOGn0kVvyiJBaJKPah4iRweg6MpAv2IZR8TDNhqj-9GUZuS17P066shVf"
	}

	// Si BaseURL no está configurada, usar sandbox por defecto
	if paypalConfig.BaseURL == "" {
		paypalConfig.BaseURL = "https://api-m.sandbox.paypal.com"
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
=======
func getOrderService() *services.OrderService {
	repo := repositories.NewOrderRepository(database.DB)
	return services.NewOrderService(repo)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// CreateOrder crea una orden de PayPal
func CreateOrder(c *gin.Context) {
	var request struct {
		SessionID string `json:"session_id" validate:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequest(err.Error())
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validar estructura
	validationErrors := validation.ValidateStruct(request)
	if len(validationErrors) > 0 {
		c.JSON(422, validationErrors.ToAPIError())
		return
	}

<<<<<<< HEAD
	// Obtener el carrito con items
	var cart models.Cart
	err := database.DB.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, request.SessionID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)

	if err == sql.ErrNoRows {
		apiErr := errors.NewNotFound("Cart", request.SessionID)
		c.JSON(apiErr.Code, apiErr)
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch cart", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Obtener items del carrito con productos
	cart.CartItems = getCartItemsWithProducts(cart.ID)

	if len(cart.CartItems) == 0 {
		apiErr := errors.NewBadRequest("Cart is empty")
		c.JSON(apiErr.Code, apiErr)
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
		apiErr := errors.NewInternal("Failed to get PayPal token")
		c.JSON(apiErr.Code, apiErr)
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
		apiErr := errors.NewInternal("Failed to communicate with PayPal")
		c.JSON(apiErr.Code, apiErr)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var paypalOrder map[string]interface{}
	if err := json.Unmarshal(body, &paypalOrder); err != nil {
		apiErr := errors.NewInternal("Failed to parse PayPal response")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		apiErr := errors.NewInternal("PayPal order creation failed")
		c.JSON(apiErr.Code, apiErr)
=======
	order, paypalOrder, err := getOrderService().CreateOrder(request.SessionID)
	if err == services.ErrOrderCartNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err == services.ErrOrderCartEmpty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}
	if err == services.ErrOrderPayPalToken {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal token"})
		return
	}
	if err == services.ErrOrderPayPalCreate {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PayPal order creation failed", "details": paypalOrder})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
<<<<<<< HEAD
		apiErr := errors.NewDatabaseError("Insert order", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		apiErr := errors.NewDatabaseError("Get order ID", err)
		c.JSON(apiErr.Code, apiErr)
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

=======
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order", "details": err.Error()})
		return
	}

	paypalID, _ := paypalOrder["id"].(string)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
	c.JSON(http.StatusOK, gin.H{
		"id":     paypalID,
		"order":  order,
		"paypal": paypalOrder,
	})
}

// CaptureOrder captura el pago de una orden de PayPal
func CaptureOrder(c *gin.Context) {
	orderid := c.Param("id")

	order, captureResult, err := getOrderService().CaptureOrder(orderid)
	if err == services.ErrOrderPayPalToken {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal token"})
		return
	}
	if err == services.ErrOrderPayPalCapture {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PayPal capture failed", "details": captureResult})
		return
	}
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found", "paypal_order_id": orderid})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":   order,
		"capture": captureResult,
	})
}

// GetOrder obtiene los detalles de una orden
func GetOrder(c *gin.Context) {
	orderid := c.Param("id")
	if orderid == "" {
		apiErr := errors.NewBadRequest("Order id is required")
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
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
		apiErr := errors.NewNotFound("Order", orderid)
		c.JSON(apiErr.Code, apiErr)
=======
	order, err := getOrderService().GetOrder(orderid)
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
		return
	}
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch order", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

<<<<<<< HEAD
	order.Finalized = finalized == 1
	order.OrderItems = getOrderItemsWithProducts(order.ID)

	c.JSON(200, order)
=======
	c.JSON(http.StatusOK, order)
>>>>>>> 56ef4a99558720e22eaa0ffde0aef19a608948d7
}

// GetOrders obtiene todas las órdenes de una sesión
func GetOrders(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		apiErr := errors.NewBadRequest("session_id query parameter is required")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	orders, err := getOrderService().GetOrdersBySession(sessionid)
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch orders", err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(200, orders)
}

// GetAllOrders obtiene todas las órdenes del sistema (solo para administradores)
func GetAllOrders(c *gin.Context) {
	orders, err := getOrderService().GetAllOrders()
	if err != nil {
		apiErr := errors.NewDatabaseError("Fetch all orders", err)
		c.JSON(apiErr.Code, apiErr)
		return
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

	order, err := getOrderService().FinalizeOrder(orderID)
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al finalizar orden"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder actualiza el estado de una orden (solo administradores)
func UpdateOrder(c *gin.Context) {
	orderid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
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

	order, err := getOrderService().UpdateOrder(orderid, req.Status, req.TotalAmount)
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder elimina una orden (solo administradores)
func DeleteOrder(c *gin.Context) {
	orderid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}

	err = getOrderService().DeleteOrder(orderid)
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// GetPayPalConfig devuelve el Client id para el frontend
func GetPayPalConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"clientId": getOrderService().GetPayPalClientID(),
	})
}
