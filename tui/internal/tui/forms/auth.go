package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/linuxunsw/vote/tui/internal/tui/validation"
)

var ZIDForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Key("zid").
			Title("what's your zid?").
			Placeholder("z5555555").
			Validate(validation.ZID),
	),
).WithTheme(styles.FormTheme())

var OTPForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("enter verification code").
			Key("otp").
			Description("a code has been sent to the email associated with your zid").
			Validate(validation.OTP),
	),
).WithTheme(styles.FormTheme())
