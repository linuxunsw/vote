package styles

import "github.com/charmbracelet/lipgloss"

var HeaderStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("11")).
	PaddingTop(1).
	PaddingBottom(2).
	Bold(true)

func Header(width int) lipgloss.Style {
	return HeaderStyle.Width(width)
}
