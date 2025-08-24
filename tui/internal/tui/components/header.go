package components

import (
	"fmt"

	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

// TODO: make header title(s) not hardcoded
const societyName = "$ linux society"
const eventName = "2025 agm nomination form"

func ShowHeader() string {
	headerContent := fmt.Sprintf("%s\n%s", societyName, eventName)
	return styles.HeaderStyle.Render(headerContent)
}
