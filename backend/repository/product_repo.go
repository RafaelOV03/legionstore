package repository

import (
	"database/sql"
	"smartech/backend/models"
)

// ProductRepository maneja toda la lógica de acceso a datos de productos
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository crea una nueva instancia de ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ProductWithStock es el modelo de producto con stock total
type ProductWithStock struct {
	models.Product
	StockTotal int `json:"stock_total"`
}

// GetAll obtiene todos los productos activos con stock
func (r *ProductRepository) GetAll() ([]ProductWithStock, error) {
	query := `
		SELECT p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description, 
		       p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
		       COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.producto_id
		WHERE p.activo = 1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, DBError(err, "GetAll products")
	}
	defer rows.Close()

	var products []ProductWithStock
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name,
			&product.Description, &product.PrecioCompra, &product.PrecioVenta, &product.Category,
			&product.Brand, &product.ImageURL, &product.Images, &activo, &product.StockTotal)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, rows.Err()
}

// GetByID obtiene un producto por su ID
func (r *ProductRepository) GetByID(id int64) (*models.Product, error) {
	var product models.Product
	var activo int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, 
		       category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name,
		&product.Description, &product.PrecioCompra, &product.PrecioVenta, &product.Category,
		&product.Brand, &product.ImageURL, &product.Images, &activo)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "GetByID product")
	}

	product.Activo = activo == 1
	return &product, nil
}

// Create crea un nuevo producto
func (r *ProductRepository) Create(product *models.Product) error {
	result, err := r.db.Exec(`
		INSERT INTO products (codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, product.Codigo, product.Name, product.Description, product.PrecioCompra, product.PrecioVenta,
		product.Category, product.Brand, product.ImageURL, product.Images, 1)

	if err != nil {
		return DBError(err, "Create product")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return DBError(err, "Get last insert id")
	}

	product.ID = id
	return nil
}

// Update actualiza un producto existente
func (r *ProductRepository) Update(id int64, product *models.Product) error {
	result, err := r.db.Exec(`
		UPDATE products
		SET codigo = ?, name = ?, description = ?, precio_compra = ?, precio_venta = ?,
		    category = ?, brand = ?, image_url = ?, images = ?
		WHERE id = ?
	`, product.Codigo, product.Name, product.Description, product.PrecioCompra, product.PrecioVenta,
		product.Category, product.Brand, product.ImageURL, product.Images, id)

	if err != nil {
		return DBError(err, "Update product")
	}

	return CheckRowsAffected(result, "Update product")
}

// Delete marca un producto como inactivo (soft delete)
func (r *ProductRepository) Delete(id int64) error {
	result, err := r.db.Exec(`
		UPDATE products SET activo = 0 WHERE id = ?
	`, id)

	if err != nil {
		return DBError(err, "Delete product")
	}

	return CheckRowsAffected(result, "Delete product")
}

// GetByCategory obtiene productos de una categoría específica
func (r *ProductRepository) GetByCategory(category string) ([]ProductWithStock, error) {
	query := `
		SELECT p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description, 
		       p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
		       COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.producto_id
		WHERE p.activo = 1 AND p.category = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, DBError(err, "GetByCategory products")
	}
	defer rows.Close()

	var products []ProductWithStock
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name,
			&product.Description, &product.PrecioCompra, &product.PrecioVenta, &product.Category,
			&product.Brand, &product.ImageURL, &product.Images, &activo, &product.StockTotal)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, rows.Err()
}

// Search busca productos por nombre o código
func (r *ProductRepository) Search(searchTerm string) ([]ProductWithStock, error) {
	query := `
		SELECT p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description, 
		       p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
		       COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.producto_id
		WHERE p.activo = 1 AND (p.name LIKE ? OR p.codigo LIKE ?)
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, DBError(err, "Search products")
	}
	defer rows.Close()

	var products []ProductWithStock
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name,
			&product.Description, &product.PrecioCompra, &product.PrecioVenta, &product.Category,
			&product.Brand, &product.ImageURL, &product.Images, &activo, &product.StockTotal)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, rows.Err()
}

// GetRandomProducts devuelve productos aleatorios
func (r *ProductRepository) GetRandomProducts(limit int) ([]ProductWithStock, error) {
	query := `
		SELECT 
			p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description,
			p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
			COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.product_id
		WHERE p.activo = 1
		GROUP BY p.id
		ORDER BY RANDOM()
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, DBError(err, "Get random products")
	}
	defer rows.Close()

	var products []ProductWithStock
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name,
			&product.Description, &product.PrecioCompra, &product.PrecioVenta, &product.Category,
			&product.Brand, &product.ImageURL, &product.Images, &activo, &product.StockTotal)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, rows.Err()
}
