package models

type ElectionState int

const (
	StateClosed ElectionState = iota
	StateNominationsOpen
	StateNominationsClosed
	StateVotingOpen
	StateVotingClosed
	StateResults
)

var stateName = map[ElectionState]string{
	StateClosed:            "CLOSED",
	StateNominationsOpen:   "NOMINATIONS_OPEN",
	StateNominationsClosed: "NOMINATIONS_CLOSED",
	StateVotingOpen:        "VOTING_OPEN",
	StateVotingClosed:      "VOTING_CLOSED",
	StateResults:           "RESULTS",
}

func (ss ElectionState) String() string {
	return stateName[ss]
}

type StateChangeEvent struct {
	NewState string `json:"new_state" enum:"CLOSED,NOMINATIONS_OPEN,NOMINATIONS_CLOSED,VOTING_OPEN,VOTING_CLOSED,RESULTS"`
}
