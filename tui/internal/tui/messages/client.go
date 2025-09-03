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

type VerifyOTPResultMsg struct {
	Error error
}

// Submission is sent as a message to the root model
type Submission struct {
	Name    string
	Email   string
	Discord string

	Roles     []string
	Statement string
	Url       string
}

type SubmitFormResultMsg struct {
	RefCode string
	Error   error
}

// Message sent to the submission page to show result data
type PublicSubmitFormResultMsg struct {
	RefCode string
	Error   error
}

// Sends a request to the root model to generate an OTP
func SendRequestOTP(zid string) tea.Cmd {
	msg := RequestOTPMsg{
		ZID: zid,
	}

	return func() tea.Msg { return msg }
}

// Sends request to the root model to verify the given OTP
func SendVerifyOTP(otp string) tea.Cmd {
	msg := VerifyOTPMsg{
		OTP: otp,
	}

	return func() tea.Msg { return msg }
}

// Sends request to the root model to submit the form
func SendSubmission(data Submission) tea.Cmd {
	return func() tea.Msg { return data }
}

func SendPublicSubmitFormResult(refCode string, error error) tea.Cmd {
	msg := PublicSubmitFormResultMsg{
		RefCode: refCode,
		Error:   error,
	}

	return func() tea.Msg { return msg }
}
