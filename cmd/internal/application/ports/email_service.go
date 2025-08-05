package ports

import "swapp-go/cmd/internal/domain"

type EmailService interface {
	SendEmail(message *domain.EmailMessage) error
}
