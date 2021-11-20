package process

import "testing"

func TestStart(t *testing.T) {
	p := New("echo", "hello")

	p.Start()

	if p.Pid == 0 {
		t.Fatalf("Process PID should not be 0")
	}

	p.Wait()

	expected := "hello"
	if p.Output() != expected {
		t.Fatalf("Process output should be %s, got %s", expected, p.Output())
	}
}