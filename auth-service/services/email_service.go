package services

import (
	"auth-server/utils"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

type EmailService interface {
	SendPasswordReset(to, resetURL string) error
}

type emailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService() EmailService {
	host := utils.GetEnv("SMTP_HOST", "")
	portStr := utils.GetEnv("SMTP_PORT", "587")
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587
	}
	username := utils.GetEnv("SMTP_USERNAME", "")
	from := utils.GetEnv("SMTP_FROM", username)

	return &emailService{
		host:     host,
		port:     port,
		username: username,
		password: utils.GetEnv("SMTP_PASSWORD", ""),
		from:     from,
	}
}

func (s *emailService) SendPasswordReset(to, resetURL string) error {
	if s.host == "" || s.username == "" {
		return fmt.Errorf("smtp not configured")
	}

	subject := "Сброс пароля — Smart Home"
	body := strings.Join([]string{
		"<!doctype html><html><body style=\"font-family: -apple-system, Segoe UI, Roboto, sans-serif; color: #0f172a;\">",
		"<div style=\"max-width: 480px; margin: 24px auto; padding: 24px; border: 1px solid #e2e8f0; border-radius: 12px;\">",
		"<h2 style=\"margin: 0 0 12px;\">Сброс пароля</h2>",
		"<p style=\"color: #475569; line-height: 1.55;\">Вы запросили сброс пароля. Перейдите по ссылке, чтобы задать новый пароль. Ссылка действует 30 минут.</p>",
		fmt.Sprintf("<p><a href=\"%s\" style=\"display:inline-block;padding:10px 18px;background:#6366f1;color:white;text-decoration:none;border-radius:8px;font-weight:500;\">Сбросить пароль</a></p>", resetURL),
		fmt.Sprintf("<p style=\"color: #94a3b8; font-size: 12px;\">Если кнопка не работает, скопируйте: %s</p>", resetURL),
		"<p style=\"color: #94a3b8; font-size: 12px;\">Если вы не запрашивали сброс — просто проигнорируйте письмо.</p>",
		"</div></body></html>",
	}, "")

	msg := []byte(strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n"))

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if s.port == 465 {
		return s.sendImplicitTLS(addr, auth, to, msg)
	}
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}

func (s *emailService) sendImplicitTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.host})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}
