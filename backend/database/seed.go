package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase inserta datos iniciales en la base de datos
func SeedDatabase() {
	// Verificar si ya hay datos
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	if count > 0 {
		log.Println("Database already has data, skipping seed.")
		return
	}

	// Crear sedes
	seedSedes()

	// Crear permisos
	seedPermisos()

	// Crear roles
	seedRoles()

	// Asignar permisos a roles
	seedRolePermisos()

	// Crear usuarios de prueba
	seedUsuarios()

	// Crear productos de ejemplo
	seedProductos()

	// Crear proveedores de ejemplo
	seedProveedores()

	// Crear insumos de ejemplo
	seedInsumos()

	// Crear órdenes de trabajo de ejemplo
	seedOrdenesTrabajo()

	// Crear traspasos de ejemplo
	//seedTraspasos()

	// Crear cotizaciones de ejemplo
	seedCotizaciones()

	// Crear deudas de ejemplo
	seedDeudas()

	// Crear promociones de ejemplo
	seedPromociones()

	log.Println("Seed completed successfully.")
}

func seedSedes() {
	sedes := []struct {
		nombre    string
		direccion string
		telefono  string
	}{
		{"Sede Central", "Av. Principal #123, Centro", "591-4-4444444"},
		{"Sucursal Norte", "Av. Norte #456, Zona Norte", "591-4-4444445"},
		{"Sucursal Sur", "Av. Sur #789, Zona Sur", "591-4-4444446"},
	}

	for _, s := range sedes {
		DB.Exec(`INSERT INTO sedes (nombre, direccion, telefono, activa) VALUES (?, ?, ?, 1)`,
			s.nombre, s.direccion, s.telefono)
	}
	log.Println("Sedes created.")
}

