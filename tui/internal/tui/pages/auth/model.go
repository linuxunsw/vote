package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authModel struct {
	logger *log.Logger

	cWidth  int
	cHeight int

	form *huh.Form

	isSubmitted bool
}

func New(logger *log.Logger) tea.Model {
	model := &authModel{
		logger:      logger,
		form:        forms.ZID(),
		isSubmitted: false,
	}

	return model
}

// Initialise form
func (m *authModel) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		tea.WindowSize(),
	)
}

func (m *authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Pass the message to the form model
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		cmds = append(cmds, cmd)
		m.form = f
	}

	// If we've completed the form, send the form data to root and change the page
	if m.form.State == huh.StateCompleted && !m.isSubmitted {
		zID := m.form.GetString("zid")

		m.isSubmitted = true

		return m, messages.SendRequestOTP(zID)
	}

	// Handle remaining bubble tea commands
	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		m.cHeight = msg.Height
		m.cWidth = msg.Width
		m.form = m.form.WithHeight(m.cHeight).WithWidth(m.cWidth)

		formHeight := lipgloss.Height(m.form.View())
		formWidth := lipgloss.Width(m.form.View())
		m.logger.Debug("Form Size", "height", formHeight, "width", formWidth)
		return m, nil

	}

	return m, tea.Batch(cmds...)
}

// Display form
func (m *authModel) View() string {
	return styles.FormStyle.Render(m.form.View())

}
