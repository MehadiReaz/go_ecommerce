package notification

import (
	"fmt"
)

// generateOrderConfirmationEmail generates order confirmation email body
func generateOrderConfirmationEmail(orderNumber string, amount float64) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Order Confirmation</h2>
			<p>Thank you for your order!</p>
			<p><strong>Order Number:</strong> %s</p>
			<p><strong>Total Amount:</strong> $%.2f</p>
			<p>We'll send you a shipping confirmation email as soon as your order ships.</p>
		</body>
		</html>
	`, orderNumber, amount)
}

// generatePasswordResetEmail generates password reset email body
func generatePasswordResetEmail(resetToken string) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>You requested a password reset for your account.</p>
			<p>Click the link below to reset your password:</p>
			<a href="https://example.com/reset-password?token=%s">Reset Password</a>
			<p>If you didn't request this, please ignore this email.</p>
			<p>This link will expire in 24 hours.</p>
		</body>
		</html>
	`, resetToken)
}

// generateWelcomeEmail generates welcome email body
func generateWelcomeEmail(name string) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to E-Commerce!</h2>
			<p>Hi %s,</p>
			<p>Thank you for signing up. We're excited to have you!</p>
			<p>Start shopping now and enjoy exclusive deals.</p>
			<a href="https://example.com/shop">Start Shopping</a>
		</body>
		</html>
	`, name)
}