func seedPermisos() {
	permisos := []struct {
		name        string
		description string
		resource    string
		action      string
	}{
		// Permisos de Gerente
		{"segmentacion.read", "Ver segmentación de productos/clientes", "segmentacion", "read"},
		{"segmentacion.create", "Crear segmentación", "segmentacion", "create"},
		{"promociones.read", "Ver promociones", "promociones", "read"},
		{"promociones.create", "Crear promociones", "promociones", "create"},
		{"promociones.update", "Editar promociones", "promociones", "update"},
		{"promociones.delete", "Eliminar promociones", "promociones", "delete"},
		{"reportes.read", "Ver reportes de ganancias y ventas", "reportes", "read"},
		{"auditoria.read", "Ver logs de auditoría", "auditoria", "read"},

		// Permisos de Vendedor
		{"compatibilidad.read", "Consultar compatibilidad de productos", "compatibilidad", "read"},
		{"compatibilidad.create", "Crear compatibilidades", "compatibilidad", "create"},
		{"compatibilidad.delete", "Eliminar compatibilidades", "compatibilidad", "delete"},
		{"cotizaciones.read", "Ver cotizaciones", "cotizaciones", "read"},
		{"cotizaciones.create", "Crear cotizaciones", "cotizaciones", "create"},
		{"cotizaciones.update", "Editar/convertir cotizaciones", "cotizaciones", "update"},
		{"cotizaciones.delete", "Eliminar cotizaciones", "cotizaciones", "delete"},
		{"stock.read", "Consultar stock multisede", "stock", "read"},
		{"stock.update", "Actualizar stock", "stock", "update"},

		// Permisos de Técnico
		{"ordenes.read", "Ver órdenes de trabajo", "ordenes", "read"},
		{"ordenes.create", "Crear órdenes de trabajo", "ordenes", "create"},
		{"ordenes.update", "Gestionar órdenes de trabajo", "ordenes", "update"},
		{"ordenes.delete", "Eliminar órdenes de trabajo", "ordenes", "delete"},
		{"insumos.read", "Ver insumos/repuestos", "insumos", "read"},
		{"insumos.create", "Crear insumos", "insumos", "create"},
		{"insumos.update", "Gestionar insumos", "insumos", "update"},
		{"insumos.delete", "Eliminar insumos", "insumos", "delete"},

		// Permisos de Administrador
		{"rmas.read", "Ver RMA/Garantías", "rmas", "read"},
		{"rmas.create", "Crear RMA/Garantías", "rmas", "create"},
		{"rmas.update", "Gestionar RMA/Garantías", "rmas", "update"},
		{"rmas.delete", "Eliminar RMA", "rmas", "delete"},
		{"deudas.read", "Ver deudas de proveedores", "deudas", "read"},
		{"deudas.create", "Crear deudas", "deudas", "create"},
		{"deudas.update", "Registrar pagos a proveedores", "deudas", "update"},
		{"traspasos.read", "Ver traspasos entre sedes", "traspasos", "read"},
		{"traspasos.create", "Crear traspasos", "traspasos", "create"},
		{"traspasos.update", "Gestionar traspasos", "traspasos", "update"},
		{"traspasos.delete", "Eliminar traspasos", "traspasos", "delete"},
		{"proveedores.read", "Ver proveedores", "proveedores", "read"},
		{"proveedores.create", "Crear proveedores", "proveedores", "create"},
		{"proveedores.update", "Gestionar proveedores", "proveedores", "update"},
		{"proveedores.delete", "Eliminar proveedores", "proveedores", "delete"},

		// Permisos generales de productos
		{"products.read", "Ver productos", "products", "read"},
		{"products.create", "Crear productos", "products", "create"},
		{"products.update", "Editar productos", "products", "update"},
		{"products.delete", "Eliminar productos", "products", "delete"},
		{"precios.update", "Actualizar precios de productos", "precios", "update"},

		// Permisos de usuarios y roles (solo admin)
		{"users.read", "Ver usuarios", "users", "read"},
		{"users.create", "Crear usuarios", "users", "create"},
		{"users.update", "Gestionar usuarios", "users", "update"},
		{"users.delete", "Eliminar usuarios", "users", "delete"},
		{"roles.read", "Ver roles", "roles", "read"},
		{"roles.create", "Crear roles", "roles", "create"},
		{"roles.update", "Gestionar roles", "roles", "update"},
		{"roles.delete", "Eliminar roles", "roles", "delete"},

		// Permisos de sedes
		{"sedes.read", "Ver sedes", "sedes", "read"},
		{"sedes.create", "Crear sedes", "sedes", "create"},
		{"sedes.update", "Gestionar sedes", "sedes", "update"},
		{"sedes.delete", "Eliminar sedes", "sedes", "delete"},

		// Dashboard
		{"dashboard.read", "Ver dashboard", "dashboard", "read"},
	}

	for _, p := range permisos {
		DB.Exec(`INSERT INTO permissions (name, description, resource, action) VALUES (?, ?, ?, ?)`,
			p.name, p.description, p.resource, p.action)
	}
	log.Println("Permissions created.")
}

func seedRoles() {
	roles := []struct {
		name        string
		description string
		isSystem    int
	}{
		{"gerente", "Gerente - Acceso a reportes, segmentación, promociones y auditoría", 1},
		{"vendedor", "Vendedor - Acceso a cotizaciones, compatibilidad y stock", 1},
		{"tecnico", "Técnico - Acceso a órdenes de trabajo, insumos y trazabilidad", 1},
		{"administrador", "Administrador - Acceso completo al sistema", 1},
	}

	for _, r := range roles {
		DB.Exec(`INSERT INTO roles (name, description, is_system) VALUES (?, ?, ?)`,
			r.name, r.description, r.isSystem)
	}
	log.Println("Roles created.")
}

