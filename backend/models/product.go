package models

import (
	"time"
)

// Product representa un producto en el inventario
type Product struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Codigo       string    `json:"codigo"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PrecioCompra float64   `json:"precio_compra"`
	PrecioVenta  float64   `json:"precio_venta"`
	Category     string    `json:"category"`
	Brand        string    `json:"brand"`
	ImageURL     string    `json:"image_url"`
	Images       string    `json:"images"` // JSON string con array de imágenes adicionales
	Activo       bool      `json:"activo"`
}
