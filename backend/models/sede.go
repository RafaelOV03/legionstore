package models

import "time"

// Sede representa una ubicación física del inventario
type Sede struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Nombre    string    `json:"nombre"`
	Direccion string    `json:"direccion"`
	Telefono  string    `json:"telefono"`
	Activa    bool      `json:"activa"`
}

// StockSede representa el stock de un producto en una sede específica
type StockSede struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SedeID      int64     `json:"sede_id"`
	ProductoID  int64     `json:"producto_id"`
	Cantidad    int       `json:"cantidad"`
	StockMinimo int       `json:"stock_minimo"`
	StockMaximo int       `json:"stock_maximo"`
	Sede        *Sede     `json:"sede,omitempty"`
	Producto    *Product  `json:"producto,omitempty"`
}
