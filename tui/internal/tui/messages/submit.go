package messages

import tea "github.com/charmbracelet/bubbletea"

// Sends data to root model to trigger submit
type Submission struct {
	Name    string
	Email   string
	Discord string

	Roles     []string
	Statement string
	Url       string
}

// Tells the UI result of submission - if an error is
// present the submission was unsuccessful
type FormSubmissionResultMsg struct {
	RefCode string
	Error   error
}

func SendSubmission(data Submission) tea.Cmd {
	return func() tea.Msg { return data }
}

func SendSubmissionResult(refCode string, result error) tea.Cmd {
	msg := FormSubmissionResultMsg{
		RefCode: refCode,
		Error:   result,
	}

	return func() tea.Msg { return msg }
}
