package models

import "time"

// RMA representa una solicitud de devolución/garantía
type RMA struct {
	ID               int64      `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	NumeroRMA        string     `json:"numero_rma"`
	ProductoID       int64      `json:"producto_id"`
	ClienteNombre    string     `json:"cliente_nombre"`
	ClienteTelefono  string     `json:"cliente_telefono"`
	ClienteEmail     string     `json:"cliente_email"`
	NumSerie         string     `json:"num_serie"`
	FechaCompra      time.Time  `json:"fecha_compra"`
	MotivoDevolucion string     `json:"motivo_devolucion"`
	Diagnostico      string     `json:"diagnostico"`
	Estado           string     `json:"estado"`   // recibido, en_revision, aprobado, rechazado, enviado_proveedor, resuelto
	Solucion         string     `json:"solucion"` // reparacion, reemplazo, devolucion_dinero
	FechaResolucion  *time.Time `json:"fecha_resolucion,omitempty"`
	UsuarioID        int64      `json:"usuario_id"`
	SedeID           int64      `json:"sede_id"`
	Notas            string     `json:"notas"`
	Producto         *Product   `json:"producto,omitempty"`
	Usuario          *User      `json:"usuario,omitempty"`
	Sede             *Sede      `json:"sede,omitempty"`
}

// HistorialRMA representa el historial de cambios de un RMA
type HistorialRMA struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	RMAID       int64     `json:"rma_id"`
	EstadoAnt   string    `json:"estado_anterior"`
	EstadoNuevo string    `json:"estado_nuevo"`
	Comentario  string    `json:"comentario"`
	UsuarioID   int64     `json:"usuario_id"`
}
