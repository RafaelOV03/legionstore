package repositories

import (
	"database/sql"
	"fmt"
	"smartech/backend/models"
	"time"
)

type OrdenTrabajoView struct {
	models.OrdenTrabajo
	TecnicoNombre string `json:"tecnico_nombre"`
	SedeNombre    string `json:"sede_nombre"`
}

type InsumoOrdenView struct {
	ID           int64  `json:"id"`
	InsumoID     int64  `json:"insumo_id"`
	Cantidad     int    `json:"cantidad"`
	InsumoNombre string `json:"insumo_nombre"`
	InsumoCodigo string `json:"insumo_codigo"`
}

type TrazabilidadItem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Accion    string    `json:"accion"`
	Detalle   string    `json:"detalle"`
	Usuario   string    `json:"usuario"`
}

type OrdenesStats struct {
	Total         int `json:"total"`
	Recibidos     int `json:"recibidos"`
	EnDiagnostico int `json:"en_diagnostico"`
	EnReparacion  int `json:"en_reparacion"`
	Terminados    int `json:"terminados"`
	Entregados    int `json:"entregados"`
	Urgentes      int `json:"urgentes"`
}

type TecnicoView struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	SedeID     *int64  `json:"sede_id"`
	SedeNombre *string `json:"sede_nombre"`
}

type CreateOrdenTrabajoParams struct {
	NumeroOrden       string
	ClienteNombre     string
	ClienteTelefono   string
	Equipo            string
	Marca             string
	Modelo            string
	NumSerie          string
	ProblemaReportado string
	Prioridad         string
	FechaPromesa      interface{}
	TecnicoID         *int64
	SedeID            int64
	Notas             string
}

type UpdateOrdenTrabajoParams struct {
	DiagnosticoTecnico string
	SolucionAplicada   string
	Estado             string
	Prioridad          string
	TecnicoID          *int64
	CostoManoObra      float64
	CostoRepuestos     float64
	Notas              string
	FechaPromesa       interface{}
	FechaEntrega       interface{}
}

type OrdenTrabajoRepository struct {
	db *sql.DB
}

func NewOrdenTrabajoRepository(db *sql.DB) *OrdenTrabajoRepository {
	return &OrdenTrabajoRepository{db: db}
}

func (r *OrdenTrabajoRepository) ListOrdenesTrabajo(estado, prioridad, tecnicoID, sedeID string) ([]OrdenTrabajoView, error) {
	query := `
		SELECT ot.id, ot.created_at, ot.updated_at, ot.numero_orden, ot.cliente_nombre,
		       ot.cliente_telefono, ot.equipo, ot.num_serie, ot.marca, ot.modelo,
		       ot.problema_reportado, ot.diagnostico_tecnico, ot.solucion_aplicada, ot.estado,
		       ot.prioridad, ot.fecha_ingreso, ot.fecha_promesa, ot.fecha_entrega, ot.costo_mano_obra,
		       ot.costo_repuestos, ot.tecnico_id, ot.sede_id, ot.notas,
		       COALESCE(u.name, 'Sin asignar') as tecnico_nombre, s.nombre as sede_nombre
		FROM ordenes_trabajo ot
		LEFT JOIN users u ON ot.tecnico_id = u.id
		INNER JOIN sedes s ON ot.sede_id = s.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)

	if estado != "" {
		query += " AND ot.estado = ?"
		args = append(args, estado)
	}
	if prioridad != "" {
		query += " AND ot.prioridad = ?"
		args = append(args, prioridad)
	}
	if tecnicoID != "" {
		query += " AND ot.tecnico_id = ?"
		args = append(args, tecnicoID)
	}
	if sedeID != "" {
		query += " AND ot.sede_id = ?"
		args = append(args, sedeID)
	}

	query += " ORDER BY CASE ot.prioridad WHEN 'urgente' THEN 1 WHEN 'alta' THEN 2 WHEN 'media' THEN 3 ELSE 4 END, ot.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordenes := make([]OrdenTrabajoView, 0)
	for rows.Next() {
		var o OrdenTrabajoView
		var fechaPromesa, fechaEntrega sql.NullTime
		var tecnicoIDN sql.NullInt64
		var diagnostico, solucion, notas sql.NullString
		err := rows.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.NumeroOrden, &o.ClienteNombre,
			&o.ClienteTelefono, &o.Equipo, &o.NumSerie, &o.Marca, &o.Modelo,
			&o.ProblemaReportado, &diagnostico, &solucion, &o.Estado,
			&o.Prioridad, &o.FechaIngreso, &fechaPromesa, &fechaEntrega, &o.CostoManoObra,
			&o.CostoRepuestos, &tecnicoIDN, &o.SedeID, &notas,
			&o.TecnicoNombre, &o.SedeNombre)
		if err != nil {
			continue
		}
		if diagnostico.Valid {
			o.DiagnosticoTecnico = diagnostico.String
		}
		if solucion.Valid {
			o.SolucionAplicada = solucion.String
		}
		if notas.Valid {
			o.Notas = notas.String
		}
		if fechaPromesa.Valid {
			o.FechaPromesa = &fechaPromesa.Time
		}
		if fechaEntrega.Valid {
			o.FechaEntrega = &fechaEntrega.Time
		}
		if tecnicoIDN.Valid {
			o.TecnicoID = &tecnicoIDN.Int64
		}
		ordenes = append(ordenes, o)
	}

	return ordenes, nil
}

func (r *OrdenTrabajoRepository) GetOrdenTrabajoByID(id int64) (models.OrdenTrabajo, error) {
	var o models.OrdenTrabajo
	var fechaPromesa, fechaEntrega sql.NullTime
	var tecnicoIDN sql.NullInt64
	var diagnostico, solucion, notas sql.NullString

	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, numero_orden, cliente_nombre,
		       cliente_telefono, equipo, num_serie, marca, modelo,
		       problema_reportado, diagnostico_tecnico, solucion_aplicada, estado,
		       prioridad, fecha_ingreso, fecha_promesa, fecha_entrega, costo_mano_obra,
		       costo_repuestos, tecnico_id, sede_id, notas
		FROM ordenes_trabajo WHERE id = ?
	`, id).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.NumeroOrden, &o.ClienteNombre,
		&o.ClienteTelefono, &o.Equipo, &o.NumSerie, &o.Marca, &o.Modelo,
		&o.ProblemaReportado, &diagnostico, &solucion, &o.Estado,
		&o.Prioridad, &o.FechaIngreso, &fechaPromesa, &fechaEntrega, &o.CostoManoObra,
		&o.CostoRepuestos, &tecnicoIDN, &o.SedeID, &notas)
	if err != nil {
		return models.OrdenTrabajo{}, err
	}

	if diagnostico.Valid {
		o.DiagnosticoTecnico = diagnostico.String
	}
	if solucion.Valid {
		o.SolucionAplicada = solucion.String
	}
	if notas.Valid {
		o.Notas = notas.String
	}
	if fechaPromesa.Valid {
		o.FechaPromesa = &fechaPromesa.Time
	}
	if fechaEntrega.Valid {
		o.FechaEntrega = &fechaEntrega.Time
	}
	if tecnicoIDN.Valid {
		o.TecnicoID = &tecnicoIDN.Int64
	}

	return o, nil
}

