package root

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"

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
	user string
	log  *log.Logger

	wWidth  int
	wHeight int
	cWidth  int
	cHeight int
	keyMap  keys.KeyMap

	pages  map[pages.PageID]tea.Model
	loaded map[pages.PageID]bool

	current pages.PageID

	isAuthenticated bool

	loadingSpinner spinner.Model
	loading        bool

	data formData
}

// INFO: testing stuff for "fake" api call so we can
// test spinner - will be replaced when api complete
type testCompleteMsg struct {
	err error
}

// Simulate API call
// API calls will need to be done in tea.Cmds like so
func testAPICall() tea.Cmd {
	return func() tea.Msg {
		// Simulate API call delay
		time.Sleep(3 * time.Second)
		return testCompleteMsg{err: nil}
	}
}

func New(user string) tea.Model {
	keyMap := keys.DefaultKeyMap()

	// Create logger
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(true)
	prefix := fmt.Sprintf("app (%s)", user)
	logger.SetPrefix(prefix)

	logDebug := viper.GetBool("tui.debug")
	if logDebug {
		logger.SetLevel(log.DebugLevel)
	}

	// Load each page
	pageMap := map[pages.PageID]tea.Model{
		pages.PageAuth:     auth.New(logger),
		pages.PageAuthCode: authcode.New(logger),
		pages.PageForm:     form.New(),
	}

	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Ellipsis

	model := &rootModel{
		log:             logger,
		keyMap:          keyMap,
		pages:           pageMap,
		isAuthenticated: false,
		loaded:          make(map[pages.PageID]bool),
		loading:         false,
		loadingSpinner:  loadingSpinner,
		current:         pages.PageAuth,
	}

	model.log.Info("Starting app...")

	return model

}
func (m *rootModel) Init() tea.Cmd {
	m.loaded[m.current] = true

	return tea.Batch(
		m.pages[m.current].Init(),
		m.loadingSpinner.Tick,
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
		m.log.Debug("PageChangeMsg", "msg", msg)
		return m, m.movePage(msg.ID)
	case messages.AuthMsg:
		m.log.Debug("AuthMsg", "msg", msg)
		m.data.zID = msg.ZID
		m.loading = true

		return m, testAPICall()
	case messages.CheckOTPMsg:
		m.log.Debug("CheckOTPMsg", "msg", msg)
		// TODO: send req to api with otp, set authenticated to true on this condition
		// log when successfully/unsuccessfully authenticated
		m.isAuthenticated = true
		// log when someone successfully authenticates
		return m, tea.Batch(messages.SendIsAuthenticated(nil))
	case testCompleteMsg:
		m.log.Debug("testCompleteMsg", "msg", msg)
		m.loading = false
		return m, messages.SendPageChange(pages.PageAuthCode)
	case messages.Submission:
		m.log.Debug("SubmissionMsg", "msg", msg)
		m.data.submission = msg
		return m, tea.Quit
	case spinner.TickMsg:
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd
	}

	// Pass any remaining messages to the current model
	updated, cmd := m.pages[m.current].Update(msg)
	m.pages[m.current] = updated

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Displays the header, current model's content and a footer if the user is authenticated
func (m *rootModel) View() string {
	// footer only when authed
	var footer string
	if m.isAuthenticated {
		footer = components.ShowFooter(m.data.zID, m.wWidth)
	}

	var content string
	var loadingSpinner string

	if m.loading {
		w, h := m.findContentSize()

		style := lipgloss.NewStyle().Align(lipgloss.Center).Width(w).Height(h)
		loadingSpinner = style.Render("requesting otp")
		content = strings.TrimRightFunc(loadingSpinner, unicode.IsSpace) + m.loadingSpinner.View()
		content = lipgloss.NewStyle().AlignVertical(lipgloss.Center).Height(h).Render(content)

	} else {
		content = m.pages[m.current].View()
	}

	return lipgloss.JoinVertical(lipgloss.Top, components.ShowHeader(m.wWidth), content, footer)
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
	m.log.Debug("WindowSizeMsg", "msg", msg)
	var cmd tea.Cmd

	m.wWidth = msg.Width
	m.wHeight = msg.Height

	w, h := m.findContentSize()

	cmd = messages.SendPageContentSize(w, h)
	m.log.Debug("SendPageContentSize", "height", h, "width", w)

	return cmd
}

func (m *rootModel) findContentSize() (w int, h int) {
	footerHeight := 0
	if m.isAuthenticated {
		footerHeight = lipgloss.Height(components.ShowFooter(m.data.zID, m.wWidth))
	}
	headerHeight := lipgloss.Height(components.ShowHeader(m.wWidth))

	h = m.wHeight - footerHeight - headerHeight - 1
	w = m.wWidth - 4

	return
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

	w, h := m.findContentSize()
	cmd := messages.SendPageContentSize(w, h)
	log.Debug("SendPageContentSize", "height", m.cHeight, "width", m.cWidth)

	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}
