// Contains the tea.Cmds which are used to interact with the API

package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/oapi-codegen/runtime/types"
)

var (
	ErrUnauthorised error = errors.New("your session has expired, please log in again")
)

func createIPRequestEditor(ip string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Real-IP", ip)
		return nil
	}
}

// Sends request to generate OTP, sends response back to root model as
// ServerErrMsg or a success message
func GenerateOTPCmd(c *ClientWithIP, zID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		body := GenerateOTPInputBody{
			Zid: zID,
		}

		resp, err := c.Client.GenerateOtpWithResponse(ctx, body, createIPRequestEditor(c.IP))
		if err != nil {
			return messages.ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		if resp.StatusCode() != http.StatusNoContent {
			var respID string
			if resp != nil && resp.HTTPResponse != nil {
				respID = resp.HTTPResponse.Header.Get("X-Request-ID")
			}

			var err error
			if resp != nil && resp.ApplicationproblemJSONDefault != nil {
				err = buildError(*resp.ApplicationproblemJSONDefault)
			}

			return messages.ServerErrMsg{
				RespID: respID,
				Error:  err,
			}
		}

		return messages.GenerateOTPSuccessMsg{}
	}
}

// Sends request to submit OTP, sends response back to root model as
// ServerErrMsg or a success message
func SubmitOTPCmd(c *ClientWithIP, zID string, otp string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		body := SubmitOTPInputBody{
			Zid: zID,
			Otp: otp,
		}

		resp, err := c.Client.SubmitOtpWithResponse(ctx, body, createIPRequestEditor(c.IP))
		if err != nil {
			return messages.ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		if resp.StatusCode() != http.StatusOK {
			respID := resp.HTTPResponse.Header.Get("X-Request-ID")
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return messages.ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return messages.SubmitOTPSuccessMsg{}
	}

}

// Sends request to submit a nomination, sends response back to root model as
// ServerErrMsg or a success message
func SubmitNominationCmd(c *ClientWithIP, data messages.Submission) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var execRoles []SubmitNominationExecutiveRoles
		for _, role := range data.Roles {
			execRoles = append(execRoles, SubmitNominationExecutiveRoles(role))
		}

		body := SubmitNomination{
			CandidateName:      data.Name,
			CandidateStatement: data.Statement,
			ContactEmail:       types.Email(data.Email),
			DiscordUsername:    data.Discord,
			ExecutiveRoles:     &execRoles,
		}

		// Prevents submitting an empty url (messes with server validation)
		if data.Url != "" {
			body.Url = &data.Url
		}

		resp, err := c.Client.SubmitNominationWithResponse(ctx, body, createIPRequestEditor(c.IP))
		if err != nil {
			return messages.ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		// Add request ID as a reference code
		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusUnauthorized {
			return messages.ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrUnauthorised,
			}
		}
		if resp.StatusCode() != http.StatusNoContent {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return messages.ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return messages.SubmitNominationSuccessMsg{
			RefCode: respID,
		}

	}
}

// Build error message from an error model
func buildError(em ErrorModel) error {
	var sb strings.Builder

	// Copy title
	if em.Title != nil && *em.Title != "" {
		sb.WriteString(*em.Title)
	} else {
		sb.WriteString("Error")
	}

	if em.Detail != nil && *em.Detail != "" {
		sb.WriteString(": " + *em.Detail)
	}

	if em.Errors == nil || len(*em.Errors) == 0 {
		return errors.New(sb.String())
	}

	for _, er := range *em.Errors {
		sb.WriteString("\n- ")

		if er.Location != nil && *er.Location != "" {
			sb.WriteString(*er.Location + ": ")
		}
		if er.Message != nil && *er.Message != "" {
			sb.WriteString(*er.Message + ": ")
		}
		if er.Value != nil {
			valueStr := fmt.Sprintf("(%v)", er.Value)
			sb.WriteString(valueStr)
		}
	}

	return errors.New(sb.String())
}
