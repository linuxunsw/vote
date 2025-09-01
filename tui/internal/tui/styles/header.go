package styles

import "github.com/charmbracelet/lipgloss"

var HeaderStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingTop(1).
	PaddingBottom(1)

func Header(width int) lipgloss.Style {
	return HeaderStyle.Width(width)
}
