package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase inicializa la conexión con la base de datos SQLite
func InitDatabase() {
	var err error
	log.Println("Connecting to database...")

	// Configurar SQLite para modo WAL para mejor concurrencia
	DB, err = sql.Open("sqlite3", "./inventario.db?cache=shared&mode=rwc&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configurar pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(0)

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connection successful.")

	log.Println("Creating tables...")
	if err := CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	log.Println("Tables created successfully.")

	log.Println("Seeding database...")
	SeedDatabase()
	log.Println("Database seeded successfully.")
}

// CreateTables crea las tablas del sistema de inventario
func CreateTables() error {
	queries := []string{
		// ===================== TABLAS BASE =====================

		// Permisos
		`CREATE TABLE IF NOT EXISTS permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			resource TEXT NOT NULL,
			action TEXT NOT NULL
		)`,

		// Roles
		`CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			is_system INTEGER DEFAULT 0
		)`,

		// Relación roles-permisos
		`CREATE TABLE IF NOT EXISTS role_permissions (
			role_id INTEGER NOT NULL,
			permission_id INTEGER NOT NULL,
			PRIMARY KEY (role_id, permission_id),
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
		)`,

		// Sedes/Sucursales
		`CREATE TABLE IF NOT EXISTS sedes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			nombre TEXT NOT NULL,
			direccion TEXT,
			telefono TEXT,
			activa INTEGER DEFAULT 1
		)`,

		// Usuarios
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role_id INTEGER NOT NULL,
			sede_id INTEGER,
			FOREIGN KEY (role_id) REFERENCES roles(id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id)
		)`,

		// ===================== PRODUCTOS E INVENTARIO =====================

		// Productos
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			codigo TEXT UNIQUE,
			name TEXT NOT NULL,
			description TEXT,
			precio_compra REAL DEFAULT 0,
			precio_venta REAL NOT NULL,
			category TEXT NOT NULL,
			brand TEXT,
			image_url TEXT,
			images TEXT DEFAULT '[]',
			activo INTEGER DEFAULT 1
		)`,

		// Stock por sede
		`CREATE TABLE IF NOT EXISTS stock_sedes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			sede_id INTEGER NOT NULL,
			producto_id INTEGER NOT NULL,
			cantidad INTEGER DEFAULT 0,
			stock_minimo INTEGER DEFAULT 5,
			stock_maximo INTEGER DEFAULT 100,
			UNIQUE(sede_id, producto_id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id),
			FOREIGN KEY (producto_id) REFERENCES products(id)
		)`,

		// Compatibilidad entre productos
		`CREATE TABLE IF NOT EXISTS compatibilidades (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			producto_id INTEGER NOT NULL,
			compatible_con_id INTEGER NOT NULL,
			tipo_relacion TEXT DEFAULT 'compatible',
			notas TEXT,
			FOREIGN KEY (producto_id) REFERENCES products(id),
			FOREIGN KEY (compatible_con_id) REFERENCES products(id)
		)`,

		// ===================== PROVEEDORES Y DEUDAS =====================

		// Proveedores
		`CREATE TABLE IF NOT EXISTS proveedores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			nombre TEXT NOT NULL,
			ruc TEXT UNIQUE,
			direccion TEXT,
			telefono TEXT,
			email TEXT,
			contacto TEXT,
			activo INTEGER DEFAULT 1
		)`,

		// Deudas con proveedores
		`CREATE TABLE IF NOT EXISTS deudas_proveedores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			proveedor_id INTEGER NOT NULL,
			monto REAL NOT NULL,
			monto_pagado REAL DEFAULT 0,
			fecha_vence DATETIME,
			estado TEXT DEFAULT 'pendiente',
			descripcion TEXT,
			num_factura TEXT,
			FOREIGN KEY (proveedor_id) REFERENCES proveedores(id)
		)`,

		// Pagos a proveedores
		`CREATE TABLE IF NOT EXISTS pagos_proveedores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deuda_id INTEGER NOT NULL,
			monto REAL NOT NULL,
			metodo TEXT,
			referencia TEXT,
			usuario_id INTEGER NOT NULL,
			FOREIGN KEY (deuda_id) REFERENCES deudas_proveedores(id),
			FOREIGN KEY (usuario_id) REFERENCES users(id)
		)`,

		// ===================== RMA / GARANTÍAS =====================

		// RMA/Garantías
		`CREATE TABLE IF NOT EXISTS rmas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			numero_rma TEXT UNIQUE NOT NULL,
			producto_id INTEGER NOT NULL,
			cliente_nombre TEXT,
			cliente_telefono TEXT,
			cliente_email TEXT,
			num_serie TEXT,
			fecha_compra DATETIME,
			motivo_devolucion TEXT,
			diagnostico TEXT,
			estado TEXT DEFAULT 'recibido',
			solucion TEXT,
			fecha_resolucion DATETIME,
			usuario_id INTEGER NOT NULL,
			sede_id INTEGER NOT NULL,
			notas TEXT,
			FOREIGN KEY (producto_id) REFERENCES products(id),
			FOREIGN KEY (usuario_id) REFERENCES users(id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id)
		)`,

		// Historial de RMA
		`CREATE TABLE IF NOT EXISTS historial_rmas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			rma_id INTEGER NOT NULL,
			estado_anterior TEXT,
			estado_nuevo TEXT,
			comentario TEXT,
			usuario_id INTEGER NOT NULL,
			FOREIGN KEY (rma_id) REFERENCES rmas(id),
			FOREIGN KEY (usuario_id) REFERENCES users(id)
		)`,

		// ===================== COTIZACIONES =====================

		// Cotizaciones
		`CREATE TABLE IF NOT EXISTS cotizaciones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			numero_cotizacion TEXT UNIQUE NOT NULL,
			cliente_nombre TEXT,
			cliente_telefono TEXT,
			cliente_email TEXT,
			cliente_empresa TEXT,
			subtotal REAL DEFAULT 0,
			descuento REAL DEFAULT 0,
			impuesto REAL DEFAULT 0,
			total REAL DEFAULT 0,
			estado TEXT DEFAULT 'borrador',
			validez INTEGER DEFAULT 30,
			fecha_vence DATETIME,
			notas TEXT,
			usuario_id INTEGER NOT NULL,
			sede_id INTEGER NOT NULL,
			FOREIGN KEY (usuario_id) REFERENCES users(id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id)
		)`,

		// Items de cotización
		`CREATE TABLE IF NOT EXISTS cotizacion_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			cotizacion_id INTEGER NOT NULL,
			producto_id INTEGER NOT NULL,
			cantidad INTEGER NOT NULL,
			precio_unit REAL NOT NULL,
			descuento REAL DEFAULT 0,
			subtotal REAL NOT NULL,
			FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(id) ON DELETE CASCADE,
			FOREIGN KEY (producto_id) REFERENCES products(id)
		)`,

		// ===================== TRASPASOS =====================

		// Traspasos entre sedes
		`CREATE TABLE IF NOT EXISTS traspasos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			numero_traspaso TEXT UNIQUE NOT NULL,
			sede_origen_id INTEGER NOT NULL,
			sede_destino_id INTEGER NOT NULL,
			estado TEXT DEFAULT 'pendiente',
			fecha_envio DATETIME,
			fecha_recepcion DATETIME,
			notas TEXT,
			usuario_envia_id INTEGER NOT NULL,
			usuario_recibe_id INTEGER,
			FOREIGN KEY (sede_origen_id) REFERENCES sedes(id),
			FOREIGN KEY (sede_destino_id) REFERENCES sedes(id),
			FOREIGN KEY (usuario_envia_id) REFERENCES users(id),
			FOREIGN KEY (usuario_recibe_id) REFERENCES users(id)
		)`,

		// Items de traspaso
		`CREATE TABLE IF NOT EXISTS traspaso_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			traspaso_id INTEGER NOT NULL,
			producto_id INTEGER NOT NULL,
			cantidad INTEGER NOT NULL,
			cantidad_recibida INTEGER DEFAULT 0,
			notas TEXT,
			FOREIGN KEY (traspaso_id) REFERENCES traspasos(id) ON DELETE CASCADE,
			FOREIGN KEY (producto_id) REFERENCES products(id)
		)`,

		// ===================== TÉCNICO - ÓRDENES DE TRABAJO =====================

		// Órdenes de trabajo
		`CREATE TABLE IF NOT EXISTS ordenes_trabajo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			numero_orden TEXT UNIQUE NOT NULL,
			cliente_nombre TEXT,
			cliente_telefono TEXT,
			cliente_email TEXT,
			equipo TEXT,
			num_serie TEXT,
			marca TEXT,
			modelo TEXT,
			problema_reportado TEXT,
			diagnostico_tecnico TEXT,
			solucion_aplicada TEXT,
			prioridad TEXT DEFAULT 'media',
			estado TEXT DEFAULT 'recibido',
			fecha_ingreso DATETIME DEFAULT CURRENT_TIMESTAMP,
			fecha_promesa DATETIME,
			fecha_entrega DATETIME,
			costo_mano_obra REAL DEFAULT 0,
			costo_repuestos REAL DEFAULT 0,
			costo_total REAL DEFAULT 0,
			tecnico_id INTEGER NOT NULL,
			sede_id INTEGER NOT NULL,
			notas TEXT,
			FOREIGN KEY (tecnico_id) REFERENCES users(id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id)
		)`,

		// Insumos/Repuestos
		`CREATE TABLE IF NOT EXISTS insumos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			codigo TEXT UNIQUE,
			nombre TEXT NOT NULL,
			descripcion TEXT,
			categoria TEXT,
			unidad_medida TEXT DEFAULT 'unidad',
			stock INTEGER DEFAULT 0,
			stock_minimo INTEGER DEFAULT 5,
			costo REAL DEFAULT 0,
			sede_id INTEGER NOT NULL,
			activo INTEGER DEFAULT 1,
			FOREIGN KEY (sede_id) REFERENCES sedes(id)
		)`,

		// Insumos usados en órdenes
		`CREATE TABLE IF NOT EXISTS insumos_orden (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			orden_id INTEGER NOT NULL,
			insumo_id INTEGER NOT NULL,
			cantidad INTEGER NOT NULL,
			FOREIGN KEY (orden_id) REFERENCES ordenes_trabajo(id) ON DELETE CASCADE,
			FOREIGN KEY (insumo_id) REFERENCES insumos(id)
		)`,

		// Trazabilidad de órdenes de trabajo
		`CREATE TABLE IF NOT EXISTS trazabilidad (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			orden_trabajo_id INTEGER,
			accion TEXT NOT NULL,
			detalle TEXT,
			usuario_id INTEGER NOT NULL,
			FOREIGN KEY (orden_trabajo_id) REFERENCES ordenes_trabajo(id),
			FOREIGN KEY (usuario_id) REFERENCES users(id)
		)`,

		// ===================== GERENCIA - PROMOCIONES Y REPORTES =====================

		// Promociones
		`CREATE TABLE IF NOT EXISTS promociones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			nombre TEXT NOT NULL,
			descripcion TEXT,
			tipo_descuento TEXT DEFAULT 'porcentaje',
			valor REAL NOT NULL,
			fecha_inicio DATETIME,
			fecha_fin DATETIME,
			activa INTEGER DEFAULT 1,
			categorias TEXT,
			producto_ids TEXT,
			usuario_id INTEGER NOT NULL,
			FOREIGN KEY (usuario_id) REFERENCES users(id)
		)`,

		// Logs de auditoría
		`CREATE TABLE IF NOT EXISTS logs_auditoria (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			usuario_id INTEGER NOT NULL,
			accion TEXT NOT NULL,
			entidad TEXT NOT NULL,
			entidad_id INTEGER,
			valor_anterior TEXT,
			valor_nuevo TEXT,
			ip_address TEXT,
			user_agent TEXT,
			detalles TEXT,
			FOREIGN KEY (usuario_id) REFERENCES users(id)
		)`,

		// Ventas (para reportes de ganancias)
		`CREATE TABLE IF NOT EXISTS ventas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			numero_venta TEXT UNIQUE NOT NULL,
			cliente_nombre TEXT,
			subtotal REAL DEFAULT 0,
			descuento REAL DEFAULT 0,
			impuesto REAL DEFAULT 0,
			total REAL NOT NULL,
			costo_total REAL DEFAULT 0,
			ganancia REAL DEFAULT 0,
			metodo_pago TEXT,
			usuario_id INTEGER NOT NULL,
			sede_id INTEGER NOT NULL,
			cotizacion_id INTEGER,
			FOREIGN KEY (usuario_id) REFERENCES users(id),
			FOREIGN KEY (sede_id) REFERENCES sedes(id),
			FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(id)
		)`,

		// Items de venta
		`CREATE TABLE IF NOT EXISTS venta_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			venta_id INTEGER NOT NULL,
			producto_id INTEGER NOT NULL,
			cantidad INTEGER NOT NULL,
			precio_unit REAL NOT NULL,
			costo_unit REAL DEFAULT 0,
			descuento REAL DEFAULT 0,
			subtotal REAL NOT NULL,
			FOREIGN KEY (venta_id) REFERENCES ventas(id) ON DELETE CASCADE,
			FOREIGN KEY (producto_id) REFERENCES products(id)
		)`,

		// ===================== ÍNDICES =====================
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id)`,
		`CREATE INDEX IF NOT EXISTS idx_users_sede_id ON users(sede_id)`,
		`CREATE INDEX IF NOT EXISTS idx_products_category ON products(category)`,
		`CREATE INDEX IF NOT EXISTS idx_products_codigo ON products(codigo)`,
		`CREATE INDEX IF NOT EXISTS idx_stock_sedes_producto ON stock_sedes(producto_id)`,
		`CREATE INDEX IF NOT EXISTS idx_stock_sedes_sede ON stock_sedes(sede_id)`,
		`CREATE INDEX IF NOT EXISTS idx_rmas_estado ON rmas(estado)`,
		`CREATE INDEX IF NOT EXISTS idx_cotizaciones_estado ON cotizaciones(estado)`,
		`CREATE INDEX IF NOT EXISTS idx_traspasos_estado ON traspasos(estado)`,
		`CREATE INDEX IF NOT EXISTS idx_ordenes_trabajo_estado ON ordenes_trabajo(estado)`,
		`CREATE INDEX IF NOT EXISTS idx_deudas_estado ON deudas_proveedores(estado)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_usuario ON logs_auditoria(usuario_id)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_entidad ON logs_auditoria(entidad)`,
		`CREATE INDEX IF NOT EXISTS idx_trazabilidad_orden ON trazabilidad(orden_trabajo_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ventas_fecha ON ventas(created_at)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
// cambios agregados 23/06/2024
// PingDatabase verifica la conectividad con la base de datos
func PingDatabase() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.Ping()
}

// GetDBStats retorna estadísticas del pool de conexiones
func GetDBStats() sql.DBStats {
	if DB == nil {
		return sql.DBStats{}
	}
	return DB.Stats()
}

// IsTableExists verifica si una tabla existe en la base de datos
func IsTableExists(tableName string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("database not initialized")
	}
	
	var count int
	query := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?"
	err := DB.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
// CloseDB cierra la conexión a la base de datos
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
