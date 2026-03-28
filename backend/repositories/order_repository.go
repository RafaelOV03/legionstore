package repositories

import (
	"database/sql"
	"smartech/backend/models"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetCartBySession(sessionID string) (models.Cart, error) {
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

func (r *OrderRepository) GetCartItemsWithProducts(cartID int64) ([]models.CartItem, error) {
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

func (r *OrderRepository) CreateOrderWithItems(sessionID, paypalOrderID string, total float64, items []models.CartItem) (models.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Order{}, err
	}

	result, err := tx.Exec(`
		INSERT INTO orders (session_id, paypal_order_id, status, total_amount, currency, finalized, created_at, updated_at)
		VALUES (?, ?, 'PENDING', ?, 'USD', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, sessionID, paypalOrderID, total)
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	for _, item := range items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, product_name, quantity, price, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, orderID, item.ProductID, item.Product.Name, item.Quantity, item.Product.PrecioVenta)
		if err != nil {
			tx.Rollback()
			return models.Order{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Order{}, err
	}

	return models.Order{
		ID:            orderID,
		SessionID:     sessionID,
		PayPalOrderID: paypalOrderID,
		Status:        "PENDING",
		TotalAmount:   total,
		Currency:      "USD",
		Finalized:     false,
	}, nil
}

func (r *OrderRepository) GetOrderByPayPalID(paypalOrderID string) (models.Order, error) {
	var order models.Order
	var finalized int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE paypal_order_id = ?
	`, paypalOrderID).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &finalized)
	if err != nil {
		return models.Order{}, err
	}
	order.Finalized = finalized == 1
	return order, nil
}

func (r *OrderRepository) GetOrderByID(orderID int64) (models.Order, error) {
	var order models.Order
	var finalized int
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE id = ?
	`, orderID).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
		&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
		&order.CompletedAt, &finalized)
	if err != nil {
		return models.Order{}, err
	}
	order.Finalized = finalized == 1
	return order, nil
}

func (r *OrderRepository) GetOrderItems(orderID int64) ([]models.OrderItem, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, order_id, product_id, product_name, quantity, price
		FROM order_items
		WHERE order_id = ?
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.OrderID,
			&item.ProductID, &item.ProductName, &item.Quantity, &item.Price)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) ListOrdersBySession(sessionID string) ([]models.Order, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE session_id = ?
		ORDER BY created_at DESC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		var finalized int
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
			&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
			&order.CompletedAt, &finalized)
		if err != nil {
			continue
		}
		order.Finalized = finalized == 1
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) ListAllOrders() ([]models.Order, error) {
	rows, err := r.db.Query(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, total_amount,
		       currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		var finalized int
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID, &order.PayPalOrderID,
			&order.Status, &order.TotalAmount, &order.Currency, &order.PayerEmail, &order.PayerName,
			&order.CompletedAt, &finalized)
		if err != nil {
			continue
		}
		order.Finalized = finalized == 1
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) MarkOrderCompleted(orderID int64, completedAt time.Time, payerEmail, payerName string) error {
	_, err := r.db.Exec(`
		UPDATE orders
		SET status = 'COMPLETED', completed_at = ?, payer_email = ?, payer_name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, completedAt, payerEmail, payerName, orderID)
	return err
}

func (r *OrderRepository) DecreaseStockByOrder(orderID int64) error {
	rows, err := r.db.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int64
		var quantity int
		if err := rows.Scan(&productID, &quantity); err != nil {
			continue
		}
		if _, err := r.db.Exec(`
			UPDATE products
			SET stock_quantity = stock_quantity - ?
			WHERE id = ?
		`, quantity, productID); err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepository) ClearCartBySession(sessionID string) error {
	var cartID int64
	err := r.db.QueryRow("SELECT id FROM carts WHERE session_id = ?", sessionID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err = r.db.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
	return err
}

func (r *OrderRepository) CountOrderByID(orderID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", orderID).Scan(&count)
	return count, err
}

func (r *OrderRepository) FinalizeOrder(orderID int64) error {
	_, err := r.db.Exec(`
		UPDATE orders SET finalized = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, orderID)
	return err
}

func (r *OrderRepository) UpdateOrderFields(orderID int64, status string, totalAmount float64) error {
	updates := make([]string, 0)
	args := make([]interface{}, 0)
	if status != "" {
		updates = append(updates, "status = ?")
		args = append(args, status)
	}
	if totalAmount > 0 {
		updates = append(updates, "total_amount = ?")
		args = append(args, totalAmount)
	}
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE orders SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += ", updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	args = append(args, orderID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *OrderRepository) DeleteOrder(orderID int64) error {
	if _, err := r.db.Exec("DELETE FROM order_items WHERE order_id = ?", orderID); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM orders WHERE id = ?", orderID)
	return err
}
