package repositories

import (
	"database/sql"
	"fmt"
	"smartech/backend/models"
	"time"
)

type RMAView struct {
	models.RMA
	ProductoNombre string `json:"producto_nombre"`
	ProductoMarca  string `json:"producto_marca"`
	UsuarioNombre  string `json:"usuario_nombre"`
	SedeNombre     string `json:"sede_nombre"`
}

type HistorialItem struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	EstadoAnterior string    `json:"estado_anterior"`
	EstadoNuevo    string    `json:"estado_nuevo"`
	Comentario     string    `json:"comentario"`
	Usuario        string    `json:"usuario"`
}

type RMAStats struct {
	Total      int `json:"total"`
	Recibidos  int `json:"recibidos"`
	EnRevision int `json:"en_revision"`
	Resueltos  int `json:"resueltos"`
	Rechazados int `json:"rechazados"`
}

type RMARepository struct {
	db *sql.DB
}

func NewRMARepository(db *sql.DB) *RMARepository {
	return &RMARepository{db: db}
}

func (r *RMARepository) ListRMAs(estado, sedeID string) ([]RMAView, error) {
	query := `
		SELECT r.id, r.created_at, r.updated_at, r.numero_rma, r.producto_id, r.cliente_nombre,
		       r.cliente_telefono, r.cliente_email, r.num_serie, r.fecha_compra, r.motivo_devolucion,
		       r.diagnostico, r.estado, r.solucion, r.fecha_resolucion, r.usuario_id, r.sede_id, r.notas,
		       p.name as producto_nombre, p.brand as producto_marca,
		       u.name as usuario_nombre,
		       s.nombre as sede_nombre
		FROM rmas r
		INNER JOIN products p ON r.producto_id = p.id
		INNER JOIN users u ON r.usuario_id = u.id
		INNER JOIN sedes s ON r.sede_id = s.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	if estado != "" {
		query += " AND r.estado = ?"
		args = append(args, estado)
	}
	if sedeID != "" {
		query += " AND r.sede_id = ?"
		args = append(args, sedeID)
	}
	query += " ORDER BY r.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rmas := make([]RMAView, 0)
	for rows.Next() {
		var item RMAView
		var fechaCompra, fechaResolucion sql.NullTime
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.NumeroRMA, &item.ProductoID,
			&item.ClienteNombre, &item.ClienteTelefono, &item.ClienteEmail, &item.NumSerie, &fechaCompra,
			&item.MotivoDevolucion, &item.Diagnostico, &item.Estado, &item.Solucion, &fechaResolucion,
			&item.UsuarioID, &item.SedeID, &item.Notas,
			&item.ProductoNombre, &item.ProductoMarca, &item.UsuarioNombre, &item.SedeNombre)
		if err != nil {
			continue
		}
		if fechaCompra.Valid {
			item.FechaCompra = fechaCompra.Time
		}
		if fechaResolucion.Valid {
			item.FechaResolucion = &fechaResolucion.Time
		}
		rmas = append(rmas, item)
	}
	return rmas, nil
}

func (r *RMARepository) GetRMAByID(id int64) (models.RMA, error) {
	var item models.RMA
	var fechaCompra, fechaResolucion sql.NullTime
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, numero_rma, producto_id, cliente_nombre,
		       cliente_telefono, cliente_email, num_serie, fecha_compra, motivo_devolucion,
		       diagnostico, estado, solucion, fecha_resolucion, usuario_id, sede_id, notas
		FROM rmas WHERE id = ?
	`, id).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.NumeroRMA, &item.ProductoID,
		&item.ClienteNombre, &item.ClienteTelefono, &item.ClienteEmail, &item.NumSerie, &fechaCompra,
		&item.MotivoDevolucion, &item.Diagnostico, &item.Estado, &item.Solucion, &fechaResolucion,
		&item.UsuarioID, &item.SedeID, &item.Notas)
	if err != nil {
		return models.RMA{}, err
	}
	if fechaCompra.Valid {
		item.FechaCompra = fechaCompra.Time
	}
	if fechaResolucion.Valid {
		item.FechaResolucion = &fechaResolucion.Time
	}
	return item, nil
}

func (r *RMARepository) ListHistorial(rmaID int64) ([]HistorialItem, error) {
	rows, err := r.db.Query(`
		SELECT h.id, h.created_at, h.estado_anterior, h.estado_nuevo, h.comentario, u.name
		FROM historial_rmas h
		INNER JOIN users u ON h.usuario_id = u.id
		WHERE h.rma_id = ?
		ORDER BY h.created_at DESC
	`, rmaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	h := make([]HistorialItem, 0)
	for rows.Next() {
		var item HistorialItem
		if err := rows.Scan(&item.ID, &item.CreatedAt, &item.EstadoAnterior, &item.EstadoNuevo, &item.Comentario, &item.Usuario); err != nil {
			continue
		}
		h = append(h, item)
	}
	return h, nil
}

func (r *RMARepository) CountRMAs() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&count)
	return count, err
}

func (r *RMARepository) BuildRMANumber(count int) string {
	return fmt.Sprintf("RMA-%d-%04d", time.Now().Year(), count+1)
}

func (r *RMARepository) InsertRMA(numeroRMA string, req models.RMA, fechaCompra interface{}, userID interface{}) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO rmas (numero_rma, producto_id, cliente_nombre, cliente_telefono, cliente_email,
		                  num_serie, fecha_compra, motivo_devolucion, estado, usuario_id, sede_id, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'recibido', ?, ?, ?)
	`, numeroRMA, req.ProductoID, req.ClienteNombre, req.ClienteTelefono, req.ClienteEmail,
		req.NumSerie, fechaCompra, req.MotivoDevolucion, userID, req.SedeID, req.Notas)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *RMARepository) InsertHistorial(rmaID int64, estadoAnt, estadoNuevo, comentario string, userID interface{}) error {
	_, err := r.db.Exec(`INSERT INTO historial_rmas (rma_id, estado_anterior, estado_nuevo, comentario, usuario_id)
		VALUES (?, ?, ?, ?, ?)`, rmaID, estadoAnt, estadoNuevo, comentario, userID)
	return err
}

func (r *RMARepository) GetEstadoByID(id int64) (string, error) {
	var estado string
	err := r.db.QueryRow("SELECT estado FROM rmas WHERE id = ?", id).Scan(&estado)
	return estado, err
}

func (r *RMARepository) UpdateRMA(id int64, diagnostico, estado, solucion, notas string, fechaResolucion interface{}) error {
	_, err := r.db.Exec(`
		UPDATE rmas SET diagnostico = ?, estado = ?, solucion = ?, notas = ?,
		                fecha_resolucion = COALESCE(?, fecha_resolucion), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, diagnostico, estado, solucion, notas, fechaResolucion, id)
	return err
}

func (r *RMARepository) DeleteRMA(id int64) error {
	if _, err := r.db.Exec("DELETE FROM historial_rmas WHERE rma_id = ?", id); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM rmas WHERE id = ?", id)
	return err
}

func (r *RMARepository) Stats() (RMAStats, error) {
	var s RMAStats
	if err := r.db.QueryRow("SELECT COUNT(*) FROM rmas").Scan(&s.Total); err != nil {
		return RMAStats{}, err
	}
	r.db.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'recibido'").Scan(&s.Recibidos)
	r.db.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'en_revision'").Scan(&s.EnRevision)
	r.db.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'resuelto'").Scan(&s.Resueltos)
	r.db.QueryRow("SELECT COUNT(*) FROM rmas WHERE estado = 'rechazado'").Scan(&s.Rechazados)
	return s, nil
}
