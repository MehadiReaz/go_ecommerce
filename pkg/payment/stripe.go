package payment

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type StripeClient struct {
	secretKey string
}

// NewStripeClient creates a new Stripe client
func NewStripeClient(secretKey string) *StripeClient {
	stripe.Key = secretKey
	return &StripeClient{secretKey: secretKey}
}

// CreatePaymentIntent creates a Stripe payment intent
func (c *StripeClient) CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

// GetPaymentIntent retrieves a Stripe payment intent
func (c *StripeClient) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

// ConfirmPaymentIntent confirms a Stripe payment intent
func (c *StripeClient) ConfirmPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Confirm(id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return pi, nil
}

// CancelPaymentIntent cancels a Stripe payment intent
func (c *StripeClient) CancelPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Cancel(id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	return pi, nil
}
