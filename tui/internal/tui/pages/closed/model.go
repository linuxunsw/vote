package closed

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

const (
	closedMessage = "voting/nominations are currently closed, please come back later!"
	exitMessage   = "exit with ctrl+c"
)

type submitModel struct {
	logger *log.Logger

	// Maximum size allowed for content
	cWidth  int
	cHeight int
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
	}
	return m, nil
}

// Displays message which changes depending on whether the submission was successful
// or there was a server error
func (m *submitModel) View() string {
	exit := styles.ExitMessageStyle.Render(exitMessage)

	message := fmt.Sprintf("%s\n\n%s", closedMessage, exit)
	return styles.SubmitText(m.cHeight, m.cWidth).Render(message)
}
