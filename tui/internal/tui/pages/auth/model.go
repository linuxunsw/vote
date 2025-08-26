// TODO: save value from form to root's data
package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/tui/components"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

type authModel struct {
	form *huh.Form

	zid string
}

func New() tea.Model {
	model := &auth{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("what's your zid?").
				Value(&model.zid),
		),
	).WithTheme(huh.ThemeBase())

	model.form = form

	return model
}

func (a *auth) Init() tea.Cmd {
	return a.form.Init()
}

func (a *auth) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := a.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		a.form = f
	}

	return a, cmd
}

func (a *auth) View() string {
	styledForm := styles.FormStyle.Render(a.form.View())
	return styles.AppStyle.Render(components.ShowHeader() + "\n" + styledForm)
}
