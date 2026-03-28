package repositories

import (
	"database/sql"
	"fmt"
	"smartech/backend/models"
	"time"
)

type TraspasoView struct {
	models.Traspaso
	SedeOrigenNombre    string `json:"sede_origen_nombre"`
	SedeDestinoNombre   string `json:"sede_destino_nombre"`
	UsuarioOrigenNombre string `json:"usuario_origen_nombre"`
}

type TraspasoItemView struct {
	ID               int64  `json:"id"`
	ProductoID       int64  `json:"producto_id"`
	Cantidad         int    `json:"cantidad"`
	CantidadRecibida int    `json:"cantidad_recibida"`
	ProductoNombre   string `json:"producto_nombre"`
	ProductoMarca    string `json:"producto_marca"`
	ProductoCodigo   string `json:"producto_codigo"`
}

type TraspasoCreateItemInput struct {
	ProductoID int64
	Cantidad   int
}

type TraspasoRecibirItemInput struct {
	ItemID           int64
	CantidadRecibida int
}

type TraspasoRepository struct {
	db *sql.DB
}

func NewTraspasoRepository(db *sql.DB) *TraspasoRepository {
	return &TraspasoRepository{db: db}
}

func (r *TraspasoRepository) ListTraspasos(estado, sedeOrigenID, sedeDestinoID string) ([]TraspasoView, error) {
	query := `
		SELECT t.id, t.created_at, t.updated_at, t.numero_traspaso, t.sede_origen_id, t.sede_destino_id,
		       t.estado, t.fecha_envio, t.fecha_recepcion, t.notas, t.usuario_envia_id, t.usuario_recibe_id,
		       so.nombre as sede_origen_nombre, sd.nombre as sede_destino_nombre,
		       uo.name as usuario_origen_nombre
		FROM traspasos t
		INNER JOIN sedes so ON t.sede_origen_id = so.id
		INNER JOIN sedes sd ON t.sede_destino_id = sd.id
		INNER JOIN users uo ON t.usuario_envia_id = uo.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)

	if estado != "" {
		query += " AND t.estado = ?"
		args = append(args, estado)
	}
	if sedeOrigenID != "" {
		query += " AND t.sede_origen_id = ?"
		args = append(args, sedeOrigenID)
	}
	if sedeDestinoID != "" {
		query += " AND t.sede_destino_id = ?"
		args = append(args, sedeDestinoID)
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	traspasos := make([]TraspasoView, 0)
	for rows.Next() {
		var t TraspasoView
		var fechaEnvio, fechaRecepcion sql.NullTime
		var usuarioRecibeID sql.NullInt64
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.NumeroTraspaso, &t.SedeOrigenID, &t.SedeDestinoID,
			&t.Estado, &fechaEnvio, &fechaRecepcion, &t.Notas, &t.UsuarioEnviaID, &usuarioRecibeID,
			&t.SedeOrigenNombre, &t.SedeDestinoNombre, &t.UsuarioOrigenNombre)
		if err != nil {
			continue
		}
		if fechaEnvio.Valid {
			t.FechaEnvio = &fechaEnvio.Time
		}
		if fechaRecepcion.Valid {
			t.FechaRecepcion = &fechaRecepcion.Time
		}
		if usuarioRecibeID.Valid {
			t.UsuarioRecibeID = &usuarioRecibeID.Int64
		}
		traspasos = append(traspasos, t)
	}

	return traspasos, nil
}

func (r *TraspasoRepository) GetTraspasoByID(id int64) (models.Traspaso, error) {
	var t models.Traspaso
	var fechaEnvio, fechaRecepcion sql.NullTime
	var usuarioRecibeID sql.NullInt64

	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, numero_traspaso, sede_origen_id, sede_destino_id,
		       estado, fecha_envio, fecha_recepcion, notas, usuario_envia_id, usuario_recibe_id
		FROM traspasos WHERE id = ?
	`, id).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.NumeroTraspaso, &t.SedeOrigenID, &t.SedeDestinoID,
		&t.Estado, &fechaEnvio, &fechaRecepcion, &t.Notas, &t.UsuarioEnviaID, &usuarioRecibeID)
	if err != nil {
		return models.Traspaso{}, err
	}

	if fechaEnvio.Valid {
		t.FechaEnvio = &fechaEnvio.Time
	}
	if fechaRecepcion.Valid {
		t.FechaRecepcion = &fechaRecepcion.Time
	}
	if usuarioRecibeID.Valid {
		t.UsuarioRecibeID = &usuarioRecibeID.Int64
	}

	return t, nil
}

