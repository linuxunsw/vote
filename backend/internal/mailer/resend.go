package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/resend/resend-go/v2"
)

type ResendMailer struct {
	fromEmail string
	apiKey    string
	client    *resend.Client
	otpExpiry time.Duration
}

func NewResendMailer(cfg config.Config) Mailer {
	client := resend.NewClient(cfg.Mailer.ResendAPIKey)
	return &ResendMailer{
		fromEmail: cfg.Mailer.FromEmail,
		apiKey:    cfg.Mailer.ResendAPIKey,
		client:    client,
		otpExpiry: cfg.OTP.Duration,
	}
}

func (m *ResendMailer) SendOTP(toEmail, otpCode string) error {
	tmpl, err := template.ParseFS(FS, "templates/otp.html")
	if err != nil {
		return err
	}

	data := struct {
		OTP           string
		ExpiryMinutes int
	}{
		OTP:           otpCode,
		ExpiryMinutes: int(m.otpExpiry.Minutes()),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	plainText := fmt.Sprintf("Your OTP code is: %s\n\nThis code will expire in %d minutes.\n\nRegards,\nLinux Society UNSW", otpCode, data.ExpiryMinutes)

	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{toEmail},
		Html:    buf.String(),
		Text:    plainText,
		Subject: "Your Linux Society Vote OTP",
	}

	_, err = m.client.Emails.Send(params)
	return err
}
