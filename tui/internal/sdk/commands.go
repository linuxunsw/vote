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
	"github.com/oapi-codegen/runtime/types"
)

var (
	ErrUnauthorised error = errors.New("your session has expired, please log in again")
	ErrForbidden    error = errors.New("you aren't authorised to vote! please check if you are a society member via rubric and contact a society executive for more help")
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
			return ServerErrMsg{
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

			return ServerErrMsg{
				RespID: respID,
				Error:  err,
			}
		}

		return GenerateOTPSuccessMsg{}
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
			return ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusForbidden {
			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrForbidden,
			}
		} else if resp.StatusCode() != http.StatusOK && resp.ApplicationproblemJSONDefault != nil {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return SubmitOTPSuccessMsg{}
	}

}

func GetElectionStateCmd(c *ClientWithIP) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := c.Client.GetElectionStateWithResponse(ctx, createIPRequestEditor(c.IP))
		if err != nil {
			return ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		// Add request ID as a reference code
		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusUnauthorized {
			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrUnauthorised,
			}
		}
		if resp.StatusCode() != http.StatusOK && resp.ApplicationproblemJSONDefault != nil {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return GetElectionStateSuccessMsg{
			State: string(resp.JSON200.State),
		}

	}
}

func GetBallotCmd(c *ClientWithIP) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := c.Client.GetBallotWithResponse(ctx, createIPRequestEditor(c.IP))
		if err != nil {
			return ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		// Add request ID as a reference code
		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusUnauthorized {
			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrUnauthorised,
			}
		}
		if resp.StatusCode() != http.StatusOK && resp.ApplicationproblemJSONDefault != nil {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return GetBallotSuccessMsg{
			Ballot: resp.JSON200,
		}

	}
}

// Sends request to submit a nomination, sends response back to root model as
// ServerErrMsg or a success message
func SubmitNominationCmd(c *ClientWithIP, data Submission) tea.Cmd {
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
			return ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		// Add request ID as a reference code
		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusUnauthorized {
			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrUnauthorised,
			}
		}
		if resp.StatusCode() != http.StatusOK && resp.ApplicationproblemJSONDefault != nil {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return SubmitNominationSuccessMsg{
			RefCode: respID,
		}

	}
}

func SubmitVoteCmd(c *ClientWithIP, data map[string]string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		body := SubmitVoteBody{
			Positions: data,
		}

		resp, err := c.Client.SubmitVoteWithResponse(ctx, body, createIPRequestEditor(c.IP))
		if err != nil {
			return ServerErrMsg{
				RespID: "",
				Error:  err,
			}
		}

		// Add request ID as a reference code
		respID := resp.HTTPResponse.Header.Get("X-Request-ID")
		if resp.StatusCode() == http.StatusUnauthorized {
			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      ErrUnauthorised,
			}
		}
		if resp.StatusCode() != http.StatusNoContent && resp.ApplicationproblemJSONDefault != nil {
			err := buildError(*resp.ApplicationproblemJSONDefault)

			return ServerErrMsg{
				StatusCode: resp.StatusCode(),
				RespID:     respID,
				Error:      err,
			}
		}

		// Build success message
		return SubmitVoteSuccessMsg{
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
