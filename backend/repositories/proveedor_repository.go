package repositories

import (
	"database/sql"
	"smartech/backend/models"
	"time"
)

type DeudaSimple struct {
	ID               int64      `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	NumeroFactura    string     `json:"numero_factura"`
	MontoTotal       float64    `json:"monto_total"`
	MontoPagado      float64    `json:"monto_pagado"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Estado           string     `json:"estado"`
	Notas            string     `json:"notas"`
	ProveedorID      int64      `json:"proveedor_id"`
}

type DeudaView struct {
	ID               int64      `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ProveedorID      int64      `json:"proveedor_id"`
	NumeroFactura    string     `json:"numero_factura"`
	MontoTotal       float64    `json:"monto_total"`
	MontoPagado      float64    `json:"monto_pagado"`
	MontoPendiente   float64    `json:"monto_pendiente"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Estado           string     `json:"estado"`
	Notas            string     `json:"notas"`
	Proveedor        struct {
		Nombre string `json:"nombre"`
	} `json:"proveedor"`
}

type PagoView struct {
	ID               int64   `json:"id"`
	CreatedAt        string  `json:"created_at"`
	Monto            float64 `json:"monto"`
	FechaPago        string  `json:"fecha_pago"`
	MetodoPago       string  `json:"metodo_pago"`
	NumeroReferencia string  `json:"numero_referencia"`
	UsuarioNombre    string  `json:"usuario_nombre"`
}

type ResumenProveedor struct {
	ProveedorID  int64   `json:"proveedor_id"`
	Nombre       string  `json:"nombre"`
	NumFacturas  int     `json:"num_facturas"`
	SaldoTotal   float64 `json:"saldo_total"`
	ProximaFecha *string `json:"proxima_fecha"`
}

type ProveedorRepository struct {
	db *sql.DB
}

func NewProveedorRepository(db *sql.DB) *ProveedorRepository {
	return &ProveedorRepository{db: db}
}

func (r *ProveedorRepository) ListProveedores() ([]models.Proveedor, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, nombre, ruc, direccion, telefono, email, contacto, activo
		FROM proveedores ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proveedores := make([]models.Proveedor, 0)
	for rows.Next() {
		var p models.Proveedor
		var activo int
		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.RucNit, &p.Direccion,
			&p.Telefono, &p.Email, &p.Contacto, &activo); err != nil {
			continue
		}
		p.Activo = activo == 1
		proveedores = append(proveedores, p)
	}
	return proveedores, nil
}

func (r *ProveedorRepository) GetProveedorByID(id int64) (models.Proveedor, error) {
	var p models.Proveedor
	var activo int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, nombre, ruc, direccion, telefono, email, contacto, activo
		FROM proveedores WHERE id = ?
	`, id).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.RucNit, &p.Direccion,
		&p.Telefono, &p.Email, &p.Contacto, &activo)
	if err != nil {
		return models.Proveedor{}, err
	}
	p.Activo = activo == 1
	return p, nil
}

func (r *ProveedorRepository) ListPendingDeudasByProveedor(id int64) ([]DeudaSimple, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, num_factura, monto, monto_pagado, fecha_vence, estado, descripcion
		FROM deudas_proveedores WHERE proveedor_id = ? AND estado != 'pagada'
		ORDER BY fecha_vence
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deudas := make([]DeudaSimple, 0)
	for rows.Next() {
		var d DeudaSimple
		var fechaVenc sql.NullTime
		if err := rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.NumeroFactura, &d.MontoTotal,
			&d.MontoPagado, &fechaVenc, &d.Estado, &d.Notas); err != nil {
			continue
		}
		if fechaVenc.Valid {
			d.FechaVencimiento = &fechaVenc.Time
		}
		d.ProveedorID = id
		deudas = append(deudas, d)
	}
	return deudas, nil
}

func (r *ProveedorRepository) InsertProveedor(p models.Proveedor) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO proveedores (nombre, ruc, direccion, telefono, email, contacto, activo)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, p.Nombre, p.RucNit, p.Direccion, p.Telefono, p.Email, p.Contacto)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ProveedorRepository) UpdateProveedor(id int64, p models.Proveedor) error {
	activo := 0
	if p.Activo {
		activo = 1
	}
	_, err := r.db.Exec(`
		UPDATE proveedores SET nombre = ?, ruc = ?, direccion = ?, telefono = ?,
		                   email = ?, contacto = ?, activo = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, p.Nombre, p.RucNit, p.Direccion, p.Telefono, p.Email, p.Contacto, activo, id)
	return err
}

func (r *ProveedorRepository) CountPendingDeudas(id int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM deudas_proveedores WHERE proveedor_id = ? AND estado != 'pagada'", id).Scan(&count)
	return count, err
}

func (r *ProveedorRepository) DeleteProveedorCascade(id int64) error {
	if _, err := r.db.Exec("DELETE FROM pagos_proveedores WHERE deuda_id IN (SELECT id FROM deudas_proveedores WHERE proveedor_id = ?)", id); err != nil {
		return err
	}
	if _, err := r.db.Exec("DELETE FROM deudas_proveedores WHERE proveedor_id = ?", id); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM proveedores WHERE id = ?", id)
	return err
}

