package repositories

import (
	"database/sql"
	"smartech/backend/models"
)

type CompatibilidadView struct {
	ID               int64  `json:"id"`
	ProductoID       int64  `json:"producto_id"`
	CompatibleCon    int64  `json:"compatible_con"`
	Notas            string `json:"notas"`
	ProductoNombre   string `json:"producto_nombre"`
	ProductoMarca    string `json:"producto_marca"`
	CompatibleNombre string `json:"compatible_nombre"`
	CompatibleMarca  string `json:"compatible_marca"`
}

type ProductoBase struct {
	ID       int64
	Name     string
	Brand    string
	Category string
}

type CompatibleProductoView struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Brand     string  `json:"brand"`
	Category  string  `json:"category"`
	Precio    float64 `json:"precio"`
	Notas     string  `json:"notas"`
	TipoMatch string  `json:"tipo_match"`
}

type InsumoStats struct {
	TotalInsumos    int     `json:"total_insumos"`
	BajoStock       int     `json:"bajo_stock"`
	SinStock        int     `json:"sin_stock"`
	ValorInventario float64 `json:"valor_inventario"`
}

type InsumoRepository struct {
	db *sql.DB
}

func NewInsumoRepository(db *sql.DB) *InsumoRepository {
	return &InsumoRepository{db: db}
}

func (r *InsumoRepository) ListInsumos(categoria string, bajoStock bool) ([]models.Insumo, error) {
	query := `
		SELECT id, created_at, updated_at, codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo
		FROM insumos WHERE 1=1
	`
	args := make([]interface{}, 0)
	if categoria != "" {
		query += " AND categoria = ?"
		args = append(args, categoria)
	}
	if bajoStock {
		query += " AND stock <= stock_minimo"
	}
	query += " ORDER BY nombre"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Insumo, 0)
	for rows.Next() {
		var i models.Insumo
		var activo int
		if err := rows.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Codigo, &i.Nombre, &i.Descripcion,
			&i.Categoria, &i.UnidadMedida, &i.Stock, &i.StockMinimo, &i.Costo, &i.SedeID, &activo); err != nil {
			continue
		}
		i.Activo = activo == 1
		items = append(items, i)
	}
	return items, nil
}

