package models

import "time"

// Proveedor representa un proveedor de productos
type Proveedor struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Nombre    string    `json:"nombre"`
	RucNit    string    `json:"ruc_nit"`
	Direccion string    `json:"direccion"`
	Telefono  string    `json:"telefono"`
	Email     string    `json:"email"`
	Contacto  string    `json:"contacto"`
	Activo    bool      `json:"activo"`
}

// DeudaProveedor representa una deuda pendiente con un proveedor
type DeudaProveedor struct {
	ID               int64      `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ProveedorID      int64      `json:"proveedor_id"`
	NumeroFactura    string     `json:"numero_factura"`
	MontoTotal       float64    `json:"monto_total"`
	MontoPagado      float64    `json:"monto_pagado"`
	FechaFactura     time.Time  `json:"fecha_factura"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Estado           string     `json:"estado"` // pendiente, parcial, pagada
	Notas            string     `json:"notas"`
	Proveedor        *Proveedor `json:"proveedor,omitempty"`
}

// PagoProveedor representa un pago realizado a un proveedor
type PagoProveedor struct {
	ID               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	DeudaID          int64     `json:"deuda_id"`
	Monto            float64   `json:"monto"`
	FechaPago        time.Time `json:"fecha_pago"`
	MetodoPago       string    `json:"metodo_pago"` // efectivo, transferencia, cheque
	NumeroReferencia string    `json:"numero_referencia"`
	Notas            string    `json:"notas"`
	UsuarioID        int64     `json:"usuario_id"`
}
