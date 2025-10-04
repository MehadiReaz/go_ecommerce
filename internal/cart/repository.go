package cart

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

// GetOrCreate gets or creates a cart for a user
func (r *Repository) GetOrCreate(userID int64) (*Cart, error) {
	// Try to get existing cart
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	
	cart := &Cart{}
	err := r.db.QueryRow(query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	
	if err == sql.ErrNoRows {
		// Create new cart
		createQuery := `INSERT INTO carts (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
		err = r.db.QueryRow(createQuery, userID, time.Now(), time.Now()).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
		cart.UserID = userID
		return cart, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	
	return cart, nil
}

// GetItems retrieves all items in a cart
func (r *Repository) GetItems(cartID int64) ([]CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, quantity, price, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
	`
	
	rows, err := r.db.Query(query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()
	
	items := []CartItem{}
	for rows.Next() {
		item := CartItem{}
		err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		item.Subtotal = item.Price * float64(item.Quantity)
		items = append(items, item)
	}
	
	return items, nil
}

// AddItem adds an item to the cart
func (r *Repository) AddItem(cartID, productID int64, quantity int, price float64) error {
	// Check if item already exists
	var existingID int64
	var existingQuantity int
	
	checkQuery := `SELECT id, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	err := r.db.QueryRow(checkQuery, cartID, productID).Scan(&existingID, &existingQuantity)
	
	if err == sql.ErrNoRows {
		// Insert new item
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity, price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = r.db.Exec(insertQuery, cartID, productID, quantity, price, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}
		return nil
	}
	
	if err != nil {
		return fmt.Errorf("failed to check existing item: %w", err)
	}
	
	// Update existing item
	updateQuery := `UPDATE cart_items SET quantity = $1, updated_at = $2 WHERE id = $3`
	_, err = r.db.Exec(updateQuery, existingQuantity+quantity, time.Now(), existingID)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}
	
	return nil
}

// UpdateItem updates a cart item quantity
func (r *Repository) UpdateItem(itemID int64, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, quantity, time.Now(), itemID)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}
	return nil
}

// RemoveItem removes an item from the cart
func (r *Repository) RemoveItem(itemID int64) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Exec(query, itemID)
	if err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	return nil
}

// Clear removes all items from a cart
func (r *Repository) Clear(cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(query, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}
