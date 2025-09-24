package pages

type PageID string

const (
	Auth             PageID = "auth"
	AuthCode         PageID = "authCode"
	NominationForm   PageID = "nominationForm"
	NominationSubmit PageID = "nominationSubmit"
	VotingForm       PageID = "votingForm"
	VotingSubmit     PageID = "votingSubmit"
	Closed           PageID = "closed"
)
