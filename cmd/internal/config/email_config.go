package config

import (
	"log"
	"os"
	"strconv"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	Sender   string
}

func LoadEmailConfig() EmailConfig {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	return EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Sender:   os.Getenv("SMTP_FROM_ADDRESS"),
	}
}
