package product

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	ComparePrice float64  `json:"compare_price,omitempty" db:"compare_price"`
	CategoryID  int64     `json:"category_id" db:"category_id"`
	SKU         string    `json:"sku" db:"sku"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsFeatured  bool      `json:"is_featured" db:"is_featured"`
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateProductRequest represents the create product request
type CreateProductRequest struct {
	Name         string  `json:"name" validate:"required"`
	Slug         string  `json:"slug" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"required,gt=0"`
	ComparePrice float64 `json:"compare_price,omitempty"`
	CategoryID   int64   `json:"category_id" validate:"required"`
	SKU          string  `json:"sku" validate:"required"`
	ImageURL     string  `json:"image_url,omitempty"`
	IsFeatured   bool    `json:"is_featured"`
}

// UpdateProductRequest represents the update product request
type UpdateProductRequest struct {
	Name         string  `json:"name,omitempty"`
	Slug         string  `json:"slug,omitempty"`
	Description  string  `json:"description,omitempty"`
	Price        float64 `json:"price,omitempty"`
	ComparePrice float64 `json:"compare_price,omitempty"`
	CategoryID   int64   `json:"category_id,omitempty"`
	SKU          string  `json:"sku,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
	IsFeatured   *bool   `json:"is_featured,omitempty"`
}

// ProductFilter represents filtering options
type ProductFilter struct {
	CategoryID int64
	MinPrice   float64
	MaxPrice   float64
	IsFeatured *bool
	Search     string
	Limit      int
	Offset     int
}
