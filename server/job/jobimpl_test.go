package job

import (
	"testing"
)

func TestStart(t *testing.T) {
	j := New("echo", "hello")

	checkStatus(t, j, "idle")

	j.Start(nil)

	checkStatus(t, j, "started")

	if j.Id() == "" {
		t.Fatalf("Job PID should not be empty")
	}

	j.Wait()
	checkStatus(t, j, "exited")

	expected := "hello"

	outputs := j.Output()
	if outputs[0].Text != expected {
		t.Fatalf("Job output should be %s, got %s", expected, outputs[0].Text)
	}
}

func TestStop(t *testing.T) {
	j := New("sleep", "1")

	checkStatus(t, j, "idle")

	j.Start(nil)

	checkStatus(t, j, "started")

	if j.Id() == "" {
		t.Fatalf("Job PID should not be empty")
	}

	_ = j.Stop()
	checkStatus(t, j, "killed")
}

func checkStatus(t *testing.T, p *JobImpl, status string) {
	if p.Status() != status {
		t.Fatalf("Job status should be \"%s\", got \"%s\"", status, p.Status())
	}
}
