package root

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"

	"github.com/linuxunsw/vote/tui/internal/client"
	"github.com/linuxunsw/vote/tui/internal/tui/components"
	"github.com/linuxunsw/vote/tui/internal/tui/keys"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/auth"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/authcode"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/form"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/submit"
)

const helpHeight = 1

// Form data
type formData struct {
	zID        string
	submission messages.Submission
}

type rootModel struct {
	log *log.Logger

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

func New(user string) tea.Model {
	keyMap := keys.DefaultKeyMap()

	// Create logger
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(true)
	prefix := fmt.Sprintf("app (%s)", user)
	logger.SetPrefix(prefix)

	// Only debug logs if debug set in config
	logDebug := viper.GetBool("tui.debug")
	if logDebug {
		logger.SetLevel(log.DebugLevel)
	}

	// Load each page
	// TODO: add submission result page
	pageMap := map[pages.PageID]tea.Model{
		pages.PageAuth:     auth.New(logger),
		pages.PageAuthCode: authcode.New(logger),
		pages.PageForm:     form.New(logger),
		pages.PageSubmit:   submit.New(logger),
	}

	// Create spinner
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

	// Handle messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m, m.handleWindowSizeMsg(msg)
	case messages.PageChangeMsg:
		m.log.Debug("PageChangeMsg", "msg", msg)
		return m, m.movePage(msg.ID)
	case messages.RequestOTPMsg:
		m.log.Debug("RequestOTPMsg", "zID", msg.ZID)
		m.data.zID = msg.ZID

		m.loading = true

		return m, client.RequestOTP(m.data.zID)
	case messages.RequestOTPResultMsg:
		m.log.Debug("RequestOTPResultMsg", "error", msg.Error)
		// only change page if we didn't error
		m.loading = false

		if msg.Error == nil {
			return m, messages.SendPageChange(pages.PageAuthCode)
		}
	case messages.VerifyOTPMsg:
		m.log.Debug("VerifyOTPMsg", "OTP", msg.OTP)

		m.loading = true

		return m, client.VerifyOTP(msg.OTP)
	case messages.VerifyOTPResultMsg:
		m.log.Debug("VerifyOTPResultMsg", "error", msg.Error)

		m.loading = false

		if msg.Error == nil {
			m.isAuthenticated = true
			m.log.Info("User authenticated", "zID", m.data.zID)
			return m, messages.SendPageChange(pages.PageForm)
		}
	case messages.Submission:
		m.log.Debug(
			"Submission",
			"zid", m.data.zID,
			"name", msg.Name,
			"roles", msg.Roles,
		)

		m.loading = true

		m.data.submission = msg

		return m, client.SubmitForm(msg)
	case messages.SubmitFormResultMsg:
		m.log.Debug("SubmitFormResultMsg", "refCode", msg.RefCode, "error", msg.Error)

		m.loading = false

		if msg.Error == nil {
			m.log.Info("Form submitted", "zID", m.data.zID)
			return m, tea.Sequence(
				messages.SendPageChange(pages.PageSubmit),
				messages.SendPublicSubmitFormResult(msg.RefCode, msg.Error),
			)
		}
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
	var footer string
	if m.isAuthenticated {
		footer = components.ShowFooter(m.data.zID, m.wWidth)
	}

	var content string
	var loadingSpinner string

	// Content changes depending on whether we are loading
	if m.loading {
		// Find the current content size to allow us to center the
		// loading spinner
		w, h := m.findContentSize()

		// Align spinner horizontally, remove spaces to the right
		// to prevent text from re-centering as the spinner increases
		// in width
		style := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(w).
			Height(h)

		loadingSpinner = style.Render(components.GetPageMsg(m.current))
		content = strings.TrimRightFunc(loadingSpinner, unicode.IsSpace) + m.loadingSpinner.View()

		// Align vertically
		content = lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			Height(h).
			Render(content)

	} else {
		content = m.pages[m.current].View()
	}

	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Top,
		components.ShowHeader(m.wWidth),
		content,
		footer,
	)
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

	// Change root model size
	m.wWidth = msg.Width
	m.wHeight = msg.Height

	// Recalculate the content size and send it to the currently shown model
	w, h := m.findContentSize()
	m.cWidth, m.cHeight = w, h
	cmd = messages.SendPageContentSize(w, h)
	m.log.Debug("SendPageContentSize", "height", h, "width", w)

	return cmd
}

func (m *rootModel) findContentSize() (w int, h int) {
	// If the user is authenticated, the footer is shown, so we should include it
	// in the height calculation only when the user is authenticated
	footerHeight := 0
	if m.isAuthenticated {
		footerHeight = lipgloss.Height(components.ShowFooter(m.data.zID, m.wWidth))
	}

	headerHeight := lipgloss.Height(components.ShowHeader(m.wWidth))

	// huh adds space onto the form for the "help" view, which we
	// should also include in size calculation
	h = m.wHeight - footerHeight - headerHeight - helpHeight
	w = m.wWidth

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

	// Ensure that the content will be the correct width - may have
	// changed due to footer size changing, etc.
	w, h := m.findContentSize()
	cmd := messages.SendPageContentSize(w, h)

	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}
