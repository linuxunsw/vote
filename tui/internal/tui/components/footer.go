package components

import (
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

func ShowFooter(zid string, width int) string {
	content := ""
	if zid != "" {
		content = "currently signed in as: " + zid
	}
	return styles.FooterStyle.Width(width).Render(content)
}

func ShowErrorFooter(error error, width int) string {
	content := ""
	if error != nil {
		content = error.Error()
	}
	return styles.ErrorFooterStyle.Width(width).Render(content)
}
