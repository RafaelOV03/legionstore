package models

import "time"

// Traspaso representa un movimiento de inventario entre sedes
type Traspaso struct {
	ID              int64          `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	NumeroTraspaso  string         `json:"numero_traspaso"`
	SedeOrigenID    int64          `json:"sede_origen_id"`
	SedeDestinoID   int64          `json:"sede_destino_id"`
	Estado          string         `json:"estado"` // pendiente, enviado, recibido, cancelado
	FechaEnvio      *time.Time     `json:"fecha_envio,omitempty"`
	FechaRecepcion  *time.Time     `json:"fecha_recepcion,omitempty"`
	Notas           string         `json:"notas"`
	UsuarioEnviaID  int64          `json:"usuario_envia_id"`
	UsuarioRecibeID *int64         `json:"usuario_recibe_id,omitempty"`
	Items           []TraspasoItem `json:"items,omitempty"`
	SedeOrigen      *Sede          `json:"sede_origen,omitempty"`
	SedeDestino     *Sede          `json:"sede_destino,omitempty"`
	UsuarioOrigen   *User          `json:"usuario_origen,omitempty"`
	UsuarioDestino  *User          `json:"usuario_destino,omitempty"`
}

// TraspasoItem representa un producto en un traspaso
type TraspasoItem struct {
	ID               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	TraspasoID       int64     `json:"traspaso_id"`
	ProductoID       int64     `json:"producto_id"`
	Cantidad         int       `json:"cantidad"`
	CantidadRecibida int       `json:"cantidad_recibida"`
	Producto         *Product  `json:"producto,omitempty"`
}
