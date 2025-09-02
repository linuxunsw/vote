package models

import "net/http"

type NominationStatusInput struct {
	Session http.Cookie `cookie:"session"`
	ID      string      `path:"id"`
}

type NominationStatusOutput struct {
	HasNominated bool       `json:"has_nominated" example:"true"`
	Nomination   Nomination `json:"nomination" required:"false"`
}

type NominationStatusResponse struct {
	Body NominationStatusOutput
}

type NominationSubmissionInput struct {
	Session http.Cookie `cookie:"session"`
	ID      string      `path:"id"`
	Body    Nomination  `json:"nomination"`
}

type Nomination struct {
	Name         string   `json:"name" doc:"The user's full name" example:"John Tux"`
	ContactEmail string   `json:"contact_email" doc:"The user's preferred contact email" example:"johntux@gmail.com"`
	Statement    string   `json:"statement" doc:"The user's candidate statement" example:""`
	Roles        []string `json:"roles" doc:"The roles the user is nominating for" enum:"president,secretary,treasurer,arc_delegate,edi_officer,grievance_officer" example:"arc_delegate,edi_officer"`
	Discord      string   `json:"discord" doc:"The user's Discord username" example:"johntux"`
	URL          string   `json:"url" doc:"An optional URL provided by the user" required:"false" format:"email" example:"linuxunsw.org"`
}

type NominationSubmissionResponse struct {
}

type NominationDeleteInput struct {
	Session http.Cookie `cookie:"session"`
	ID      string      `path:"id"`
}

type NominationDeleteResponse struct {
}