func (r *OrdenTrabajoRepository) ListInsumosByOrdenID(id int64) ([]InsumoOrdenView, error) {
	rows, err := r.db.Query(`
		SELECT io.id, io.insumo_id, io.cantidad, i.nombre, i.codigo
		FROM insumos_orden io
		INNER JOIN insumos i ON io.insumo_id = i.id
		WHERE io.orden_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	insumos := make([]InsumoOrdenView, 0)
	for rows.Next() {
		var ins InsumoOrdenView
		if err := rows.Scan(&ins.ID, &ins.InsumoID, &ins.Cantidad, &ins.InsumoNombre, &ins.InsumoCodigo); err != nil {
			continue
		}
		insumos = append(insumos, ins)
	}
	return insumos, nil
}

func (r *OrdenTrabajoRepository) ListTrazabilidadByOrdenID(id int64) ([]TrazabilidadItem, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.created_at, t.accion, t.detalle, u.name
		FROM trazabilidad t
		INNER JOIN users u ON t.usuario_id = u.id
		WHERE t.orden_trabajo_id = ?
		ORDER BY t.created_at DESC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]TrazabilidadItem, 0)
	for rows.Next() {
		var t TrazabilidadItem
		if err := rows.Scan(&t.ID, &t.CreatedAt, &t.Accion, &t.Detalle, &t.Usuario); err != nil {
			continue
		}
		items = append(items, t)
	}
	return items, nil
}

func (r *OrdenTrabajoRepository) CountOrdenesTrabajo() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo").Scan(&count)
	return count, err
}

func (r *OrdenTrabajoRepository) BuildNumeroOrden(count int) string {
	return fmt.Sprintf("OT-%d-%04d", time.Now().Year(), count+1)
}

func (r *OrdenTrabajoRepository) CreateOrdenTrabajo(params CreateOrdenTrabajoParams) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO ordenes_trabajo (numero_orden, cliente_nombre, cliente_telefono, equipo,
		                             marca, modelo, num_serie, problema_reportado,
		                             estado, prioridad, fecha_ingreso, fecha_promesa, tecnico_id,
		                             sede_id, notas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'recibido', ?, CURRENT_TIMESTAMP, ?, ?, ?, ?)
	`, params.NumeroOrden, params.ClienteNombre, params.ClienteTelefono, params.Equipo,
		params.Marca, params.Modelo, params.NumSerie, params.ProblemaReportado,
		params.Prioridad, params.FechaPromesa, params.TecnicoID, params.SedeID, params.Notas)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *OrdenTrabajoRepository) InsertTrazabilidad(ordenID int64, accion, detalle string, usuarioID interface{}) error {
	_, err := r.db.Exec(`
		INSERT INTO trazabilidad (orden_trabajo_id, accion, detalle, usuario_id)
		VALUES (?, ?, ?, ?)
	`, ordenID, accion, detalle, usuarioID)
	return err
}

