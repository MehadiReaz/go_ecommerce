package auth

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

// StoreRefreshToken stores a refresh token
func (r *Repository) StoreRefreshToken(userID int64, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, userID, token, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

// ValidateRefreshToken checks if a refresh token is valid
func (r *Repository) ValidateRefreshToken(token string) (int64, error) {
	query := `
		SELECT user_id, expires_at
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false
	`

	var userID int64
	var expiresAt time.Time

	err := r.db.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("invalid refresh token")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("refresh token expired")
	}

	return userID, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *Repository) RevokeRefreshToken(token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`

	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// CleanExpiredTokens removes expired tokens
func (r *Repository) CleanExpiredTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to clean expired tokens: %w", err)
	}

	return nil
}
