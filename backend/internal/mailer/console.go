package mailer

import (
	"log/slog"
)

type ConsoleMailer struct {
	Logger *slog.Logger
}

func (c *ConsoleMailer) SendOTP(toEmail, otpCode string) error {
	c.Logger.Info("Sending OTP", "email", toEmail, "otp", otpCode)
	return nil
}

func NewConsoleMailer(logger *slog.Logger) Mailer {
	return &ConsoleMailer{Logger: logger}
}