func seedRolePermisos() {
	// Mapeo de roles a permisos
	rolePermisos := map[string][]string{
		"gerente": {
			// segmentacion.read/create, promociones.*, reportes.read, auditoria.read, products.read, stock.read
			"segmentacion.read", "segmentacion.create",
			"promociones.read", "promociones.create", "promociones.update", "promociones.delete",
			"reportes.read",
			"auditoria.read",
			"products.read",
			"stock.read",
			"dashboard.read", "sedes.read",
		},
		"vendedor": {
			// compatibilidad.*, cotizaciones.*, stock.read, insumos.read, products.read
			"compatibilidad.read", "compatibilidad.create", "compatibilidad.delete",
			"cotizaciones.read", "cotizaciones.create", "cotizaciones.update", "cotizaciones.delete",
			"stock.read",
			"insumos.read",
			"products.read",
			"dashboard.read", "sedes.read",
		},
		"tecnico": {
			// ordenes.*, insumos.*, products.read, stock.read
			"ordenes.read", "ordenes.create", "ordenes.update", "ordenes.delete",
			"insumos.read", "insumos.create", "insumos.update", "insumos.delete",
			"products.read",
			"stock.read",
			"dashboard.read", "sedes.read",
		},
		"administrador": {
			// rmas.*, deudas.*, traspasos.*, proveedores.*, products.*, users.*, roles.*, sedes.*, stock.*, reportes.read, auditoria.read
			"rmas.read", "rmas.create", "rmas.update", "rmas.delete",
			"deudas.read", "deudas.create", "deudas.update",
			"traspasos.read", "traspasos.create", "traspasos.update", "traspasos.delete",
			"proveedores.read", "proveedores.create", "proveedores.update", "proveedores.delete",
			"products.read", "products.create", "products.update", "products.delete",
			"users.read", "users.create", "users.update", "users.delete",
			"roles.read", "roles.create", "roles.update", "roles.delete",
			"sedes.read", "sedes.create", "sedes.update", "sedes.delete",
			"stock.read", "stock.update",
			"reportes.read",
			"auditoria.read",
			"dashboard.read",
			"precios.update",
		},
	}

	for roleName, permisos := range rolePermisos {
		var roleID int64
		DB.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)

		for _, permName := range permisos {
			var permID int64
			DB.QueryRow("SELECT id FROM permissions WHERE name = ?", permName).Scan(&permID)
			if permID > 0 && roleID > 0 {
				DB.Exec("INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permID)
			}
		}
	}
	log.Println("Role permissions assigned.")
}

func seedUsuarios() {
	// Hash de contraseñas
	hashPassword := func(password string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hash)
	}

	usuarios := []struct {
		name     string
		email    string
		password string
		role     string
		sedeID   int
	}{
		{"Admin Sistema", "admin@inventario.com", "admin123", "administrador", 1},
		{"Juan Gerente", "gerente@inventario.com", "gerente123", "gerente", 1},
		{"María Vendedora", "vendedor@inventario.com", "vendedor123", "vendedor", 1},
		{"Carlos Técnico", "tecnico@inventario.com", "tecnico123", "tecnico", 1},
		{"Ana Vendedora Norte", "vendedor.norte@inventario.com", "vendedor123", "vendedor", 2},
		{"Pedro Técnico Sur", "tecnico.sur@inventario.com", "tecnico123", "tecnico", 3},
	}

	for _, u := range usuarios {
		var roleID int64
		DB.QueryRow("SELECT id FROM roles WHERE name = ?", u.role).Scan(&roleID)

		DB.Exec(`INSERT INTO users (name, email, password, role_id, sede_id) VALUES (?, ?, ?, ?, ?)`,
			u.name, u.email, hashPassword(u.password), roleID, u.sedeID)
	}
	log.Println("Users created.")
}

