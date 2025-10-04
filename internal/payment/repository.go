package payment

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

// Create creates a new payment record
func (r *Repository) Create(payment *Payment) error {
	query := `
		INSERT INTO payments (order_id, user_id, amount, currency, payment_method, transaction_id, status, payment_gateway, gateway_response, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		payment.OrderID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.Status,
		payment.PaymentGateway,
		payment.GatewayResponse,
		time.Now(),
		time.Now(),
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetByID retrieves a payment by ID
func (r *Repository) GetByID(id int64) (*Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, payment_method, transaction_id, status, payment_gateway, gateway_response, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.Status,
		&payment.PaymentGateway,
		&payment.GatewayResponse,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// GetByOrderID retrieves a payment by order ID
func (r *Repository) GetByOrderID(orderID int64) (*Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, payment_method, transaction_id, status, payment_gateway, gateway_response, created_at, updated_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	payment := &Payment{}
	err := r.db.QueryRow(query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.Status,
		&payment.PaymentGateway,
		&payment.GatewayResponse,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// UpdateStatus updates a payment's status
func (r *Repository) UpdateStatus(id int64, status, transactionID, gatewayResponse string) error {
	query := `
		UPDATE payments
		SET status = $1, transaction_id = $2, gateway_response = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, status, transactionID, gatewayResponse, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}
