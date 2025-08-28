package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"

	"github.com/linuxunsw/vote/tui/internal/tui"
)

// Default host and port
const (
	host = "0.0.0.0"
	port = "2222"
)

func main() {

	local := flag.Bool("local", false, "run locally (no SSH server)")
	sshHost := flag.String("host", host, "host")
	sshPort := flag.String("port", port, "port")
	flag.Parse()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing log file: %v", err)
		}
	}()

	if *local != true {
		tui.SSH(*sshHost, *sshPort)
	} else {
		tui.Local()
	}

}
