package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/linuxunsw/vote/tui/internal/tui/validation"
)

// Creates a new form to prompt the user for their zID
func ZID() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("zid").
				Title("what's your zid?").
				Placeholder("z5555555").
				Validate(validation.ZID),
		),
	).WithTheme(styles.FormTheme())
}

// Creates a new form to prompt the user for the OTP sent
// to the email associated with their zID
func OTP() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("enter verification code").
				Key("otp").
				Description("a code has been sent to the email associated with your zid").
				Validate(validation.OTP),
		),
	).WithTheme(styles.FormTheme())
}
