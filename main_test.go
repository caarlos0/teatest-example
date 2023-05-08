package main

import (
	"bytes"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestFullOutput(t *testing.T) {
	m := initialModel(time.Second)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))
	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*2)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestFinalModel(t *testing.T) {
	tm := teatest.NewTestModel(t, initialModel(time.Second), teatest.WithInitialTermSize(300, 100))
	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second*2))
	m, ok := fm.(model)
	if !ok {
		t.Fatalf("final model have the wrong type: %T", fm)
	}
	if m.duration != time.Second {
		t.Errorf("m.duration != 1s: %s", m.duration)
	}
	if m.start.After(time.Now().Add(-1 * time.Second)) {
		t.Errorf("m.start should be more than 1 second ago: %s", m.start)
	}
}

func TestOuput(t *testing.T) {
	tm := teatest.NewTestModel(t, initialModel(10*time.Second), teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("sleeping 8s"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
