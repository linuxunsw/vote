package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
)

// Triggers page change to given ID
type PageChangeMsg struct {
	ID pages.PageID
}

func SendPageChange(id pages.PageID) tea.Cmd {
	msg := PageChangeMsg{
		ID: id,
	}

	return func() tea.Msg { return msg }
}
