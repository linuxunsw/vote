package models

import (
	"fmt"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

type NominationRequest struct {
	Body struct {
		CandidateName      string   `json:"candidate_name" minLength:"2" maxLength:"100" example:"John Doe"`
		ContactEmail       string   `json:"contact_email" format:"email" example:"john@example.com"`
		DiscordUsername    string   `json:"discord_username" maxLength:"32" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" minItems:"1" maxItems:"6" uniqueItems:"true" enum:"president,secretary,treasurer,arc_delegate,edi_officer,greivance_officer" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" required:"true" minLength:"50" maxLength:"2000" example:"I am running for president because..."`
		URL                *string  `json:"url" format:"uri" required:"false" example:"https://johndoe.com"`
	}
}

func (b *NominationRequest) Resolve(ctx huma.Context) []error {
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

type NominationResponse struct {
	Body struct {
		ID int64 `json:"id" example:"1"`
	}
}

type PublicNominationResponse struct {
	Body struct {
		ID                 int64    `json:"id" example:"1"`
		CandidateName      string   `json:"candidate_name" example:"John Doe"`
		DiscordUsername    string   `json:"discord_username,omitempty" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" example:"I am running for president because..."`
		URL                *string  `json:"url,omitempty" format:"uri" example:"https://johndoe.com"`
	}
}
