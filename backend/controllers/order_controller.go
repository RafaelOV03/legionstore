package controllers

import (
	"net/http"
	"smartech/backend/database"
	"smartech/backend/repositories"
	"smartech/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getOrderService() *services.OrderService {
	repo := repositories.NewOrderRepository(database.DB)
	return services.NewOrderService(repo)
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
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order", "details": err.Error()})
		return
	}

	paypalID, _ := paypalOrder["id"].(string)
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

	order, err := getOrderService().GetOrder(orderid)
	if err == services.ErrOrderNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders obtiene todas las órdenes de una sesión
func GetOrders(c *gin.Context) {
	sessionid := c.Query("session_id")
	if sessionid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	orders, err := getOrderService().GetOrdersBySession(sessionid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetAllOrders obtiene todas las órdenes del sistema (solo para administradores)
func GetAllOrders(c *gin.Context) {
	orders, err := getOrderService().GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
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
