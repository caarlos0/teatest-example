package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	spinner  spinner.Model
	quitting bool
	duration time.Duration
	start    time.Time
	err      error
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func initialModel(duration time.Duration) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:  s,
		duration: duration,
		start:    time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) timeLeft() time.Duration {
	return m.duration - time.Since(m.start)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		if m.timeLeft() <= 0 {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf(
		"\n\n   %s sleeping %s... %s\n\n",
		m.spinner.View(),
		m.timeLeft().Round(time.Second),
		quitKeys.Help().Desc,
	)
	if m.quitting {
		return str + "\n"
	}
	return str
}

func main() {
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel(duration))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
