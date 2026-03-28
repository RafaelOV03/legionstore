package repositories

import (
	"database/sql"
	"smartech/backend/models"
)

type LogView struct {
	models.LogAuditoria
	UsuarioNombre string `json:"usuario_nombre"`
}

type AccionStats struct {
	Accion string `json:"accion"`
	Total  int    `json:"total"`
}

type UserStats struct {
	Usuario string `json:"usuario"`
	Total   int    `json:"total"`
}

type DayStats struct {
	Fecha string `json:"fecha"`
	Total int    `json:"total"`
}

type CatStats struct {
	Categoria string  `json:"categoria"`
	Total     float64 `json:"total"`
}

type SedeStats struct {
	Sede  string  `json:"sede"`
	Total float64 `json:"total"`
}

type PromocionCreateInput struct {
	Nombre         string
	Descripcion    string
	Tipo           string
	Valor          float64
	FechaInicio    string
	FechaFin       string
	ProductosIDs   string
	SegmentacionID *int64
}

type PromocionUpdateInput struct {
	Nombre         string
	Descripcion    string
	Tipo           string
	Valor          float64
	FechaInicio    string
	FechaFin       string
	ProductosIDs   string
	SegmentacionID *int64
	Activa         bool
}

type AuditoriaRepository struct {
	db *sql.DB
}

func NewAuditoriaRepository(db *sql.DB) *AuditoriaRepository {
	return &AuditoriaRepository{db: db}
}