func (r *OrdenTrabajoRepository) GetEstadoByOrdenID(id int64) (string, error) {
	var estado string
	err := r.db.QueryRow("SELECT estado FROM ordenes_trabajo WHERE id = ?", id).Scan(&estado)
	return estado, err
}

func (r *OrdenTrabajoRepository) UpdateOrdenTrabajo(id int64, params UpdateOrdenTrabajoParams) error {
	_, err := r.db.Exec(`
		UPDATE ordenes_trabajo SET
		    diagnostico_tecnico = ?, solucion_aplicada = ?, estado = ?, prioridad = ?, tecnico_id = ?,
		    costo_mano_obra = ?, costo_repuestos = ?, notas = ?, fecha_promesa = COALESCE(?, fecha_promesa),
		    fecha_entrega = COALESCE(?, fecha_entrega), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, params.DiagnosticoTecnico, params.SolucionAplicada, params.Estado, params.Prioridad, params.TecnicoID,
		params.CostoManoObra, params.CostoRepuestos, params.Notas, params.FechaPromesa, params.FechaEntrega, id)
	return err
}

func (r *OrdenTrabajoRepository) GetTecnicoNombre(tecnicoID int64) (string, error) {
	var name string
	err := r.db.QueryRow("SELECT name FROM users WHERE id = ?", tecnicoID).Scan(&name)
	return name, err
}

func (r *OrdenTrabajoRepository) AssignTecnico(id, tecnicoID int64) error {
	_, err := r.db.Exec(`
		UPDATE ordenes_trabajo SET tecnico_id = ?,
		    estado = CASE WHEN estado = 'recibido' THEN 'en_diagnostico' ELSE estado END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, tecnicoID, id)
	return err
}

func (r *OrdenTrabajoRepository) GetInsumoStockAndNombre(insumoID int64) (string, int, error) {
	var nombre string
	var stock int
	err := r.db.QueryRow("SELECT nombre, stock FROM insumos WHERE id = ?", insumoID).Scan(&nombre, &stock)
	return nombre, stock, err
}

func (r *OrdenTrabajoRepository) AddInsumoToOrden(ordenID, insumoID int64, cantidad int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO insumos_orden (orden_id, insumo_id, cantidad)
		VALUES (?, ?, ?)
	`, ordenID, insumoID, cantidad)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE insumos SET stock = stock - ? WHERE id = ?", cantidad, insumoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *OrdenTrabajoRepository) DeleteOrdenTrabajo(id int64) error {
	if _, err := r.db.Exec("DELETE FROM trazabilidad WHERE orden_trabajo_id = ?", id); err != nil {
		return err
	}
	if _, err := r.db.Exec("DELETE FROM insumos_orden WHERE orden_id = ?", id); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM ordenes_trabajo WHERE id = ?", id)
	return err
}

func (r *OrdenTrabajoRepository) Stats() (OrdenesStats, error) {
	var stats OrdenesStats
	if err := r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo").Scan(&stats.Total); err != nil {
		return OrdenesStats{}, err
	}
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'recibido'").Scan(&stats.Recibidos)
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'en_diagnostico'").Scan(&stats.EnDiagnostico)
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'en_reparacion'").Scan(&stats.EnReparacion)
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'terminado'").Scan(&stats.Terminados)
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE estado = 'entregado'").Scan(&stats.Entregados)
	r.db.QueryRow("SELECT COUNT(*) FROM ordenes_trabajo WHERE prioridad = 'urgente' AND estado NOT IN ('entregado', 'cancelado')").Scan(&stats.Urgentes)
	return stats, nil
}

func (r *OrdenTrabajoRepository) ListTecnicos() ([]TecnicoView, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.email, u.sede_id, s.nombre as sede_nombre
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		LEFT JOIN sedes s ON u.sede_id = s.id
		WHERE r.name = 'tecnico'
		ORDER BY u.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tecnicos := make([]TecnicoView, 0)
	for rows.Next() {
		var t TecnicoView
		var sedeID sql.NullInt64
		var sedeNombre sql.NullString
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &sedeID, &sedeNombre); err != nil {
			continue
		}
		if sedeID.Valid {
			t.SedeID = &sedeID.Int64
		}
		if sedeNombre.Valid {
			t.SedeNombre = &sedeNombre.String
		}
		tecnicos = append(tecnicos, t)
	}
	return tecnicos, nil
}
