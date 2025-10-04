package shipping

import (
	"time"
)

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	FullName    string    `json:"full_name" db:"full_name"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	AddressLine1 string   `json:"address_line1" db:"address_line1"`
	AddressLine2 string   `json:"address_line2,omitempty" db:"address_line2"`
	City        string    `json:"city" db:"city"`
	State       string    `json:"state" db:"state"`
	PostalCode  string    `json:"postal_code" db:"postal_code"`
	Country     string    `json:"country" db:"country"`
	IsDefault   bool      `json:"is_default" db:"is_default"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAddressRequest represents creating an address
type CreateAddressRequest struct {
	FullName     string `json:"full_name" validate:"required"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	AddressLine1 string `json:"address_line1" validate:"required"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	PostalCode   string `json:"postal_code" validate:"required"`
	Country      string `json:"country" validate:"required"`
	IsDefault    bool   `json:"is_default"`
}

// UpdateAddressRequest represents updating an address
type UpdateAddressRequest struct {
	FullName     string `json:"full_name,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	IsDefault    *bool  `json:"is_default,omitempty"`
}
