package email

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

const (
	smtpServerAddress = "smtp.gmail.com:587"
	smtpAuthAddress   = "smtp.gmail.com"
)

type EmailSender interface {
	SendEmail(
		to []string,
		cc []string,
		bcc []string,
		subject string,
		content string,
		attachedFiles []string,
	) error
}

type GmailSender struct {
	name           string
	senderEmail    string
	senderPassword string
	pool           *email.Pool
}

func NewGmailSender(name, senderEmail, senderPassword string) EmailSender {
	pool, err := email.NewPool(
		smtpServerAddress,
		4,
		smtp.PlainAuth("", senderEmail, senderPassword, smtpAuthAddress),
	)

	if err != nil {
		log.Panic().Err(err).Msg("failed to create email sender")
	}

	return &GmailSender{
		pool:           pool,
		name:           name,
		senderEmail:    senderEmail,
		senderPassword: senderPassword,
	}
}

func (sender *GmailSender) SendEmail(
	to []string,
	cc []string,
	bcc []string,
	subject string,
	content string,
	attachedFiles []string,
) error {
	e := &email.Email{
		To:      to,
		Cc:      cc,
		Bcc:     bcc,
		From:    fmt.Sprintf("%s <%s>", sender.name, sender.senderEmail),
		Subject: subject,
		HTML:    []byte(content),
	}
	for _, file := range attachedFiles {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", file, err)
		}
	}

	err := sender.pool.Send(e, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
