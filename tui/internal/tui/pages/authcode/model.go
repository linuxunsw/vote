package authcode

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authCodeModel struct {
	logger *log.Logger

	wWidth  int
	wHeight int

	form *huh.Form

	isSubmitted bool
}

func New(logger *log.Logger) tea.Model {
	model := &authCodeModel{
		logger:      logger,
		form:        forms.OTP(),
		isSubmitted: false,
	}

	return model
}

func (m *authCodeModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *authCodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		log.Debug("PageContentSizeMsg", "msg", msg)
		m.wHeight = msg.Height
		m.wWidth = msg.Width
		m.form = m.form.WithHeight(m.wHeight).WithWidth(m.wWidth)

		formHeight := lipgloss.Height(m.form.View())
		formWidth := lipgloss.Width(m.form.View())
		m.logger.Debug("Form Size", "height", formHeight, "width", formWidth)
		return m, nil
	case tea.WindowSizeMsg:
		return m, nil
	}
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted && !m.isSubmitted {
		m.isSubmitted = true
		otp := m.form.GetString("otp")
		return m, messages.SendVerifyOTP(otp)
	}

	return m, cmd
}

func (m *authCodeModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
