// Messages to interact with the client
package messages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ServerErrMsg struct {
	StatusCode int
	RespID     string
	Error      error
}

type ResetFormMsg struct{}

type GenerateOTPMsg struct {
	ZID string
}
type GenerateOTPSuccessMsg struct{}

type SubmitOTPMsg struct {
	OTP string
}
type SubmitOTPSuccessMsg struct{}

// Submission is sent as a message to the root model
type Submission struct {
	Name      string   `json:"candidate_name"`
	Email     string   `json:"contact_email"`
	Discord   string   `json:"discord_username"`
	Roles     []string `json:"executive_roles"`
	Statement string   `json:"candidate_statement"`
	Url       string   `json:"url,omitempty"`
}

type SubmitNominationSuccessMsg struct {
	RefCode string
}

// Message sent to the submission page to show result data
type PublicSubmitFormResultMsg struct {
	RefCode string
	Error   error
}

// Sends a request to the root model to generate an OTP
func SendGenerateOTP(zid string) tea.Cmd {
	msg := GenerateOTPMsg{
		ZID: zid,
	}

	return func() tea.Msg { return msg }
}

// Sends request to the root model to verify the given OTP
func SendSubmitOTP(otp string) tea.Cmd {
	msg := SubmitOTPMsg{
		OTP: otp,
	}

	return func() tea.Msg { return msg }
}

// Sends request to the root model to submit the form
func SendNomination(data Submission) tea.Cmd {
	return func() tea.Msg { return data }
}

func SendPublicSubmitFormResult(refCode string, error error) tea.Cmd {
	msg := PublicSubmitFormResultMsg{
		RefCode: refCode,
		Error:   error,
	}

	return func() tea.Msg { return msg }
}

func SendResetForm() tea.Cmd {
	msg := ResetFormMsg{}

	return func() tea.Msg { return msg }
}
