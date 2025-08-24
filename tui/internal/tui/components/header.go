package components

import (
	"fmt"
	"os"

	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

func ShowHeader() string {
	societyName := os.Getenv("SOCIETY_NAME")
	eventName := os.Getenv("EVENT_NAME")
	headerContent := fmt.Sprintf("%s\n%s", societyName, eventName)
	return styles.HeaderStyle.Render(headerContent)
}
