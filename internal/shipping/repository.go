package shipping

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

// Create creates a new shipping address
func (r *Repository) Create(address *ShippingAddress) error {
	// If this is the default address, unset any existing default
	if address.IsDefault {
		if err := r.UnsetDefault(address.UserID); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO shipping_addresses (user_id, full_name, phone_number, address_line1, address_line2, city, state, postal_code, country, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		address.UserID,
		address.FullName,
		address.PhoneNumber,
		address.AddressLine1,
		address.AddressLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		time.Now(),
		time.Now(),
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create address: %w", err)
	}

	return nil
}

// GetByID retrieves an address by ID
func (r *Repository) GetByID(id int64) (*ShippingAddress, error) {
	query := `
		SELECT id, user_id, full_name, phone_number, address_line1, address_line2, city, state, postal_code, country, is_default, created_at, updated_at
		FROM shipping_addresses
		WHERE id = $1
	`

	address := &ShippingAddress{}
	err := r.db.QueryRow(query, id).Scan(
		&address.ID,
		&address.UserID,
		&address.FullName,
		&address.PhoneNumber,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.City,
		&address.State,
		&address.PostalCode,
		&address.Country,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("address not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return address, nil
}

// ListByUserID retrieves all addresses for a user
func (r *Repository) ListByUserID(userID int64) ([]*ShippingAddress, error) {
	query := `
		SELECT id, user_id, full_name, phone_number, address_line1, address_line2, city, state, postal_code, country, is_default, created_at, updated_at
		FROM shipping_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}
	defer rows.Close()

	addresses := []*ShippingAddress{}
	for rows.Next() {
		address := &ShippingAddress{}
		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.FullName,
			&address.PhoneNumber,
			&address.AddressLine1,
			&address.AddressLine2,
			&address.City,
			&address.State,
			&address.PostalCode,
			&address.Country,
			&address.IsDefault,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Update updates an address
func (r *Repository) Update(address *ShippingAddress) error {
	// If this is the default address, unset any existing default
	if address.IsDefault {
		if err := r.UnsetDefault(address.UserID); err != nil {
			return err
		}
	}

	query := `
		UPDATE shipping_addresses
		SET full_name = $1, phone_number = $2, address_line1 = $3, address_line2 = $4, 
		    city = $5, state = $6, postal_code = $7, country = $8, is_default = $9, updated_at = $10
		WHERE id = $11
	`

	_, err := r.db.Exec(
		query,
		address.FullName,
		address.PhoneNumber,
		address.AddressLine1,
		address.AddressLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		time.Now(),
		address.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	return nil
}

// Delete deletes an address
func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM shipping_addresses WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	return nil
}

// UnsetDefault removes default flag from all addresses for a user
func (r *Repository) UnsetDefault(userID int64) error {
	query := `UPDATE shipping_addresses SET is_default = false WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to unset default: %w", err)
	}

	return nil
}
