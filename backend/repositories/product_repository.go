package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"smartech/backend/models"
	"strings"
)

var ErrProductNotFound = errors.New("product not found")

type ProductWithStock struct {
	models.Product
	StockTotal int `json:"stock_total"`
}

type ProductCreateInput struct {
	Codigo       string
	Name         string
	Description  string
	PrecioCompra float64
	PrecioVenta  float64
	Category     string
	Brand        string
	ImageURL     string
	Images       []string
}

type ProductUpdateInput struct {
	Codigo       *string
	Name         *string
	Description  *string
	PrecioCompra *float64
	PrecioVenta  *float64
	Category     *string
	Brand        *string
	ImageURL     *string
	Images       *string
	Activo       *bool
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListActiveWithStock() ([]ProductWithStock, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.codigo, p.name, p.description,
		       p.precio_compra, p.precio_venta, p.category, p.brand, p.image_url, p.images, p.activo,
		       COALESCE(SUM(ss.cantidad), 0) as stock_total
		FROM products p
		LEFT JOIN stock_sedes ss ON p.id = ss.producto_id
		WHERE p.activo = 1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]ProductWithStock, 0)
	for rows.Next() {
		var product ProductWithStock
		var activo int
		err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo,
			&product.StockTotal,
		)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id int64) (*models.Product, error) {
	var product models.Product
	var activo int

	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, id).Scan(
		&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
		&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo,
	)
	if err == sql.ErrNoRows {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	product.Activo = activo == 1
	return &product, nil
}

func (r *ProductRepository) Create(input ProductCreateInput) (*models.Product, error) {
	imagesJSON := "[]"
	if len(input.Images) > 0 {
		imagesBytes, _ := json.Marshal(input.Images)
		imagesJSON = string(imagesBytes)
	}

	result, err := r.db.Exec(`
		INSERT INTO products (codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, input.Codigo, input.Name, input.Description, input.PrecioCompra, input.PrecioVenta, input.Category, input.Brand, input.ImageURL, imagesJSON)
	if err != nil {
		return nil, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		ID:           productID,
		Codigo:       input.Codigo,
		Name:         input.Name,
		Description:  input.Description,
		PrecioCompra: input.PrecioCompra,
		PrecioVenta:  input.PrecioVenta,
		Category:     input.Category,
		Brand:        input.Brand,
		ImageURL:     input.ImageURL,
		Images:       imagesJSON,
		Activo:       true,
	}
	return product, nil
}

func (r *ProductRepository) ExistsByID(id int64) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProductRepository) UpdatePartial(id int64, input ProductUpdateInput) error {
	updates := []string{}
	args := []interface{}{}

	if input.Codigo != nil {
		updates = append(updates, "codigo = ?")
		args = append(args, *input.Codigo)
	}
	if input.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *input.Description)
	}
	if input.PrecioCompra != nil {
		updates = append(updates, "precio_compra = ?")
		args = append(args, *input.PrecioCompra)
	}
	if input.PrecioVenta != nil {
		updates = append(updates, "precio_venta = ?")
		args = append(args, *input.PrecioVenta)
	}
	if input.Category != nil {
		updates = append(updates, "category = ?")
		args = append(args, *input.Category)
	}
	if input.Brand != nil {
		updates = append(updates, "brand = ?")
		args = append(args, *input.Brand)
	}
	if input.ImageURL != nil {
		updates = append(updates, "image_url = ?")
		args = append(args, *input.ImageURL)
	}
	if input.Images != nil {
		updates = append(updates, "images = ?")
		args = append(args, *input.Images)
	}
	if input.Activo != nil {
		activo := 0
		if *input.Activo {
			activo = 1
		}
		updates = append(updates, "activo = ?")
		args = append(args, activo)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)
	query := "UPDATE products SET " + strings.Join(updates, ", ") + " WHERE id = ?"

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProductRepository) SoftDelete(id int64) error {
	_, err := r.db.Exec("UPDATE products SET activo = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *ProductRepository) ListRandom(limit int) ([]models.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE activo = 1
		ORDER BY RANDOM()
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		var activo int
		err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo,
		)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) ListByCategory(category string) ([]models.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE category = ? AND activo = 1
		ORDER BY created_at DESC
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		var activo int
		err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
			&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo,
		)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		products = append(products, product)
	}

	return products, nil
}