func seedProductos() {
	productos := []struct {
		codigo       string
		name         string
		description  string
		precioCompra float64
		precioVenta  float64
		category     string
		brand        string
		imageURL     string
	}{
		{"LAPTOP-001", "Laptop HP Pavilion 15", "Laptop con Intel Core i5, 8GB RAM, 512GB SSD", 450.00, 599.99, "Laptops", "HP", "https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500"},
		{"LAPTOP-002", "MacBook Air M2", "Apple MacBook Air con chip M2, 8GB RAM, 256GB SSD", 850.00, 1199.99, "Laptops", "Apple", "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500"},
		{"LAPTOP-003", "Lenovo ThinkPad X1", "ThinkPad X1 Carbon, Intel Core i7, 16GB RAM", 780.00, 1099.99, "Laptops", "Lenovo", "https://images.unsplash.com/photo-1588872657578-7efd1f1555ed?w=500"},
		{"DESKTOP-001", "PC Gamer RTX 4060", "PC Gaming con RTX 4060, Ryzen 5 5600X, 16GB RAM", 650.00, 899.99, "Desktops", "Custom", "https://images.unsplash.com/photo-1587202372775-e229f172b9d7?w=500"},
		{"MONITOR-001", "Monitor LG 27'' 4K", "Monitor UltraFine 4K IPS, 60Hz, USB-C", 280.00, 399.99, "Monitores", "LG", "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=500"},
		{"MONITOR-002", "Monitor Samsung Curvo 32''", "Monitor curvo gaming 165Hz, 1ms, QHD", 220.00, 349.99, "Monitores", "Samsung", "https://images.unsplash.com/photo-1585792180666-f7347c490ee2?w=500"},
		{"TECLADO-001", "Teclado Mecánico Logitech G915", "Teclado mecánico inalámbrico RGB", 120.00, 199.99, "Periféricos", "Logitech", "https://images.unsplash.com/photo-1511467687858-23d96c32e4ae?w=500"},
		{"MOUSE-001", "Mouse Logitech MX Master 3", "Mouse ergonómico inalámbrico premium", 60.00, 99.99, "Periféricos", "Logitech", "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=500"},
		{"RAM-001", "Memoria RAM DDR4 16GB", "Kit 2x8GB DDR4 3200MHz Corsair Vengeance", 35.00, 59.99, "Componentes", "Corsair", "https://images.unsplash.com/photo-1562976540-1502c2145186?w=500"},
		{"SSD-001", "SSD Samsung 970 EVO 1TB", "SSD NVMe M.2 alta velocidad", 70.00, 119.99, "Almacenamiento", "Samsung", "https://images.unsplash.com/photo-1597872200969-2b65d56bd16b?w=500"},
		{"GPU-001", "NVIDIA RTX 4070", "Tarjeta gráfica RTX 4070 12GB GDDR6X", 450.00, 599.99, "Componentes", "NVIDIA", "https://images.unsplash.com/photo-1591488320449-011701bb6704?w=500"},
		{"IMPRESORA-001", "Impresora HP LaserJet Pro", "Impresora láser monocromática WiFi", 180.00, 279.99, "Impresoras", "HP", "https://images.unsplash.com/photo-1612815154858-60aa4c59eaa6?w=500"},
	}

	for _, p := range productos {
		var prodID int64
		result, err := DB.Exec(`INSERT INTO products (codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, '[]', 1)`,
			p.codigo, p.name, p.description, p.precioCompra, p.precioVenta, p.category, p.brand, p.imageURL)

		if err == nil {
			prodID, _ = result.LastInsertId()
			// Agregar stock en todas las sedes
			for sedeID := 1; sedeID <= 3; sedeID++ {
				stock := 10 + (sedeID * 5) // Stock variable por sede
				DB.Exec(`INSERT INTO stock_sedes (sede_id, producto_id, cantidad, stock_minimo, stock_maximo) VALUES (?, ?, ?, 5, 50)`,
					sedeID, prodID, stock)
			}
		}
	}

	// Agregar algunas compatibilidades
	compatibilidades := []struct {
		prod1 string
		prod2 string
		tipo  string
		notas string
	}{
		{"RAM-001", "LAPTOP-001", "upgrade", "Compatible para upgrade de memoria"},
		{"RAM-001", "DESKTOP-001", "upgrade", "Compatible para upgrade de memoria"},
		{"SSD-001", "LAPTOP-001", "upgrade", "Compatible para almacenamiento adicional"},
		{"SSD-001", "LAPTOP-002", "upgrade", "Compatible vía adaptador USB-C"},
		{"MONITOR-001", "LAPTOP-002", "accesorio", "Conexión USB-C directa"},
		{"TECLADO-001", "LAPTOP-001", "accesorio", "Compatible Bluetooth/USB"},
		{"MOUSE-001", "LAPTOP-001", "accesorio", "Compatible Bluetooth/USB"},
	}

	for _, c := range compatibilidades {
		var prod1ID, prod2ID int64
		DB.QueryRow("SELECT id FROM products WHERE codigo = ?", c.prod1).Scan(&prod1ID)
		DB.QueryRow("SELECT id FROM products WHERE codigo = ?", c.prod2).Scan(&prod2ID)
		if prod1ID > 0 && prod2ID > 0 {
			DB.Exec(`INSERT INTO compatibilidades (producto_id, compatible_con_id, tipo_relacion, notas) VALUES (?, ?, ?, ?)`,
				prod1ID, prod2ID, c.tipo, c.notas)
		}
	}

	log.Println("Products and stock created.")
}

