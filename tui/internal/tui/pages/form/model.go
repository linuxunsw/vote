package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/linuxunsw/vote/tui/internal/tui/validation"
)

type formModel struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

func New() tea.Model {
	model := &formModel{}

	// TODO: take in a config for fields displayed
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("full name").
				Validate(validation.NotEmpty),
			huh.NewInput().
				Key("email").
				Title("preferred contact email").
				Validate(validation.Email),
			huh.NewInput().
				Key("discord").
				Title("discord username").
				Validate(validation.NotEmpty),
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
				Validate(validation.NotEmpty),
			huh.NewInput().
				Key("url").
				Title("url"),
		),
	).WithTheme(huh.ThemeBase16())

	model.form = form

	return model
}

func (m *formModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wHeight = msg.Height
		m.wWidth = msg.Width
	}

	return m, cmd
}

func (m *formModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
