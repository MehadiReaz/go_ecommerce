package category

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description,omitempty" db:"description"`
	ParentID    *int64    `json:"parent_id,omitempty" db:"parent_id"`
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateCategoryRequest represents the create category request
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description,omitempty"`
	ParentID    *int64 `json:"parent_id,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// UpdateCategoryRequest represents the update category request
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    *int64 `json:"parent_id,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
