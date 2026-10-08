package auth

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type smtpMailer struct {
	host     string
	port     int
	username string
	password string
	appURL   string
}

type disabledMailer struct{}

func (disabledMailer) SendResetPassword(string, string) error {
	return fmt.Errorf("SMTP is not configured")
}

func NewSMTPMailer(host string, port int, username, password, appURL string) Mailer {
	if strings.TrimSpace(host) == "" || port < 1 {
		return disabledMailer{}
	}
	return smtpMailer{host: host, port: port, username: username, password: password, appURL: strings.TrimRight(appURL, "/")}
}

func (m smtpMailer) SendResetPassword(to, token string) error {
	from := m.username
	if from == "" {
		return fmt.Errorf("SMTP_USER is required to send password reset email")
	}
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}
	if _, err := mail.ParseAddress(from); err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}

	address := net.JoinHostPort(m.host, strconv.Itoa(m.port))
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	var conn net.Conn
	var err error
	if m.port == 465 {
		conn, err = tls.DialWithDialer(dialer, "tcp", address, &tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = dialer.Dial("tcp", address)
	}
	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return fmt.Errorf("start SMTP client: %w", err)
	}
	defer client.Close()
	if m.port != 465 {
		if err := client.StartTLS(&tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("start SMTP TLS: %w", err)
		}
	}
	if m.username != "" {
		if err := client.Auth(smtp.PlainAuth("", m.username, m.password, m.host)); err != nil {
			return fmt.Errorf("authenticate SMTP client: %w", err)
		}
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start SMTP message: %w", err)
	}
	resetURL := m.appURL + "/reset-password?token=" + url.QueryEscape(token)
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nЧтобы задать новый пароль, откройте ссылку:\r\n%s\r\n\r\nСсылка действует один час. Если вы не запрашивали сброс пароля, проигнорируйте это письмо.\r\n",
		from, to, mime.QEncoding.Encode("UTF-8", "Сброс пароля УНИФЛОУ"), resetURL,
	)
	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("finish SMTP session: %w", err)
	}
	return nil
}

func newResetToken() string { return uuid.NewString() }
