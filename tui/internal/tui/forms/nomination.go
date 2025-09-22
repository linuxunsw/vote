package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/linuxunsw/vote/tui/internal/tui/validation"
)

// Creates a form for the user's nomination information
func Nomination() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("full name").
				Validate(huh.ValidateLength(2, 100)),
			huh.NewInput().
				Key("email").
				Title("preferred contact email").
				Validate(validation.Email),
			huh.NewInput().
				Key("discord").
				Title("discord username").
				Validate(huh.ValidateLength(2, 32)),
			huh.NewMultiSelect[string]().
				Key("roles").
				Title("roles you are nominating for").
				Options(
					huh.NewOption("president", "president"),
					huh.NewOption("secretary", "secretary"),
					huh.NewOption("treasurer", "treasurer"),
					huh.NewOption("arc delegate", "arc_delegate"),
					huh.NewOption("edi officer", "edi_officer"),
					huh.NewOption("grievance officer", "grievance_officer"),
				).
				Validate(validation.Role),
			huh.NewText().
				Key("statement").
				Title("please provide a candidate statement").
				ExternalEditor(false).
				Validate(huh.ValidateLength(50, 2000)),
			huh.NewInput().
				Key("url").
				Title("url (optional)").
				Validate(validation.URL),
		),
	).WithTheme(styles.FormTheme())
}
