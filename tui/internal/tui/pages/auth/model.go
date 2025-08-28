package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authModel struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

func New() tea.Model {
	model := &authModel{
		form: forms.ZID(),
	}

	return model
}

// Initialise form
func (m *authModel) Init() tea.Cmd {
	return m.form.Init()
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
	if m.form.State == huh.StateCompleted {
		zID := m.form.GetString("zid")
		return m, tea.Batch(
			messages.SendAuth(zID),
			messages.SendPageChange(pages.PageAuthCode),
		)
	}

	// Handle remaining bubble tea commands
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wHeight = msg.Height
		m.wWidth = msg.Width
	}

	return m, tea.Batch(cmds...)
}

// Display form
func (m *authModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
