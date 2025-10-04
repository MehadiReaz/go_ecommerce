package payment

import (
	"ecommerce_project/internal/config"
	"fmt"
)

type Service struct {
	repo   *Repository
	config *config.PaymentConfig
}

func NewService(repo *Repository, config *config.PaymentConfig) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

// CreatePayment creates a new payment
func (s *Service) CreatePayment(userID int64, req *CreatePaymentRequest) (*Payment, error) {
	// Validate payment method
	if req.PaymentMethod != "stripe" && req.PaymentMethod != "bkash" {
		return nil, fmt.Errorf("invalid payment method")
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	payment := &Payment{
		OrderID:       req.OrderID,
		UserID:        userID,
		Amount:        0, // This should be fetched from order
		Currency:      currency,
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
	}

	// Process payment with gateway
	var transactionID string
	var err error

	switch req.PaymentMethod {
	case "stripe":
		transactionID, err = s.processStripePayment(payment)
	case "bkash":
		transactionID, err = s.processBkashPayment(payment)
	default:
		return nil, fmt.Errorf("unsupported payment method")
	}

	if err != nil {
		payment.Status = "failed"
		payment.GatewayResponse = err.Error()
	} else {
		payment.TransactionID = transactionID
		payment.Status = "completed"
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPayment retrieves a payment
func (s *Service) GetPayment(paymentID, userID int64) (*Payment, error) {
	payment, err := s.repo.GetByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if payment.UserID != userID {
		return nil, fmt.Errorf("payment not found")
	}

	return payment, nil
}

// ProcessWebhook processes payment webhook
func (s *Service) ProcessWebhook(gateway string, payload *PaymentWebhookPayload) error {
	payment, err := s.repo.GetByOrderID(payload.OrderID)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(payment.ID, payload.Status, payload.TransactionID, "Webhook processed")
}

func (s *Service) processStripePayment(payment *Payment) (string, error) {
	// Integration with Stripe
	// This is a placeholder implementation
	if s.config.StripeSecretKey == "" {
		return "", fmt.Errorf("stripe not configured")
	}

	// Simulate successful payment
	return fmt.Sprintf("stripe_txn_%d", payment.OrderID), nil
}

func (s *Service) processBkashPayment(payment *Payment) (string, error) {
	// Integration with bKash
	// This is a placeholder implementation
	if s.config.BkashAppKey == "" {
		return "", fmt.Errorf("bkash not configured")
	}

	// Simulate successful payment
	return fmt.Sprintf("bkash_txn_%d", payment.OrderID), nil
}
