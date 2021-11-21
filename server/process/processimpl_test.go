package process

import (
	"testing"
)

func TestStart(t *testing.T) {
	p := New("echo", "hello")

	checkStatus(t, p, "idle")

	p.Start()

	checkStatus(t, p, "started")

	if p.Id() == "" {
		t.Fatalf("Process PID should not be empty")
	}

	p.Wait()
	checkStatus(t, p, "exited")

	expected := "hello"

	outputs := p.Output()
	if outputs[0].Text != expected {
		t.Fatalf("Process output should be %s, got %s", expected, outputs[0].Text)
	}
}

func TestStop(t *testing.T) {
	p := New("sleep", "1")

	checkStatus(t, p, "idle")

	p.Start()

	checkStatus(t, p, "started")

	if p.Id() == "" {
		t.Fatalf("Process PID should not be empty")
	}

	_ = p.Stop()
	checkStatus(t, p, "killed")
}

func checkStatus(t *testing.T, p *ProcessImpl, status string) {
	if p.Status() != status {
		t.Fatalf("Process status should be \"%s\", got \"%s\"", status, p.Status())
	}
}
