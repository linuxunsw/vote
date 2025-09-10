package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type formModel struct {
	logger *log.Logger

	cWidth  int
	cHeight int

	form *huh.Form

	isSubmitted bool
}

// Creates model
// TODO: determine which form to display depending on
// server state (e.g. nomination vs voting)
func New(logger *log.Logger) tea.Model {
	model := &formModel{
		logger:      logger,
		form:        forms.Nomination(),
		isSubmitted: false,
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
	if m.form.State == huh.StateCompleted && !m.isSubmitted {
		m.isSubmitted = true

		// TODO: get roles from form

		data := messages.Submission{
			Name:      m.form.GetString("name"),
			Email:     m.form.GetString("email"),
			Discord:   m.form.GetString("discord"),
			Statement: m.form.GetString("statement"),
			Url:       m.form.GetString("url"),
		}

		formRoles := m.form.Get("roles")
		roles, ok := formRoles.([]string)
		if !ok {
			// FIX: remove
			m.logger.Debug("what the fuck")
		}
		data.Roles = roles

		return m, tea.Batch(
			messages.SendSubmission(data),
		)
	}

	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		log.Debug("PageContentSizeMsg", "height", msg.Height, "width", msg.Width)
		m.cHeight = msg.Height
		m.cWidth = msg.Width

		m.form = m.form.WithHeight(m.cHeight).WithWidth(m.cWidth)
	}

	return m, cmd
}

// Display the form
func (m *formModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
