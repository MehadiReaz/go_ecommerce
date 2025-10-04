package category

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new category
func (s *Service) Create(req *CreateCategoryRequest) (*Category, error) {
	category := &Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetByID retrieves a category by ID
func (s *Service) GetByID(id int64) (*Category, error) {
	return s.repo.GetByID(id)
}

// List retrieves all categories
func (s *Service) List() ([]*Category, error) {
	return s.repo.List()
}

// Update updates a category
func (s *Service) Update(id int64, req *UpdateCategoryRequest) (*Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.ImageURL != "" {
		category.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a category
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