func (r *AuditoriaRepository) ListLogs(accion, entidad, usuarioID, fechaDesde, fechaHasta, limit string) ([]LogView, error) {
	query := `
		SELECT l.id, l.created_at, l.usuario_id, l.accion, l.entidad, l.entidad_id,
		       l.valor_anterior, l.valor_nuevo, l.ip_address, u.name as usuario_nombre
		FROM logs_auditoria l
		INNER JOIN users u ON l.usuario_id = u.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)

	if accion != "" {
		query += " AND l.accion = ?"
		args = append(args, accion)
	}
	if entidad != "" {
		query += " AND l.entidad = ?"
		args = append(args, entidad)
	}
	if usuarioID != "" {
		query += " AND l.usuario_id = ?"
		args = append(args, usuarioID)
	}
	if fechaDesde != "" {
		query += " AND l.created_at >= ?"
		args = append(args, fechaDesde)
	}
	if fechaHasta != "" {
		query += " AND l.created_at <= ?"
		args = append(args, fechaHasta+" 23:59:59")
	}

	query += " ORDER BY l.created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]LogView, 0)
	for rows.Next() {
		var item LogView
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.UsuarioID, &item.Accion, &item.Entidad, &item.EntidadID,
			&item.ValorAnterior, &item.ValorNuevo, &item.IPAddress, &item.UsuarioNombre)
		if err != nil {
			continue
		}
		logs = append(logs, item)
	}

	return logs, nil
}

func (r *AuditoriaRepository) ListAccionStats() ([]AccionStats, error) {
	rows, err := r.db.Query(`
		SELECT accion, COUNT(*) as total
		FROM logs_auditoria
		WHERE created_at >= date('now', '-30 days')
		GROUP BY accion
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]AccionStats, 0)
	for rows.Next() {
		var item AccionStats
		if err := rows.Scan(&item.Accion, &item.Total); err != nil {
			continue
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (r *AuditoriaRepository) ListActividadPorUsuario() ([]UserStats, error) {
	rows, err := r.db.Query(`
		SELECT u.name, COUNT(*) as total
		FROM logs_auditoria l
		INNER JOIN users u ON l.usuario_id = u.id
		WHERE l.created_at >= date('now', '-30 days')
		GROUP BY l.usuario_id
		ORDER BY total DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]UserStats, 0)
	for rows.Next() {
		var item UserStats
		if err := rows.Scan(&item.Usuario, &item.Total); err != nil {
			continue
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (r *AuditoriaRepository) ListActividadPorDia() ([]DayStats, error) {
	rows, err := r.db.Query(`
		SELECT date(created_at) as fecha, COUNT(*) as total
		FROM logs_auditoria
		WHERE created_at >= date('now', '-7 days')
		GROUP BY date(created_at)
		ORDER BY fecha
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]DayStats, 0)
	for rows.Next() {
		var item DayStats
		if err := rows.Scan(&item.Fecha, &item.Total); err != nil {
			continue
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (r *AuditoriaRepository) TotalVentas(fechaDesde, fechaHasta, sedeID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total), 0) FROM ventas WHERE 1=1
	`
	args := make([]interface{}, 0)

	if fechaDesde != "" {
		query += " AND created_at >= ?"
		args = append(args, fechaDesde)
	}
	if fechaHasta != "" {
		query += " AND created_at <= ?"
		args = append(args, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		query += " AND sede_id = ?"
		args = append(args, sedeID)
	}

	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func (r *AuditoriaRepository) CostoProductos(fechaDesde, fechaHasta, sedeID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(vi.cantidad * p.precio_compra), 0)
		FROM venta_items vi
		INNER JOIN products p ON vi.producto_id = p.id
		INNER JOIN ventas v ON vi.venta_id = v.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)

	if fechaDesde != "" {
		query += " AND v.created_at >= ?"
		args = append(args, fechaDesde)
	}
	if fechaHasta != "" {
		query += " AND v.created_at <= ?"
		args = append(args, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		query += " AND v.sede_id = ?"
		args = append(args, sedeID)
	}

	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func (r *AuditoriaRepository) TotalServicios(fechaDesde, fechaHasta, sedeID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(costo_servicio + costo_repuestos), 0)
		FROM ordenes_trabajo
		WHERE estado = 'entregado'
	`
	args := make([]interface{}, 0)

	if fechaDesde != "" {
		query += " AND fecha_entrega >= ?"
		args = append(args, fechaDesde)
	}
	if fechaHasta != "" {
		query += " AND fecha_entrega <= ?"
		args = append(args, fechaHasta+" 23:59:59")
	}
	if sedeID != "" {
		query += " AND sede_id = ?"
		args = append(args, sedeID)
	}

	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func (r *AuditoriaRepository) ListVentasPorCategoria() ([]CatStats, error) {
	rows, err := r.db.Query(`
		SELECT p.category, COALESCE(SUM(vi.subtotal), 0) as total
		FROM venta_items vi
		INNER JOIN products p ON vi.producto_id = p.id
		INNER JOIN ventas v ON vi.venta_id = v.id
		GROUP BY p.category
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]CatStats, 0)
	for rows.Next() {
		var item CatStats
		if err := rows.Scan(&item.Categoria, &item.Total); err != nil {
			continue
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (r *AuditoriaRepository) ListVentasPorSede() ([]SedeStats, error) {
	rows, err := r.db.Query(`
		SELECT s.nombre, COALESCE(SUM(v.total), 0) as total
		FROM ventas v
		INNER JOIN sedes s ON v.sede_id = s.id
		GROUP BY v.sede_id
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]SedeStats, 0)
	for rows.Next() {
		var item SedeStats
		if err := rows.Scan(&item.Sede, &item.Total); err != nil {
			continue
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (r *AuditoriaRepository) ListSegmentaciones() ([]models.Segmentacion, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, nombre, descripcion, criterios, activa
		FROM segmentaciones ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	segmentaciones := make([]models.Segmentacion, 0)
	for rows.Next() {
		var s models.Segmentacion
		var activa int
		err := rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Nombre, &s.Descripcion, &s.Criterios, &activa)
		if err != nil {
			continue
		}
		s.Activa = activa == 1
		segmentaciones = append(segmentaciones, s)
	}
	return segmentaciones, nil
}

func (r *AuditoriaRepository) CreateSegmentacion(nombre, descripcion, criterios string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO segmentaciones (nombre, descripcion, criterios, activa)
		VALUES (?, ?, ?, 1)
	`, nombre, descripcion, criterios)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AuditoriaRepository) ListPromociones(activas bool) ([]models.Promocion, error) {
	query := `
		SELECT id, created_at, updated_at, nombre, descripcion, tipo_descuento, valor, fecha_inicio,
		       fecha_fin, producto_ids, categorias, activa
		FROM promociones
	`
	if activas {
		query += ` WHERE activa = 1 AND fecha_inicio <= date('now') AND fecha_fin >= date('now')`
	}
	query += " ORDER BY fecha_inicio DESC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	promociones := make([]models.Promocion, 0)
	for rows.Next() {
		var p models.Promocion
		var activaInt int
		var prodIDs, categorias sql.NullString
		var fechaInicio, fechaFin sql.NullTime
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Nombre, &p.Descripcion, &p.Tipo,
			&p.Valor, &fechaInicio, &fechaFin, &prodIDs, &categorias, &activaInt)
		if err != nil {
			continue
		}
		p.Activa = activaInt == 1
		if prodIDs.Valid {
			p.ProductosIDs = prodIDs.String
		}
		if fechaInicio.Valid {
			p.FechaInicio = fechaInicio.Time
		}
		if fechaFin.Valid {
			p.FechaFin = fechaFin.Time
		}
		promociones = append(promociones, p)
	}
	return promociones, nil
}

func (r *AuditoriaRepository) CreatePromocion(input PromocionCreateInput) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO promociones (nombre, descripcion, tipo, valor, fecha_inicio, fecha_fin,
		                         productos_ids, segmentacion_id, activa)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
	`, input.Nombre, input.Descripcion, input.Tipo, input.Valor, input.FechaInicio, input.FechaFin,
		input.ProductosIDs, input.SegmentacionID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AuditoriaRepository) UpdatePromocion(id int64, input PromocionUpdateInput) error {
	activa := 0
	if input.Activa {
		activa = 1
	}

	_, err := r.db.Exec(`
		UPDATE promociones SET nombre = ?, descripcion = ?, tipo = ?, valor = ?, fecha_inicio = ?,
		                   fecha_fin = ?, productos_ids = ?, segmentacion_id = ?, activa = ?,
		                   updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, input.Nombre, input.Descripcion, input.Tipo, input.Valor, input.FechaInicio, input.FechaFin,
		input.ProductosIDs, input.SegmentacionID, activa, id)
	return err
}

func (r *AuditoriaRepository) DeletePromocion(id int64) error {
	_, err := r.db.Exec("DELETE FROM promociones WHERE id = ?", id)
	return err
}
