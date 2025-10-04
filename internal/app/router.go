package app

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"ecommerce_project/internal/auth"
	"ecommerce_project/internal/cart"
	"ecommerce_project/internal/category"
	"ecommerce_project/internal/config"
	"ecommerce_project/internal/inventory"
	"ecommerce_project/internal/notification"
	"ecommerce_project/internal/order"
	"ecommerce_project/internal/payment"
	"ecommerce_project/internal/product"
	"ecommerce_project/internal/review"
	"ecommerce_project/internal/shipping"
	"ecommerce_project/internal/user"
)

// SetupRouter initializes all routes and dependencies
func SetupRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(LoggingMiddleware)
	router.Use(CORSMiddleware)
	router.Use(RecoveryMiddleware)

	// Initialize repositories
	userRepo := user.NewRepository(db)
	productRepo := product.NewRepository(db)
	categoryRepo := category.NewRepository(db)
	cartRepo := cart.NewRepository(db)
	orderRepo := order.NewRepository(db)
	paymentRepo := payment.NewRepository(db)
	inventoryRepo := inventory.NewRepository(db)
	reviewRepo := review.NewRepository(db)
	shippingRepo := shipping.NewRepository(db)

	// Initialize services
	authService := auth.NewService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	notificationService := notification.NewService(&cfg.Email)
	userService := user.NewService(userRepo, authService)
	productService := product.NewService(productRepo)
	categoryService := category.NewService(categoryRepo)
	cartService := cart.NewService(cartRepo, productRepo)
	orderService := order.NewService(orderRepo, cartRepo, inventoryRepo)
	paymentService := payment.NewService(paymentRepo, &cfg.Payment)
	inventoryService := inventory.NewService(inventoryRepo)
	reviewService := review.NewService(reviewRepo)
	shippingService := shipping.NewService(shippingRepo)

	// Initialize handlers
	userHandler := user.NewHandler(userService)
	productHandler := product.NewHandler(productService)
	categoryHandler := category.NewHandler(categoryService)
	cartHandler := cart.NewHandler(cartService)
	orderHandler := order.NewHandler(orderService)
	paymentHandler := payment.NewHandler(paymentService)
	inventoryHandler := inventory.NewHandler(inventoryService)
	reviewHandler := review.NewHandler(reviewService)
	shippingHandler := shipping.NewHandler(shippingService)

	// Auth middleware
	authMiddleware := auth.NewMiddleware(authService)

	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// Public routes
	api.HandleFunc("/auth/signup", userHandler.Signup).Methods("POST")
	api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", userHandler.RefreshToken).Methods("POST")

	// Product routes (public)
	api.HandleFunc("/products", productHandler.List).Methods("GET")
	api.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")
	api.HandleFunc("/products/search", productHandler.Search).Methods("GET")

	// Category routes (public)
	api.HandleFunc("/categories", categoryHandler.List).Methods("GET")
	api.HandleFunc("/categories/{id}", categoryHandler.GetByID).Methods("GET")

	// Review routes (public read)
	api.HandleFunc("/products/{id}/reviews", reviewHandler.GetProductReviews).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.RequireAuth)

	// User routes
	protected.HandleFunc("/users/me", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/me/password", userHandler.ChangePassword).Methods("PUT")

	// Cart routes
	protected.HandleFunc("/cart", cartHandler.Get).Methods("GET")
	protected.HandleFunc("/cart/items", cartHandler.AddItem).Methods("POST")
	protected.HandleFunc("/cart/items/{id}", cartHandler.UpdateItem).Methods("PUT")
	protected.HandleFunc("/cart/items/{id}", cartHandler.RemoveItem).Methods("DELETE")
	protected.HandleFunc("/cart/clear", cartHandler.Clear).Methods("DELETE")

	// Order routes
	protected.HandleFunc("/orders", orderHandler.List).Methods("GET")
	protected.HandleFunc("/orders", orderHandler.Create).Methods("POST")
	protected.HandleFunc("/orders/{id}", orderHandler.GetByID).Methods("GET")
	protected.HandleFunc("/orders/{id}/cancel", orderHandler.Cancel).Methods("POST")

	// Payment routes
	protected.HandleFunc("/payments", paymentHandler.CreatePayment).Methods("POST")
	protected.HandleFunc("/payments/{id}", paymentHandler.GetPayment).Methods("GET")
	api.HandleFunc("/payments/webhook/stripe", paymentHandler.StripeWebhook).Methods("POST")
	api.HandleFunc("/payments/webhook/bkash", paymentHandler.BkashWebhook).Methods("POST")

	// Review routes (protected write)
	protected.HandleFunc("/reviews", reviewHandler.Create).Methods("POST")
	protected.HandleFunc("/reviews/{id}", reviewHandler.Update).Methods("PUT")
	protected.HandleFunc("/reviews/{id}", reviewHandler.Delete).Methods("DELETE")

	// Shipping routes
	protected.HandleFunc("/shipping/addresses", shippingHandler.ListAddresses).Methods("GET")
	protected.HandleFunc("/shipping/addresses", shippingHandler.CreateAddress).Methods("POST")
	protected.HandleFunc("/shipping/addresses/{id}", shippingHandler.UpdateAddress).Methods("PUT")
	protected.HandleFunc("/shipping/addresses/{id}", shippingHandler.DeleteAddress).Methods("DELETE")

	// Admin routes (add admin middleware later)
	admin := protected.PathPrefix("/admin").Subrouter()
	// admin.Use(authMiddleware.RequireAdmin) // Implement admin check

	admin.HandleFunc("/products", productHandler.Create).Methods("POST")
	admin.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
	admin.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")

	admin.HandleFunc("/categories", categoryHandler.Create).Methods("POST")
	admin.HandleFunc("/categories/{id}", categoryHandler.Update).Methods("PUT")
	admin.HandleFunc("/categories/{id}", categoryHandler.Delete).Methods("DELETE")

	admin.HandleFunc("/inventory", inventoryHandler.List).Methods("GET")
	admin.HandleFunc("/inventory/{id}", inventoryHandler.Update).Methods("PUT")

	return router
}

// HealthCheckHandler returns the health status of the API
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","message":"E-Commerce API is running"}`))
}
