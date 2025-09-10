// Contains the tea.Cmds which are used to interact with the API
package client

import (
	"context"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
)

// TODO: replace sleeps with api calls
func GenerateOTPCmd(client *http.Client, zID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := generateOTP(client, ctx, zID)

		return messages.RequestOTPResultMsg{
			Error: err,
		}
	}
}

func SubmitOTPCmd(client *http.Client, zID string, otp string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := submitOTP(client, ctx, zID, otp)

		return messages.VerifyOTPResultMsg{
			Error: err,
		}
	}

}

func SubmitNominationCmd(client *http.Client, data messages.Submission) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		respID, err := submitNomination(client, ctx, data)

		return messages.SubmitFormResultMsg{
			Error:   err,
			RefCode: respID,
		}
	}
}
