package styles

import "github.com/charmbracelet/lipgloss"

var FooterStyle = lipgloss.NewStyle().
	AlignHorizontal(lipgloss.Center).
	Italic(true).
	Faint(true).
	PaddingTop(1).
	PaddingBottom(1)

var ErrorFooterStyle = lipgloss.NewStyle().
	AlignHorizontal(lipgloss.Center).
	Italic(true).
	PaddingTop(1).
	PaddingBottom(1).
	Foreground(lipgloss.ANSIColor(9))
