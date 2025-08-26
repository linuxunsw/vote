// TODO: move between pages

package root

import (
	"os"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/linuxunsw/vote/tui/internal/tui/components"
	"github.com/linuxunsw/vote/tui/internal/tui/keys"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/auth"
	authcode "github.com/linuxunsw/vote/tui/internal/tui/pages/authCode"
)

/*
type data struct {
	zID     string
	name    string
	email   string
	discord string
	roles   []bool

	statement string
	url       string
}*/

type root struct {
	keyMap keys.KeyMap

	pages  map[pages.PageID]tea.Model
	loaded map[pages.PageID]bool

	current  pages.PageID
	previous pages.PageID

	// data data
}

func New() tea.Model {
	keyMap := keys.DefaultKeyMap()
	pageMap := map[pages.PageID]tea.Model{
		pages.PageAuth:     auth.New(),
		pages.PageAuthCode: authcode.New(),
	}

	model := &root{
		keyMap:  keyMap,
		pages:   pageMap,
		loaded:  make(map[pages.PageID]bool),
		current: pages.PageAuth,
	}

	return model

}

func (r *root) Init() tea.Cmd {
	r.loaded[r.current] = true

	windowTitle := os.Getenv("EVENT_NAME")

	return tea.Batch(
		r.pages[r.current].Init(),
		tea.SetWindowTitle(windowTitle),
	)
}

func (r *root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return r, r.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return r, r.handleWindowSizeMsg(msg)
	case pages.PageChangeMsg:
		return r, r.movePage(msg.ID)
	case auth.SendAuthMsg:
		r.data.zID = msg.ZID
		// TODO: send req to api with zid
		return r, cmd
	}

	updated, cmd := r.pages[r.current].Update(msg)
	r.pages[r.current] = updated

	cmds = append(cmds, cmd)
	return r, tea.Batch(cmds...)
}

func (r *root) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, r.keyMap.Quit):
		return tea.Quit
	default:
		updated, cmd := r.pages[r.current].Update(msg)
		r.pages[r.current] = updated
		return cmd
	}
}

func (r *root) handleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	r.wWidth = msg.Width
	r.wHeight = msg.Height

	// also pass to the current model
	updated, cmd := r.pages[r.current].Update(msg)
	r.pages[r.current] = updated
	return cmd
}

func (r *root) movePage(pageID pages.PageID) tea.Cmd {
	r.current = pageID

	var cmds []tea.Cmd
	if !r.loaded[r.current] {
		cmd := r.pages[r.current].Init()
		cmds = append(cmds, cmd)
		r.loaded[r.current] = true
	}

	return tea.Batch(cmds...)
}

func (r *root) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, components.ShowHeader(r.wWidth), r.pages[r.current].View())
}
