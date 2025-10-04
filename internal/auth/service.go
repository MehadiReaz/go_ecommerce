package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret      string
	expiryHours int
}

func NewService(secret string, expiryHours int) *Service {
	return &Service{
		secret:      secret,
		expiryHours: expiryHours,
	}
}

// GenerateToken generates a new JWT access token
func (s *Service) GenerateToken(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.expiryHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *Service) GenerateRefreshToken(userID int64) (string, error) {
	// Generate a UUID as refresh token
	refreshToken := uuid.New().String()
	return refreshToken, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID := int64(claims["user_id"].(float64))
	email := claims["email"].(string)
	role := claims["role"].(string)

	return &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

// ValidateRefreshToken validates a refresh token
func (s *Service) ValidateRefreshToken(token string) (int64, error) {
	// This should be implemented with database validation
	// For now, return an error indicating it needs DB implementation
	return 0, fmt.Errorf("refresh token validation requires database")
}
