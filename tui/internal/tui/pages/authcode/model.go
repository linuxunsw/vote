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

	// Maximum size allowed for content
	cWidth  int
	cHeight int

	// OTP form
	form *huh.Form

	// Prevent resubmission of an already submitted form
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
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		log.Debug("PageContentSizeMsg", "height", msg.Height, "width", msg.Width)
		m.cHeight = msg.Height
		m.cWidth = msg.Width
		m.form = m.form.WithHeight(m.cHeight).WithWidth(m.cWidth)

		formHeight := lipgloss.Height(m.form.View())
		formWidth := lipgloss.Width(m.form.View())
		m.logger.Debug("Form Size", "height", formHeight, "width", formWidth)
		return m, nil
	case messages.ServerErrMsg:
		m.isSubmitted = false
		m.form = forms.OTP().WithHeight(m.cHeight).WithWidth(m.cWidth)
		m.form.Init()
		return m, nil
	}

	// If the user has completed the form, send the OTP
	// and prevent them from resubmitting the form until
	// the page either changes or an error occurs and we
	// reset the form + display an error message
	if m.form.State == huh.StateCompleted && !m.isSubmitted {
		m.isSubmitted = true
		otp := m.form.GetString("otp")
		return m, messages.SendSubmitOTP(otp)
	}

	return m, cmd
}

func (m *authCodeModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
