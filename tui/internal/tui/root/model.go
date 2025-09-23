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

	"github.com/linuxunsw/vote/tui/internal/sdk"
	"github.com/linuxunsw/vote/tui/internal/tui/components"
	"github.com/linuxunsw/vote/tui/internal/tui/keys"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/auth"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/authcode"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/closed"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/nominationform"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/nominationsubmit"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/voting"
	"github.com/linuxunsw/vote/tui/internal/tui/pages/votingsubmit"
)

const helpHeight = 1

// Form data
type formData struct {
	zID        string
	submission sdk.Submission
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

	client          *sdk.ClientWithIP
	isAuthenticated bool
	error           error

	loadingSpinner spinner.Model
	loading        bool

	needsSizeUpdate bool

	data formData
}

func New(user, ip string) tea.Model {
	keyMap := keys.DefaultKeyMap()

	logger := createLogger(user)

	// Load each page
	pageMap := map[pages.PageID]tea.Model{
		pages.Auth:             auth.New(logger),
		pages.AuthCode:         authcode.New(logger),
		pages.NominationForm:   nominationform.New(logger),
		pages.NominationSubmit: nominationsubmit.New(logger),
		pages.Closed:           closed.New(logger),
		pages.VotingSubmit:     votingsubmit.New(logger),
	}

	// Create spinner
	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Ellipsis

	client := sdk.CreateClient(logger, ip)

	model := &rootModel{
		log:             logger,
		keyMap:          keyMap,
		pages:           pageMap,
		client:          client,
		isAuthenticated: false,
		loaded:          make(map[pages.PageID]bool),
		loading:         false,
		loadingSpinner:  loadingSpinner,
		current:         pages.Auth,
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
	case sdk.GenerateOTPMsg:
		m.log.Debug("GenerateOTPMsg", "zID", msg.ZID)

		m.data.zID = msg.ZID
		m.loading = true
		m.error = nil

		return m, sdk.GenerateOTPCmd(m.client, m.data.zID)
	case sdk.SubmitOTPMsg:
		m.log.Debug("SubmitOTPMsg", "OTP", msg.OTP)

		m.loading = true
		m.error = nil

		return m, sdk.SubmitOTPCmd(m.client, m.data.zID, msg.OTP)
	case sdk.Submission:
		m.log.Debug(
			"Submission",
			"zid", m.data.zID,
			"name", msg.Name,
			"roles", msg.Roles,
		)

		m.loading = true
		m.error = nil

		m.data.submission = msg

		return m, sdk.SubmitNominationCmd(m.client, m.data.submission)
	case sdk.GenerateOTPSuccessMsg:
		m.log.Debug("GenerateOTPSuccessMsg")
		m.loading = false
		return m, messages.SendPageChange(pages.AuthCode)
	case sdk.SubmitOTPSuccessMsg:
		m.log.Debug("SubmitOTPSuccessMsg")

		m.isAuthenticated = true

		return m, sdk.GetElectionStateCmd(m.client)
	case sdk.GetElectionStateSuccessMsg:
		m.log.Debug("GetElectionStateSuccessMsg", "state", msg.State)

		m.loading = false

		// change form depending on state
		if msg.State == string(sdk.GetElectionStateResponseBodyStateNOMINATIONSOPEN) {
			return m, messages.SendPageChange(pages.NominationForm)
		} else if msg.State == string(sdk.GetElectionStateResponseBodyStateVOTINGOPEN) {
			m.loading = true
			return m, sdk.GetBallotCmd(m.client)
		} else {
			return m, messages.SendPageChange(pages.Closed)
		}
	case sdk.GetBallotSuccessMsg:
		m.loading = false
		m.pages[pages.VotingForm] = voting.New(m.log, *msg.Ballot)
		return m, messages.SendPageChange(pages.VotingForm)
	case sdk.SubmitVoteMsg:
		m.loading = true
		return m, sdk.SubmitVoteCmd(m.client, msg.Votes)
	case sdk.SubmitVoteSuccessMsg:
		m.loading = false
		return m, tea.Sequence(
			messages.SendPageChange(pages.VotingSubmit),
			sdk.SendPublicSubmitFormResult(msg.RefCode, nil),
		)
	case sdk.SubmitNominationSuccessMsg:
		m.log.Debug("SubmitNominationSuccessMsg", "refCode", msg.RefCode)
		m.loading = false
		return m, tea.Sequence(
			messages.SendPageChange(pages.NominationSubmit),
			sdk.SendPublicSubmitFormResult(msg.RefCode, nil),
		)
	case sdk.ServerErrMsg:
		// INFO: we must handle the ServerErrMsg here as well as in the current
		// model as failing to disable loading will prevent the current page
		// from displaying (thus not displaying the error)
		m.loading = false
		m.error = msg.Error
		m.needsSizeUpdate = true

		// Reset everything if unauthorised (means cookie has expired)
		if msg.Error == sdk.ErrUnauthorised {
			// reset pages
			m.pages[pages.Auth] = auth.New(m.log)
			m.pages[pages.AuthCode] = authcode.New(m.log)
			m.pages[pages.NominationForm] = nominationform.New(m.log)
			m.loaded[pages.Auth] = false
			m.loaded[pages.AuthCode] = false
			m.loaded[pages.NominationForm] = false
			m.loaded[pages.VotingForm] = false

			m.isAuthenticated = false
			return m, messages.SendPageChange(pages.Auth)

		}
	case spinner.TickMsg:
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd
	}

	// Pass any remaining messages to the current model
	updated, cmd := m.pages[m.current].Update(msg)
	m.pages[m.current] = updated

	// INFO: this allows us to trigger a content size change to account
	// for the error message being included in the footer while allowing
	// the individual page model and the root model to both handle
	// a ServerErrMsg
	if m.needsSizeUpdate {
		m.needsSizeUpdate = false
		w, h := m.findContentSize()
		cmds = append(cmds, messages.SendPageContentSize(w, h))
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Displays the header, current model's content and a footer if the user is authenticated
func (m *rootModel) View() string {
	var footer string
	if m.error != nil {
		footer = components.ShowErrorFooter(m.error, m.wWidth)
	}
	if m.isAuthenticated {
		footer += components.ShowFooter(m.data.zID, m.wWidth)
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
	if m.error != nil {
		footerHeight += lipgloss.Height(components.ShowErrorFooter(m.error, m.wWidth))
	}
	if m.isAuthenticated {
		footerHeight += lipgloss.Height(components.ShowFooter(m.data.zID, m.wWidth))
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

func createLogger(user string) *log.Logger {
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

	return logger
}
