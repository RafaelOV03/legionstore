package repository

import (
	"database/sql"
	"smartech/backend/models"
)

// OrderRepository maneja toda la lógica de acceso a datos de órdenes
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository crea una nueva instancia de OrderRepository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// GetAll obtiene todas las órdenes
func (r *OrderRepository) GetAll() ([]models.Order, error) {
	query := `
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, 
		       total_amount, currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, DBError(err, "Get all orders")
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var completedAt sql.NullTime
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID,
			&order.PayPalOrderID, &order.Status, &order.TotalAmount, &order.Currency,
			&order.PayerEmail, &order.PayerName, &completedAt, &order.Finalized)
		if err != nil {
			continue
		}
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// GetByID obtiene una orden por su ID
func (r *OrderRepository) GetByID(id int64) (*models.Order, error) {
	var order models.Order
	var completedAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, 
		       total_amount, currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE id = ?
	`, id).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID,
		&order.PayPalOrderID, &order.Status, &order.TotalAmount, &order.Currency,
		&order.PayerEmail, &order.PayerName, &completedAt, &order.Finalized)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "Get order by ID")
	}

	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}

	return &order, nil
}

// GetBySessionID obtiene una orden por ID de sesión
func (r *OrderRepository) GetBySessionID(sessionID string) (*models.Order, error) {
	var order models.Order
	var completedAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, 
		       total_amount, currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE session_id = ?
	`, sessionID).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID,
		&order.PayPalOrderID, &order.Status, &order.TotalAmount, &order.Currency,
		&order.PayerEmail, &order.PayerName, &completedAt, &order.Finalized)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, DBError(err, "Get order by session ID")
	}

	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}

	return &order, nil
}

// Create crea una nueva orden
func (r *OrderRepository) Create(order *models.Order) error {
	result, err := r.db.Exec(`
		INSERT INTO orders (session_id, paypal_order_id, status, total_amount, currency, 
		                    payer_email, payer_name, finalized, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, order.SessionID, order.PayPalOrderID, order.Status, order.TotalAmount, order.Currency,
		order.PayerEmail, order.PayerName, order.Finalized)

	if err != nil {
		return DBError(err, "Create order")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return DBError(err, "Get last insert ID")
	}

	order.ID = id
	return nil
}

// Update actualiza una orden existente
func (r *OrderRepository) Update(id int64, order *models.Order) error {
	var completedAtVal interface{} = nil
	if order.CompletedAt != nil {
		completedAtVal = order.CompletedAt
	}

	result, err := r.db.Exec(`
		UPDATE orders
		SET paypal_order_id = ?, status = ?, total_amount = ?, currency = ?,
		    payer_email = ?, payer_name = ?, completed_at = ?, finalized = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, order.PayPalOrderID, order.Status, order.TotalAmount, order.Currency,
		order.PayerEmail, order.PayerName, completedAtVal, order.Finalized, id)

	if err != nil {
		return DBError(err, "Update order")
	}
	return CheckRowsAffected(result, "Update order")
}

// UpdateStatus actualiza solo el estado de una orden
func (r *OrderRepository) UpdateStatus(id int64, status string) error {
	result, err := r.db.Exec(`
		UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, status, id)

	if err != nil {
		return DBError(err, "Update order status")
	}
	return CheckRowsAffected(result, "Update order status")
}

// Delete elimina una orden
func (r *OrderRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM orders WHERE id = ?`, id)
	if err != nil {
		return DBError(err, "Delete order")
	}
	return CheckRowsAffected(result, "Delete order")
}

// GetByStatus obtiene órdenes por estado
func (r *OrderRepository) GetByStatus(status string) ([]models.Order, error) {
	query := `
		SELECT id, created_at, updated_at, session_id, paypal_order_id, status, 
		       total_amount, currency, payer_email, payer_name, completed_at, finalized
		FROM orders
		WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, DBError(err, "Get orders by status")
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var completedAt sql.NullTime
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.SessionID,
			&order.PayPalOrderID, &order.Status, &order.TotalAmount, &order.Currency,
			&order.PayerEmail, &order.PayerName, &completedAt, &order.Finalized)
		if err != nil {
			continue
		}
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// Count devuelve el número total de órdenes
func (r *OrderRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&count)
	return count, ScanError(err, "count orders")
}

// GetTotalRevenue calcula el ingreso total de órdenes completadas
func (r *OrderRepository) GetTotalRevenue() (float64, error) {
	var revenue sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0) FROM orders 
		WHERE status IN ('COMPLETED', 'DELIVERED')
	`).Scan(&revenue)

	if err != nil {
		return 0, DBError(err, "Get total revenue")
	}

	if revenue.Valid {
		return revenue.Float64, nil
	}
	return 0, nil
}