func seedProveedores() {
	proveedores := []struct {
		nombre    string
		ruc       string
		direccion string
		telefono  string
		email     string
		contacto  string
	}{
		{"Tech Distribuidores S.A.", "1234567890", "Av. Industrial #100", "591-4-4441111", "ventas@techdist.com", "Roberto García"},
		{"Componentes Plus", "0987654321", "Calle Comercio #50", "591-4-4442222", "info@compplus.com", "Laura Méndez"},
		{"Import Tech Bolivia", "1122334455", "Zona Franca #25", "591-4-4443333", "compras@importech.bo", "Carlos Quispe"},
	}

	for _, p := range proveedores {
		DB.Exec(`INSERT INTO proveedores (nombre, ruc, direccion, telefono, email, contacto, activo) VALUES (?, ?, ?, ?, ?, ?, 1)`,
			p.nombre, p.ruc, p.direccion, p.telefono, p.email, p.contacto)
	}
	log.Println("Proveedores created.")
}

func seedInsumos() {
	insumos := []struct {
		codigo       string
		nombre       string
		descripcion  string
		categoria    string
		unidadMedida string
		stock        int
		stockMinimo  int
		costo        float64
		sedeID       int
	}{
		{"INS-001", "Pasta térmica Arctic MX-4", "Pasta térmica de alta calidad 4g", "Consumibles", "unidad", 25, 5, 8.50, 1},
		{"INS-002", "Cable SATA III", "Cable SATA 6Gb/s 50cm", "Cables y Conectores", "unidad", 50, 10, 2.50, 1},
		{"INS-003", "Destornillador Phillips #2", "Destornillador magnético precisión", "Herramientas", "unidad", 10, 2, 5.00, 1},
		{"INS-004", "Memoria RAM DDR4 4GB", "Repuesto DDR4 2666MHz", "Repuestos", "unidad", 8, 3, 18.00, 1},
		{"INS-005", "Ventilador 120mm", "Ventilador refrigeración 120mm PWM", "Componentes Electrónicos", "unidad", 15, 5, 12.00, 2},
		{"INS-006", "Teclado USB básico", "Teclado español USB", "Repuestos", "unidad", 20, 5, 8.00, 2},
	}

	for _, i := range insumos {
		DB.Exec(`INSERT INTO insumos (codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
			i.codigo, i.nombre, i.descripcion, i.categoria, i.unidadMedida, i.stock, i.stockMinimo, i.costo, i.sedeID)
	}
	log.Println("Insumos created.")
}

func seedOrdenesTrabajo() {
	// Obtener ID del técnico
	var tecnicoID int64
	DB.QueryRow("SELECT id FROM users WHERE email = 'tecnico@inventario.com'").Scan(&tecnicoID)

	ordenes := []struct {
		numeroOrden      string
		clienteNombre    string
		clienteTelefono  string
		equipo           string
		marca            string
		modelo           string
		numSerie         string
		problemaReportado string
		estado           string
		prioridad        string
		sedeID           int
		tecnicoID        int64
	}{
		{"OT-2024-0001", "Juan Pérez", "591-71234567", "Laptop", "HP", "Pavilion 15", "SN123456", "No enciende, se quedó sin batería", "en_reparacion", "alta", 1, tecnicoID},
		{"OT-2024-0002", "María García", "591-72345678", "Desktop", "Dell", "Optiplex 3080", "SN789012", "Lentitud general, posible virus", "en_diagnostico", "media", 1, tecnicoID},
		{"OT-2024-0003", "Carlos López", "591-73456789", "Laptop", "Lenovo", "ThinkPad T480", "SN345678", "Pantalla parpadea", "recibido", "baja", 1, 0},
	}

	for _, o := range ordenes {
		if o.tecnicoID > 0 {
			DB.Exec(`INSERT INTO ordenes_trabajo (numero_orden, cliente_nombre, cliente_telefono, equipo, marca, 
				modelo, num_serie, problema_reportado, estado, prioridad, sede_id, tecnico_id, fecha_ingreso) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
				o.numeroOrden, o.clienteNombre, o.clienteTelefono, o.equipo, o.marca, o.modelo,
				o.numSerie, o.problemaReportado, o.estado, o.prioridad, o.sedeID, o.tecnicoID)
		} else {
			DB.Exec(`INSERT INTO ordenes_trabajo (numero_orden, cliente_nombre, cliente_telefono, equipo, marca, 
				modelo, num_serie, problema_reportado, estado, prioridad, sede_id, fecha_ingreso) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
				o.numeroOrden, o.clienteNombre, o.clienteTelefono, o.equipo, o.marca, o.modelo,
				o.numSerie, o.problemaReportado, o.estado, o.prioridad, o.sedeID)
		}
	}
	log.Println("Ordenes de trabajo created.")
}

func seedTraspasos() {
	// Obtener ID del admin
	var adminID int64
	DB.QueryRow("SELECT id FROM users WHERE email = 'admin@inventario.com'").Scan(&adminID)

	// Obtener un producto
	var productoID int64
	DB.QueryRow("SELECT id FROM products LIMIT 1").Scan(&productoID)

	if adminID > 0 && productoID > 0 {
		// Crear traspaso
		result, err := DB.Exec(`INSERT INTO traspasos (numero_traspaso, sede_origen_id, sede_destino_id, estado, notas, usuario_envia_id, fecha_envio) 
			VALUES ('TRP-2024-0001', 1, 2, 'enviado', 'Traspaso de ejemplo para demostración', ?, datetime('now'))`, adminID)
		if err == nil {
			traspasoID, _ := result.LastInsertId()
			// Agregar item al traspaso
			DB.Exec(`INSERT INTO traspaso_items (traspaso_id, producto_id, cantidad) VALUES (?, ?, 5)`, traspasoID, productoID)
		}
	}
	log.Println("Traspasos created.")
}

func seedCotizaciones() {
	// Obtener ID del vendedor
	var vendedorID int64
	DB.QueryRow("SELECT id FROM users WHERE email = 'vendedor@inventario.com'").Scan(&vendedorID)

	// Obtener productos
	var producto1ID, producto2ID int64
	rows, _ := DB.Query("SELECT id FROM products LIMIT 2")
	if rows.Next() {
		rows.Scan(&producto1ID)
	}
	if rows.Next() {
		rows.Scan(&producto2ID)
	}
	rows.Close()

	if vendedorID > 0 && producto1ID > 0 {
		// Crear cotización
		result, err := DB.Exec(`INSERT INTO cotizaciones (numero_cotizacion, cliente_nombre, cliente_email, cliente_telefono, 
			estado, subtotal, descuento, total, notas, validez, usuario_id, sede_id) 
			VALUES ('COT-2024-0001', 'Empresa ABC S.R.L.', 'contacto@empresaabc.com', '591-74567890', 
			'pendiente', 1000.00, 50.00, 950.00, 'Cotización de equipos de oficina', 15, ?, 1)`, vendedorID)
		if err == nil {
			cotizacionID, _ := result.LastInsertId()
			// Agregar items
			DB.Exec(`INSERT INTO cotizacion_items (cotizacion_id, producto_id, cantidad, precio_unit, descuento, subtotal) 
				VALUES (?, ?, 1, 599.99, 0, 599.99)`, cotizacionID, producto1ID)
			if producto2ID > 0 {
				DB.Exec(`INSERT INTO cotizacion_items (cotizacion_id, producto_id, cantidad, precio_unit, descuento, subtotal) 
					VALUES (?, ?, 2, 199.99, 0, 399.98)`, cotizacionID, producto2ID)
			}
		}
	}
	log.Println("Cotizaciones created.")
}

func seedDeudas() {
	// Obtener proveedor
	var proveedorID int64
	DB.QueryRow("SELECT id FROM proveedores LIMIT 1").Scan(&proveedorID)

	if proveedorID > 0 {
		DB.Exec(`INSERT INTO deudas_proveedores (proveedor_id, num_factura, monto, monto_pagado, 
			fecha_vence, estado, descripcion) 
			VALUES (?, 'FAC-2024-001', 5000.00, 2000.00, date('now', '+30 days'), 'pendiente', 
			'Compra de equipos para stock')`, proveedorID)
		DB.Exec(`INSERT INTO deudas_proveedores (proveedor_id, num_factura, monto, monto_pagado, 
			fecha_vence, estado, descripcion) 
			VALUES (?, 'FAC-2024-002', 1500.00, 1500.00, date('now', '-30 days'), 'pagada', 
			'Compra de accesorios')`, proveedorID)
	}
	log.Println("Deudas created.")
}

func seedPromociones() {
	// Obtener un usuario gerente para las promociones
	var usuarioID int64
	err := DB.QueryRow("SELECT u.id FROM users u INNER JOIN roles r ON u.role_id = r.id WHERE r.name = 'gerente' LIMIT 1").Scan(&usuarioID)
	if err != nil {
		// Si no hay gerente, usar el primer usuario
		DB.QueryRow("SELECT id FROM users LIMIT 1").Scan(&usuarioID)
	}
	
	if usuarioID == 0 {
		log.Println("No users found, skipping promociones seed")
		return
	}

	// Obtener un producto para la promoción específica
	var productoID int64
	DB.QueryRow("SELECT id FROM products LIMIT 1").Scan(&productoID)

	// Promoción general
	_, err = DB.Exec(`INSERT INTO promociones (nombre, descripcion, tipo_descuento, valor, fecha_inicio, fecha_fin, activa, usuario_id) 
		VALUES (?, ?, ?, ?, date('now'), date('now', '+30 days'), 1, ?)`,
		"Descuento Bienvenida", "Descuento del 10% en primera compra", "porcentaje", 10.0, usuarioID)
	if err != nil {
		log.Printf("Error creating general promocion: %v", err)
	}

	// Promoción por producto específico (si hay productos)
	if productoID > 0 {
		productIDs := fmt.Sprintf("[%d]", productoID) // JSON array format
		_, err = DB.Exec(`INSERT INTO promociones (nombre, descripcion, tipo_descuento, valor, fecha_inicio, fecha_fin, activa, producto_ids, usuario_id) 
			VALUES (?, ?, ?, ?, date('now'), date('now', '+15 days'), 1, ?, ?)`,
			"Oferta Laptop", "Descuento especial en laptops seleccionadas", "porcentaje", 15.0, productIDs, usuarioID)
		if err != nil {
			log.Printf("Error creating product promocion: %v", err)
		}
	}

	log.Println("Promociones created.")
}
