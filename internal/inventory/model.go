package inventory

import (
	"time"
)

// Inventory represents product inventory
type Inventory struct {
	ID        int64     `json:"id" db:"id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Reserved  int       `json:"reserved" db:"reserved"`
	Available int       `json:"available"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateInventoryRequest represents updating inventory
type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=0"`
}
