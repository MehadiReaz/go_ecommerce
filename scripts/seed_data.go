package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Database connection
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Seed data
	seedUsers(db)
	seedCategories(db)
	seedProducts(db)
	seedInventory(db)

	log.Println("Seed data inserted successfully!")
}

func seedUsers(db *sql.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []struct {
		Email     string
		FirstName string
		LastName  string
		Role      string
	}{
		{"admin@ecommerce.com", "Admin", "User", "admin"},
		{"john@example.com", "John", "Doe", "customer"},
		{"jane@example.com", "Jane", "Smith", "customer"},
	}

	for _, u := range users {
		_, err := db.Exec(`
			INSERT INTO users (email, password, first_name, last_name, role, is_active, email_verified)
			VALUES ($1, $2, $3, $4, $5, true, true)
			ON CONFLICT (email) DO NOTHING
		`, u.Email, string(password), u.FirstName, u.LastName, u.Role)

		if err != nil {
			log.Printf("Error inserting user %s: %v", u.Email, err)
		}
	}

	log.Println("Users seeded")
}

func seedCategories(db *sql.DB) {
	categories := []struct {
		Name string
		Slug string
		Desc string
	}{
		{"Electronics", "electronics", "Electronic devices and accessories"},
		{"Clothing", "clothing", "Fashion and apparel"},
		{"Books", "books", "Books and reading materials"},
		{"Home & Garden", "home-garden", "Home and garden products"},
		{"Sports", "sports", "Sports and fitness equipment"},
	}

	for _, c := range categories {
		_, err := db.Exec(`
			INSERT INTO categories (name, slug, description, is_active)
			VALUES ($1, $2, $3, true)
			ON CONFLICT (slug) DO NOTHING
		`, c.Name, c.Slug, c.Desc)

		if err != nil {
			log.Printf("Error inserting category %s: %v", c.Name, err)
		}
	}

	log.Println("Categories seeded")
}

func seedProducts(db *sql.DB) {
	products := []struct {
		Name     string
		Slug     string
		Desc     string
		Price    float64
		Category int64
		SKU      string
	}{
		{"Laptop", "laptop", "High-performance laptop", 999.99, 1, "LAP001"},
		{"Smartphone", "smartphone", "Latest smartphone", 699.99, 1, "PHN001"},
		{"T-Shirt", "t-shirt", "Cotton t-shirt", 19.99, 2, "TSH001"},
		{"Jeans", "jeans", "Denim jeans", 49.99, 2, "JNS001"},
		{"Programming Book", "programming-book", "Learn programming", 39.99, 3, "BK001"},
		{"Coffee Maker", "coffee-maker", "Automatic coffee maker", 79.99, 4, "HOM001"},
		{"Yoga Mat", "yoga-mat", "Exercise yoga mat", 29.99, 5, "SPT001"},
	}

	for _, p := range products {
		_, err := db.Exec(`
			INSERT INTO products (name, slug, description, price, category_id, sku, is_active, is_featured)
			VALUES ($1, $2, $3, $4, $5, $6, true, $7)
			ON CONFLICT (slug) DO NOTHING
		`, p.Name, p.Slug, p.Desc, p.Price, p.Category, p.SKU, rand.Intn(2) == 0)

		if err != nil {
			log.Printf("Error inserting product %s: %v", p.Name, err)
		}
	}

	log.Println("Products seeded")
}

func seedInventory(db *sql.DB) {
	// Add inventory for all products
	_, err := db.Exec(`
		INSERT INTO inventory (product_id, quantity, reserved)
		SELECT id, $1, 0 FROM products
		ON CONFLICT (product_id) DO NOTHING
	`, rand.Intn(100)+10)

	if err != nil {
		log.Printf("Error inserting inventory: %v", err)
	}

	log.Println("Inventory seeded")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
