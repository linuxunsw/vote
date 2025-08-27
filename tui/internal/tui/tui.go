package tui

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/comment"
	"github.com/charmbracelet/wish/logging"
	"github.com/linuxunsw/vote/tui/internal/tui/root"
)

func Local() {
	p := tea.NewProgram(root.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("TUI error:", err)
	}
}

func SSH(address string) {
	// Make a new server
	server, err := wish.NewServer(
		wish.WithAddress(address),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(handler),
			logging.Middleware(),
			comment.Middleware("session ended."),
		),
	)
	if err != nil {
		log.Fatal("could not start server", err)
	}

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("shutting down SSH server…")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatal("shutdown error", err)
		}
		os.Exit(0)
	}()

	log.Println("SSH server listening", address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Fatal("server error", "err", err)
	}
}

func handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	_, winCh, _ := s.Pty()

	m := root.New()

	go func() {
		for win := range winCh {
			m.Update(tea.WindowSizeMsg{Width: win.Width, Height: win.Height})
		}
	}()

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
