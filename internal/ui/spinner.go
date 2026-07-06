package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

const defaultSpinnerSpeed = 100 * time.Millisecond

type Spinner struct {
	writer     io.Writer
	model      spinner.Model
	speed      time.Duration
	isTerminal bool

	mu      sync.Mutex
	prefix  string
	suffix  string
	done    chan struct{}
	wg      sync.WaitGroup
	started bool
}

func NewSpinner(writer io.Writer, speed time.Duration) *Spinner {
	if writer == nil {
		writer = os.Stdout
	}
	if speed <= 0 {
		speed = defaultSpinnerSpeed
	}

	model := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(lipgloss.NewStyle()),
	)

	s := &Spinner{
		writer: writer,
		model:  model,
		speed:  speed,
	}

	if file, ok := writer.(*os.File); ok {
		s.isTerminal = term.IsTerminal(file.Fd())
	}

	return s
}

func (s *Spinner) SetPrefix(prefix string) {
	s.mu.Lock()
	s.prefix = prefix
	s.mu.Unlock()
}

func (s *Spinner) SetSuffix(suffix string) {
	s.mu.Lock()
	s.suffix = suffix
	s.mu.Unlock()
}

func (s *Spinner) Start() {
	if !s.isTerminal {
		return
	}

	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return
	}
	s.started = true
	s.done = make(chan struct{})
	s.mu.Unlock()

	s.writeFrame()

	s.wg.Add(1)
	go s.run()
}

func (s *Spinner) Stop() {
	if !s.isTerminal {
		return
	}

	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	s.started = false
	close(s.done)
	s.mu.Unlock()

	s.wg.Wait()

	fmt.Fprint(s.writer, "\r\033[K")
}

func (s *Spinner) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.speed)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.writeFrame()
		}
	}
}

func (s *Spinner) writeFrame() {
	s.mu.Lock()
	prefix := s.prefix
	suffix := s.suffix
	s.mu.Unlock()

	frame := s.model.View()
	updated, _ := s.model.Update(spinner.TickMsg{Time: time.Now(), ID: s.model.ID()})
	s.model = updated

	fmt.Fprintf(s.writer, "\r\033[K%s%s %s", prefix, frame, suffix)
}
