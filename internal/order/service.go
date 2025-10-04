package order

import (
	"fmt"
	"time"
)

type Service struct {
	repo          *Repository
	cartRepo      CartRepository
	inventoryRepo InventoryRepository
}

type CartRepository interface {
	GetOrCreate(userID int64) (*Cart, error)
	GetItems(cartID int64) ([]CartItem, error)
	Clear(cartID int64) error
}

type Cart struct {
	ID    int64
	Items []CartItem
}

type CartItem struct {
	ProductID int64
	Quantity  int
	Price     float64
}

type InventoryRepository interface {
	CheckStock(productID int64, quantity int) (bool, error)
	ReduceStock(productID int64, quantity int) error
}

func NewService(repo *Repository, cartRepo CartRepository, inventoryRepo InventoryRepository) *Service {
	return &Service{
		repo:          repo,
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
	}
}

// Create creates a new order from cart
func (s *Service) Create(userID int64, req *CreateOrderRequest) (*Order, error) {
	// Get cart
	cart, err := s.cartRepo.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}

	items, err := s.cartRepo.GetItems(cart.ID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Check inventory
	for _, item := range items {
		hasStock, err := s.inventoryRepo.CheckStock(item.ProductID, item.Quantity)
		if err != nil || !hasStock {
			return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}
	}

	// Calculate totals
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.Price * float64(item.Quantity)
	}
	tax := subtotal * 0.1 // 10% tax
	shippingCost := 10.0
	total := subtotal + tax + shippingCost

	// Create order
	order := &Order{
		UserID:          userID,
		OrderNumber:     generateOrderNumber(),
		Status:          "pending",
		PaymentStatus:   "pending",
		Subtotal:        subtotal,
		Tax:             tax,
		ShippingCost:    shippingCost,
		Total:           total,
		ShippingAddress: fmt.Sprintf("Address ID: %d", req.ShippingAddressID),
		BillingAddress:  fmt.Sprintf("Address ID: %d", req.BillingAddressID),
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Create order items and reduce inventory
	for _, item := range items {
		orderItem := &OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Price * float64(item.Quantity),
		}

		if err := s.repo.CreateItem(orderItem); err != nil {
			return nil, err
		}

		// Reduce inventory
		if err := s.inventoryRepo.ReduceStock(item.ProductID, item.Quantity); err != nil {
			return nil, err
		}
	}

	// Clear cart
	if err := s.cartRepo.Clear(cart.ID); err != nil {
		return nil, err
	}

	// Get full order with items
	return s.repo.GetByID(order.ID)
}

// GetByID retrieves an order by ID
func (s *Service) GetByID(orderID, userID int64) (*Order, error) {
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// List retrieves user's orders
func (s *Service) List(userID int64, limit, offset int) ([]*Order, error) {
	filter := &OrderFilter{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	if filter.Limit == 0 {
		filter.Limit = 20
	}

	return s.repo.List(filter)
}

// Cancel cancels an order
func (s *Service) Cancel(orderID, userID int64) error {
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Verify ownership
	if order.UserID != userID {
		return fmt.Errorf("order not found")
	}

	// Check if order can be cancelled
	if order.Status == "shipped" || order.Status == "delivered" || order.Status == "cancelled" {
		return fmt.Errorf("order cannot be cancelled")
	}

	return s.repo.UpdateStatus(orderID, "cancelled")
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().Unix())
}
