package email

import (
	"fmt"
	"net/smtp"

	"ecommerce_project/internal/config"
)

type SMTPClient struct {
	config *config.EmailConfig
}

// NewSMTPClient creates a new SMTP client
func NewSMTPClient(config *config.EmailConfig) *SMTPClient {
	return &SMTPClient{config: config}
}

// SendEmail sends an email via SMTP
func (c *SMTPClient) SendEmail(to, subject, body string) error {
	from := c.config.FromEmail
	password := c.config.SMTPPassword

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", c.config.FromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Compose message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Authentication
	auth := smtp.PlainAuth("", from, password, c.config.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", c.config.SMTPHost, c.config.SMTPPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendBulkEmail sends emails to multiple recipients
func (c *SMTPClient) SendBulkEmail(recipients []string, subject, body string) error {
	for _, to := range recipients {
		if err := c.SendEmail(to, subject, body); err != nil {
			return err
		}
	}
	return nil
}
