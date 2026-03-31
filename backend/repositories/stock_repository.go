package repositories

import "database/sql"

type StockItem struct {
	ID          int64  `json:"id"`
	SedeID      int64  `json:"sede_id"`
	ProductoID  int64  `json:"producto_id"`
	Cantidad    int    `json:"cantidad"`
	StockMinimo int    `json:"stock_minimo"`
	StockMaximo int    `json:"stock_maximo"`
	SedeNombre  string `json:"sede_nombre"`
	Codigo      string `json:"codigo"`
	Producto    string `json:"producto"`
	Categoria   string `json:"categoria"`
	Marca       string `json:"marca"`
}

type StockProducto struct {
	ID          int64   `json:"id"`
	ProductoID  int64   `json:"producto_id"`
	Cantidad    int     `json:"cantidad"`
	StockMinimo int     `json:"stock_minimo"`
	StockMaximo int     `json:"stock_maximo"`
	Codigo      string  `json:"codigo"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Categoria   string  `json:"categoria"`
	Marca       string  `json:"marca"`
	ImageURL    string  `json:"image_url"`
}

type StockUpdateInput struct {
	SedeID      int64
	ProductoID  int64
	Cantidad    int
	StockMinimo int
	StockMaximo int
}

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) ListMultisede() ([]StockItem, error) {
	rows, err := r.db.Query(`
		SELECT ss.id, ss.sede_id, ss.producto_id, ss.cantidad, ss.stock_minimo, ss.stock_maximo,
		       s.nombre as sede_nombre,
		       p.codigo, p.name, p.category, p.brand
		FROM stock_sedes ss
		INNER JOIN sedes s ON ss.sede_id = s.id
		INNER JOIN products p ON ss.producto_id = p.id
		WHERE p.activo = 1 AND s.activa = 1
		ORDER BY p.name, s.nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]StockItem, 0)
	for rows.Next() {
		var item StockItem
		err := rows.Scan(&item.ID, &item.SedeID, &item.ProductoID, &item.Cantidad, &item.StockMinimo, &item.StockMaximo,
			&item.SedeNombre, &item.Codigo, &item.Producto, &item.Categoria, &item.Marca)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *StockRepository) ListBySede(sedeID int64) ([]StockProducto, error) {
	rows, err := r.db.Query(`
		SELECT ss.id, ss.producto_id, ss.cantidad, ss.stock_minimo, ss.stock_maximo,
		       p.codigo, p.name, p.description, p.precio_venta, p.category, p.brand, p.image_url
		FROM stock_sedes ss
		INNER JOIN products p ON ss.producto_id = p.id
		WHERE ss.sede_id = ? AND p.activo = 1
		ORDER BY p.name
	`, sedeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]StockProducto, 0)
	for rows.Next() {
		var item StockProducto
		err := rows.Scan(&item.ID, &item.ProductoID, &item.Cantidad, &item.StockMinimo, &item.StockMaximo,
			&item.Codigo, &item.Nombre, &item.Descripcion, &item.Precio, &item.Categoria, &item.Marca, &item.ImageURL)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *StockRepository) Upsert(input StockUpdateInput) error {
	var existingID int64
	err := r.db.QueryRow("SELECT id FROM stock_sedes WHERE sede_id = ? AND producto_id = ?", input.SedeID, input.ProductoID).Scan(&existingID)

	if err == sql.ErrNoRows {
		_, err = r.db.Exec(`INSERT INTO stock_sedes (sede_id, producto_id, cantidad, stock_minimo, stock_maximo) VALUES (?, ?, ?, ?, ?)`,
			input.SedeID, input.ProductoID, input.Cantidad, input.StockMinimo, input.StockMaximo)
		return err
	}
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`UPDATE stock_sedes SET cantidad = ?, stock_minimo = ?, stock_maximo = ?, updated_at = CURRENT_TIMESTAMP WHERE sede_id = ? AND producto_id = ?`,
		input.Cantidad, input.StockMinimo, input.StockMaximo, input.SedeID, input.ProductoID)
	return err
}
