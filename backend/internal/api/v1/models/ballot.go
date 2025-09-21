package models

type PublicBallot struct {
	ElectionID string                        `json:"election_id" example:"1" doc:"Election ID"`
	Candidates map[string][]PublicNomination `json:"candidates" doc:"Map of executive role to list of candidates running for that role" example:"{\"president\": [], \"secretary\": []}"`
	HasVoted   bool                          `json:"has_voted" doc:"Whether the current user has already voted in this election"`
}

type GetBallotResponse struct {
	Body PublicBallot
}
