package repositories

import (
	"database/sql"
	"smartech/backend/models"
)

type UploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) GetProductImages(productID int64) (string, error) {
	var images string
	err := r.db.QueryRow("SELECT images FROM products WHERE id = ?", productID).Scan(&images)
	return images, err
}

func (r *UploadRepository) CountProductByID(productID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productID).Scan(&count)
	return count, err
}

func (r *UploadRepository) UpdateProductImages(productID int64, images string) error {
	_, err := r.db.Exec(`
		UPDATE products SET images = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, images, productID)
	return err
}

func (r *UploadRepository) GetProductByID(productID int64) (models.Product, error) {
	var product models.Product
	var activo int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, codigo, name, description, precio_compra, precio_venta, category, brand, image_url, images, activo
		FROM products
		WHERE id = ?
	`, productID).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Codigo, &product.Name, &product.Description,
		&product.PrecioCompra, &product.PrecioVenta, &product.Category, &product.Brand, &product.ImageURL, &product.Images, &activo)
	if err != nil {
		return models.Product{}, err
	}
	product.Activo = activo == 1
	return product, nil
}