func (r *InsumoRepository) GetInsumoByID(id int64) (models.Insumo, error) {
	var i models.Insumo
	var activo int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo
		FROM insumos WHERE id = ?
	`, id).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Codigo, &i.Nombre, &i.Descripcion,
		&i.Categoria, &i.UnidadMedida, &i.Stock, &i.StockMinimo, &i.Costo, &i.SedeID, &activo)
	if err != nil {
		return models.Insumo{}, err
	}
	i.Activo = activo == 1
	return i, nil
}

func (r *InsumoRepository) InsertInsumo(i models.Insumo) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO insumos (codigo, nombre, descripcion, categoria, unidad_medida, stock, stock_minimo, costo, sede_id, activo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
	`, i.Codigo, i.Nombre, i.Descripcion, i.Categoria, i.UnidadMedida, i.Stock, i.StockMinimo, i.Costo, i.SedeID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *InsumoRepository) UpdateInsumo(id int64, i models.Insumo) error {
	activo := 0
	if i.Activo {
		activo = 1
	}
	_, err := r.db.Exec(`
		UPDATE insumos SET codigo = ?, nombre = ?, descripcion = ?, categoria = ?, unidad_medida = ?, stock = ?,
		               stock_minimo = ?, costo = ?, activo = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, i.Codigo, i.Nombre, i.Descripcion, i.Categoria, i.UnidadMedida, i.Stock,
		i.StockMinimo, i.Costo, activo, id)
	return err
}

func (r *InsumoRepository) GetStock(id int64) (int, error) {
	var stock int
	err := r.db.QueryRow("SELECT stock FROM insumos WHERE id = ?", id).Scan(&stock)
	return stock, err
}

func (r *InsumoRepository) UpdateStock(id int64, newStock int) error {
	_, err := r.db.Exec("UPDATE insumos SET stock = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newStock, id)
	return err
}

func (r *InsumoRepository) DeleteInsumo(id int64) error {
	_, err := r.db.Exec("DELETE FROM insumos WHERE id = ?", id)
	return err
}

func (r *InsumoRepository) ListCompatibilidades(productoID string) ([]CompatibilidadView, error) {
	query := `
		SELECT c.id, c.producto_id, c.compatible_con_id, c.notas,
		       p1.name as producto_nombre, p1.brand as producto_marca,
		       p2.name as compatible_nombre, p2.brand as compatible_marca
		FROM compatibilidades c
		INNER JOIN products p1 ON c.producto_id = p1.id
		INNER JOIN products p2 ON c.compatible_con_id = p2.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	if productoID != "" {
		query += " AND (c.producto_id = ? OR c.compatible_con_id = ?)"
		args = append(args, productoID, productoID)
	}
	query += " ORDER BY p1.name"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CompatibilidadView, 0)
	for rows.Next() {
		var c CompatibilidadView
		if err := rows.Scan(&c.ID, &c.ProductoID, &c.CompatibleCon, &c.Notas,
			&c.ProductoNombre, &c.ProductoMarca, &c.CompatibleNombre, &c.CompatibleMarca); err != nil {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *InsumoRepository) GetProductoBase(productoID string) (ProductoBase, error) {
	var p ProductoBase
	err := r.db.QueryRow("SELECT id, name, brand, category FROM products WHERE id = ?", productoID).Scan(&p.ID, &p.Name, &p.Brand, &p.Category)
	return p, err
}

func (r *InsumoRepository) ListCompatiblesDirectos(productoID string) ([]CompatibleProductoView, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.brand, p.category, p.precio_venta, COALESCE(c.notas, '')
		FROM compatibilidades c
		INNER JOIN products p ON p.id = CASE WHEN c.producto_id = ? THEN c.compatible_con_id ELSE c.producto_id END
		WHERE c.producto_id = ? OR c.compatible_con_id = ?
	`, productoID, productoID, productoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CompatibleProductoView, 0)
	for rows.Next() {
		var c CompatibleProductoView
		if err := rows.Scan(&c.ID, &c.Name, &c.Brand, &c.Category, &c.Precio, &c.Notas); err != nil {
			continue
		}
		c.TipoMatch = "directo"
		out = append(out, c)
	}
	return out, nil
}

func (r *InsumoRepository) ListCompatiblesByCategoriaMarca(productoID string, category string, brand string) ([]CompatibleProductoView, error) {
	rows, err := r.db.Query(`
		SELECT id, name, brand, category, precio_venta
		FROM products
		WHERE id != ? AND activo = 1 AND (category = ? OR brand = ?)
		LIMIT 20
	`, productoID, category, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CompatibleProductoView, 0)
	for rows.Next() {
		var c CompatibleProductoView
		if err := rows.Scan(&c.ID, &c.Name, &c.Brand, &c.Category, &c.Precio); err != nil {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *InsumoRepository) CountCompatibilidadPair(a, b int64) (int, error) {
	var exists int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM compatibilidades
		WHERE (producto_id = ? AND compatible_con_id = ?) OR (producto_id = ? AND compatible_con_id = ?)
	`, a, b, b, a).Scan(&exists)
	return exists, err
}

func (r *InsumoRepository) InsertCompatibilidad(a, b int64, notas string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO compatibilidades (producto_id, compatible_con_id, notas)
		VALUES (?, ?, ?)
	`, a, b, notas)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *InsumoRepository) DeleteCompatibilidad(id int64) error {
	_, err := r.db.Exec("DELETE FROM compatibilidades WHERE id = ?", id)
	return err
}

func (r *InsumoRepository) Stats() (InsumoStats, error) {
	var s InsumoStats
	if err := r.db.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1").Scan(&s.TotalInsumos); err != nil {
		return InsumoStats{}, err
	}
	r.db.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1 AND stock <= stock_minimo AND stock > 0").Scan(&s.BajoStock)
	r.db.QueryRow("SELECT COUNT(*) FROM insumos WHERE activo = 1 AND stock = 0").Scan(&s.SinStock)
	r.db.QueryRow("SELECT COALESCE(SUM(stock * costo), 0) FROM insumos WHERE activo = 1").Scan(&s.ValorInventario)
	return s, nil
}
