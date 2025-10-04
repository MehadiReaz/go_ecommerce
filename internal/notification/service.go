package notification

import (
	"ecommerce_project/internal/config"
	"ecommerce_project/pkg/logger"
)

type Service struct {
	emailConfig *config.EmailConfig
}

func NewService(emailConfig *config.EmailConfig) *Service {
	return &Service{
		emailConfig: emailConfig,
	}
}

// SendEmail sends an email notification
func (s *Service) SendEmail(to, subject, body string) error {
	logger.Info("Sending email",
		"to", to,
		"subject", subject,
	)

	// Placeholder for actual email sending implementation
	// In production, use SMTP or a service like SendGrid

	return nil
}

// SendOrderConfirmation sends order confirmation email
func (s *Service) SendOrderConfirmation(to, orderNumber string, amount float64) error {
	subject := "Order Confirmation"
	body := generateOrderConfirmationEmail(orderNumber, amount)
	return s.SendEmail(to, subject, body)
}

// SendPasswordReset sends password reset email
func (s *Service) SendPasswordReset(to, resetToken string) error {
	subject := "Password Reset Request"
	body := generatePasswordResetEmail(resetToken)
	return s.SendEmail(to, subject, body)
}

// SendWelcome sends welcome email
func (s *Service) SendWelcome(to, name string) error {
	subject := "Welcome to E-Commerce"
	body := generateWelcomeEmail(name)
	return s.SendEmail(to, subject, body)
}

// SendSMS sends an SMS notification (placeholder)
func (s *Service) SendSMS(phone, message string) error {
	logger.Info("Sending SMS",
		"phone", phone,
		"message", message,
	)

	// Placeholder for actual SMS sending implementation
	// In production, use a service like Twilio

	return nil
}

// SendPushNotification sends a push notification (placeholder)
func (s *Service) SendPushNotification(userID int64, title, message string) error {
	logger.Info("Sending push notification",
		"user_id", userID,
		"title", title,
		"message", message,
	)

	// Placeholder for actual push notification implementation
	// In production, use a service like Firebase Cloud Messaging

	return nil
}
