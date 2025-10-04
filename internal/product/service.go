package product

import (
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new product
func (s *Service) Create(req *CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CategoryID:   req.CategoryID,
		SKU:          req.SKU,
		ImageURL:     req.ImageURL,
		IsActive:     true,
		IsFeatured:   req.IsFeatured,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetByID retrieves a product by ID
func (s *Service) GetByID(id int64) (*Product, error) {
	return s.repo.GetByID(id)
}

// List retrieves products with filtering
func (s *Service) List(filter *ProductFilter) ([]*Product, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	return s.repo.List(filter)
}

// Update updates a product
func (s *Service) Update(id int64, req *UpdateProductRequest) (*Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.ComparePrice > 0 {
		product.ComparePrice = req.ComparePrice
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Search searches products
func (s *Service) Search(searchTerm string, limit, offset int) ([]*Product, error) {
	if limit == 0 {
		limit = 20
	}

	if searchTerm == "" {
		return nil, fmt.Errorf("search term is required")
	}

	return s.repo.Search(searchTerm, limit, offset)
}
