package email

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
}

func NewEmailService() *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		panic("RESEND_API_KEY environment variable is required")
	}

	client := resend.NewClient(apiKey)
	return &EmailService{
		client: client,
	}
}

type VerificationEmailData struct {
	Email             string
	FullName          string
	VerificationToken string
	VerificationURL   string
}

func (e *EmailService) SendVerificationEmail(ctx context.Context, data VerificationEmailData) error {
	// Generate verification URL
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	verificationURL := fmt.Sprintf("%s/auth/verify?token=%s", baseURL, data.VerificationToken)

	// Create HTML email content
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Verify Your Email - Syncer</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #4f46e5;
            margin-bottom: 10px;
        }
        .title {
            font-size: 24px;
            font-weight: 600;
            color: #1f2937;
            margin-bottom: 20px;
        }
        .content {
            font-size: 16px;
            color: #6b7280;
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            background: #4f46e5;
            color: white;
            padding: 12px 30px;
            border-radius: 6px;
            text-decoration: none;
            font-weight: 600;
            margin: 20px 0;
        }
        .button:hover {
            background: #4338ca;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e5e7eb;
            font-size: 14px;
            color: #9ca3af;
        }
        .expiry {
            background: #fef3c7;
            border: 1px solid #f59e0b;
            border-radius: 6px;
            padding: 12px;
            margin: 20px 0;
            font-size: 14px;
            color: #92400e;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Syncer</div>
            <h1 class="title">Verify Your Email Address</h1>
        </div>
        
        <div class="content">
            <p>Hi %s,</p>
            
            <p>Welcome to Syncer! We're excited to have you as part of our synchronization platform.</p>
            
            <p>To complete your registration and start using Syncer, please verify your email address by clicking the button below:</p>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s" style="display:inline-block;background-color:#4f46e5;color:#ffffff !important;padding:12px 30px;border-radius:6px;text-decoration:none;font-weight:600;">Verify Email Address</a>
            </div>
            
            <div class="expiry">
                <strong>⏰ Important:</strong> This verification link will expire in 24 hours for security reasons.
            </div>
            
            <p>If you can't click the button above, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #4f46e5;">%s</p>
            
            <p>If you didn't create an account with Syncer, you can safely ignore this email.</p>
        </div>
        
        <div class="footer">
            <p>Best regards,<br>The Syncer Team</p>
            <p>This email was sent to %s. If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`, data.FullName, verificationURL, verificationURL, data.Email)

	// Create plain text version
	textContent := fmt.Sprintf(`
Hi %s,

Welcome to Syncer! We're excited to have you as part of our synchronization platform.

To complete your registration and start using Syncer, please verify your email address by clicking the link below:

%s

⏰ Important: This verification link will expire in 24 hours for security reasons.

If you didn't create an account with Syncer, you can safely ignore this email.

Best regards,
The Syncer Team

This email was sent to %s.
`, data.FullName, verificationURL, data.Email)

	// Send email using Resend
	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM_EMAIL"), // e.g., "noreply@syncer.net"
		To:      []string{data.Email},
		Subject: "Verify Your Email Address - Syncer",
		Html:    htmlContent,
		Text:    textContent,
	}

	// Set default from email if not provided
	if params.From == "" {
		params.From = "noreply@syncer.net"
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (e *EmailService) SendPasswordResetEmail(ctx context.Context, email, fullName, resetToken string) error {
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, resetToken)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Reset Your Password - Syncer</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #4f46e5;
            margin-bottom: 10px;
        }
        .title {
            font-size: 24px;
            font-weight: 600;
            color: #1f2937;
            margin-bottom: 20px;
        }
        .content {
            font-size: 16px;
            color: #6b7280;
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            background: #dc2626;
            color: white;
            padding: 12px 30px;
            border-radius: 6px;
            text-decoration: none;
            font-weight: 600;
            margin: 20px 0;
        }
        .button:hover {
            background: #b91c1c;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e5e7eb;
            font-size: 14px;
            color: #9ca3af;
        }
        .expiry {
            background: #fef2f2;
            border: 1px solid #f87171;
            border-radius: 6px;
            padding: 12px;
            margin: 20px 0;
            font-size: 14px;
            color: #dc2626;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Syncer</div>
            <h1 class="title">Reset Your Password</h1>
        </div>
        
        <div class="content">
            <p>Hi %s,</p>
            
            <p>We received a request to reset your password for your Syncer account.</p>
            
            <p>If you made this request, click the button below to reset your password:</p>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Reset Password</a>
            </div>
            
            <div class="expiry">
                <strong>⏰ Important:</strong> This password reset link will expire in 1 hour for security reasons.
            </div>
            
            <p>If you can't click the button above, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #4f46e5;">%s</p>
            
            <p>If you didn't request a password reset, you can safely ignore this email. Your password will remain unchanged.</p>
        </div>
        
        <div class="footer">
            <p>Best regards,<br>The Syncer Team</p>
            <p>This email was sent to %s. If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`, fullName, resetURL, resetURL, email)

	textContent := fmt.Sprintf(`
Hi %s,

We received a request to reset your password for your Syncer account.

If you made this request, click the link below to reset your password:

%s

⏰ Important: This password reset link will expire in 1 hour for security reasons.

If you didn't request a password reset, you can safely ignore this email. Your password will remain unchanged.

Best regards,
The Syncer Team

This email was sent to %s.
`, fullName, resetURL, email)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM_EMAIL"),
		To:      []string{email},
		Subject: "Reset Your Password - Syncer",
		Html:    htmlContent,
		Text:    textContent,
	}

	if params.From == "" {
		params.From = "noreply@syncer.net"
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (e *EmailService) SendWelcomeEmail(ctx context.Context, email, fullName string) error {
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to Syncer!</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #4f46e5;
            margin-bottom: 10px;
        }
        .title {
            font-size: 24px;
            font-weight: 600;
            color: #1f2937;
            margin-bottom: 20px;
        }
        .content {
            font-size: 16px;
            color: #6b7280;
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            background: #10b981;
            color: white;
            padding: 12px 30px;
            border-radius: 6px;
            text-decoration: none;
            font-weight: 600;
            margin: 20px 0;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e5e7eb;
            font-size: 14px;
            color: #9ca3af;
        }
        .feature {
            background: #f9fafb;
            border-radius: 6px;
            padding: 20px;
            margin: 20px 0;
        }
        .feature-title {
            font-weight: 600;
            color: #1f2937;
            margin-bottom: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Syncer</div>
            <h1 class="title">Welcome to Syncer! 🎉</h1>
        </div>
        
        <div class="content">
            <p>Hi %s,</p>
            
            <p>Congratulations! Your email has been verified and you're now part of the Syncer community.</p>
            
            <p>Syncer is your all-in-one synchronization platform that helps you keep your digital life organized across multiple services.</p>
            
            <div class="feature">
                <div class="feature-title">🗓️ Calendar Sync</div>
                <p>Seamlessly sync your calendars across Google, Outlook, and other services.</p>
            </div>
            
            <div class="feature">
                <div class="feature-title">🎵 Music Sync</div>
                <p>Keep your playlists and music preferences synchronized across Spotify, Apple Music, and more.</p>
            </div>
            
            <div class="feature">
                <div class="feature-title">🔄 Real-time Updates</div>
                <p>Get instant notifications when your data syncs across platforms.</p>
            </div>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s/dashboard" class="button">Get Started</a>
            </div>
            
            <p>If you have any questions or need help getting started, don't hesitate to reach out to our support team.</p>
        </div>
        
        <div class="footer">
            <p>Best regards,<br>The Syncer Team</p>
            <p>This email was sent to %s.</p>
        </div>
    </div>
</body>
</html>
`, fullName, os.Getenv("FRONTEND_URL"), email)

	textContent := fmt.Sprintf(`
Hi %s,

Congratulations! Your email has been verified and you're now part of the Syncer community.

Syncer is your all-in-one synchronization platform that helps you keep your digital life organized across multiple services.

🗓️ Calendar Sync - Seamlessly sync your calendars across Google, Outlook, and other services.
🎵 Music Sync - Keep your playlists and music preferences synchronized across Spotify, Apple Music, and more.
🔄 Real-time Updates - Get instant notifications when your data syncs across platforms.

Get started: %s/dashboard

If you have any questions or need help getting started, don't hesitate to reach out to our support team.

Best regards,
The Syncer Team

This email was sent to %s.
`, fullName, os.Getenv("FRONTEND_URL"), email)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM_EMAIL"),
		To:      []string{email},
		Subject: "Welcome to Syncer! 🎉",
		Html:    htmlContent,
		Text:    textContent,
	}

	if params.From == "" {
		params.From = "noreply@syncer.net"
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}
