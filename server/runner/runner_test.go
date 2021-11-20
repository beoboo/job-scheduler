package runner

import (
	"github.com/beoboo/job-worker-service/server/process"
	"strings"
	"testing"
)

type DummyProcess struct {
	executable string
	args       []string
	logs       []string
}

func (p *DummyProcess) log(msg string) {
	p.logs = append(p.logs, msg)
}

func (p *DummyProcess) Start() int {
	p.log("start")

	return 0
}

func (p *DummyProcess) Stop() error {
	p.log("stop")

	return nil
}

func (p *DummyProcess) Wait() {
	p.log("wait")
}

func (p *DummyProcess) Output() string {
	if p.executable == "echo" {
		return strings.Join(p.args, " ")
	}

	return ""
}

func (p *DummyProcess) Error() string {
	return ""
}

func (p *DummyProcess) Status() string {
	return ""
}

type DummyProcessFactory struct {
}

func (f *DummyProcessFactory) Create(executable string, args ...string) process.Process {
	return &DummyProcess{
		executable: executable,
		args:       args,
	}
}

func TestStart(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	pid, _ := runner.Start("sleep", "1")

	if len(runner.Processes) != 1 {
		t.Fatalf("Process not started")
	}

	runner.Stop(pid)
}

func TestStop(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	pid, _ := runner.Start("sleep", "1")
	_, _ = runner.Stop(pid)

	if len(runner.Processes) != 0 {
		t.Fatalf("Process not stopped")
	}
}

func TestOutput(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	expected := "hello"

	pid, _ := runner.Start("echo", expected)
	output := runner.Output(pid)

	if output != expected {
		t.Fatalf("Wrong output, want %s, got %s", output, expected)
	}

	_, _ = runner.Stop(pid)
}
