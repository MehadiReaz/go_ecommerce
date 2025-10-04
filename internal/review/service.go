package review

import (
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new review
func (s *Service) Create(userID int64, req *CreateReviewRequest) (*Review, error) {
	// Check if user has already reviewed this product
	hasReviewed, err := s.repo.UserHasReviewed(userID, req.ProductID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, fmt.Errorf("you have already reviewed this product")
	}

	review := &Review{
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Title:     req.Title,
		Comment:   req.Comment,
		Verified:  false, // Should check if user purchased this product
		Helpful:   0,
	}

	if err := s.repo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

// GetProductReviews retrieves reviews for a product
func (s *Service) GetProductReviews(productID int64, limit, offset int) ([]*Review, error) {
	if limit == 0 {
		limit = 20
	}

	return s.repo.GetByProductID(productID, limit, offset)
}

// Update updates a review
func (s *Service) Update(userID, reviewID int64, req *UpdateReviewRequest) (*Review, error) {
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if review.UserID != userID {
		return nil, fmt.Errorf("review not found")
	}

	// Update fields
	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Title != "" {
		review.Title = req.Title
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

// Delete deletes a review
func (s *Service) Delete(userID, reviewID int64) error {
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return err
	}

	// Verify ownership
	if review.UserID != userID {
		return fmt.Errorf("review not found")
	}

	return s.repo.Delete(reviewID)
}
