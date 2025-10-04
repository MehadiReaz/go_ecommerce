package shipping

import (
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateAddress creates a new shipping address
func (s *Service) CreateAddress(userID int64, req *CreateAddressRequest) (*ShippingAddress, error) {
	address := &ShippingAddress{
		UserID:       userID,
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		IsDefault:    req.IsDefault,
	}

	if err := s.repo.Create(address); err != nil {
		return nil, err
	}

	return address, nil
}

// ListAddresses retrieves all addresses for a user
func (s *Service) ListAddresses(userID int64) ([]*ShippingAddress, error) {
	return s.repo.ListByUserID(userID)
}

// GetAddress retrieves an address
func (s *Service) GetAddress(userID, addressID int64) (*ShippingAddress, error) {
	address, err := s.repo.GetByID(addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, fmt.Errorf("address not found")
	}

	return address, nil
}

// UpdateAddress updates an address
func (s *Service) UpdateAddress(userID, addressID int64, req *UpdateAddressRequest) (*ShippingAddress, error) {
	address, err := s.repo.GetByID(addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, fmt.Errorf("address not found")
	}

	// Update fields
	if req.FullName != "" {
		address.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		address.PhoneNumber = req.PhoneNumber
	}
	if req.AddressLine1 != "" {
		address.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		address.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.State != "" {
		address.State = req.State
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		address.Country = req.Country
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	if err := s.repo.Update(address); err != nil {
		return nil, err
	}

	return address, nil
}

// DeleteAddress deletes an address
func (s *Service) DeleteAddress(userID, addressID int64) error {
	address, err := s.repo.GetByID(addressID)
	if err != nil {
		return err
	}

	// Verify ownership
	if address.UserID != userID {
		return fmt.Errorf("address not found")
	}

	return s.repo.Delete(addressID)
}
