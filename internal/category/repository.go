package category

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

// Create creates a new category
func (r *Repository) Create(category *Category) error {
	query := `
		INSERT INTO categories (name, slug, description, parent_id, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.ImageURL,
		category.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// GetByID retrieves a category by ID
func (r *Repository) GetByID(id int64) (*Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, image_url, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1 AND is_active = true
	`

	category := &Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ParentID,
		&category.ImageURL,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// List retrieves all active categories
func (r *Repository) List() ([]*Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, image_url, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		category := &Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentID,
			&category.ImageURL,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Update updates a category
func (r *Repository) Update(category *Category) error {
	query := `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, parent_id = $4, image_url = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.ImageURL,
		category.IsActive,
		time.Now(),
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// Delete soft deletes a category
func (r *Repository) Delete(id int64) error {
	query := `UPDATE categories SET is_active = false, updated_at = $1 WHERE id = $2`

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
