package cart

import (
	"database/sql"
	"fmt"
)

type Service struct {
	repo        *Repository
	productRepo ProductRepository
}

type ProductRepository interface {
	GetByID(id int64) (*Product, error)
}

type Product struct {
	ID    int64
	Price float64
}

func NewService(repo *Repository, productRepo ProductRepository) *Service {
	return &Service{
		repo:        repo,
		productRepo: productRepo,
	}
}

// Get retrieves a user's cart
func (s *Service) Get(userID int64) (*Cart, error) {
	cart, err := s.repo.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}
	
	items, err := s.repo.GetItems(cart.ID)
	if err != nil {
		return nil, err
	}
	
	cart.Items = items
	
	// Calculate total
	total := 0.0
	for _, item := range items {
		total += item.Subtotal
	}
	cart.Total = total
	
	return cart, nil
}

// AddItem adds an item to the cart
func (s *Service) AddItem(userID int64, req *AddItemRequest) error {
	// Get or create cart
	cart, err := s.repo.GetOrCreate(userID)
	if err != nil {
		return err
	}
	
	// Get product to verify existence and get price
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found")
	}
	
	// Add item to cart
	return s.repo.AddItem(cart.ID, req.ProductID, req.Quantity, product.Price)
}

// UpdateItem updates a cart item
func (s *Service) UpdateItem(userID, itemID int64, req *UpdateItemRequest) error {
	// Verify cart ownership (simplified, should verify item belongs to user's cart)
	return s.repo.UpdateItem(itemID, req.Quantity)
}

// RemoveItem removes an item from the cart
func (s *Service) RemoveItem(userID, itemID int64) error {
	// Verify cart ownership (simplified, should verify item belongs to user's cart)
	return s.repo.RemoveItem(itemID)
}

// Clear clears all items from the cart
func (s *Service) Clear(userID int64) error {
	cart, err := s.repo.GetOrCreate(userID)
	if err != nil {
		return err
	}
	
	return s.repo.Clear(cart.ID)
}
