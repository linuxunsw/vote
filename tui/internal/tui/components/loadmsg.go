package components

import "github.com/linuxunsw/vote/tui/internal/tui/pages"

var messages = map[pages.PageID]string{
	pages.PageAuth:     "requesting OTP",
	pages.PageAuthCode: "verifying OTP",
	pages.PageForm:     "submitting form",
	pages.PageSubmit:   "loading",
}

func GetPageMsg(id pages.PageID) string {
	return messages[id]
}
