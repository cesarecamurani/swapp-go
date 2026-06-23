package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/config"
	"swapp-go/cmd/internal/domain"
)

type SmtpEmailService struct {
	config config.EmailConfig
}

func NewSmtpEmailService(config config.EmailConfig) ports.EmailService {
	return &SmtpEmailService{config: config}
}

func (s *SmtpEmailService) SendEmail(message *domain.EmailMessage) error {
	from := s.config.Sender
	to := message.Recipient

	msg := buildMessage(from, to, message.Subject, message.Body)

	auth := smtp.PlainAuth(
		"",
		s.config.Username,
		s.config.Password,
		s.config.SMTPHost,
	)

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	
	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	data, err := client.Data()
	if err != nil {
		return err
	}

	_, err = data.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = data.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func buildMessage(from, to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""

	var msg strings.Builder

	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	msg.WriteString("\r\n" + body)

	return msg.String()
}
