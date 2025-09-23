package components

import "github.com/linuxunsw/vote/tui/internal/tui/pages"

var messages = map[pages.PageID]string{
	pages.Auth:             "requesting OTP",
	pages.AuthCode:         "verifying OTP",
	pages.NominationForm:   "submitting nomination",
	pages.VotingForm:       "submitting vote",
	pages.NominationSubmit: "loading",
	pages.VotingSubmit:     "loading",
}

func GetPageMsg(id pages.PageID) string {
	return messages[id]
}
