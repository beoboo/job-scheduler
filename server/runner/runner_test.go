package runner

import (
	"github.com/beoboo/job-worker-service/server/process"
	"strings"
	"testing"
)

type DummyProcess struct {
	command string
	args    []string
	logs    []string
}

func (p *DummyProcess) log(msg string) {
	p.logs = append(p.logs, msg)
}

func (p *DummyProcess) Start() int {
	p.log("start")
	return 0
}

func (p *DummyProcess) Stop() {
	p.log("stop")
}

func (p *DummyProcess) Wait() {
	p.log("wait")
}

func (p *DummyProcess) Output() string {
	if p.command == "echo" {
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

func (f *DummyProcessFactory) Create(command string, args ...string) process.Process {
	return &DummyProcess{
		command: command,
		args:    args,
	}
}

func TestStart(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	pid := runner.Start("sleep", "1")

	if len(runner.Processes) != 1 {
		t.Fatalf("Process not started")
	}

	runner.Stop(pid)
}

func TestStop(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	pid := runner.Start("sleep", "1")
	runner.Stop(pid)

	if len(runner.Processes) != 0 {
		t.Fatalf("Process not stopped")
	}
}

func TestOutput(t *testing.T) {
	factory := DummyProcessFactory{}
	runner := New(&factory)

	expected := "hello"

	pid := runner.Start("echo", expected)
	output := runner.Output(pid)

	if output != expected {
		t.Fatalf("Wrong output, want %s, got %s", output, expected)
	}

	runner.Stop(pid)
}
