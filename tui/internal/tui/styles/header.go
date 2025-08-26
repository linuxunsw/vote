package styles

import "github.com/charmbracelet/lipgloss"

var HeaderStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	PaddingTop(1).
	PaddingBottom(1).
	MarginBottom(1)

func Header(width int) lipgloss.Style {
	return HeaderStyle.Width(width)
}
