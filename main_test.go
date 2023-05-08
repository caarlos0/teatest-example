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
