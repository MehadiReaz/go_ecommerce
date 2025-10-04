package inventory

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetByProductID retrieves inventory for a product
func (s *Service) GetByProductID(productID int64) (*Inventory, error) {
	return s.repo.GetByProductID(productID)
}

// List retrieves all inventory records
func (s *Service) List(limit, offset int) ([]*Inventory, error) {
	if limit == 0 {
		limit = 50
	}

	return s.repo.List(limit, offset)
}

// Update updates inventory quantity
func (s *Service) Update(id int64, req *UpdateInventoryRequest) (*Inventory, error) {
	if err := s.repo.Update(id, req.Quantity); err != nil {
		return nil, err
	}

	// Return updated inventory (simplified, should get by ID)
	return &Inventory{
		ID:       id,
		Quantity: req.Quantity,
	}, nil
}

// CheckStock checks if sufficient stock is available
func (s *Service) CheckStock(productID int64, quantity int) (bool, error) {
	return s.repo.CheckStock(productID, quantity)
}

// ReduceStock reduces inventory stock
func (s *Service) ReduceStock(productID int64, quantity int) error {
	return s.repo.ReduceStock(productID, quantity)
}
