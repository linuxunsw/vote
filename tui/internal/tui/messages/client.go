// Messages to interact with the client
package messages

import tea "github.com/charmbracelet/bubbletea"

type RequestOTPMsg struct {
	ZID string
}

type RequestOTPResultMsg struct {
	Error error
}

type VerifyOTPMsg struct {
	OTP string
}

// TODO: incl session token?
type VerifyOTPResultMsg struct {
	Error error
}

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
type SubmitFormResultMsg struct {
	RefCode string
	Error   error
}

func SendRequestOTP(zid string) tea.Cmd {
	msg := RequestOTPMsg{
		ZID: zid,
	}

	return func() tea.Msg { return msg }
}

func SendRequestOTPResult(err error) tea.Cmd {
	msg := RequestOTPResultMsg{
		Error: err,
	}

	return func() tea.Msg { return msg }
}

func SendVerifyOTP(otp string) tea.Cmd {
	msg := VerifyOTPMsg{
		OTP: otp,
	}

	return func() tea.Msg { return msg }
}

func SendVerifyOTPResult(err error) tea.Cmd {
	msg := VerifyOTPResultMsg{
		Error: err,
	}

	return func() tea.Msg { return msg }
}

func SendSubmission(data Submission) tea.Cmd {
	return func() tea.Msg { return data }
}

func SendSubmitFormResult(refCode string, result error) tea.Cmd {
	msg := SubmitFormResultMsg{
		RefCode: refCode,
		Error:   result,
	}

	return func() tea.Msg { return msg }
}
