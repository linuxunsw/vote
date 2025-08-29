package components

import (
	"fmt"

	"github.com/linuxunsw/vote/tui/internal/tui/styles"
	"github.com/spf13/viper"
)

func ShowHeader(width int) string {
	societyName := viper.GetString("society")
	eventName := viper.GetString("event")
	headerContent := fmt.Sprintf("%s\n%s", societyName, eventName)
	return styles.Header(width).Render(headerContent)
}
