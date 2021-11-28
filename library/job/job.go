package job

import (
	"bufio"
	"fmt"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/status"
	"github.com/beoboo/job-scheduler/library/stream"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// Job wraps the execution of a process, capturing its stdout and stderr streams,
// and providing the process status
// TODO: it should also report the exit status code
type Job struct {
	id     string
	cmd    *exec.Cmd
	output *stream.Stream
	done   chan bool
	status string
	m      sync.Mutex
}

// New creates a new Job
func New() *Job {
	p := &Job{
		id:     generateRandomId(),
		done:   make(chan bool),
		output: stream.New(),
		status: status.IDLE,
	}

	return p
}

// TODO: This could be a UUID or some other random generated value
func generateRandomId() string {
	now := time.Now()
	return fmt.Sprintf("%d", now.UnixNano())
}

// Id returns the random job ID
func (j *Job) Id() string {
	return j.id
}

// StartChild starts the execution of a process through a parent/child mechanism
func (j *Job) StartChild(executable string, args ...string) error {
	args = append([]string{
		"child",    // Main subcommand
		j.id,       // The job ID
		executable, // The original executable
	}, args...)

	return j.Start("/proc/self/exe", args...)
}

// Start starts the execution of a process, capturing its output
func (j *Job) Start(executable string, args ...string) error {
	cmd := exec.Command(executable, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC,
	}

	j.cmd = cmd

	errCh := make(chan error, 1)
	log.Debugf("Starting: %s\n", j.cmd.Path)

	go func() {
		j.run(errCh)
	}()

	// Waits for the job to be started successfully
	err := <-errCh

	return err
}

// Stop stops a running process
func (j *Job) Stop() error {
	j.lock("Stop")
	err := j.cmd.Process.Kill()
	j.unlock("Stop")

	if err != nil {
		return fmt.Errorf("cannot kill job %d: (%s)", j.pid(), err)
	}

	j.updateStatus(status.KILLED)
	return nil
}

// Output returns the stream of captured stdout/stderr of the running process.
func (j *Job) Output() *stream.Stream {
	j.lock("Output")
	defer j.unlock("Output")

	j.output.Rewind()

	return j.output
}

// Status returns the current status of the Job
func (j *Job) Status() string {
	j.lock("Status")
	defer j.unlock("Status")

	return j.status
}

// Wait blocks until the process is completed
func (j *Job) Wait() {
	<-j.done
}

func (j *Job) run(started chan error) {
	log.Debugf("Running: %s\n", j.cmd.Path)

	stdout, _ := j.cmd.StdoutPipe()
	stderr, _ := j.cmd.StderrPipe()

	err := j.cmd.Start()

	started <- err
	if err != nil {
		return
	}

	j.updateStatus(status.RUNNING)

	j.pipe("output", stdout)
	j.pipe("error", stderr)

	_ = j.cmd.Wait()

	j.updateStatus(status.EXITED)

	j.done <- true
}

func (j *Job) pipe(channel string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		go func() {
			j.lock("pipe")
			defer j.unlock("pipe")

			log.Debugf("[%s] %s\n", channel, text)
			_ = j.output.Write(stream.Line{
				Channel: channel,
				Time:    time.Now(),
				Text:    text,
			})
		}()
	}
}

func (j *Job) pid() int {
	return j.cmd.Process.Pid
}

func (j *Job) updateStatus(st string) {
	j.lock("updateStatus")
	defer j.unlock("updateStatus")

	switch j.status {
	case status.IDLE:
		if st == status.RUNNING {
			j.status = st
		}
	case status.RUNNING:
		if st == status.EXITED || st == status.KILLED {
			j.status = st
			j.output.Close()
		}
	default:
		return
	}
}

func (j *Job) lock(id string) {
	log.Tracef("Job locking %s\n", id)
	j.m.Lock()
}

func (j *Job) unlock(id string) {
	log.Tracef("Job unlocking %s\n", id)
	j.m.Unlock()
}
