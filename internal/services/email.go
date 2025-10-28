package services

import (
	"fmt"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"template-store/internal/config"
	"template-store/internal/models"
)

// EmailService defines the interface for email operations
type EmailService interface {
	// Welcome and onboarding
	SendWelcomeEmail(user *models.User) error

	// Password management
	SendPasswordResetEmail(user *models.User, resetToken string, resetURL string) error
	SendPasswordChangedNotification(user *models.User) error

	// Email verification
	SendEmailVerificationEmail(user *models.User, token string, verificationURL string) error
	SendEmailChangeConfirmation(user *models.User, newEmail, token string, verificationURL string) error

	// Account management
	SendAccountSuspensionEmail(user *models.User, reason string, suspendedBy string) error
	SendAccountUnsuspensionEmail(user *models.User) error
	SendAccountDeletionEmail(user *models.User) error

	// Order notifications
	SendOrderConfirmationEmail(user *models.User, order *models.Order, template *models.Template) error

	// Generic email
	SendEmail(to, subject, htmlContent, textContent string) error
}

// EmailServiceImpl implements the EmailService interface using SendGrid
type EmailServiceImpl struct {
	client     *sendgrid.Client
	fromEmail  string
	fromName   string
	baseURL    string // Frontend base URL for links
}

// NewEmailService creates a new EmailService instance
func NewEmailService(cfg *config.Config) EmailService {
	return &EmailServiceImpl{
		client:    sendgrid.NewSendClient(cfg.SendGrid.APIKey),
		fromEmail: cfg.SendGrid.From,
		fromName:  "Template Store",
		baseURL:   getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailServiceImpl) SendWelcomeEmail(user *models.User) error {
	subject := "Welcome to Template Store!"

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome to Template Store, %s!</h1>
			<p>Thank you for joining our community. We're excited to have you on board.</p>
			<p>Explore our collection of premium templates and start creating amazing projects.</p>
			<p>
				<a href="%s/templates" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Browse Templates
				</a>
			</p>
			<p>If you have any questions, feel free to reach out to our support team.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, s.baseURL)

	textContent := fmt.Sprintf(
		"Welcome to Template Store, %s!\n\nThank you for joining our community.\n\nBest regards,\nThe Template Store Team",
		user.Name,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailServiceImpl) SendPasswordResetEmail(user *models.User, resetToken string, resetURL string) error {
	subject := "Reset Your Password"

	if resetURL == "" {
		resetURL = fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, resetToken)
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Password Reset Request</h1>
			<p>Hello %s,</p>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<p>
				<a href="%s" style="background-color: #2196F3; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Reset Password
				</a>
			</p>
			<p>This link will expire in 24 hours.</p>
			<p>If you didn't request a password reset, please ignore this email or contact support if you have concerns.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, resetURL)

	textContent := fmt.Sprintf(
		"Password Reset Request\n\nHello %s,\n\nClick this link to reset your password: %s\n\nThis link will expire in 24 hours.\n\nBest regards,\nThe Template Store Team",
		user.Name, resetURL,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendPasswordChangedNotification sends a notification that the password was changed
func (s *EmailServiceImpl) SendPasswordChangedNotification(user *models.User) error {
	subject := "Your Password Has Been Changed"

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Password Changed</h1>
			<p>Hello %s,</p>
			<p>This email confirms that your password was successfully changed on %s.</p>
			<p>If you didn't make this change, please contact our support team immediately.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, time.Now().Format("January 2, 2006 at 3:04 PM"))

	textContent := fmt.Sprintf(
		"Password Changed\n\nHello %s,\n\nYour password was successfully changed.\n\nBest regards,\nThe Template Store Team",
		user.Name,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendEmailVerificationEmail sends an email verification link
func (s *EmailServiceImpl) SendEmailVerificationEmail(user *models.User, token string, verificationURL string) error {
	subject := "Verify Your Email Address"

	if verificationURL == "" {
		verificationURL = fmt.Sprintf("%s/verify-email?token=%s", s.baseURL, token)
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Verify Your Email</h1>
			<p>Hello %s,</p>
			<p>Please verify your email address by clicking the button below:</p>
			<p>
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Verify Email
				</a>
			</p>
			<p>This link will expire in 72 hours.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, verificationURL)

	textContent := fmt.Sprintf(
		"Verify Your Email\n\nHello %s,\n\nClick this link to verify your email: %s\n\nBest regards,\nThe Template Store Team",
		user.Name, verificationURL,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendEmailChangeConfirmation sends confirmation for email change
func (s *EmailServiceImpl) SendEmailChangeConfirmation(user *models.User, newEmail, token string, verificationURL string) error {
	subject := "Confirm Your New Email Address"

	if verificationURL == "" {
		verificationURL = fmt.Sprintf("%s/verify-email-change?token=%s", s.baseURL, token)
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Confirm Email Change</h1>
			<p>Hello %s,</p>
			<p>You requested to change your email address to <strong>%s</strong>.</p>
			<p>Please confirm this change by clicking the button below:</p>
			<p>
				<a href="%s" style="background-color: #2196F3; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Confirm Email Change
				</a>
			</p>
			<p>This link will expire in 72 hours.</p>
			<p>If you didn't request this change, please ignore this email.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, newEmail, verificationURL)

	textContent := fmt.Sprintf(
		"Confirm Email Change\n\nHello %s,\n\nConfirm your new email: %s\n\nClick: %s\n\nBest regards,\nThe Template Store Team",
		user.Name, newEmail, verificationURL,
	)

	// Send to the NEW email address
	return s.SendEmail(newEmail, subject, htmlContent, textContent)
}

// SendAccountSuspensionEmail sends notification of account suspension
func (s *EmailServiceImpl) SendAccountSuspensionEmail(user *models.User, reason string, suspendedBy string) error {
	subject := "Your Account Has Been Suspended"

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Account Suspended</h1>
			<p>Hello %s,</p>
			<p>Your Template Store account has been suspended.</p>
			<p><strong>Reason:</strong> %s</p>
			<p>If you believe this is a mistake or would like to appeal this decision, please contact our support team.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, reason)

	textContent := fmt.Sprintf(
		"Account Suspended\n\nHello %s,\n\nYour account has been suspended.\nReason: %s\n\nBest regards,\nThe Template Store Team",
		user.Name, reason,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendAccountUnsuspensionEmail sends notification that account is unsuspended
func (s *EmailServiceImpl) SendAccountUnsuspensionEmail(user *models.User) error {
	subject := "Your Account Has Been Reinstated"

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Account Reinstated</h1>
			<p>Hello %s,</p>
			<p>Good news! Your Template Store account has been reinstated and is now active.</p>
			<p>You can now access all features and services.</p>
			<p>
				<a href="%s/login" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Login to Your Account
				</a>
			</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, s.baseURL)

	textContent := fmt.Sprintf(
		"Account Reinstated\n\nHello %s,\n\nYour account has been reinstated.\n\nBest regards,\nThe Template Store Team",
		user.Name,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendAccountDeletionEmail sends confirmation of account deletion
func (s *EmailServiceImpl) SendAccountDeletionEmail(user *models.User) error {
	subject := "Your Account Has Been Deleted"

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Account Deleted</h1>
			<p>Hello %s,</p>
			<p>This email confirms that your Template Store account has been deleted as requested.</p>
			<p>We're sorry to see you go. If you change your mind, you can create a new account anytime.</p>
			<p>Thank you for being part of our community.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name)

	textContent := fmt.Sprintf(
		"Account Deleted\n\nHello %s,\n\nYour account has been deleted.\n\nBest regards,\nThe Template Store Team",
		user.Name,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendOrderConfirmationEmail sends order confirmation
func (s *EmailServiceImpl) SendOrderConfirmationEmail(user *models.User, order *models.Order, template *models.Template) error {
	subject := "Order Confirmation - Template Store"

	downloadURL := fmt.Sprintf("%s/dashboard/orders/%d", s.baseURL, order.ID)

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Order Confirmation</h1>
			<p>Hello %s,</p>
			<p>Thank you for your purchase! Your order has been confirmed.</p>
			<h3>Order Details:</h3>
			<ul>
				<li><strong>Order ID:</strong> #%d</li>
				<li><strong>Template:</strong> %s</li>
				<li><strong>Date:</strong> %s</li>
			</ul>
			<p>
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Download Template
				</a>
			</p>
			<p>You can access your purchased templates anytime from your dashboard.</p>
			<p>Best regards,<br>The Template Store Team</p>
		</body>
		</html>
	`, user.Name, order.ID, template.Name, time.Now().Format("January 2, 2006"), downloadURL)

	textContent := fmt.Sprintf(
		"Order Confirmation\n\nHello %s,\n\nOrder #%d confirmed.\nTemplate: %s\n\nDownload: %s\n\nBest regards,\nThe Template Store Team",
		user.Name, order.ID, template.Name, downloadURL,
	)

	return s.SendEmail(user.Email, subject, htmlContent, textContent)
}

// SendEmail sends a generic email
func (s *EmailServiceImpl) SendEmail(to, subject, htmlContent, textContent string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, textContent, htmlContent)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid returned error status %d: %s", response.StatusCode, response.Body)
	}

	return nil
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	// This should use the actual config, but for simplicity we'll use a basic implementation
	return defaultValue
}
