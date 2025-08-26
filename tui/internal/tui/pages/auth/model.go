// TODO: save value from form to root's data
package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type state int

const (
	initalising state = iota
	ready
)

type auth struct {
	form *huh.Form

	state state
}

type SendAuthMsg struct {
	ZID string
}

func New() tea.Model {
	model := &auth{
		state: initalising,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("zid").
				Title("what's your zid?"),
		),
	).WithTheme(huh.ThemeBase16())

	model.form = form

	return model
}

func (a *auth) Init() tea.Cmd {
	return a.form.Init()
}

func (a *auth) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := a.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		cmds = append(cmds, cmd)
		a.form = f
	}

	if a.form.State == huh.StateCompleted {
		zid := a.form.GetString("zid")
		return a, tea.Batch(
			func() tea.Msg { return SendAuthMsg{ZID: zid} },
			func() tea.Msg { return pages.PageChangeMsg{ID: pages.PageAuthCode} },
		)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.wHeight = msg.Height
		a.wWidth = msg.Width
		a.state = ready
	}

	return a, tea.Batch(cmds...)
}

func (a *auth) View() string {
	// set content to spinner
	return styles.FormStyle.Render(a.form.View())
}
