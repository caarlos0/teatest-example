package main

import (
	"io"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest"
)

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
