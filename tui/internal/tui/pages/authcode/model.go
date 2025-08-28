package authcode

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authCodeModel struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

func New() tea.Model {
	model := &authCodeModel{
		form: forms.OTP(),
	}

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
	case messages.IsAuthenticatedMsg:
		return m, messages.SendPageChange(pages.PageForm)
	}

	if m.form.State == huh.StateCompleted {
		otp := m.form.GetString("otp")
		return m, messages.SendCheckOTP(otp)
	}

	return m, cmd
}

func (m *authCodeModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
