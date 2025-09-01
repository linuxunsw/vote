// Contains the tea.Cmds which are used to interact with the API
package client

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
)

// TODO: replace sleeps with api calls
func RequestOTP(zid string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)

		return messages.RequestOTPResultMsg{
			Error: nil,
		}
	}
}

func VerifyOTP(otp string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)

		return messages.VerifyOTPResultMsg{
			Error: nil,
		}
	}

}

func SubmitForm(data messages.Submission) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)

		return messages.SubmitFormResultMsg{
			Error:   nil,
			RefCode: "blank",
		}
	}
}
