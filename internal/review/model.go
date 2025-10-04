package review

import (
	"time"
)

// Review represents a product review
type Review struct {
	ID        int64     `json:"id" db:"id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"` // 1-5 stars
	Title     string    `json:"title,omitempty" db:"title"`
	Comment   string    `json:"comment,omitempty" db:"comment"`
	Verified  bool      `json:"verified" db:"verified"` // verified purchase
	Helpful   int       `json:"helpful" db:"helpful"`   // helpful count
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateReviewRequest represents creating a review
type CreateReviewRequest struct {
	ProductID int64  `json:"product_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Title     string `json:"title,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// UpdateReviewRequest represents updating a review
type UpdateReviewRequest struct {
	Rating  int    `json:"rating,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
}
