package tui

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/linuxunsw/vote/tui/internal/tui/root"
)

// Runs the program locally
func Local() {
	p := tea.NewProgram(root.New(""), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("TUI error:", err)
	}
}

// Runs the program served via wish (SSH)
func SSH(host string, port string) {
	log := log.New(os.Stderr)
	log.SetReportTimestamp(true)
	log.SetPrefix("server")

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.StructuredMiddlewareWithLogger(log, log.GetLevel()),
		),

		ssh.AllocatePty(),
	)
	if err != nil {
		log.Fatal("could not start server,", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server...", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	// pty, _, _ := s.Pty()
	m := root.New(s.User())
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
