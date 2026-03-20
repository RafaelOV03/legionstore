package repository

import (
	"database/sql"
	"smartech/backend/models"
)

// CartRepository maneja toda la lógica de acceso a datos del carrito
type CartRepository struct {
	db *sql.DB
}

// NewCartRepository crea una nueva instancia de CartRepository
func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetBySessionID obtiene el carrito de una sesión (o lo crea si no existe)
func (r *CartRepository) GetBySessionID(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id
		FROM carts
		WHERE session_id = ?
	`, sessionID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.SessionID)

	if err == sql.ErrNoRows {
		// Crear carrito nuevo
		result, err := r.db.Exec(`
			INSERT INTO carts (session_id, created_at, updated_at)
			VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, sessionID)

		if err != nil {
			return nil, DBError(err, "Create cart")
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, DBError(err, "Get last insert id")
		}

		cart.ID = id
		cart.SessionID = sessionID
		return &cart, nil
	}

	if err != nil {
		return nil, DBError(err, "Get cart by session")
	}

	return &cart, nil
}

// GetItems obtiene todos los items del carrito
func (r *CartRepository) GetItems(cartID int64) ([]models.CartItem, error) {
	query := `
		SELECT id, created_at, updated_at, cart_id, product_id, quantity
		FROM cart_items
		WHERE cart_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, cartID)
	if err != nil {
		return nil, DBError(err, "Get cart items")
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.CartID, &item.ProductID, &item.Quantity)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// AddItem agrega un item al carrito
func (r *CartRepository) AddItem(cartID int64, productID int64, quantity int) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := r.db.Exec(query, cartID, productID, quantity)
	if err != nil {
		return DBError(err, "Add cart item")
	}
	return CheckRowsAffected(result, "Add cart item")
}

// UpdateItem actualiza la cantidad de un item
func (r *CartRepository) UpdateItem(itemID int64, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, quantity, itemID)
	if err != nil {
		return DBError(err, "Update cart item")
	}
	return CheckRowsAffected(result, "Update cart item")
}

// RemoveItem elimina un item del carrito
func (r *CartRepository) RemoveItem(itemID int64) error {
	query := `DELETE FROM cart_items WHERE id = ?`
	result, err := r.db.Exec(query, itemID)
	if err != nil {
		return DBError(err, "Remove cart item")
	}
	return CheckRowsAffected(result, "Remove cart item")
}

// Clear vacía todos los items del carrito
func (r *CartRepository) Clear(cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = ?`
	result, err := r.db.Exec(query, cartID)
	if err != nil {
		return DBError(err, "Clear cart")
	}
	return CheckRowsAffected(result, "Clear cart")
}

// CountItems cuenta los items en el carrito
func (r *CartRepository) CountItems(cartID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM cart_items WHERE cart_id = ?
	`, cartID).Scan(&count)
	return count, ScanError(err, "count cart items")
}
