package models

import "time"

// OrdenTrabajo representa una orden de trabajo técnico
type OrdenTrabajo struct {
	ID                int64         `json:"id"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	NumeroOrden       string        `json:"numero_orden"`
	ClienteNombre     string        `json:"cliente_nombre"`
	ClienteTelefono   string        `json:"cliente_telefono"`
	Equipo            string        `json:"equipo"`
	Marca             string        `json:"marca"`
	Modelo            string        `json:"modelo"`
	NumSerie          string        `json:"num_serie"`
	ProblemaReportado string        `json:"problema_reportado"`
	DiagnosticoTecnico string       `json:"diagnostico_tecnico"`
	SolucionAplicada  string        `json:"solucion_aplicada"`
	Estado            string        `json:"estado"`    // recibido, en_diagnostico, en_reparacion, terminado, entregado, cancelado
	Prioridad         string        `json:"prioridad"` // baja, media, alta, urgente
	FechaIngreso      time.Time     `json:"fecha_ingreso"`
	FechaPromesa      *time.Time    `json:"fecha_promesa,omitempty"`
	FechaEntrega      *time.Time    `json:"fecha_entrega,omitempty"`
	CostoManoObra     float64       `json:"costo_mano_obra"`
	CostoRepuestos    float64       `json:"costo_repuestos"`
	TecnicoID         *int64        `json:"tecnico_id,omitempty"`
	SedeID            int64         `json:"sede_id"`
	Notas             string        `json:"notas"`
	Insumos           []InsumoOrden `json:"insumos,omitempty"`
	Tecnico           *User         `json:"tecnico,omitempty"`
	Sede              *Sede         `json:"sede,omitempty"`
}

// Insumo representa un repuesto o material del inventario técnico
type Insumo struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Codigo       string    `json:"codigo"`
	Nombre       string    `json:"nombre"`
	Descripcion  string    `json:"descripcion"`
	Categoria    string    `json:"categoria"`
	UnidadMedida string    `json:"unidad_medida"`
	Stock        int       `json:"stock"`
	StockMinimo  int       `json:"stock_minimo"`
	Costo        float64   `json:"costo"`
	SedeID       int64     `json:"sede_id"`
	Activo       bool      `json:"activo"`
}

// InsumoOrden representa un insumo usado en una orden de trabajo
type InsumoOrden struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	OrdenTrabajoID int64     `json:"orden_trabajo_id"`
	InsumoID       int64     `json:"insumo_id"`
	Cantidad       int       `json:"cantidad"`
	Insumo         *Insumo   `json:"insumo,omitempty"`
}

// Trazabilidad representa el historial de eventos de una orden de trabajo
type Trazabilidad struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	OrdenTrabajoID int64     `json:"orden_trabajo_id"`
	Accion         string    `json:"accion"`
	Detalle        string    `json:"detalle"`
	UsuarioID      int64     `json:"usuario_id"`
	Usuario        *User     `json:"usuario,omitempty"`
}

// Compatibilidad representa la compatibilidad entre productos
type Compatibilidad struct {
	ID            int64    `json:"id"`
	ProductoID    int64    `json:"producto_id"`
	CompatibleCon int64    `json:"compatible_con"`
	Notas         string   `json:"notas"`
	Producto      *Product `json:"producto,omitempty"`
	ProductoComp  *Product `json:"producto_compatible,omitempty"`
}
