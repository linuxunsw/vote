package authcode

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/linuxunsw/vote/tui/internal/tui/validation"
)

type authCodeModel struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

func New() tea.Model {
	model := &authCodeModel{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("enter verification code").
				Key("otp").
				Description("a code has been sent to the email associated with your zid").
				Validate(validation.OTP),
		),
	).WithTheme(huh.ThemeBase16())

	model.form = form

	return model
}

func (m *authCodeModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *authCodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wHeight = msg.Height
		m.wWidth = msg.Width
	case messages.AuthenticatedMsg:
		return m, tea.Batch(func() tea.Msg { return messages.PageChangeMsg{ID: pages.PageForm} })
	}

	if m.form.State == huh.StateCompleted {
		otp := m.form.GetString("otp")
		return m, tea.Batch(
			func() tea.Msg { return messages.CheckOTPMsg{OTP: otp} },
		)
	}

	return m, cmd
}

func (m *authCodeModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
