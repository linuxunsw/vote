package messages

import tea "github.com/charmbracelet/bubbletea"

// Sends the user's form input (their zID) to the root model
type AuthMsg struct {
	ZID string
}

// Triggers the root model to check if otp is valid
type CheckOTPMsg struct {
	OTP string
}

// Determines authenticated state - if there was an error we are
// not authenticated
type IsAuthenticatedMsg struct {
	Error error
}

type ZIDSubmittedMsg struct {
	Error error
}

// Sends a message containing the user's zID for authentication
func SendAuth(zID string) tea.Cmd {
	msg := AuthMsg{
		ZID: zID,
	}

	return func() tea.Msg { return msg }
}

// Sends a message containing the user's OTP input
func SendCheckOTP(otp string) tea.Cmd {
	msg := CheckOTPMsg{
		OTP: otp,
	}

	return func() tea.Msg { return msg }
}

// Sends a message showing if the user was authenticated or
// if there was an error in authentication
func SendIsAuthenticated(err error) tea.Cmd {
	msg := IsAuthenticatedMsg{
		Error: err,
	}

	return func() tea.Msg { return msg }
}

func SendZIDSumbitted(err error) tea.Cmd {
	msg := ZIDSubmittedMsg{
		Error: err,
	}

	return func() tea.Msg { return msg }
}
