package service

import (
	"auth-server/models"
	"fmt"
	"math/rand/v2"
	"net/smtp"
	"strings"
	"time"
)

func (s *AuthService) sendResetPasswordCode(code string, email string) error {
	emailBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>Password Reset</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f4f4f4;
			padding: 20px;
		}
	.container {
		max-width: 600px;
		margin: auto;
		background: white;
		padding: 30px;
		border-radius: 8px;
		box-shadow: 0 0 10px rgba(0,0,0,0.1);
	}
	h1 {
		color: #333;
	}
	p {
		font-size: 16px;
		line-height: 1.5;
	}
	.button {
		display: inline-block;
		margin-top: 20px;
		padding: 12px 24px;
		background-color: #335f8f;
		color: white;
		text-decoration: none;
		font-weight: bold;
		border-radius: 5px;
	}
	a {
      color: white;
    }
	.footer {
		margin-top: 30px;
		font-size: 14px;
		color: #777;
	}
	</style>
	</head>
	<body>
	<div class="container">
	<h1>Password Reset Request</h1>
	<p>We received a request to reset your account password. If this was you, please click the button below to continue:</p>

	<a href="https://smarthome.hipahopa.ru/reset-password?code=%s&email=%s" class="button">Reset Password</a>

	<p>If you did not request a password reset, you can safely ignore this email.</p>

	<div class="footer">
		This link will expire in 15 minutes for security reasons.
	</div>
	</div>
	</body>
	</html>`, code, email)

	err := s.sendEmail([]string{email}, "Password Reset Request", emailBody)
	if err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return nil
}

func (s *AuthService) sendEmailConfirmationCode(user *models.User) error {
	// generate code
	confirmationCode := fmt.Sprintf("%06d", int(rand.Float32()*1000000))
	expiresAt := 15 * time.Minute

	err := s.repo.SetEmailConfirmationCode(confirmationCode, user.Email, expiresAt)
	if err != nil {
		return err
	}

	// send code
	emailBody := fmt.Sprintf(`
       <!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Verification Code</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f9f9f9;
					padding: 20px;
				}
				.email-container {
					max-width: 500px;
					margin: auto;
					background-color: #ffffff;
					padding: 30px;
					border-radius: 8px;
					box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
				}
				h2 {
					color: #333333;
					font-size: 24px;
					margin-bottom: 20px;
				}
				p {
					font-size: 16px;
					color: #555555;
					line-height: 1.5;
				}
				.code {
					display: inline-block;
					margin: 20px 0;
					padding: 12px 24px;
					font-size: 22px;
					letter-spacing: 2px;
					font-weight: bold;
					color: #333333;
					background-color: #f0f0f0;
					border-radius: 6px;
					word-break: break-all;
				}
				.footer {
					margin-top: 20px;
					font-size: 14px;
					color: #aaaaaa;
				}
			</style>
		</head>
		<body>
			<div class="email-container">
				<h2>Your Verification Code</h2>
				<p>Please use the following code to verify your email:</p>
				<div class="code">%s</div>
				<p>This code will expire in %s minutes.</p>
				<div class="footer">
					If you did not request this code, please ignore this email.
				</div>
			</div>
		</body>
		</html>
    `, confirmationCode, fmt.Sprintf("%.0f", expiresAt.Minutes()))

	err = s.sendEmail([]string{user.Email}, "Your Verification Code", emailBody)
	if err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return nil
}

func (s *AuthService) sendEmail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.smtpConfig.Username, s.smtpConfig.Password, s.smtpConfig.Host)

	message := []byte(
		"From: " + subject + "\r\n" +
			"To: " + strings.Join(to, ", ") + "\r\n" +
			"From: " + "Smarthome" + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n" +
			"\r\n" + // Empty line to separate headers from body
			body,
	)

	err := smtp.SendMail(
		s.smtpConfig.Host+":"+s.smtpConfig.Port,
		auth,
		s.smtpConfig.Username,
		to,
		message,
	)

	return err
}
