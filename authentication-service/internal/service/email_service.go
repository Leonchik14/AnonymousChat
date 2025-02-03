package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

// EmailService - структура сервиса отправки email
type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewEmailService - конструктор сервиса Email
func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

// SendEmail - Отправка email через SMTP
func (e *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.host, e.port, e.username, e.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("Ошибка отправки email: %v", err)
	}
	return nil
}

func (e *EmailService) SendVerificationEmail(to string, token string) error {
	subject := "Подтверждение регистрации"
	body := fmt.Sprintf(`
		<h2>Подтвердите свою регистрацию</h2>
		<p>Нажмите <a href="http://localhost:8080/auth/verify?token=%s">сюда</a>, чтобы подтвердить email.</p>`, token)

	return e.SendEmail(to, subject, body)
}

// SendPasswordResetEmail - Отправка письма для сброса пароля
func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Сброс пароля"
	body := fmt.Sprintf(`
		<h2>Восстановление пароля</h2>
		<p>Для сброса пароля нажмите <a href="http://localhost:8080/auth/reset-password?token=%s">сюда</a>.</p>`, token)

	return e.SendEmail(to, subject, body)
}
