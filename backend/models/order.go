package models

import (
	"time"
)

// Order representa una orden de compra
type Order struct {
	ID            int64       `json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	SessionID     string      `json:"session_id"`
	PayPalOrderID string      `json:"paypal_order_id"`
	Status        string      `json:"status"` // PENDING, COMPLETED, CANCELLED
	TotalAmount   float64     `json:"total_amount"`
	Currency      string      `json:"currency"`
	PayerEmail    string      `json:"payer_email"`
	PayerName     string      `json:"payer_name"`
	CompletedAt   *time.Time  `json:"completed_at,omitempty"`
	Finalized     bool        `json:"finalized"`
	OrderItems    []OrderItem `json:"order_items,omitempty"`
}

// OrderItem representa un producto dentro de una orden
type OrderItem struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderID     int64     `json:"order_id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
}
