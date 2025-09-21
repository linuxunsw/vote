package models

import (
	"fmt"
	"net/url"
	"time"

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

type SubmitNominationResponse struct {
	Body SubmitNominationResponseBody
}

type SubmitNominationResponseBody struct {
	NominationID string `json:"nomination_id" example:"abc123" doc:"Public nomination ID"`
}

// Public nominations don't have zID or contact email
type PublicNomination struct {
	NominationId       string    `json:"nomination_id"`
	ElectionID         string    `json:"election_id"`
	CandidateName      string    `json:"candidate_name" example:"John Doe"`
	DiscordUsername    string    `json:"discord_username" example:"johndoe"`
	ExecutiveRoles     []string  `json:"executive_roles" example:"[\"president\", \"secretary\"]"`
	CandidateStatement string    `json:"candidate_statement" example:"I am running for president because..."`
	URL                *string   `json:"url,omitempty" example:"https://johndoe.com"`
	CreatedAt          time.Time `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	UpdatedAt          time.Time `json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

func FromStoreNomination(nom store.Nomination) PublicNomination {
	return PublicNomination{
		NominationId:       nom.NominationId,
		ElectionID:         nom.ElectionID,
		CandidateName:      nom.CandidateName,
		DiscordUsername:    nom.DiscordUsername,
		ExecutiveRoles:     nom.ExecutiveRoles,
		CandidateStatement: nom.CandidateStatement,
		URL:                nom.URL,
		CreatedAt:          nom.CreatedAt,
		UpdatedAt:          nom.UpdatedAt,
	}
}

type GetNominationResponse struct {
	Body store.Nomination
}

type GetPublicNominationInput struct {
	NominationId string `path:"nomination_id" doc:"Public nomination ID"`
}

type GetPublicNominationResponse struct {
	Body PublicNomination
}
