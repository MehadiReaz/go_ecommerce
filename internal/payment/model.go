package payment

import (
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID              int64     `json:"id" db:"id"`
	OrderID         int64     `json:"order_id" db:"order_id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Currency        string    `json:"currency" db:"currency"`
	PaymentMethod   string    `json:"payment_method" db:"payment_method"` // stripe, bkash
	TransactionID   string    `json:"transaction_id,omitempty" db:"transaction_id"`
	Status          string    `json:"status" db:"status"` // pending, completed, failed, refunded
	PaymentGateway  string    `json:"payment_gateway,omitempty" db:"payment_gateway"`
	GatewayResponse string    `json:"gateway_response,omitempty" db:"gateway_response"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePaymentRequest represents creating a payment
type CreatePaymentRequest struct {
	OrderID       int64  `json:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	Currency      string `json:"currency,omitempty"`
}

// PaymentWebhookPayload represents payment webhook payload
type PaymentWebhookPayload struct {
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	OrderID       int64   `json:"order_id"`
}
