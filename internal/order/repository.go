package order

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new order
func (r *Repository) Create(order *Order) error {
	query := `
		INSERT INTO orders (user_id, order_number, status, payment_status, subtotal, tax, shipping_cost, total, shipping_address, billing_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.PaymentStatus,
		order.Subtotal,
		order.Tax,
		order.ShippingCost,
		order.Total,
		order.ShippingAddress,
		order.BillingAddress,
		time.Now(),
		time.Now(),
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// CreateItem creates an order item
func (r *Repository) CreateItem(item *OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price, subtotal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
		item.Subtotal,
		time.Now(),
	).Scan(&item.ID, &item.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}

	return nil
}

// GetByID retrieves an order by ID
func (r *Repository) GetByID(id int64) (*Order, error) {
	query := `
		SELECT id, user_id, order_number, status, payment_status, subtotal, tax, shipping_cost, total, shipping_address, billing_address, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &Order{}
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.PaymentStatus,
		&order.Subtotal,
		&order.Tax,
		&order.ShippingCost,
		&order.Total,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get order items
	items, err := r.GetItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

// GetItems retrieves all items for an order
func (r *Repository) GetItems(orderID int64) ([]OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price, subtotal, created_at
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	items := []OrderItem{}
	for rows.Next() {
		item := OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.Subtotal, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// List retrieves orders with filtering
func (r *Repository) List(filter *OrderFilter) ([]*Order, error) {
	query := `SELECT id, user_id, order_number, status, payment_status, subtotal, tax, shipping_cost, total, shipping_address, billing_address, created_at, updated_at FROM orders WHERE 1=1`
	args := []interface{}{}
	argPosition := 1

	if filter.UserID > 0 {
		query += fmt.Sprintf(" AND user_id = $%d", argPosition)
		args = append(args, filter.UserID)
		argPosition++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPosition)
		args = append(args, filter.Status)
		argPosition++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPosition)
		args = append(args, filter.Limit)
		argPosition++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPosition)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.PaymentStatus,
			&order.Subtotal,
			&order.Tax,
			&order.ShippingCost,
			&order.Total,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateStatus updates an order's status
func (r *Repository) UpdateStatus(orderID int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

// UpdatePaymentStatus updates an order's payment status
func (r *Repository) UpdatePaymentStatus(orderID int64, status string) error {
	query := `UPDATE orders SET payment_status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}
