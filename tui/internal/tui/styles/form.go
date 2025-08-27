package styles

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var FormStyle = lipgloss.NewStyle().
	Align(lipgloss.Left).
	MarginRight(2).
	MarginLeft(2)

func FormTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("8"))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("11")).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color("6"))
	t.Focused.Directory = t.Focused.Directory.Foreground(lipgloss.Color("6"))
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color("8"))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lipgloss.Color("9"))
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lipgloss.Color("9")).Italic(true)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color("11"))
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(lipgloss.Color("11"))
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(lipgloss.Color("11"))
	t.Focused.Option = t.Focused.Option.Foreground(lipgloss.Color("11"))
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color("15"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color("15")).Bold(true)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color("11"))
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(lipgloss.Color("7"))
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("5"))
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0"))

	t.Focused.TextInput.Cursor.Foreground(lipgloss.Color("5"))
	t.Focused.TextInput.Placeholder.Foreground(lipgloss.Color("8"))
	t.Focused.TextInput.Prompt.Foreground(lipgloss.Color("3"))

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(lipgloss.Color("8"))
	t.Blurred.Title = t.Blurred.NoteTitle.Foreground(lipgloss.Color("8"))

	t.Blurred.TextInput.Prompt = t.Blurred.TextInput.Prompt.Foreground(lipgloss.Color("8"))
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(lipgloss.Color("7"))

	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
