// TODO: move between pages

package root

import (
	"os"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/keys"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/auth"
)

// Form data
type formData struct {
	zID string

	// TODO: add submission to remaining fields

	// name    string
	// email   string
	// discord string
	// roles   []bool
	// statement string
	// url       string
}

type root struct {
	wWidth  int
	wHeight int
	keyMap  keys.KeyMap

	pages    map[PageID]tea.Model
	loaded   map[PageID]bool
	current  PageID
	previous PageID

	data formData
}

func New() tea.Model {
	keyMap := keys.DefaultKeyMap()
	pages := map[PageID]tea.Model{
		PageAuth: auth.New(),
	}

	model := &root{
		keyMap:  keyMap,
		pages:   pages,
		loaded:  map[PageID]bool{},
		current: PageAuth,
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
		return r, r.handleKeypressMsg(msg)
	}

	// INFO: this repetition is temporary to make the form work
	// as the form needs more messages than just the single keypress
	// msg passed to the model - so we just temporarily send through
	// the remaining 'not global' messages to the currently shown model
	updated, cmd := r.pages[r.current].Update(msg)
	r.pages[r.current] = updated

	cmds = append(cmds, cmd)
	return r, tea.Batch(cmds...)
}

func (r *root) handleKeypressMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, r.keyMap.Quit):
		return tea.Quit
	default:
		updated, cmd := r.pages[r.current].Update(msg)
		r.pages[r.current] = updated
		return cmd
	}
}

func (r *root) View() string {
	return r.pages[r.current].View()
}
