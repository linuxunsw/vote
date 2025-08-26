package messages

import "github.com/linuxunsw/vote/tui/internal/tui/pages"

// Triggers page change to given ID
type PageChangeMsg struct {
	ID pages.PageID
}
