package models

import (
	"regexp"
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

var zIDRegex = regexp.MustCompile(`^z[0-9]{7}$`)

type SubmitVoteBody struct {
	Positions map[string]string `json:"positions" example:"{\"president\":\"z1234567\",\"secretary\":\"z7654321\"}"`
}

func (b *SubmitVoteInput) Resolve(ctx huma.Context) []error {
	var errors []error

	for position, zid := range b.Body.Positions {
		if _, ok := validPositions[position]; !ok {
			errors = append(errors, &huma.ErrorDetail{
				Message:  "invalid position",
				Location: "body.positions",
				Value:    position,
			})
		}
		if !zIDRegex.MatchString(zid) {
			errors = append(errors, &huma.ErrorDetail{
				Message:  "invalid zid format",
				Location: "body.positions." + position,
				Value:    zid,
			})
		}
	}

	return errors
}

type GetVoteResponse struct {
	Body Vote
}

type Vote struct {
	Positions map[string]string `json:"positions" example:"{\"president\":\"z1234567\",\"secretary\":\"z7654321\"}"`
	CreatedAt time.Time            `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time            `json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
}
