package models

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
		Zids []string `json:"zids" doc:"User zIDs"`
	}
}

type ElectionMemberListSetResponse struct {
}
