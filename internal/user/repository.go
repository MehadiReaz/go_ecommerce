package user

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

// Create creates a new user
func (r *Repository) Create(user *User) error {
	query := `
		INSERT INTO users (email, password, first_name, last_name, phone_number, role, is_active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Role,
		user.IsActive,
		user.EmailVerified,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *Repository) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, role, is_active, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *Repository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, role, is_active, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user's information
func (r *Repository) Update(user *User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, phone_number = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.PhoneNumber, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *Repository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Delete soft deletes a user
func (r *Repository) Delete(id int64) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// EmailExists checks if an email already exists
func (r *Repository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
