package user

import (
	"fmt"

	"ecommerce_project/internal/auth"
	"ecommerce_project/pkg/utils"
)

type Service struct {
	repo        *Repository
	authService *auth.Service
}

func NewService(repo *Repository, authService *auth.Service) *Service {
	return &Service{
		repo:        repo,
		authService: authService,
	}
}

// Signup creates a new user account
func (s *Service) Signup(req *SignupRequest) (*User, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		Email:         req.Email,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		PhoneNumber:   req.PhoneNumber,
		Role:          "customer",
		IsActive:      true,
		EmailVerified: false,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	token, err := s.authService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := s.authService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password before returning
	user.Password = ""

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// GetProfile retrieves a user's profile
func (s *Service) GetProfile(userID int64) (*User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}

// UpdateProfile updates a user's profile
func (s *Service) UpdateProfile(userID int64, req *UpdateProfileRequest) (*User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(userID int64, req *ChangePasswordRequest) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(userID, hashedPassword)
}

// RefreshToken generates a new access token from refresh token
func (s *Service) RefreshToken(refreshToken string) (string, error) {
	userID, err := s.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		return "", err
	}

	token, err := s.authService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
