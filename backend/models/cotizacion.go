package models

import "time"

// Cotizacion representa una cotización para un cliente
type Cotizacion struct {
	ID               int64            `json:"id"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	NumeroCotizacion string           `json:"numero_cotizacion"`
	ClienteNombre    string           `json:"cliente_nombre"`
	ClienteTelefono  string           `json:"cliente_telefono"`
	ClienteEmail     string           `json:"cliente_email"`
	ValidezDias      int              `json:"validez_dias"` // días de validez
	Estado           string           `json:"estado"`       // pendiente, aprobada, rechazada, vencida, convertida
	Total            float64          `json:"total"`
	Descuento        float64          `json:"descuento"`
	Notas            string           `json:"notas"`
	UsuarioID        int64            `json:"usuario_id"`
	SedeID           int64            `json:"sede_id"`
	Items            []CotizacionItem `json:"items,omitempty"`
	Usuario          *User            `json:"usuario,omitempty"`
	Sede             *Sede            `json:"sede,omitempty"`
}

// CotizacionItem representa un item dentro de una cotización
type CotizacionItem struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	CotizacionID   int64     `json:"cotizacion_id"`
	ProductoID     int64     `json:"producto_id"`
	Cantidad       int       `json:"cantidad"`
	PrecioUnitario float64   `json:"precio_unitario"`
	Subtotal       float64   `json:"subtotal"`
	Producto       *Product  `json:"producto,omitempty"`
}
