package job

import (
	"github.com/beoboo/job-scheduler/library/assert"
	"github.com/beoboo/job-scheduler/library/status"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	j := NewJob("sleep", "0.1")

	assertStatus(t, j, status.IDLE)

	_, _ = j.Start()

	assertStatus(t, j, status.RUNNING)

	if j.Id() == "" {
		t.Fatalf("Job PID should not be empty")
	}

	j.Wait()

	assertStatus(t, j, status.EXITED)
}

func TestStop(t *testing.T) {
	j := NewJob("sleep", "1")

	assertStatus(t, j, status.IDLE)

	_, _ = j.Start()

	assertStatus(t, j, status.RUNNING)

	if j.Id() == "" {
		t.Fatalf("Job PID should not be empty")
	}

	_ = j.Stop()
	assertStatus(t, j, status.KILLED)
}

func TestOutput(t *testing.T) {
	j := NewJob("../../test.sh", "2 0.1")

	assertStatus(t, j, status.IDLE)

	expectedLines := []string{
		"Running for 2 times, sleeping for 0.1",
		"#1",
		"#2",
	}

	_, err := j.Start()
	if err != nil {
		t.Fatal(err)
	}

	assertOutput(t, j, expectedLines)

	j.Wait()

	assertStatus(t, j, status.EXITED)

	assertOutput(t, j, expectedLines)
}

// * Add resource control for CPU, Memory and Disk IO per job using cgroups.
// * Add resource isolation for using PID, mount, and networking namespaces.

func TestNamespaces(t *testing.T) {
	j := NewJob("sleep", "1")

	assertStatus(t, j, "idle")

	_, _ = j.Start()

	assertStatus(t, j, status.RUNNING)

	if j.Id() == "" {
		t.Fatalf("Job PID should not be empty")
	}

	_ = j.Stop()
	assertStatus(t, j, status.KILLED)
}

func assertStatus(t *testing.T, j *Job, expected string) {
	time.Sleep(10 * time.Millisecond)
	assert.AssertStatus(t, j.Status(), expected)
}

func assertOutput(t *testing.T, j *Job, expected []string) {
	assert.AssertOutput(t, j.Output(), expected)
}
