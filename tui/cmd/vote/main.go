package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"

	"github.com/linuxunsw/vote/tui/internal/tui/root"
)

func main() {
	// TODO: set up propper logging to file w/ charm logs
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(root.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("fatal:", err)
		os.Exit(1)
	}
}
