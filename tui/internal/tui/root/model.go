package root

import (
	"os"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/linuxunsw/vote/tui/internal/tui/components"
	"github.com/linuxunsw/vote/tui/internal/tui/keys"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/auth"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/authcode"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/form"
)

// Form data
type formData struct {
	zID        string
	submission messages.Submission
}

type rootModel struct {
	wWidth  int
	wHeight int
	keyMap  keys.KeyMap

	pages  map[pages.PageID]tea.Model
	loaded map[pages.PageID]bool

	current pages.PageID

	isAuthenticated bool

	data formData
}

func New() tea.Model {
	keyMap := keys.DefaultKeyMap()

	// Load each page
	pageMap := map[pages.PageID]tea.Model{
		pages.PageAuth:     auth.New(),
		pages.PageAuthCode: authcode.New(),
		pages.PageForm:     form.New(),
	}

	model := &rootModel{
		keyMap:          keyMap,
		pages:           pageMap,
		isAuthenticated: false,
		loaded:          make(map[pages.PageID]bool),
		current:         pages.PageAuth,
	}

	return model

}
func (m *rootModel) Init() tea.Cmd {
	m.loaded[m.current] = true

	windowTitle := os.Getenv("EVENT_NAME")

	return tea.Batch(
		m.pages[m.current].Init(),
		tea.SetWindowTitle(windowTitle),
	)
}

// Handles all messages recieved by the app
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m, m.handleWindowSizeMsg(msg)
	case messages.PageChangeMsg:
		return m, m.movePage(msg.ID)
	case messages.AuthMsg:
		m.data.zID = msg.ZID
		// TODO: send req to api with zid
		return m, cmd
	case messages.CheckOTPMsg:
		// TODO: send req to api with otp, set authenticated to true on this condition
		m.isAuthenticated = true
		return m, tea.Batch(messages.SendIsAuthenticated(nil))
	case messages.Submission:
		m.data.submission = msg
		return m, tea.Batch(messages.SendSubmissionResult("insert code here", nil))
	}

	// Pass any remaining messages to the current model
	updated, cmd := m.pages[m.current].Update(msg)
	m.pages[m.current] = updated

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Displays the header, current model's content and a footer if the user is authenticated
func (m *rootModel) View() string {
	// Create footer with the user's zID if authenticated
	var footer string
	if m.isAuthenticated {
		footer = components.ShowFooter(m.data.zID, m.wWidth)
	}

	return lipgloss.JoinVertical(lipgloss.Top, components.ShowHeader(m.wWidth), m.pages[m.current].View(), footer)
}

// Handles any global keybinds as defined in `m.keyMap`, then passes down
// to the current 'page' (model)
func (m *rootModel) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return tea.Quit
	default:
		updated, cmd := m.pages[m.current].Update(msg)
		m.pages[m.current] = updated
		return cmd
	}
}

// Sets the window size in the rootModel, then passes down to the current
// 'page' (model)
func (m *rootModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	m.wWidth = msg.Width
	m.wHeight = msg.Height

	// Pass the WindowSizeMsg to the current model
	updated, cmd := m.pages[m.current].Update(msg)
	m.pages[m.current] = updated
	return cmd
}

// Switches current page given a pageID
func (m *rootModel) movePage(pageID pages.PageID) tea.Cmd {
	m.current = pageID

	// If we haven't initialised the model before, do so
	var cmds []tea.Cmd
	if !m.loaded[m.current] {
		cmd := m.pages[m.current].Init()
		cmds = append(cmds, cmd)
		m.loaded[m.current] = true
	}

	return tea.Batch(cmds...)
}
