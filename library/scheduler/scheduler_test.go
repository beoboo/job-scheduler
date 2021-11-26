package scheduler

import (
	"github.com/beoboo/job-scheduler/library/job"
	"strings"
	"sync"
	"testing"
	"time"
)

type DummyJob struct {
	executable string
	args       []string
	logs       []string
	m          sync.Mutex
}

func (j *DummyJob) Id() string {
	return "123"
}

func (j *DummyJob) Start(listener job.OnJobListener) string {
	j.log("start")

	go func() {
		time.Sleep(50 * time.Millisecond)
		listener.OnFinishedJob(j)
	}()

	return j.Id()
}

func (j *DummyJob) log(msg string) {
	j.m.Lock()
	defer j.m.Unlock()

	j.logs = append(j.logs, msg)
}

func (j *DummyJob) Stop() error {
	j.log("stop")

	return nil
}

func (j *DummyJob) Wait() {
	j.log("wait")
}

func (j *DummyJob) Output() []job.OutputStream {
	if j.executable == "echo" {
		return []job.OutputStream{
			{Text: strings.Join(j.args, " ")},
		}
	}

	return nil
}

func (j *DummyJob) Error() []job.OutputStream {
	panic("implement me")
}

func (j *DummyJob) Status() string {
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
	s := New(&factory)

	_, _ = s.Start("echo", "world")

	if s.Size() != 1 {
		t.Fatalf("Job not started")
	}

	time.Sleep(100 * time.Millisecond)

	if s.Size() != 0 {
		t.Fatalf("Job not deleted")
	}
}

func TestStop(t *testing.T) {
	factory := DummyJobFactory{}
	s := New(&factory)

	id, _ := s.Start("sleep", "1")
	_, _ = s.Stop(id)

	if s.Size() != 0 {
		t.Fatalf("Job not stopped")
	}
}

func TestOutput(t *testing.T) {
	factory := DummyJobFactory{}
	s := New(&factory)

	expected := "hello"

	id, _ := s.Start("echo", expected)
	output, _ := s.Output(id)

	if output[0].Text != expected {
		t.Fatalf("Wrong output, want %s, got %s", output[0].Text, expected)
	}

	_, _ = s.Stop(id)
}
