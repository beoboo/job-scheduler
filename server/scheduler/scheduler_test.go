package scheduler

import (
	"github.com/beoboo/job-worker-service/server/job"
	"strings"
	"testing"
	"time"
)

type DummyJob struct {
	executable string
	args       []string
	logs       []string
}

func (p *DummyJob) Id() string {
	return "123"
}

func (p *DummyJob) Start(listener job.OnJobListener) string {
	p.log("start")

	go func() {
		time.Sleep(50 * time.Millisecond)
		listener.OnFinishedJob(p)
	}()

	return p.Id()
}

func (p *DummyJob) log(msg string) {
	p.logs = append(p.logs, msg)
}

func (p *DummyJob) Stop() error {
	p.log("stop")

	return nil
}

func (p *DummyJob) Wait() {
	p.log("wait")
}

func (p *DummyJob) Output() []job.OutputStream {
	if p.executable == "echo" {
		return []job.OutputStream{
			{Text: strings.Join(p.args, " ")},
		}
	}

	return nil
}

func (p *DummyJob) Error() []job.OutputStream {
	panic("implement me")
}

func (p *DummyJob) Status() string {
	return ""
}

type DummyJobFactory struct {
}

func (f *DummyJobFactory) Create(executable string, args ...string) job.Job {
	return &DummyJob{
		executable: executable,
		args:       args,
	}
}

func TestStart(t *testing.T) {
	factory := DummyJobFactory{}
	scheduler := New(&factory)

	_, _ = scheduler.Start("echo", "world")

	if len(scheduler.jobs) != 1 {
		t.Fatalf("Job not started")
	}

	time.Sleep(100 * time.Millisecond)

	if len(scheduler.jobs) != 0 {
		t.Fatalf("Job not deleted")
	}
}

/*
func TestStop(t *testing.T) {
	factory := DummyJobFactory{}
	scheduler := New(&factory)

	id, _ := scheduler.Start("sleep", "1")
	_, _ = scheduler.Stop(id)

	if len(scheduler.jobs) != 0 {
		t.Fatalf("Job not stopped")
	}
}

func TestOutput(t *testing.T) {
	factory := DummyJobFactory{}
	scheduler := New(&factory)

	expected := "hello"

	id, _ := scheduler.Start("echo", expected)
	output, _ := scheduler.Output(id)

	if output[0].Text != expected {
		t.Fatalf("Wrong output, want %s, got %s", output[0].Text, expected)
	}

	_, _ = scheduler.Stop(id)
}
*/
