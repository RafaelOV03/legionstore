package models

import (
	"time"
)

// Cart representa un carrito de compras
type Cart struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	SessionID string     `json:"session_id"`
	CartItems []CartItem `json:"items,omitempty"`
}

// CartItem representa un producto dentro del carrito
type CartItem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CartID    int64     `json:"cart_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Product   *Product  `json:"product,omitempty"`
}
