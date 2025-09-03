package submit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

const (
	successMessage = "your nomination was submitted successfully! \n\nyour reference code is %s. a copy of your nomination has been submitted to your provided email address."
	errorMessage   = "something went wrong :( \n\nplease try again later. if you are still encountering issues, please contact a society executive on discord with the following reference code: %s."
	exitMessage    = "exit with ctrl+c"
)

type submitModel struct {
	logger *log.Logger

	// Maximum size allowed for content
	cWidth  int
	cHeight int

	// Submission details
	refCode string
	error   error
}

func New(logger *log.Logger) tea.Model {
	model := &submitModel{
		logger: logger,
	}

	return model
}

func (m *submitModel) Init() tea.Cmd {
	return nil
}

func (m *submitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.PageContentSizeMsg:
		log.Debug("PageContentSizeMsg", "height", msg.Height, "width", msg.Width)
		m.cHeight = msg.Height
		m.cWidth = msg.Width

		return m, nil
	case messages.PublicSubmitFormResultMsg:
		log.Debug("PublicSubmitFormResultMsg", "refCode", msg.RefCode, "error", msg.Error)
		m.refCode = msg.RefCode
		m.error = msg.Error

		return m, nil
	}

	return m, nil
}

// Displays message which changes depending on whether the submission was successful
// or there was a server error
func (m *submitModel) View() string {
	var content string

	if m.error != nil {
		content = fmt.Sprintf(errorMessage, m.refCode)
	} else {
		content = fmt.Sprintf(successMessage, m.refCode)
	}

	exit := styles.ExitMessageStyle.Render(exitMessage)

	message := fmt.Sprintf("%s\n\n%s", content, exit)
	return styles.SubmitText(m.cHeight, m.cWidth).Render(message)
}
