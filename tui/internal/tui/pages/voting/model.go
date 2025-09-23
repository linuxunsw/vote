package voting

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/linuxunsw/vote/tui/internal/sdk"
	"github.com/linuxunsw/vote/tui/internal/tui/forms"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
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
func New(logger *log.Logger, data sdk.PublicBallot) tea.Model {
	var vote map[string]string

	model := &formModel{
		logger:      logger,
		isSubmitted: false,
	}

	model.form = forms.Voting(data, vote)

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

	if m.form.State == huh.StateCompleted && !m.isSubmitted {
		m.isSubmitted = true

		positions := make(map[string]string)

		for key, val := range map[string]string{
			"president":         m.form.GetString("president"),
			"secretary":         m.form.GetString("secretary"),
			"treasurer":         m.form.GetString("treasurer"),
			"arc_delegate":      m.form.GetString("arc_delegate"),
			"edi_officer":       m.form.GetString("edi_officer"),
			"grievance_officer": m.form.GetString("grievance_officer"),
		} {
			if val != "" {
				positions[key] = val
			}
		}

		return m, sdk.SendSubmitVote(positions)
	}

	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		log.Debug("PageContentSizeMsg", "height", msg.Height, "width", msg.Width)
		m.cHeight = msg.Height
		m.cWidth = msg.Width

		m.form = m.form.WithHeight(m.cHeight).WithWidth(m.cWidth)
	case sdk.ServerErrMsg:
		return m, tea.Sequence(
			messages.SendPageChange(pages.VotingSubmit),
			sdk.SendPublicSubmitFormResult(msg.RespID, msg.Error),
		)
	}

	return m, cmd
}

// Display the form
func (m *formModel) View() string {
	return styles.FormStyle.Render(m.form.View())
}
