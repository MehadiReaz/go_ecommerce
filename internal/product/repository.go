package product

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new product
func (r *Repository) Create(product *Product) error {
	query := `
		INSERT INTO products (name, slug, description, price, compare_price, category_id, sku, is_active, is_featured, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.ComparePrice,
		product.CategoryID,
		product.SKU,
		product.IsActive,
		product.IsFeatured,
		product.ImageURL,
		time.Now(),
		time.Now(),
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *Repository) GetByID(id int64) (*Product, error) {
	query := `
		SELECT id, name, slug, description, price, compare_price, category_id, sku, is_active, is_featured, image_url, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_active = true
	`

	product := &Product{}
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.ComparePrice,
		&product.CategoryID,
		&product.SKU,
		&product.IsActive,
		&product.IsFeatured,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// List retrieves products with filtering
func (r *Repository) List(filter *ProductFilter) ([]*Product, error) {
	query := `SELECT id, name, slug, description, price, compare_price, category_id, sku, is_active, is_featured, image_url, created_at, updated_at FROM products WHERE is_active = true`
	args := []interface{}{}
	argPosition := 1

	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argPosition)
		args = append(args, filter.CategoryID)
		argPosition++
	}

	if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argPosition)
		args = append(args, filter.MinPrice)
		argPosition++
	}

	if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argPosition)
		args = append(args, filter.MaxPrice)
		argPosition++
	}

	if filter.IsFeatured != nil {
		query += fmt.Sprintf(" AND is_featured = $%d", argPosition)
		args = append(args, *filter.IsFeatured)
		argPosition++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argPosition, argPosition)
		args = append(args, "%"+filter.Search+"%")
		argPosition++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPosition)
		args = append(args, filter.Limit)
		argPosition++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPosition)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.ComparePrice,
			&product.CategoryID,
			&product.SKU,
			&product.IsActive,
			&product.IsFeatured,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// Update updates a product
func (r *Repository) Update(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, slug = $2, description = $3, price = $4, compare_price = $5, 
		    category_id = $6, sku = $7, is_active = $8, is_featured = $9, image_url = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.Exec(
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.ComparePrice,
		product.CategoryID,
		product.SKU,
		product.IsActive,
		product.IsFeatured,
		product.ImageURL,
		time.Now(),
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// Delete soft deletes a product
func (r *Repository) Delete(id int64) error {
	query := `UPDATE products SET is_active = false, updated_at = $1 WHERE id = $2`

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// Search searches products by name or description
func (r *Repository) Search(searchTerm string, limit, offset int) ([]*Product, error) {
	query := `
		SELECT id, name, slug, description, price, compare_price, category_id, sku, is_active, is_featured, image_url, created_at, updated_at
		FROM products
		WHERE is_active = true AND (name ILIKE $1 OR description ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.ComparePrice,
			&product.CategoryID,
			&product.SKU,
			&product.IsActive,
			&product.IsFeatured,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}
