package scheduler

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

func (p *DummyProcess) Id() string {
	return "123"
}

func (p *DummyProcess) Start() string {
	p.log("start")

	return p.Id()
}

func (p *DummyProcess) log(msg string) {
	p.logs = append(p.logs, msg)
}

func (p *DummyProcess) Stop() error {
	p.log("stop")

	return nil
}

func (p *DummyProcess) Wait() {
	p.log("wait")
}

func (p *DummyProcess) Output() []process.OutputStream {
	if p.executable == "echo" {
		return []process.OutputStream{
			{Text: strings.Join(p.args, " ")},
		}
	}

	return nil
}

func (p *DummyProcess) Error() []process.OutputStream {
	panic("implement me")
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
	scheduler := New(&factory)

	pid, _ := scheduler.Start("sleep", "1")

	if len(scheduler.Processes) != 1 {
		t.Fatalf("Process not started")
	}

	_, _ = scheduler.Stop(pid)
}

func TestStop(t *testing.T) {
	factory := DummyProcessFactory{}
	scheduler := New(&factory)

	pid, _ := scheduler.Start("sleep", "1")
	_, _ = scheduler.Stop(pid)

	if len(scheduler.Processes) != 0 {
		t.Fatalf("Process not stopped")
	}
}

func TestOutput(t *testing.T) {
	factory := DummyProcessFactory{}
	scheduler := New(&factory)

	expected := "hello"

	pid, _ := scheduler.Start("echo", expected)
	output, _ := scheduler.Output(pid)

	if output[0].Text != expected {
		t.Fatalf("Wrong output, want %s, got %s", output[0].Text, expected)
	}

	_, _ = scheduler.Stop(pid)
}
