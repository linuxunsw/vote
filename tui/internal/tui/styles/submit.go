package styles

import "github.com/charmbracelet/lipgloss"

var SubmitTextStyle = lipgloss.NewStyle().
	Padding(0, 20, 0, 20).
	Align(lipgloss.Center).
	AlignHorizontal(lipgloss.Center).
	AlignVertical(lipgloss.Center)

var ExitMessageStyle = lipgloss.NewStyle().
	Faint(true)

func SubmitText(height, width int) lipgloss.Style {
	return SubmitTextStyle.Height(height).Width(width)
}