func (r *TraspasoRepository) ListTraspasoItems(traspasoID int64) ([]TraspasoItemView, error) {
	rows, err := r.db.Query(`
		SELECT ti.id, ti.producto_id, ti.cantidad, ti.cantidad_recibida,
		       p.name, p.brand, p.codigo
		FROM traspaso_items ti
		INNER JOIN products p ON ti.producto_id = p.id
		WHERE ti.traspaso_id = ?
	`, traspasoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]TraspasoItemView, 0)
	for rows.Next() {
		var item TraspasoItemView
		var cantRecibida sql.NullInt64
		if err := rows.Scan(&item.ID, &item.ProductoID, &item.Cantidad, &cantRecibida,
			&item.ProductoNombre, &item.ProductoMarca, &item.ProductoCodigo); err != nil {
			continue
		}
		if cantRecibida.Valid {
			item.CantidadRecibida = int(cantRecibida.Int64)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *TraspasoRepository) CountSedeByID(sedeID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM sedes WHERE id = ?", sedeID).Scan(&count)
	return count, err
}

func (r *TraspasoRepository) GetStockDisponible(productoID, sedeID int64) (int, error) {
	var stock int
	err := r.db.QueryRow(`
		SELECT COALESCE(cantidad, 0) FROM stock_sedes
		WHERE producto_id = ? AND sede_id = ?
	`, productoID, sedeID).Scan(&stock)
	return stock, err
}

func (r *TraspasoRepository) GetProductName(productoID int64) (string, error) {
	var name string
	err := r.db.QueryRow("SELECT name FROM products WHERE id = ?", productoID).Scan(&name)
	return name, err
}

func (r *TraspasoRepository) CountTraspasos() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM traspasos").Scan(&count)
	return count, err
}

func (r *TraspasoRepository) BuildNumeroTraspaso(count int) string {
	return fmt.Sprintf("TRP-%d-%04d", time.Now().Year(), count+1)
}

