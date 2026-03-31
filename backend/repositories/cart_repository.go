package repositories

import (
	"database/sql"
	"smartech/backend/models"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetCartBySession(sessionID string) (models.Cart, error) {
	var cart models.Cart
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, sessionID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)
	if err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}

func (r *CartRepository) CreateCart(sessionID string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO carts (session_id, created_at, updated_at)
		VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, sessionID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CartRepository) GetProductStock(productID int64) (int, error) {
	var stockQuantity int
	err := r.db.QueryRow(`SELECT stock_quantity FROM products WHERE id = ?`, productID).Scan(&stockQuantity)
	return stockQuantity, err
}

func (r *CartRepository) GetCartItemsWithProducts(cartID int64) ([]models.CartItem, error) {
	rows, err := r.db.Query(`
		SELECT ci.id, ci.created_at, ci.updated_at, ci.cart_id, ci.product_id, ci.quantity,
		       p.id, p.created_at, p.updated_at, p.name, p.description, p.precio_venta, p.precio_compra,
		       p.category, p.brand, p.image_url, p.images, p.activo
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = ?
	`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.CartItem, 0)
	for rows.Next() {
		var item models.CartItem
		var product models.Product
		var activo int
		err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.CartID, &item.ProductID, &item.Quantity,
			&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Name, &product.Description,
			&product.PrecioVenta, &product.PrecioCompra, &product.Category, &product.Brand,
			&product.ImageURL, &product.Images, &activo,
		)
		if err != nil {
			continue
		}
		product.Activo = activo == 1
		item.Product = &product
		items = append(items, item)
	}

	return items, nil
}

func (r *CartRepository) GetCartItemByCartAndProduct(cartID, productID int64) (int64, int, error) {
	var itemID int64
	var quantity int
	err := r.db.QueryRow(`
		SELECT id, quantity FROM cart_items
		WHERE cart_id = ? AND product_id = ?
	`, cartID, productID).Scan(&itemID, &quantity)
	return itemID, quantity, err
}

func (r *CartRepository) UpdateCartItemQuantity(itemID int64, quantity int) error {
	_, err := r.db.Exec(`
		UPDATE cart_items SET quantity = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, quantity, itemID)
	return err
}

func (r *CartRepository) InsertCartItem(cartID, productID int64, quantity int) error {
	_, err := r.db.Exec(`
		INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, cartID, productID, quantity)
	return err
}

func (r *CartRepository) GetCartItemWithProduct(itemID int64) (models.CartItem, error) {
	var cartItem models.CartItem
	var product models.Product
	var activo int
	err := r.db.QueryRow(`
		SELECT ci.id, ci.created_at, ci.updated_at, ci.cart_id, ci.product_id, ci.quantity,
		       p.id, p.created_at, p.updated_at, p.name, p.description, p.precio_venta, p.precio_compra,
		       p.category, p.brand, p.image_url, p.images, p.activo
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		WHERE ci.id = ?
	`, itemID).Scan(
		&cartItem.ID, &cartItem.CreatedAt, &cartItem.UpdatedAt, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity,
		&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Name, &product.Description,
		&product.PrecioVenta, &product.PrecioCompra, &product.Category, &product.Brand,
		&product.ImageURL, &product.Images, &activo,
	)
	if err != nil {
		return models.CartItem{}, err
	}

	product.Activo = activo == 1
	cartItem.Product = &product
	return cartItem, nil
}

func (r *CartRepository) GetSessionByCartID(cartID int64) (string, error) {
	var sessionID string
	err := r.db.QueryRow("SELECT session_id FROM carts WHERE id = ?", cartID).Scan(&sessionID)
	return sessionID, err
}

func (r *CartRepository) GetCartIDBySession(sessionID string) (int64, error) {
	var cartID int64
	err := r.db.QueryRow("SELECT id FROM carts WHERE session_id = ?", sessionID).Scan(&cartID)
	return cartID, err
}

func (r *CartRepository) DeleteCartItem(itemID int64) error {
	_, err := r.db.Exec("DELETE FROM cart_items WHERE id = ?", itemID)
	return err
}

func (r *CartRepository) DeleteCartItemsByCartID(cartID int64) error {
	_, err := r.db.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
	return err
}

func (r *CartRepository) GetCartIDByItem(itemID int64) (int64, error) {
	var cartID int64
	err := r.db.QueryRow("SELECT cart_id FROM cart_items WHERE id = ?", itemID).Scan(&cartID)
	return cartID, err
}
