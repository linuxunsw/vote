package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/pages"
)

// Triggers page change to given ID
type PageChangeMsg struct {
	ID pages.PageID
}

// Sends maximum size of page content (window size without header and footer)
type PageContentSizeMsg struct {
	Width  int
	Height int
}

func SendPageChange(id pages.PageID) tea.Cmd {
	msg := PageChangeMsg{
		ID: id,
	}

	return func() tea.Msg { return msg }
}

func SendPageContentSize(width int, height int) tea.Cmd {
	msg := PageContentSizeMsg{
		Width:  width,
		Height: height,
	}

	return func() tea.Msg { return msg }
}
