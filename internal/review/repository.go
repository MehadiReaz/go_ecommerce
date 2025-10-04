package review

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

// Create creates a new review
func (r *Repository) Create(review *Review) error {
	query := `
		INSERT INTO reviews (product_id, user_id, rating, title, comment, verified, helpful, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		review.ProductID,
		review.UserID,
		review.Rating,
		review.Title,
		review.Comment,
		review.Verified,
		review.Helpful,
		time.Now(),
		time.Now(),
	).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// GetByID retrieves a review by ID
func (r *Repository) GetByID(id int64) (*Review, error) {
	query := `
		SELECT id, product_id, user_id, rating, title, comment, verified, helpful, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`

	review := &Review{}
	err := r.db.QueryRow(query, id).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Rating,
		&review.Title,
		&review.Comment,
		&review.Verified,
		&review.Helpful,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return review, nil
}

// GetByProductID retrieves reviews for a product
func (r *Repository) GetByProductID(productID int64, limit, offset int) ([]*Review, error) {
	query := `
		SELECT id, product_id, user_id, rating, title, comment, verified, helpful, created_at, updated_at
		FROM reviews
		WHERE product_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	reviews := []*Review{}
	for rows.Next() {
		review := &Review{}
		err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.UserID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.Verified,
			&review.Helpful,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Update updates a review
func (r *Repository) Update(review *Review) error {
	query := `
		UPDATE reviews
		SET rating = $1, title = $2, comment = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, review.Rating, review.Title, review.Comment, time.Now(), review.ID)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

// Delete deletes a review
func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM reviews WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

// UserHasReviewed checks if user has already reviewed a product
func (r *Repository) UserHasReviewed(userID, productID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND product_id = $2)`

	var exists bool
	err := r.db.QueryRow(query, userID, productID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check review existence: %w", err)
	}

	return exists, nil
}