func (r *ProveedorRepository) ListDeudas(estado, proveedorID string) ([]DeudaView, error) {
	query := `
		SELECT d.id, d.created_at, d.updated_at, d.proveedor_id, d.num_factura, d.monto,
		       d.monto_pagado, d.fecha_vence, d.estado, d.descripcion,
		       p.nombre as proveedor_nombre
		FROM deudas_proveedores d
		INNER JOIN proveedores p ON d.proveedor_id = p.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	if estado != "" {
		query += " AND d.estado = ?"
		args = append(args, estado)
	} else {
		query += " AND d.estado != 'pagada'"
	}
	if proveedorID != "" {
		query += " AND d.proveedor_id = ?"
		args = append(args, proveedorID)
	}
	query += " ORDER BY d.fecha_vence"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deudas := make([]DeudaView, 0)
	for rows.Next() {
		var d DeudaView
		var fechaVenc sql.NullTime
		var provNombre string
		if err := rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ProveedorID, &d.NumeroFactura,
			&d.MontoTotal, &d.MontoPagado, &fechaVenc, &d.Estado, &d.Notas, &provNombre); err != nil {
			continue
		}
		if fechaVenc.Valid {
			d.FechaVencimiento = &fechaVenc.Time
		}
		d.MontoPendiente = d.MontoTotal - d.MontoPagado
		d.Proveedor.Nombre = provNombre
		deudas = append(deudas, d)
	}
	return deudas, nil
}

func (r *ProveedorRepository) InsertDeuda(proveedorID int64, numeroFactura string, montoTotal float64, fechaVenc interface{}, notas string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO deudas_proveedores (proveedor_id, num_factura, monto, monto_pagado,
		                                fecha_vence, estado, descripcion)
		VALUES (?, ?, ?, 0, ?, 'pendiente', ?)
	`, proveedorID, numeroFactura, montoTotal, fechaVenc, notas)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ProveedorRepository) GetDeudaMontos(deudaID int64) (float64, float64, error) {
	var montoTotal, montoPagado float64
	err := r.db.QueryRow("SELECT monto, monto_pagado FROM deudas_proveedores WHERE id = ?", deudaID).Scan(&montoTotal, &montoPagado)
	return montoTotal, montoPagado, err
}

func (r *ProveedorRepository) RegisterPago(deudaID int64, monto float64, metodo, referencia string, usuarioID interface{}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var montoTotal, montoPagado float64
	if err := tx.QueryRow("SELECT monto, monto_pagado FROM deudas_proveedores WHERE id = ?", deudaID).Scan(&montoTotal, &montoPagado); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO pagos_proveedores (deuda_id, monto, metodo, referencia, usuario_id)
		VALUES (?, ?, ?, ?, ?)
	`, deudaID, monto, metodo, referencia, usuarioID); err != nil {
		tx.Rollback()
		return err
	}

	nuevoMontoPagado := montoPagado + monto
	nuevoEstado := "pendiente"
	if nuevoMontoPagado >= montoTotal {
		nuevoEstado = "pagada"
	} else if nuevoMontoPagado > 0 {
		nuevoEstado = "parcial"
	}

	if _, err := tx.Exec(`
		UPDATE deudas_proveedores SET monto_pagado = ?, estado = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, nuevoMontoPagado, nuevoEstado, deudaID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ProveedorRepository) ListPagosByDeuda(deudaID int64) ([]PagoView, error) {
	rows, err := r.db.Query(`
		SELECT pp.id, pp.created_at, pp.monto, pp.metodo, pp.referencia, u.name
		FROM pagos_proveedores pp
		INNER JOIN users u ON pp.usuario_id = u.id
		WHERE pp.deuda_id = ?
		ORDER BY pp.created_at DESC
	`, deudaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := make([]PagoView, 0)
	for rows.Next() {
		var p PagoView
		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.Monto, &p.MetodoPago, &p.NumeroReferencia, &p.UsuarioNombre); err != nil {
			continue
		}
		p.FechaPago = p.CreatedAt
		pagos = append(pagos, p)
	}
	return pagos, nil
}

func (r *ProveedorRepository) ResumenTotales() (float64, int, int, int, error) {
	var totalDeuda float64
	var pendientes, vencidas, pagadas int

	if err := r.db.QueryRow(`SELECT COALESCE(SUM(monto - monto_pagado), 0) FROM deudas_proveedores WHERE estado != 'pagada'`).Scan(&totalDeuda); err != nil {
		return 0, 0, 0, 0, err
	}
	r.db.QueryRow(`SELECT COUNT(*) FROM deudas_proveedores WHERE estado = 'pendiente'`).Scan(&pendientes)
	r.db.QueryRow(`SELECT COUNT(*) FROM deudas_proveedores WHERE estado != 'pagada' AND fecha_vence < DATE('now')`).Scan(&vencidas)
	r.db.QueryRow(`SELECT COUNT(*) FROM deudas_proveedores WHERE estado = 'pagada'`).Scan(&pagadas)
	return totalDeuda, pendientes, vencidas, pagadas, nil
}

func (r *ProveedorRepository) ResumenPorProveedor() ([]ResumenProveedor, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.nombre,
		       COUNT(d.id) as num_facturas,
		       COALESCE(SUM(d.monto - d.monto_pagado), 0) as saldo_total,
		       MIN(d.fecha_vence) as proxima_fecha
		FROM proveedores p
		LEFT JOIN deudas_proveedores d ON p.id = d.proveedor_id AND d.estado != 'pagada'
		WHERE p.activo = 1
		GROUP BY p.id, p.nombre
		HAVING saldo_total > 0
		ORDER BY saldo_total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resumen := make([]ResumenProveedor, 0)
	for rows.Next() {
		var rItem ResumenProveedor
		var fecha sql.NullString
		if err := rows.Scan(&rItem.ProveedorID, &rItem.Nombre, &rItem.NumFacturas, &rItem.SaldoTotal, &fecha); err != nil {
			continue
		}
		if fecha.Valid {
			rItem.ProximaFecha = &fecha.String
		}
		resumen = append(resumen, rItem)
	}
	return resumen, nil
}
