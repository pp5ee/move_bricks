package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	Username     string
	Password     string
	From         string
	To           string
	UseTLS       bool
}

// DefaultEmailConfig creates email config from environment variables
func DefaultEmailConfig() *EmailConfig {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}

	return &EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
		To:       os.Getenv("EMAIL_TO"),
		UseTLS:   os.Getenv("SMTP_USE_TLS") != "false",
	}
}

// Validate checks if the email config is valid
func (c *EmailConfig) Validate() error {
	if c.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.Username == "" {
		return fmt.Errorf("SMTP_USERNAME is required")
	}
	if c.Password == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}
	if c.From == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}
	if c.To == "" {
		return fmt.Errorf("EMAIL_TO is required")
	}
	return nil
}

// EmailSender handles sending emails
type EmailSender struct {
	config *EmailConfig
}

// NewEmailSender creates a new email sender
func NewEmailSender(config *EmailConfig) *EmailSender {
	if config == nil {
		config = DefaultEmailConfig()
	}
	return &EmailSender{config: config}
}

// Send sends an email
func (s *EmailSender) Send(subject, body string) error {
	if err := s.config.Validate(); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", s.config.To)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.Username,
		s.config.Password,
	)

	dialer.TLSConfig.InsecureSkipVerify = true

	return dialer.SendMail(m)
}

// SendAlert sends a price alert email
func (s *EmailSender) SendAlert(symbol string, price, threshold float64, alertType string) error {
	subject := fmt.Sprintf("[%s] Price Alert: %s", alertType, strings.ToUpper(symbol))

	body := fmt.Sprintf(`
		<h2>Price Alert</h2>
		<p><strong>Symbol:</strong> %s</p>
		<p><strong>Current Price:</strong> %.8f</p>
		<p><strong>Threshold:</strong> %.8f</p>
		<p><strong>Alert Type:</strong> %s</p>
		<p><strong>Time:</strong> %s</p>
	`, strings.ToUpper(symbol), price, threshold, alertType, s.getCurrentTime())

	return s.Send(subject, body)
}

// SendTradeNotification sends a trade notification email
func (s *EmailSender) SendTradeNotification(symbol string, side string, amount, price float64, orderID string) error {
	subject := fmt.Sprintf("[Trade %s] %s %s", side, strings.ToUpper(symbol), s.getCurrentTime())

	body := fmt.Sprintf(`
		<h2>Trade Notification</h2>
		<p><strong>Symbol:</strong> %s</p>
		<p><strong>Side:</strong> %s</p>
		<p><strong>Amount:</strong> %.8f</p>
		<p><strong>Price:</strong> %.8f</p>
		<p><strong>Order ID:</strong> %s</p>
		<p><strong>Time:</strong> %s</p>
	`, strings.ToUpper(symbol), strings.ToUpper(side), amount, price, orderID, s.getCurrentTime())

	return s.Send(subject, body)
}

// SendErrorNotification sends an error notification email
func (s *EmailSender) SendErrorNotification(err error, context string) error {
	subject := fmt.Sprintf("[Error] %s", s.getCurrentTime())

	body := fmt.Sprintf(`
		<h2>Error Notification</h2>
		<p><strong>Error:</strong> %s</p>
		<p><strong>Context:</strong> %s</p>
		<p><strong>Time:</strong> %s</p>
	`, err.Error(), context, s.getCurrentTime())

	return s.Send(subject, body)
}

// SendBalanceNotification sends a balance notification email
func (s *EmailSender) SendBalanceNotification(currency string, balance float64) error {
	subject := fmt.Sprintf("[Balance Update] %s", strings.ToUpper(currency))

	body := fmt.Sprintf(`
		<h2>Balance Update</h2>
		<p><strong>Currency:</strong> %s</p>
		<p><strong>Balance:</strong> %.8f</p>
		<p><strong>Time:</strong> %s</p>
	`, strings.ToUpper(currency), balance, s.getCurrentTime())

	return s.Send(subject, body)
}

func (s *EmailSender) getCurrentTime() string {
	return "Now"
}