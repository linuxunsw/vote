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

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	// TODO: set up propper logging to file w/ charm logs

	ssh := flag.Bool("ssh", false, "serve TUI over SSH")
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

	if *ssh {
		tui.SSH(*sshHost, *sshPort)
	} else {
		tui.Local()
	}

}
