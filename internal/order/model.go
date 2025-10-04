package order

import (
	"time"
)

// Order represents a customer order
type Order struct {
	ID            int64       `json:"id" db:"id"`
	UserID        int64       `json:"user_id" db:"user_id"`
	OrderNumber   string      `json:"order_number" db:"order_number"`
	Status        string      `json:"status" db:"status"` // pending, confirmed, shipped, delivered, cancelled
	PaymentStatus string      `json:"payment_status" db:"payment_status"` // pending, paid, failed, refunded
	Subtotal      float64     `json:"subtotal" db:"subtotal"`
	Tax           float64     `json:"tax" db:"tax"`
	ShippingCost  float64     `json:"shipping_cost" db:"shipping_cost"`
	Total         float64     `json:"total" db:"total"`
	ShippingAddress string    `json:"shipping_address" db:"shipping_address"`
	BillingAddress  string    `json:"billing_address" db:"billing_address"`
	Items         []OrderItem `json:"items"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int64     `json:"id" db:"id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Subtotal  float64   `json:"subtotal" db:"subtotal"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateOrderRequest represents creating an order
type CreateOrderRequest struct {
	ShippingAddressID int64  `json:"shipping_address_id" validate:"required"`
	BillingAddressID  int64  `json:"billing_address_id" validate:"required"`
	PaymentMethod     string `json:"payment_method" validate:"required"`
}

// OrderFilter represents filtering options
type OrderFilter struct {
	UserID int64
	Status string
	Limit  int
	Offset int
}
