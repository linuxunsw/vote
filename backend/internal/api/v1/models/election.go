package models

import "github.com/linuxunsw/vote/backend/internal/store"

type CreateElectionInput struct {
	Body struct {
		Name string `json:"name" doc:"Election name"`
	}
}

type CreateElectionResponse struct {
	Body CreateElectionResponseBody
}

type CreateElectionResponseBody struct {
	ElectionId string `json:"election_id" doc:"Election ID"`
}

type ElectionMemberListSetInput struct {
	ElectionId string `path:"election_id" doc:"Election ID"`

	Body struct {
		Zids []string `json:"zids" doc:"User zIDs" example:"z0000000"`
	}
}

type ElectionMemberListSetResponse struct {
}

type StateChangeEvent struct {
	NewState string `json:"new_state" enum:"CLOSED,NOMINATIONS_OPEN,NOMINATIONS_CLOSED,VOTING_OPEN,VOTING_CLOSED,RESULTS,END"`
}

type GetElectionStateResponse struct {
	Body GetElectionStateResponseBody
}

type GetElectionStateResponseBody struct {
	State          string `json:"state" enum:"NO_ELECTION,CLOSED,NOMINATIONS_OPEN,NOMINATIONS_CLOSED,VOTING_OPEN,VOTING_CLOSED,RESULTS,END"`
	ElectionId     string `json:"election_id,omitempty" doc:"Election ID. Only set if an election is running."`
	StateCreatedAt string `json:"state_created_at,omitempty" doc:"Timestamp when the election first entered this state, in RFC3339 format. Only Set if an election is running." example:"2023-01-01T00:00:00Z"`
}

type TransitionElectionStateInput struct {
	Body TransitionElectionStateBody
}

type TransitionElectionStateBody struct {
	State store.ElectionState `json:"state" enum:"CLOSED,NOMINATIONS_OPEN,NOMINATIONS_CLOSED,VOTING_OPEN,VOTING_CLOSED,RESULTS,END" doc:"State to transition to"`
}

type TransitionElectionStateResponse struct {
}
