package authcode

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authCode struct {
	wWidth  int
	wHeight int

	form *huh.Form
}

func New() tea.Model {
	model := &authCode{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("enter code").
				Description("a code has been sent to the email associated with your zid"),
		),
	).WithTheme(huh.ThemeBase16())

	model.form = form

	return model
}

func (a *authCode) Init() tea.Cmd {
	return a.form.Init()
}

func (a *authCode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := a.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		a.form = f
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.wHeight = msg.Height
		a.wWidth = msg.Width
	}

	return a, cmd
}

func (a *authCode) View() string {
	return styles.FormStyle.Render(a.form.View())
}
