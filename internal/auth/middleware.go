package auth

import (
	"context"
	"net/http"
	"strings"

	"ecommerce_project/pkg/utils"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(service *Service) *Middleware {
	return &Middleware{service: service}
}

// RequireAuth middleware requires valid JWT token
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.service.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin middleware requires admin role
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)

		if role != "admin" {
			utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
