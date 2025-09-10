package mailer

import "embed"

//go:embed "templates"
var FS embed.FS

type Mailer interface {
	SendOTP(toEmail string, otpCode string) error
}
