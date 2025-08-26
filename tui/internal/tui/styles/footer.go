package styles

import "github.com/charmbracelet/lipgloss"

var FooterStyle = lipgloss.NewStyle().
	AlignHorizontal(lipgloss.Center).
	Italic(true).
	Faint(true).
	PaddingTop(1).
	PaddingBottom(1)
