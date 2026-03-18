package models

import "time"

// LogAuditoria representa un registro de auditoría del sistema
type LogAuditoria struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UsuarioID     int64     `json:"usuario_id"`
	Accion        string    `json:"accion"`  // crear, editar, eliminar, login, logout, etc
	Entidad       string    `json:"entidad"` // producto, rma, cotizacion, traspaso, etc
	EntidadID     int64     `json:"entidad_id"`
	ValorAnterior string    `json:"valor_anterior,omitempty"`
	ValorNuevo    string    `json:"valor_nuevo,omitempty"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Detalles      string    `json:"detalles"`
	Usuario       *User     `json:"usuario,omitempty"`
}

// ReporteGanancias representa datos para el reporte de ganancias
type ReporteGanancias struct {
	Periodo        string  `json:"periodo"`
	TotalVentas    float64 `json:"total_ventas"`
	TotalCostos    float64 `json:"total_costos"`
	GananciaBruta  float64 `json:"ganancia_bruta"`
	GananciaNeta   float64 `json:"ganancia_neta"`
	CantidadVentas int     `json:"cantidad_ventas"`
}

// Segmentacion representa una segmentación de clientes/productos
type Segmentacion struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Criterios   string    `json:"criterios"` // JSON con criterios de segmentación
	Activa      bool      `json:"activa"`
}

// Promocion representa una promoción activa
type Promocion struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Nombre         string    `json:"nombre"`
	Descripcion    string    `json:"descripcion"`
	Tipo           string    `json:"tipo"` // porcentaje, monto_fijo, 2x1
	Valor          float64   `json:"valor"`
	FechaInicio    time.Time `json:"fecha_inicio"`
	FechaFin       time.Time `json:"fecha_fin"`
	ProductosIDs   string    `json:"productos_ids"`   // JSON array de IDs de productos
	SegmentacionID *int64    `json:"segmentacion_id"` // Segmentación a la que aplica
	Activa         bool      `json:"activa"`
}
