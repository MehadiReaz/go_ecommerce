package cart

import (
	"time"
)

// Cart represents a shopping cart
type Cart struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID        int64     `json:"id" db:"id"`
	CartID    int64     `json:"cart_id" db:"cart_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AddItemRequest represents adding an item to cart
type AddItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

// UpdateItemRequest represents updating a cart item
type UpdateItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}
