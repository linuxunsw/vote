package models

import (
	"fmt"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type SubmitNominationRequest struct {
	Body store.SubmitNomination
}

func (b *SubmitNominationRequest) Resolve(ctx huma.Context) []error {
	if b.Body.URL == nil {
		return nil
	}

	urlStr := *b.Body.URL
	parsed, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return []error{
			&huma.ErrorDetail{
				Message:  fmt.Sprintf("invalid URL: %q", urlStr),
				Location: "body.url",
				Value:    urlStr,
			},
		}
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return []error{&huma.ErrorDetail{
			Message:  "URL must use http or https scheme",
			Location: "body.url",
			Value:    urlStr,
		}}
	}

	return nil
}

type GetNominationResponse struct {
	Body store.Nomination
}

type PublicNominationResponse struct {
	Body struct {
		CandidateName      string   `json:"candidate_name" example:"John Doe"`
		DiscordUsername    string   `json:"discord_username,omitempty" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" example:"I am running for president because..."`
		URL                *string  `json:"url,omitempty" format:"uri" example:"https://johndoe.com"`
	}
}
