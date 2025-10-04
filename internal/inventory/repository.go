package inventory

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

// GetByProductID retrieves inventory for a product
func (r *Repository) GetByProductID(productID int64) (*Inventory, error) {
	query := `
		SELECT id, product_id, quantity, reserved, updated_at
		FROM inventory
		WHERE product_id = $1
	`

	inventory := &Inventory{}
	err := r.db.QueryRow(query, productID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Reserved,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	inventory.Available = inventory.Quantity - inventory.Reserved
	return inventory, nil
}

// List retrieves all inventory records
func (r *Repository) List(limit, offset int) ([]*Inventory, error) {
	query := `
		SELECT id, product_id, quantity, reserved, updated_at
		FROM inventory
		ORDER BY product_id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	defer rows.Close()

	inventories := []*Inventory{}
	for rows.Next() {
		inventory := &Inventory{}
		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Reserved,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inventory.Available = inventory.Quantity - inventory.Reserved
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

// CheckStock checks if sufficient stock is available
func (r *Repository) CheckStock(productID int64, quantity int) (bool, error) {
	inventory, err := r.GetByProductID(productID)
	if err != nil {
		return false, err
	}

	return inventory.Available >= quantity, nil
}

// ReduceStock reduces inventory stock
func (r *Repository) ReduceStock(productID int64, quantity int) error {
	query := `
		UPDATE inventory
		SET quantity = quantity - $1, updated_at = $2
		WHERE product_id = $3 AND quantity >= $1
	`

	result, err := r.db.Exec(query, quantity, time.Now(), productID)
	if err != nil {
		return fmt.Errorf("failed to reduce stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock")
	}

	return nil
}

// Update updates inventory quantity
func (r *Repository) Update(id int64, quantity int) error {
	query := `
		UPDATE inventory
		SET quantity = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, quantity, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}