func (r *TraspasoRepository) CreateTraspasoWithItems(numero string, sedeOrigenID, sedeDestinoID int64, notas string, usuarioEnviaID interface{}, items []TraspasoCreateItemInput) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO traspasos (numero_traspaso, sede_origen_id, sede_destino_id, estado, notas, usuario_envia_id)
		VALUES (?, ?, ?, 'pendiente', ?, ?)
	`, numero, sedeOrigenID, sedeDestinoID, notas, usuarioEnviaID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	traspasoID, _ := result.LastInsertId()

	for _, item := range items {
		_, err = tx.Exec(`
			INSERT INTO traspaso_items (traspaso_id, producto_id, cantidad)
			VALUES (?, ?, ?)
		`, traspasoID, item.ProductoID, item.Cantidad)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return traspasoID, nil
}

func (r *TraspasoRepository) GetEstadoAndSedeOrigen(id int64) (string, int64, error) {
	var estado string
	var sedeOrigenID int64
	err := r.db.QueryRow("SELECT estado, sede_origen_id FROM traspasos WHERE id = ?", id).Scan(&estado, &sedeOrigenID)
	return estado, sedeOrigenID, err
}

func (r *TraspasoRepository) GetEstadoAndSedeDestino(id int64) (string, int64, error) {
	var estado string
	var sedeDestinoID int64
	err := r.db.QueryRow("SELECT estado, sede_destino_id FROM traspasos WHERE id = ?", id).Scan(&estado, &sedeDestinoID)
	return estado, sedeDestinoID, err
}

func (r *TraspasoRepository) GetEstado(id int64) (string, error) {
	var estado string
	err := r.db.QueryRow("SELECT estado FROM traspasos WHERE id = ?", id).Scan(&estado)
	return estado, err
}

func (r *TraspasoRepository) EnviarTraspaso(id, sedeOrigenID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query("SELECT producto_id, cantidad FROM traspaso_items WHERE traspaso_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	type item struct {
		ProductoID int64
		Cantidad   int
	}
	items := make([]item, 0)
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.ProductoID, &it.Cantidad); err != nil {
			continue
		}
		items = append(items, it)
	}
	rows.Close()

	for _, it := range items {
		_, err = tx.Exec(`
			UPDATE stock_sedes SET cantidad = cantidad - ? WHERE producto_id = ? AND sede_id = ?
		`, it.Cantidad, it.ProductoID, sedeOrigenID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	now := time.Now()
	_, err = tx.Exec("UPDATE traspasos SET estado = 'enviado', fecha_envio = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TraspasoRepository) RecibirTraspasoAuto(id, sedeDestinoID int64, usuarioRecibeID interface{}, notas string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query("SELECT id, producto_id, cantidad FROM traspaso_items WHERE traspaso_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	type autoItem struct {
		ID         int64
		ProductoID int64
		Cantidad   int
	}
	autoItems := make([]autoItem, 0)
	for rows.Next() {
		var ai autoItem
		if err := rows.Scan(&ai.ID, &ai.ProductoID, &ai.Cantidad); err != nil {
			continue
		}
		autoItems = append(autoItems, ai)
	}
	rows.Close()

	for _, ai := range autoItems {
		_, err = tx.Exec("UPDATE traspaso_items SET cantidad_recibida = ? WHERE id = ?", ai.Cantidad, ai.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var exists int
		tx.QueryRow("SELECT COUNT(*) FROM stock_sedes WHERE producto_id = ? AND sede_id = ?", ai.ProductoID, sedeDestinoID).Scan(&exists)
		if exists > 0 {
			_, err = tx.Exec(`
				UPDATE stock_sedes SET cantidad = cantidad + ? WHERE producto_id = ? AND sede_id = ?
			`, ai.Cantidad, ai.ProductoID, sedeDestinoID)
		} else {
			_, err = tx.Exec(`
				INSERT INTO stock_sedes (producto_id, sede_id, cantidad, stock_minimo) VALUES (?, ?, ?, 5)
			`, ai.ProductoID, sedeDestinoID, ai.Cantidad)
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	now := time.Now()
	_, err = tx.Exec(`
		UPDATE traspasos SET estado = 'recibido', fecha_recepcion = ?, usuario_recibe_id = ?,
		                   notas = CASE WHEN ? != '' THEN notas || ' | Recepción: ' || ? ELSE notas END,
		                   updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, now, usuarioRecibeID, notas, notas, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TraspasoRepository) RecibirTraspasoConItems(id, sedeDestinoID int64, usuarioRecibeID interface{}, notas string, items []TraspasoRecibirItemInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, item := range items {
		var productoID int64
		if err := tx.QueryRow("SELECT producto_id FROM traspaso_items WHERE id = ?", item.ItemID).Scan(&productoID); err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("UPDATE traspaso_items SET cantidad_recibida = ? WHERE id = ?", item.CantidadRecibida, item.ItemID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var exists int
		tx.QueryRow("SELECT COUNT(*) FROM stock_sedes WHERE producto_id = ? AND sede_id = ?", productoID, sedeDestinoID).Scan(&exists)
		if exists > 0 {
			_, err = tx.Exec(`
				UPDATE stock_sedes SET cantidad = cantidad + ? WHERE producto_id = ? AND sede_id = ?
			`, item.CantidadRecibida, productoID, sedeDestinoID)
		} else {
			_, err = tx.Exec(`
				INSERT INTO stock_sedes (producto_id, sede_id, cantidad, stock_minimo) VALUES (?, ?, ?, 5)
			`, productoID, sedeDestinoID, item.CantidadRecibida)
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	now := time.Now()
	_, err = tx.Exec(`
		UPDATE traspasos SET estado = 'recibido', fecha_recepcion = ?, usuario_recibe_id = ?,
		                   notas = CASE WHEN ? != '' THEN notas || ' | Recepción: ' || ? ELSE notas END,
		                   updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, now, usuarioRecibeID, notas, notas, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TraspasoRepository) CancelarTraspaso(id int64) error {
	_, err := r.db.Exec("UPDATE traspasos SET estado = 'cancelado', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *TraspasoRepository) DeleteTraspaso(id int64) error {
	if _, err := r.db.Exec("DELETE FROM traspaso_items WHERE traspaso_id = ?", id); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM traspasos WHERE id = ?", id)
	return err
}
