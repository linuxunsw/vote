package store

import (
	"context"
	"time"
)

type Nomination struct {
	ElectionID         string   `db:"election_id" json:"election_id" example:"1"`
	CandidateZID       string   `db:"candidate_zid" json:"candidate_zid" example:"z1234567"`
	CandidateName      string   `db:"candidate_name" json:"candidate_name" example:"John Doe"`
	ContactEmail       string   `db:"contact_email" json:"contact_email" example:"john@example.com"`
	DiscordUsername    string   `db:"discord_username" json:"discord_username" example:"johndoe"`
	ExecutiveRoles     []string `db:"executive_roles" json:"executive_roles" example:"[\"president\", \"secretary\"]"`
	CandidateStatement string   `db:"candidate_statement" json:"candidate_statement" example:"I am running for president because..."`
	URL                *string  `db:"url" json:"url,omitempty" example:"https://johndoe.com"`
	CreatedAt          time.Time   `db:"created_at" json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	UpdatedAt          time.Time   `db:"updated_at" json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

func (nom Nomination) IsRunningFor(role string) bool {
	for _, r := range nom.ExecutiveRoles {
		if r == role {
			return true
		}
	}
	return false
}

type SubmitNomination struct {
	CandidateName      string   `db:"candidate_name" json:"candidate_name" minLength:"2" maxLength:"100" example:"John Doe"`
	ContactEmail       string   `db:"contact_email" json:"contact_email" format:"email" example:"john@example.com"`
	DiscordUsername    string   `db:"discord_username" json:"discord_username" maxLength:"32" example:"johndoe"`
	ExecutiveRoles     []string `db:"executive_roles" json:"executive_roles" minItems:"1" maxItems:"6" uniqueItems:"true" enum:"president,secretary,treasurer,arc_delegate,edi_officer,grievance_officer" example:"[\"president\", \"secretary\"]"`
	CandidateStatement string   `db:"candidate_statement" json:"candidate_statement" required:"true" minLength:"50" maxLength:"2000" example:"I am running for president because..."`
	URL                *string  `db:"url" json:"url,omitempty" format:"uri" required:"false" example:"https://johndoe.com"`
}

type NominationStore interface {
	// Create or replace a new nomination.
	SubmitOrReplaceNomination(ctx context.Context, electionID string, candidateZid string, submission SubmitNomination) error

	// Get a nomination by election ID and candidate zID. Returns nil if not found.
	GetNomination(ctx context.Context, electionID string, candidateZid string) (*Nomination, error)

	// Delete a nomination by election ID and candidate zID. Does nothing if the nomination doesn't exist.
	TryDeleteNomination(ctx context.Context, electionID string, candidateZid string) error
}
