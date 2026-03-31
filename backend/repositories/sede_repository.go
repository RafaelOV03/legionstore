package repositories

import (
	"database/sql"
	"errors"
	"smartech/backend/models"
)

var (
	ErrSedeNotFound           = errors.New("sede not found")
	ErrSedeHasAssociatedUsers = errors.New("sede has associated users")
)

type SedeRepository struct {
	db *sql.DB
}

func NewSedeRepository(db *sql.DB) *SedeRepository {
	return &SedeRepository{db: db}
}

func (r *SedeRepository) List() ([]models.Sede, error) {
	rows, err := r.db.Query(`SELECT id, created_at, updated_at, nombre, direccion, telefono, activa FROM sedes ORDER BY nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sedes := make([]models.Sede, 0)
	for rows.Next() {
		var sede models.Sede
		var activa int
		err := rows.Scan(&sede.ID, &sede.CreatedAt, &sede.UpdatedAt, &sede.Nombre, &sede.Direccion, &sede.Telefono, &activa)
		if err != nil {
			continue
		}
		sede.Activa = activa == 1
		sedes = append(sedes, sede)
	}

	return sedes, nil
}

func (r *SedeRepository) GetByID(id int64) (*models.Sede, error) {
	var sede models.Sede
	var activa int

	err := r.db.QueryRow(`SELECT id, created_at, updated_at, nombre, direccion, telefono, activa FROM sedes WHERE id = ?`, id).
		Scan(&sede.ID, &sede.CreatedAt, &sede.UpdatedAt, &sede.Nombre, &sede.Direccion, &sede.Telefono, &activa)
	if err == sql.ErrNoRows {
		return nil, ErrSedeNotFound
	}
	if err != nil {
		return nil, err
	}

	sede.Activa = activa == 1
	return &sede, nil
}

func (r *SedeRepository) Create(input models.Sede) (*models.Sede, error) {
	result, err := r.db.Exec(`INSERT INTO sedes (nombre, direccion, telefono, activa) VALUES (?, ?, ?, ?)`,
		input.Nombre, input.Direccion, input.Telefono, 1)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	input.ID = id
	input.Activa = true
	return &input, nil
}

func (r *SedeRepository) Update(id int64, input models.Sede) (*models.Sede, error) {
	activa := 0
	if input.Activa {
		activa = 1
	}

	result, err := r.db.Exec(`UPDATE sedes SET nombre = ?, direccion = ?, telefono = ?, activa = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		input.Nombre, input.Direccion, input.Telefono, activa, id)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrSedeNotFound
	}

	input.ID = id
	return &input, nil
}

func (r *SedeRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM sedes WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrSedeNotFound
	}

	return nil
}

func (r *SedeRepository) HasAssociatedUsers(id int64) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE sede_id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
