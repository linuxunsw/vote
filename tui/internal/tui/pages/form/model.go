package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type formModel struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

// Creates model
// TODO: determine which form to display depending on
// server state (e.g. nomination vs voting)
func New() tea.Model {
	model := &formModel{
		form: forms.NominationForm(),
	}

	return model
}

func (m *formModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Quit app when form completed
	// TODO: change to submission page
	if m.form.State == huh.StateCompleted {
		// TODO: get roles from form
		data := messages.Submission{
			Name:      m.form.GetString("name"),
			Email:     m.form.GetString("email"),
			Discord:   m.form.GetString("discord"),
			Statement: m.form.GetString("statement"),
			Url:       m.form.GetString("url"),
		}

		return m, tea.Batch(
			messages.SendSubmission(data),
		)
	}

	// Handle remaining bubble tea commands
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wHeight = msg.Height
		m.wWidth = msg.Width
	}

	return m, cmd
}

// Display the form
func (m *formModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
