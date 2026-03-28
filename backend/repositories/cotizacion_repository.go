package repositories

import (
	"database/sql"
	"fmt"
	"smartech/backend/models"
	"time"
)

type CotizacionView struct {
	models.Cotizacion
	UsuarioNombre string `json:"usuario_nombre"`
	SedeNombre    string `json:"sede_nombre"`
}

type CotizacionItemView struct {
	ID             int64   `json:"id"`
	ProductoID     int64   `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
	ProductoNombre string  `json:"producto_nombre"`
	ProductoMarca  string  `json:"producto_marca"`
	ProductoCodigo string  `json:"producto_codigo"`
}

type CotizacionPDF struct {
	NumeroCotizacion string    `json:"numero_cotizacion"`
	ClienteNombre    string    `json:"cliente_nombre"`
	ClienteTelefono  string    `json:"cliente_telefono"`
	ClienteEmail     string    `json:"cliente_email"`
	ValidezDias      int       `json:"validez_dias"`
	Total            float64   `json:"total"`
	Descuento        float64   `json:"descuento"`
	Notas            string    `json:"notas"`
	FechaCreacion    time.Time `json:"fecha_creacion"`
	UsuarioNombre    string    `json:"usuario_nombre"`
	SedeNombre       string    `json:"sede_nombre"`
	SedeDireccion    string    `json:"sede_direccion"`
}

type CotizacionPDFItem struct {
	Codigo         string  `json:"codigo"`
	Nombre         string  `json:"nombre"`
	Marca          string  `json:"marca"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type CotizacionCreateItem struct {
	ProductoID     int64
	Cantidad       int
	PrecioUnitario float64
}

type VentaItem struct {
	ProductoID     int64   `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type CotizacionRepository struct {
	db *sql.DB
}

func NewCotizacionRepository(db *sql.DB) *CotizacionRepository {
	return &CotizacionRepository{db: db}
}

func (r *CotizacionRepository) ListCotizaciones(estado, sedeID string) ([]CotizacionView, error) {
	query := `
		SELECT c.id, c.created_at, c.updated_at, c.numero_cotizacion, c.cliente_nombre,
		       c.cliente_telefono, c.cliente_email, c.validez, c.estado, c.total,
		       c.descuento, c.notas, c.usuario_id, c.sede_id,
		       u.name as usuario_nombre, s.nombre as sede_nombre
		FROM cotizaciones c
		INNER JOIN users u ON c.usuario_id = u.id
		INNER JOIN sedes s ON c.sede_id = s.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	if estado != "" {
		query += " AND c.estado = ?"
		args = append(args, estado)
	}
	if sedeID != "" {
		query += " AND c.sede_id = ?"
		args = append(args, sedeID)
	}
	query += " ORDER BY c.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CotizacionView, 0)
	for rows.Next() {
		var c CotizacionView
		if err := rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.NumeroCotizacion, &c.ClienteNombre,
			&c.ClienteTelefono, &c.ClienteEmail, &c.ValidezDias, &c.Estado, &c.Total,
			&c.Descuento, &c.Notas, &c.UsuarioID, &c.SedeID, &c.UsuarioNombre, &c.SedeNombre); err != nil {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *CotizacionRepository) GetCotizacionByID(id int64) (models.Cotizacion, error) {
	var c models.Cotizacion
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, numero_cotizacion, cliente_nombre,
		       cliente_telefono, cliente_email, validez, estado, total,
		       descuento, notas, usuario_id, sede_id
		FROM cotizaciones WHERE id = ?
	`, id).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.NumeroCotizacion, &c.ClienteNombre,
		&c.ClienteTelefono, &c.ClienteEmail, &c.ValidezDias, &c.Estado, &c.Total,
		&c.Descuento, &c.Notas, &c.UsuarioID, &c.SedeID)
	if err != nil {
		return models.Cotizacion{}, err
	}
	return c, nil
}

func (r *CotizacionRepository) ListCotizacionItems(id int64) ([]CotizacionItemView, error) {
	rows, err := r.db.Query(`
		SELECT ci.id, ci.producto_id, ci.cantidad, ci.precio_unit, ci.subtotal,
		       p.name, p.brand, p.codigo
		FROM cotizacion_items ci
		INNER JOIN products p ON ci.producto_id = p.id
		WHERE ci.cotizacion_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CotizacionItemView, 0)
	for rows.Next() {
		var i CotizacionItemView
		if err := rows.Scan(&i.ID, &i.ProductoID, &i.Cantidad, &i.PrecioUnitario, &i.Subtotal,
			&i.ProductoNombre, &i.ProductoMarca, &i.ProductoCodigo); err != nil {
			continue
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *CotizacionRepository) CountCotizaciones() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM cotizaciones").Scan(&count)
	return count, err
}

func (r *CotizacionRepository) BuildNumeroCotizacion(count int) string {
	return fmt.Sprintf("COT-%d-%04d", time.Now().Year(), count+1)
}

func (r *CotizacionRepository) CreateCotizacionWithItems(numero, clienteNombre, clienteTelefono, clienteEmail string, validez int, total, descuento float64, notas string, usuarioID interface{}, sedeID int64, items []CotizacionCreateItem) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO cotizaciones (numero_cotizacion, cliente_nombre, cliente_telefono, cliente_email,
		                          validez, estado, total, descuento, notas, usuario_id, sede_id)
		VALUES (?, ?, ?, ?, ?, 'pendiente', ?, ?, ?, ?, ?)
	`, numero, clienteNombre, clienteTelefono, clienteEmail, validez, total, descuento, notas, usuarioID, sedeID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	cotID, _ := result.LastInsertId()
	for _, item := range items {
		subtotal := float64(item.Cantidad) * item.PrecioUnitario
		if _, err := tx.Exec(`
			INSERT INTO cotizacion_items (cotizacion_id, producto_id, cantidad, precio_unit, subtotal)
			VALUES (?, ?, ?, ?, ?)
		`, cotID, item.ProductoID, item.Cantidad, item.PrecioUnitario, subtotal); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return cotID, nil
}

func (r *CotizacionRepository) GetEstadoByID(id int64) (string, error) {
	var estado string
	err := r.db.QueryRow("SELECT estado FROM cotizaciones WHERE id = ?", id).Scan(&estado)
	return estado, err
}

func (r *CotizacionRepository) UpdateEstado(id int64, estado string) error {
	_, err := r.db.Exec("UPDATE cotizaciones SET estado = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", estado, id)
	return err
}

func (r *CotizacionRepository) DeleteCotizacion(id int64) error {
	if _, err := r.db.Exec("DELETE FROM cotizacion_items WHERE cotizacion_id = ?", id); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM cotizaciones WHERE id = ?", id)
	return err
}

func (r *CotizacionRepository) GetCotizacionConversionData(id int64) (string, int64, float64, string, error) {
	var estado string
	var sedeID int64
	var total float64
	var clienteNombre string
	err := r.db.QueryRow("SELECT estado, sede_id, total, cliente_nombre FROM cotizaciones WHERE id = ?", id).Scan(&estado, &sedeID, &total, &clienteNombre)
	return estado, sedeID, total, clienteNombre, err
}

func (r *CotizacionRepository) ConvertToVenta(id int64, usuarioID interface{}, sedeID int64, total float64, clienteNombre string) (int64, string, []VentaItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, "", nil, err
	}

	var count int
	tx.QueryRow("SELECT COUNT(*) FROM ventas").Scan(&count)
	numeroVenta := fmt.Sprintf("VNT-%d-%04d", time.Now().Year(), count+1)

	result, err := tx.Exec(`
		INSERT INTO ventas (numero_venta, cliente_nombre, total, usuario_id, sede_id, cotizacion_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, numeroVenta, clienteNombre, total, usuarioID, sedeID, id)
	if err != nil {
		tx.Rollback()
		return 0, "", nil, err
	}

	ventaID, _ := result.LastInsertId()

	rows, err := tx.Query(`
		SELECT producto_id, cantidad, precio_unit, subtotal
		FROM cotizacion_items WHERE cotizacion_id = ?
	`, id)
	if err != nil {
		tx.Rollback()
		return 0, "", nil, err
	}

	items := make([]VentaItem, 0)
	for rows.Next() {
		var i VentaItem
		if err := rows.Scan(&i.ProductoID, &i.Cantidad, &i.PrecioUnitario, &i.Subtotal); err != nil {
			continue
		}
		items = append(items, i)
	}
	rows.Close()

	for _, item := range items {
		if _, err := tx.Exec(`
			INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unit, subtotal)
			VALUES (?, ?, ?, ?, ?)
		`, ventaID, item.ProductoID, item.Cantidad, item.PrecioUnitario, item.Subtotal); err != nil {
			tx.Rollback()
			return 0, "", nil, err
		}
		if _, err := tx.Exec(`
			UPDATE stock_sedes SET cantidad = cantidad - ? WHERE producto_id = ? AND sede_id = ?
		`, item.Cantidad, item.ProductoID, sedeID); err != nil {
			tx.Rollback()
			return 0, "", nil, err
		}
	}

	tx.Exec("UPDATE cotizaciones SET estado = 'convertida', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)

	if err := tx.Commit(); err != nil {
		return 0, "", nil, err
	}
	return ventaID, numeroVenta, items, nil
}

func (r *CotizacionRepository) GetCotizacionPDF(id int64) (CotizacionPDF, error) {
	var c CotizacionPDF
	err := r.db.QueryRow(`
		SELECT c.numero_cotizacion, c.cliente_nombre, c.cliente_telefono, c.cliente_email,
		       c.validez, c.total, c.descuento, c.notas, c.created_at,
		       u.name, s.nombre, s.direccion
		FROM cotizaciones c
		INNER JOIN users u ON c.usuario_id = u.id
		INNER JOIN sedes s ON c.sede_id = s.id
		WHERE c.id = ?
	`, id).Scan(&c.NumeroCotizacion, &c.ClienteNombre, &c.ClienteTelefono, &c.ClienteEmail,
		&c.ValidezDias, &c.Total, &c.Descuento, &c.Notas, &c.FechaCreacion,
		&c.UsuarioNombre, &c.SedeNombre, &c.SedeDireccion)
	return c, err
}

func (r *CotizacionRepository) ListCotizacionPDFItems(id int64) ([]CotizacionPDFItem, error) {
	rows, err := r.db.Query(`
		SELECT p.codigo, p.name, p.brand, ci.cantidad, ci.precio_unit, ci.subtotal
		FROM cotizacion_items ci
		INNER JOIN products p ON ci.producto_id = p.id
		WHERE ci.cotizacion_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CotizacionPDFItem, 0)
	for rows.Next() {
		var i CotizacionPDFItem
		if err := rows.Scan(&i.Codigo, &i.Nombre, &i.Marca, &i.Cantidad, &i.PrecioUnitario, &i.Subtotal); err != nil {
			continue
		}
		items = append(items, i)
	}
	return items, nil
}
