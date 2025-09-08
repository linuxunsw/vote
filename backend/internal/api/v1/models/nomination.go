package models

import (
	"net/url"
)

type NominationRequest struct {
	Body struct {
		CandidateName      string   `json:"candidate_name" minLength:"2" maxLength:"100" example:"John Doe"`
		ContactEmail       string   `json:"contact_email" format:"email" example:"john@example.com"`
		DiscordUsername    string   `json:"discord_username" maxLength:"32" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" minItems:"1" maxItems:"6" uniqueItems:"true" enum:"president,secretary,treasurer,arc_delegate,edi_officer,greivance_officer" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" required:"true" minLength:"50" maxLength:"2000" example:"I am running for president because..."`
		URL                *string  `json:"url" required:"false" format:"uri" example:"https://johndoe.com"`
	}
}

type NominationResponse struct {
	Body struct {
		ID                 int64    `json:"id" example:"1"`
		ElectionID         int64    `json:"election_id" example:"1"`
		CandidateZID       string   `json:"candidate_zid" example:"z1234567"`
		CandidateName      string   `json:"candidate_name" example:"John Doe"`
		ContactEmail       string   `json:"contact_email" example:"john@example.com"`
		DiscordUsername    string   `json:"discord_username" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" example:"I am running for president because..."`
		URL                *url.URL `json:"url,omitempty" example:"https://johndoe.com"`
		CreatedAt          string   `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
		UpdatedAt          string   `json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	}
}

type PublicNominationResponse struct {
	Body struct {
		ID                 int64    `json:"id" example:"1"`
		CandidateName      string   `json:"candidate_name" example:"John Doe"`
		DiscordUsername    string   `json:"discord_username,omitempty" example:"johndoe"`
		ExecutiveRoles     []string `json:"executive_roles" example:"[\"president\", \"secretary\"]"`
		CandidateStatement string   `json:"candidate_statement" example:"I am running for president because..."`
		URL                *url.URL `json:"url,omitempty" format:"uri" example:"https://johndoe.com"`
	}
}
