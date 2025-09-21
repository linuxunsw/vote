package models

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type SubmitVoteInput struct {
	Body SubmitVoteBody
}

// president,secretary,treasurer,arc_delegate,edi_officer,grievance_officer

var validPositions = map[string]struct{}{
	"president":         {},
	"secretary":         {},
	"treasurer":         {},
	"arc_delegate":      {},
	"edi_officer":       {},
	"grievance_officer": {},
}

type SubmitVoteBody struct {
	Positions map[string]string `json:"positions" example:"{\"president\":\"z1234567\",\"secretary\":\"z7654321\"}"`
}

func (b *SubmitVoteInput) Resolve(ctx huma.Context) []error {
	var errors []error

	for position := range b.Body.Positions {
		if _, ok := validPositions[position]; !ok {
			errors = append(errors, &huma.ErrorDetail{
				Message:  "invalid position",
				Location: "body.positions",
				Value:    position,
			})
		}
	}

	return errors
}

type GetVoteResponse struct {
	Body Vote
}

type Vote struct {
	Positions map[string]string `json:"positions" doc:"A map from categories to public nomination IDs. Find these by accessing your ballot." example:"{\"president\":\"01996ae6-31e5-7bc6-bac4-399ffc8c80de\",\"secretary\":\"01996ae6-31e5-7bc6-bac4-399ffc8c80de\"}"`
	CreatedAt time.Time         `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time         `json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
}
